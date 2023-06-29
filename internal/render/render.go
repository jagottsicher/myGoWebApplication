package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/jagottsicher/myGoWebApplication/internal/config"
	"github.com/jagottsicher/myGoWebApplication/internal/models"
	"github.com/justinas/nosurf"
)

// AddDefaultData contains Data which will be added to data sent to templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// Template serves as a wrapper and renders
// a layout and a template from folder /templates to a desired writer
func Template(w http.ResponseWriter, r *http.Request, tpml string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get the right template from cache
	t, ok := tc[tpml]
	if !ok {
		return errors.New("template not in cache for some reason")
	}

	// store result in a buffer and double-check if it is a valid value
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render that template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	return nil
}

// CreateTemplateCache creates a map and stores the tempales in for caching.
func CreateTemplateCache() (map[string]*template.Template, error) {
	theCache := map[string]*template.Template{}

	// get all available files *-page.tpml from folder ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*-page.tpml", pathToTemplates))
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

		matches, err := filepath.Glob(fmt.Sprintf("%s/*-layout.tpml", pathToTemplates))
		if err != nil {
			return theCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*-layout.tpml", pathToTemplates))
			if err != nil {
				return theCache, err
			}
		}

		theCache[name] = ts
	}
	return theCache, nil
}
