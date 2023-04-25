package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Ziyue95/bookingandreservation/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	// Write to app.InfoLog
	app.InfoLog.Println("Client error with status of", status)
	// throw a http error using http.Error
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	// Get a trace of error: error msg + stack trace of error
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	// throw a http error
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
