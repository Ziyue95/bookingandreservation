package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Ziyue95/bookingandreservation/internal/config"
	"github.com/Ziyue95/bookingandreservation/internal/models"
	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// Setup components to be put into the session
	gob.Register(models.Reservation{})

	// change this to true when in production mode
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// initialize session config
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	// store session in Cookie
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false

	testApp.Session = sessionManager

	app = &testApp

	os.Exit(m.Run())
}

// Create a type that satisfy the http.ResponseWriter interface
type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
