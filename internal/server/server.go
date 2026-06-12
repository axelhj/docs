package server

import (
	"encoding/json"
	"fmt"

	"github.com/axelhj/docs/internal/domain"
	"github.com/axelhj/docs/internal/repository"
	"github.com/axelhj/docs/internal/service"
	"github.com/google/uuid"
)

type DocServer struct {
	db     repository.CouchDBClient
	schema repository.CouchDBClient // "alias" for schema db
	docs   repository.CouchDBClient // "alias" for documents db
}

func NewDocServer(host, port, user, pass string) (*DocServer, error) {
	c := repository.NewCouchDB("http", port, host, user, pass)
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

func (s *DocServer) GetSchemaByType(t domain.DocType) (*domain.Schema, error) {
	rows, err := s.schema.QueryView("schema", "schema", "schema-by-type", map[string]string{
		"key":          `"` + string(t) + `"`,
		"include_docs": "true",
		"limit":        "1",
	})
	if err != nil || len(rows) == 0 {
		return nil, fmt.Errorf("schema not found for type %s", t)
	}
	var sch domain.Schema
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
func (s *DocServer) SaveDocument(d *domain.Document) (*domain.Document, error) {
	if d.ID == "" {
		d.ID = domain.DocID(uuid.NewString())
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

func validateAgainstSchema(sch *domain.Schema, doc *domain.Document) *service.ValidationError {
	// Minimal required fields validation from schema (dense version)
	var schemaMap map[string]any
	json.Unmarshal(sch.JSON, &schemaMap)

	required, _ := schemaMap["required"].([]any)
	// props, _ := schemaMap["properties"].(map[string]any)

	if len(required) > 0 {
		for _, r := range required {
			if r == "title" && doc.Title == "" {
				return &service.ValidationError{Message: "missing required title"}
			}
			if r == "body" && doc.Body == "" {
				return &service.ValidationError{Message: "missing required body"}
			}
		}
	}

	// Could be extended with full JSON Schema validation library
	return nil
}

func (s *DocServer) GetDocument(id domain.DocID) (*domain.Document, error) {
	var d domain.Document
	if err := s.docs.Get("documents", string(id), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *DocServer) ListDocuments() ([]json.RawMessage, error) {
	return s.docs.List("documents", true)
}
