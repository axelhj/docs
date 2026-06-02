package test

import (
	"testing"
)

func TestAny(t *testing.T) {
	person := "Me"
	if false {
		t.Errorf(`Hello there %q`, person)
	}
}
