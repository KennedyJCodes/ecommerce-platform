// Package http implements HTTP handlers for the ecommerce-platform application.
// This file contains the MainPageHandler, responsible for serving the main page of the application.
package http

import (
	"net/http"
	"path/filepath"
	"text/template"
)

// MainPageHandler handles HTTP requests to the main page.

// It serves the main HTML page of the application, typically used as the entry point for client-side rendered applications.
type MainPageHandler struct {
	staticDir string
}

// NewMainPageHandler creates a new instance of MainPageHandler.
// It initializes the handler without a predefined static directory.
func NewMainPageHandler() *MainPageHandler {
	return &MainPageHandler{}
}

// SetStaticDir sets the directory from which static files are served.

// This method allows configuring the path to the directory containing static assets, such as HTML, CSS, and JavaScript files.
func (h *MainPageHandler) SetStaticDir(staticDir string) {
	h.staticDir = staticDir
}

// Handle processes HTTP requests to the main page.

// It determines the path to the index.html file, either from the configured static directory or a default path. It then parses and executes the template, writing the rendered HTML to the response. If an error occurs during template parsing, it responds with an HTTP 500 Internal Server Error.
func (h *MainPageHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Determine the path to index.html
	var indexPath string

	// If we have a static directory configured, use its path
	if h.staticDir != "" {
		indexPath = filepath.Join(h.staticDir, "index.html")
	} else {
		// If not, use the default route
		indexPath = "./../frontend/index.html"
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles(indexPath)
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
