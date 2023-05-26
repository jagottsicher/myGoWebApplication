package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// all fine nothing to do
	default:
		t.Error(fmt.Sprintf("Type mismatch: Expected http.Handler, got %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// all fine nothing to do
	default:
		t.Error(fmt.Sprintf("Type mismatch: Expected http.Handler, got %T", v))
	}
}
