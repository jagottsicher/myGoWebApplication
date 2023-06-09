package render

import (
	"net/http"
	"testing"

	"github.com/jagottsicher/myGoWebApplication/internal/models"
)

func TestAddDefaultData(t *testing.T) {

	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	session.Put(r.Context(), "flash", "a flash message")

	result := AddDefaultData(&td, r)
	if result.Flash != "a flash message" {
		t.Error("expected a value for key flash but flash message not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home-page.tpml", &models.TemplateData{})
	if err != nil {
		t.Error("writing template to browser failed:", err)
	}

	err = RenderTemplate(&ww, r, "does-not-exist-page.tpml", &models.TemplateData{})
	if err == nil {
		t.Error("requested template does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/an-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
