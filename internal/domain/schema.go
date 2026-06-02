package domain

import (
	"encoding/json"
)

type Schema struct {
	ID   DocID           `json:"_id,omitempty"`
	Rev  string          `json:"_rev,omitempty"`
	Type DocType         `json:"type"`
	JSON json.RawMessage `json:"-"` // full JSON Schema
	Raw  map[string]any  `json:"-"` // parsed
}
