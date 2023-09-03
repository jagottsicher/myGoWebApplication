package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jagottsicher/myGoWebApplication/internal/config"
	"github.com/jagottsicher/myGoWebApplication/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/eremite", handlers.Repo.Eremite)
	mux.Get("/couple", handlers.Repo.Couple)
	mux.Get("/family", handlers.Repo.Family)
	mux.Get("/reservation", handlers.Repo.Reservation)
	mux.Post("/reservation", handlers.Repo.PostReservation)
	mux.Post("/reservation-json", handlers.Repo.ReservationJSON)
	mux.Get("/choose-bungalow/{id}", handlers.Repo.ChooseBungalow)
	mux.Get("/book-bungalow", handlers.Repo.BookBungalow)
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservation)
	mux.Get("/reservation-overview", handlers.Repo.ReservationOverview)
	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
