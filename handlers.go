package main

import "net/http"

// Home is the handler for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	rendernTemplate(w, "home-page.tpml")
}

// About is the handler for the about page
func About(w http.ResponseWriter, r *http.Request) {
	rendernTemplate(w, "about-page.tpml")
}
