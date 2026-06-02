package service

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

// Maybe this is domain level stuff
type ValidationError struct {
	Message string `json:"message"`
	Errors  []any  `json:"errors,omitempty"`
}

func (err *ValidationError) Error() error {
	return errors.New(err.Message)
}

