package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ziyue95/bookingandreservation/internal/config"
	"github.com/Ziyue95/bookingandreservation/internal/driver"
	"github.com/Ziyue95/bookingandreservation/internal/handlers"
	"github.com/Ziyue95/bookingandreservation/internal/helpers"
	"github.com/Ziyue95/bookingandreservation/internal/models"
	"github.com/Ziyue95/bookingandreservation/internal/render"
	"github.com/alexedwards/scs/v2"
)

// Declare global object

// Port number: const(never will be changed)
const portNumber = ":8080"

var app config.AppConfig
var sessionManager *scs.SessionManager // use pointer for easy injection to other pkgs
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	// run main logic in run
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	// Close DB connection after return
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Start the web application on port %s", portNumber))

	// Setup ROUTING service using pat/chi:
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// Change this to true when in production mode
	app.InProduction = false

	// 1. Set logger and put it into app object
	// infoLog: print to terminal with prefix "INFO", and flag: log.Ldate|log.Ltime
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	// errorLog: print to terminal with prefix "ERROR\t", log.Lshortfile: info about the error
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Initialize components of session(Reservation object)
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// initialize session config
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	// store session in Cookie
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction // in production should be true

	app.Session = sessionManager

	// Connect to DB
	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5433 dbname=bookings user=home password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	log.Println("Connected to database!")

	// Add template files into template cache for app <- app config variable
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	// In development mode, set app.UseCache = false -> update app variable when application is running
	app.UseCache = false

	// set app.config for render pkg to use
	render.NewRenderer(&app)

	// set repos for handlers pkg to use
	repo := handlers.NewRepo(&app, db)
	// pass repo into pointer Repo in handlers pkg
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	/*
		// Setup basic ROUTING service:
		// HandleFunc requires a URL, a handler which is a func of 1. http.ResponseWriter, 2. a pointer *http.Request
		// http.HandleFunc("/", handlers.Repo.Home)
		http.HandleFunc("/about", handlers.Repo.About)
		_ = http.ListenAndServe(portNumber, nil)
	*/

	return db, nil
}
