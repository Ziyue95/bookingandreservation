package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// WriteToConsole is a self-defined middleware
func WriteToConsole(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		// move to next middleware
		next.ServeHTTP(w, r)
	})
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	// initialize a nodurf handler
	csrfHandler := nosurf.New(next)

	// create nosurf token using basic http cookie
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/.",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request <- let the web application to be stateful using session
func SessionLoad(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}
