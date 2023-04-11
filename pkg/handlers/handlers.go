package handlers

import (
	"fmt"
	"net/http"

	"github.com/Ziyue95/bookingandreservation/pkg/config"
	"github.com/Ziyue95/bookingandreservation/pkg/models"
	"github.com/Ziyue95/bookingandreservation/pkg/render"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo is the repository used by the handlers
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

// Build receivers for all handlers and link them with the repository
// -> All handlers have access to the repository

// Define the Handler function for home page:
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	// Put the remoteIP into session
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// Define the Handler function for about page:
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again!"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	if remoteIP == "" {
		fmt.Println("Can not retrieve remote ip address")
	}

	stringMap["remote_ip"] = remoteIP

	// send data to the template

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
