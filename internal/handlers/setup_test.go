package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"testing"

	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jagottsicher/myGoWebApplication/internal/config"
	"github.com/justinas/nosurf"

	"github.com/jagottsicher/myGoWebApplication/internal/models"
	"github.com/jagottsicher/myGoWebApplication/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

var functions = template.FuncMap{
	"humanReadableDate": render.HumanReadableDate,
	"formatDate":        render.FormatDate,
	"iterate":           render.Iterate,
	"add":               render.Add,
}

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Bungalow{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// change this to true when in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	listenForMail()

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/eremite", Repo.Eremite)
	mux.Get("/couple", Repo.Couple)
	mux.Get("/family", Repo.Family)
	mux.Get("/reservation", Repo.Reservation)
	mux.Post("/reservation", Repo.PostReservation)
	mux.Post("/reservation-json", Repo.ReservationJSON)
	mux.Get("/choose-bungalow/{id}", Repo.ChooseBungalow)
	mux.Get("/book-bungalow", Repo.BookBungalow)
	mux.Get("/make-reservation", Repo.MakeReservation)
	mux.Post("/make-reservation", Repo.PostMakeReservation)
	mux.Get("/reservation-overview", Repo.ReservationOverview)
	mux.Get("/user/login", Repo.ShowLogin)
	mux.Post("/user/login", Repo.PostShowLogin)
	mux.Get("/user/logout", Repo.Logout)

	mux.Get("/admin/dashboard", Repo.AdminDashboard)
	mux.Get("/admin/reservations-new", Repo.AdminNewReservations)
	mux.Get("/admin/reservations-all", Repo.AdminAllReservations)
	mux.Get("/admin/reservations-calendar", Repo.AdminReservationsCalendar)
	mux.Post("/admin/reservations-calendar", Repo.AdminPostReservationsCalendar)
	mux.Get("/admin/reservations/{src}/{id}/show", Repo.AdminShowReservation)
	mux.Post("/admin/reservations/{src}/{id}", Repo.AdminPostShowReservation)
	mux.Get("/admin/process-reservation/{src}/{id}/do", Repo.AdminProcessReservation)
	mux.Get("/admin/delete-reservation/{src}/{id}/do", Repo.AdminDeleteReservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf serves as a CSRF protection middleware
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads saves session data for each request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a map and stores the tempales in for caching.
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	theCache := map[string]*template.Template{}

	// get all available files *-page.tpml from folder ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*-page.tpml", pathToTemplates))
	if err != nil {
		return theCache, err
	}

	// range through the slice of *-page.tpml
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
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
