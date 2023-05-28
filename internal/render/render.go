package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Ziyue95/bookingandreservation/internal/config"
	"github.com/Ziyue95/bookingandreservation/internal/models"
	"github.com/justinas/nosurf"
)

// functions contains all function available to the golang template
var functions = template.FuncMap{}

// app is a pointer to config.AppConfig for render pkg to use
var app *config.AppConfig

var pathToTemplates = "./templates"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default value to td
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// Take success/warning/error msg from session
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template
// Complex way to implement template caching
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	// IMPORTANT: Only read tc from disk when app.UseCache is false
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		// Same as: tc := (*app).TemplateCache
		tc = app.TemplateCache
	} else {
		var err error
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		// log.Fatal("Can't find template in cache")
		return errors.New("Can't get template from cache")
	}

	// create a buffer to execute the template for finer-grained checking
	buf := new(bytes.Buffer)
	// set default value to td
	td = AddDefaultData(td, r)
	// Add td into the buffer
	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
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
		// ts: parsed template set injected with functions & parsed by the template name in page
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

/*
// Simple ways to implement template caching

// tc(template cache) is a package level variable to store rendered template
var tc = make(map[string]*template.Template)

// RenderTemplateTest: t is the template name
func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check to see if we already have the template in our cache
	// inMap equals true if t in tc, and false otherwise
	_, inMap := tc[t]
	if !inMap {
		// need to create the template
		log.Println(fmt.Sprintf("creating template: %s and caching this template", t))
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// we have the template in cache
		log.Println("using cached template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}

}

// createTemplateCache renders template and store in template cache (map)
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	// parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// add template to cache (map)
	tc[t] = tmpl

	return nil
}
*/
