package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jagottsicher/myGoWebApplication/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// all fine nothing to do
	default:
		t.Error(fmt.Sprintf("Type mismatch: Expected *chi.Mux, got %T", v))
	}
}
