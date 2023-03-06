package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// rendernTemplate serves as a wrapper and renders
// a layout and a template from folder /templates to a desired writer
func RenderTemplate(w http.ResponseWriter, tpml string) {
	// create a template cache
	tc, err := createTemplateCache()
	if err != nil {
		log.Fatalln("error creating template cache ", err)
	}

	// get the right template from cache
	t, ok := tc[tpml]
	if !ok {
		log.Fatalln("template not in cache for some reason ", err)
	}

	// store result in a buffer and double-check if it is a valid value
	buf := new(bytes.Buffer)

	err = t.Execute(buf, nil)
	if err != nil {
		log.Println(err)
	}

	// render that template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache() (map[string]*template.Template, error) {
	theCache := map[string]*template.Template{}

	// get all available files *-page.tpml from folder ./templates
	pages, err := filepath.Glob("./templates/*-page.tpml")
	if err != nil {
		return theCache, err
	}

	// range through the slice of *-page.tpml
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return theCache, err
		}

		matches, err := filepath.Glob("./templates/*-layout.tpml")
		if err != nil {
			return theCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*-layout.tpml")
			if err != nil {
				return theCache, err
			}
		}

		theCache[name] = ts
	}
	return theCache, nil
}
