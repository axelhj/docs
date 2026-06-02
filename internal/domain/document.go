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

type Document struct {
	ID    DocID          `json:"id,omitempty"`
	Type  DocType        `json:"type"`	Title string         `json:"title"`
	Body  string         `json:"body"`
	Rev   string         `json:"_rev,omitempty"`
	Extra map[string]any `json:"_extra,omitempty"` // for extensibility
}
