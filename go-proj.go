package docs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	// "github.com/santhosh-tekuri/jsonschema"
)


// Core types
type DocType string
type DocID string

type Document struct {
	ID    DocID          `json:"id,omitempty"`
	Type  DocType        `json:"type"`
	Title string         `json:"title"`
	Body  string         `json:"body"`
	Rev   string         `json:"_rev,omitempty"`
	Extra map[string]any `json:"_extra,omitempty"` // for extensibility
}

type Schema struct {
	ID   DocID          `json:"_id,omitempty"`
	Rev  string         `json:"_rev,omitempty"`
	Type DocType        `json:"type"`
	JSON json.RawMessage `json:"-"` // full JSON Schema
	Raw  map[string]any `json:"-"` // parsed
}

// Validation error
type ValidationError struct {
	Message string `json:"message"`
	Errors  []any  `json:"errors,omitempty"`
}

func (err *ValidationError) Error() error {
	return errors.New(err.Message)
}

// CouchDB client wrapper (dense, minimal)
type CouchDB struct {
	baseURL string
	client  *http.Client
	auth    string
	mu      sync.RWMutex
}

func NewCouchDB(host, user, pass string) *CouchDB {
	auth := ""
	if user != "" || pass != "" {
		auth = "Basic " + basicAuth(user, pass)
	}
	return &CouchDB{
		baseURL: strings.TrimSuffix(host, "/"),
		client:  &http.Client{Timeout: 10 * time.Second},
		auth:    auth,
	}
}

func basicAuth(user, pass string) string {
	return fmt.Sprintf("%s:%s", user, pass) // base64 in header
}

func (c *CouchDB) req(method, path string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}
	if c.auth != "" {
		req.Header.Set("Authorization", c.auth)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.client.Do(req)
}

// Core DB ops
func (c *CouchDB) EnsureDB(name string) error {
	_, err := c.req("PUT", "/"+name, nil)
	if err != nil && !strings.Contains(err.Error(), "412") { // already exists
		return err
	}
	return nil
}

type ViewRow struct {
	ID    string         `json:"id"`
	Key   any            `json:"key"`
	Value json.RawMessage `json:"value"`
	Doc   json.RawMessage `json:"doc,omitempty"`
}

func (c *CouchDB) QueryView(db, design, view string, params map[string]string) ([]ViewRow, error) {
	q := ""
	if len(params) > 0 {
		parts := []string{}
		for k, v := range params {
			parts = append(parts, k+"="+v)
		}
		q = "?" + strings.Join(parts, "&")
	}
	resp, err := c.req("GET", fmt.Sprintf("/%s/_design/%s/_view/%s%s", db, design, view, q), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("view query failed: %d", resp.StatusCode)
	}
	var res struct {
		Rows []ViewRow `json:"rows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Rows, nil
}

func (c *CouchDB) Put(db, id string, doc any) (string, error) {
	b, _ := json.Marshal(doc)
	resp, err := c.req("PUT", fmt.Sprintf("/%s/%s", db, id), bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var r struct{ Rev string `json:"rev"` }
	json.NewDecoder(resp.Body).Decode(&r)
	return r.Rev, nil
}

func (c *CouchDB) Get(db, id string, out any) error {
	resp, err := c.req("GET", fmt.Sprintf("/%s/%s", db, id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *CouchDB) List(db string, includeDocs bool) ([]json.RawMessage, error) {
	q := "?include_docs=true"
	if !includeDocs {
		q = ""
	}
	resp, err := c.req("GET", "/"+db+"/_all_docs"+q, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res struct {
		Rows []struct {
			Doc json.RawMessage `json:"doc"`
		} `json:"rows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	docs := make([]json.RawMessage, len(res.Rows))
	for i, r := range res.Rows {
		docs[i] = r.Doc
	}
	return docs, nil
}

// Main service
type DocServer struct {
	db     *CouchDB
	schema *CouchDB // alias for schema db
	docs   *CouchDB // alias for documents db
}

func NewDocServer(host, user, pass string) (*DocServer, error) {
	c := NewCouchDB(host, user, pass)
	if err := c.EnsureDB("schema"); err != nil {
		return nil, err
	}
	if err := c.EnsureDB("documents"); err != nil {
		return nil, err
	}

	s := &DocServer{db: c}
	s.schema = c // same client
	s.docs = c

	// Ensure schema view (by-type)
	viewDef := map[string]any{
		"views": map[string]any{
			"schema-by-type": map[string]any{
				"map": `function(doc) { if(doc.type) { emit(doc.type, doc); } }`,
			},
		},
	}
	_, _ = c.Put("schema", "_design/schema", viewDef)

	return s, nil
}

// Schema API
func (s *DocServer) SaveSchema(id string, schema json.RawMessage) (string, error) {
	var sch map[string]any
	if err := json.Unmarshal(schema, &sch); err != nil {
		return "", err
	}
	if id == "" {
		id = uuid.NewString()
	}
	sch["_id"] = id
	rev, err := s.schema.Put("schema", id, sch)
	if err != nil {
		return "", err
	}
	return rev, nil
}

func (s *DocServer) GetSchemaByType(t DocType) (*Schema, error) {
	rows, err := s.schema.QueryView("schema", "schema", "schema-by-type", map[string]string{
		"key":          `"` + string(t) + `"`,
		"include_docs": "true",
		"limit":        "1",
	})
	if err != nil || len(rows) == 0 {
		return nil, fmt.Errorf("schema not found for type %s", t)
	}
	var sch Schema
	if err := json.Unmarshal(rows[0].Doc, &sch); err != nil {
		return nil, err
	}
	sch.JSON = rows[0].Doc // full schema
	return &sch, nil
}

func (s *DocServer) ListSchemas() ([]json.RawMessage, error) {
	return s.schema.List("schema", true)
}

// Document API with validation
func (s *DocServer) SaveDocument(d *Document) (*Document, error) {
	if d.ID == "" {
		d.ID = DocID(uuid.NewString())
	}

	// Get schema
	sch, err := s.GetSchemaByType(d.Type)
	if err != nil {
		return nil, err
	}

	// Simple validation (extend with full AJV-like if needed)
	if err := validateAgainstSchema(sch, d); err != nil {
		return nil, err.Error()
	}

	// Store
	_, err = s.docs.Put("documents", string(d.ID), d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func validateAgainstSchema(sch *Schema, doc *Document) *ValidationError {
	// Minimal required fields validation from schema (dense version)
	var schemaMap map[string]any
	json.Unmarshal(sch.JSON, &schemaMap)

	required, _ := schemaMap["required"].([]any)
	// props, _ := schemaMap["properties"].(map[string]any)

	if len(required) > 0 {
		for _, r := range required {
			if r == "title" && doc.Title == "" {
				return &ValidationError{Message: "missing required title"}
			}
			if r == "body" && doc.Body == "" {
				return &ValidationError{Message: "missing required body"}
			}
		}
	}

	// Could be extended with full JSON Schema validation library
	return nil
}

func (s *DocServer) GetDocument(id DocID) (*Document, error) {
	var d Document
	if err := s.docs.Get("documents", string(id), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *DocServer) ListDocuments() ([]json.RawMessage, error) {
	return s.docs.List("documents", true)
}
