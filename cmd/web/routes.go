package main

import (
	"net/http"

	"github.com/Ziyue95/bookingandreservation/pkg/config"
	"github.com/Ziyue95/bookingandreservation/pkg/handlers"

	// "github.com/bmizerany/pat"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	/*
		// Setup routing using pat pkg
		// Initialize HTTP request multiplexer (mux)
		mux := pat.New()

		// set up routes + call associated handler function
		mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
		mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	*/

	// Setup routing using chi pkg
	mux := chi.NewRouter()

	// set middlewares using chi
	mux.Use(middleware.Recoverer)
	// use self-defined middleware
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
