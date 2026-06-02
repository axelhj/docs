package domain

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

type DocType string
type DocID string

type Schema struct {
	ID   DocID          `json:"_id,omitempty"`
	Rev  string         `json:"_rev,omitempty"`
	Type DocType        `json:"type"`
	JSON json.RawMessage `json:"-"` // full JSON Schema
	Raw  map[string]any `json:"-"` // parsed
}
