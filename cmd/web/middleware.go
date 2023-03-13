package main

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

func hitLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("HIT ... the road, Jack ...")
		next.ServeHTTP(w, r)
	})
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
