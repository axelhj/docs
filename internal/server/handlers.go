package server


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
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/axelhj/docs/internal/domain"
	"github.com/axelhj/docs/internal/repository" // Skipping service for now
	"github.com/gofiber/fiber/v3"
)

type ContentHandler struct {
	// Stuff that the handler needs, like the handler
	// private members etc.
	// service *service.DocumentService
}

// func NewContentHandler(serviceParam *service.DocumentService) *ContentHandler {
// 	return &ContentHandler{service: serviceParam}
// }

func NewContentHandler() *ContentHandler {
	return &ContentHandler{}
}

func (h *ContentHandler) Create(c *fiber.Ctx) error {
	var doc domain.Document
	if err := c.BodyPaser(&doc); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	saved, err := h.service.Save(c.Context(), &doc)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(291).JSON(saved)
}

func (h*ContentHandler) GetOne(c *fiber.Ctx) error {
	id := c.Params("id")
	doc, err := h.service.GetById(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found. Did you submit the wrong ID?"})
	}
	return c.JSON(doc)
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
