package handlers

import (
	"net/http"

	"github.com/jagottsicher/myGoWebApplication/pkg/config"
	"github.com/jagottsicher/myGoWebApplication/pkg/models"
	"github.com/jagottsicher/myGoWebApplication/pkg/render"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home-page.tpml", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// some data or calculation takes place

	sidekickMap := make(map[string]string)
	sidekickMap["morty"] = "Ooh, wee!"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	sidekickMap["remote_ip"] = remoteIP

	// send the result or any prepared data to the template
	render.RenderTemplate(w, "about-page.tpml", &models.TemplateData{
		StringMap: sidekickMap,
	})
}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact-page.tpml", &models.TemplateData{})
}

// Eremite is the handler for the eremite page
func (m *Repository) Eremite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "eremite-page.tpml", &models.TemplateData{})
}

// Couple is the handler for the couple page
func (m *Repository) Couple(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "couple-page.tpml", &models.TemplateData{})
}

// Family is the handler for the family page
func (m *Repository) Family(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "family-page.tpml", &models.TemplateData{})
}

// Reservation is the handler for the reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "check-availability-page.tpml", &models.TemplateData{})
}

// MakeReservation is the handler for the booking page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation-page.tpml", &models.TemplateData{})
}
