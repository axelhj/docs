package service

import (
	"errors"
)

// Maybe this is domain level stuff
type ValidationError struct {
	Message string `json:"message"`
	Errors  []any  `json:"errors,omitempty"`
}

func (err *ValidationError) Error() error {
	return errors.New(err.Message)
}
