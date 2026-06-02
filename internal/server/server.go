package server

import (
	// "github.com/santhosh-tekuri/jsonschema"

	"github.com/axelhj/docs/internal/repository"
)


type DocServer struct {
	db     *repository.CouchDB
	schema *repository.CouchDB // alias for schema db
	docs   *repository.CouchDB // alias for documents db
}

func NewDocServer(host, user, pass string) (*DocServer, error) {
	c := repository.NewCouchDB(host, user, pass)
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

