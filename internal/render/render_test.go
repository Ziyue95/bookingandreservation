package render

import (
	"net/http"
	"testing"

	"github.com/Ziyue95/bookingandreservation/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	// Get http request with session
	r, err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	// Put testing data into session
	sessionManager.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}

}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {

	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	// Build a empty http.Request pointer using http.NewRequest
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	// retrieve context from http.Request pointer
	ctx := r.Context()
	ctx, _ = sessionManager.Load(ctx, r.Header.Get("X-Session"))
	// Add session to current http.Request
	r = r.WithContext(ctx)
	return r, nil
}
