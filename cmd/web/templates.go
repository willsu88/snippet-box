package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/willsu88/snippet-box/pkg/forms"
	"github.com/willsu88/snippet-box/pkg/models"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Struct to hold the dynamic data
// that we pass to the HTML template
type templateData struct {
	CurrentYear int
	Form        *forms.Form
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file pat
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the page template file in to a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
