package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// rendernTemplate serves as a wrapper and renders
// a template from folder /templates to a desired writer
func rendernTemplate(w http.ResponseWriter, tpml string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tpml)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
	}
}
