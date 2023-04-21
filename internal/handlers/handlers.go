package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Ziyue95/bookingandreservation/internal/config"
	"github.com/Ziyue95/bookingandreservation/internal/forms"
	"github.com/Ziyue95/bookingandreservation/internal/models"
	"github.com/Ziyue95/bookingandreservation/internal/render"
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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		// initialize an empty form $ data for when first time rendering the page
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// Parse the form and check err
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	// Get client-entered data from form with id using r.Form.Get
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	// Form validation
	form := forms.New(r.PostForm)
	// 1. Check if required field has value
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email", "phone")
	// check minimum length for "first_name"
	form.MinLength("first_name", 3, r)
	// check email
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		// render the page using both client-entered data(reservation) and info stored in form(like errors, etc)
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			// Include form info: export error msg to web page
			Form: form,
			Data: data,
		})
		return
	}

	// Throw the reservation variable into session stored in App (m.App.Session)
	m.App.Session.Put(r.Context(), "reservation", reservation)
	// Direct the user post the reservation to /reservation-summary with respond code http.StatusSeeOther(303)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

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

// ReservationSummary displays the res summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	// Make sure if we get reservation from session
	if !ok {
		log.Println("can't get item from session")
		// Put error msg into r.Context() with field "error"
		// Can take msg using app.Session.PopString(r.Context(), "error");
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		// Redirect user back to homepage
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// After get reservation info -> remove reservation from Session
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
