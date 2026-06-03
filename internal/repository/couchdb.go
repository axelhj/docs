package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

//mockery:generate: true
type CouchDBClient interface {
	EnsureDB(name string) error;
	QueryView(db, design, view string, params map[string]string) ([]ViewRow, error);
	Put(db, id string, doc any) (string, error);
	Get(db, id string, out any) error;
	List(db string, includeDocs bool) ([]json.RawMessage, error);
}

type CouchDB struct {
	dber CouchDBClient
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
	// TODO: Make sure an actual base-url and appropriate protocol is set/defined.
	_, err := c.req("PUT", "/"+name, nil)
	if err != nil && !strings.Contains(err.Error(), "412") { // already exists
		return err
	}
	return nil
}

type ViewRow struct {
	ID    string          `json:"id"`
	Key   any             `json:"key"`
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
	var r struct {
		Rev string `json:"rev"`
	}
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
