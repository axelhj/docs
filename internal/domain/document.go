package domain

import ()

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
