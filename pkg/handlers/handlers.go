package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
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

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// start & end are string pulled from the template form
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

// `json:<feature_name>`
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// Availability JSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// Build a JSON response
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	// Marshal resp into json(out) with indent
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	// Create a header to tell the web browser wut kind of response I am sending:
	// application/json <- standard header for json files
	w.Header().Set("Content-Type", "application/json") // set Content-Type in the header as application/json
	// Write to the web browser
	w.Write(out)

}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
