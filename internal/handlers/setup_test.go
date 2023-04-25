package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Ziyue95/bookingandreservation/internal/config"
	"github.com/Ziyue95/bookingandreservation/internal/models"
	"github.com/Ziyue95/bookingandreservation/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var sessionManager *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {

	// Copy from main.go(func run)
	gob.Register(models.Reservation{})
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction
	app.Session = sessionManager

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = true // Avoid calling pathToTemplates in render.go which is inaccessible here
	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	// Copy from routes.go(func routes)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}

// Copy from middleware.go
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

func SessionLoad(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}

// Copy from render.go
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all files named *.page.tmpl from ./template
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// get all layout files (*.layout.tmpl) from ./template
	matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	// page: the full filepath of each *.page.tmpl file
	for _, page := range pages {
		// name: name of template file (*.page.tmpl)
		name := filepath.Base(page)
		// ts: parsed template set by parsing the template name in page
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
