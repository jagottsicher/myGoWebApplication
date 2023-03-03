package handlers

import (
	"net/http"

	"github.com/jagottsicher/myGoWebApplication/pkg/render"
)

// Home is the handler for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home-page.tpml")
}

// About is the handler for the about page
func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about-page.tpml")
}
