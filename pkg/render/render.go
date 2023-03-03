package render

import (
	"fmt"
	"html/template"
	"net/http"
)

// rendernTemplate serves as a wrapper and renders
// a layout and a template from folder /templates to a desired writer
func RenderTemplate(w http.ResponseWriter, tpml string) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+tpml, "./templates/base-layout.tpml")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
	}
}
