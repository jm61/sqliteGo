package main

import (
	"chinook/internal/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Users           []models.User
	Employees       []models.Employee
	Albums          []models.Album
	Records         []models.Record
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	Artist          string
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// template function to convert milliseconds column of tracks table in seconds
	funcMap := template.FuncMap{
		"div": func(a, b int) int {
			if b == 0 {
				return 0 // Handle division by zero
			}
			return a / b
		},
	}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	// Loop through the page filepaths one-by-one.
	for _, page := range pages {

		name := filepath.Base(page)
		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			"./ui/html/partials/albums.html",
			page,
		}
		// Parse the files into a template set to include the funcmap.
		ts, err := template.New("test").Funcs(funcMap).ParseFiles(files...)
		//ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}
