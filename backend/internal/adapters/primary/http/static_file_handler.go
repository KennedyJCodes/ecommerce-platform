// Package http implements HTTP handlers for the ecommerce-platform application.
// This file contains the StaticFileHandler which is responsible for serving static files.
// It restricts the served files based on allowed extensions and ensures that only valid paths are served.
package http

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/gorilla/mux"
)

// StaticFileHandler is an HTTP handler for serving static files.

// It uses a static file service to retrieve file handlers and verifies that files have allowed extensions before serving them. This helps ensure security by preventing unauthorized file types from being accessed.
type StaticFileHandler struct {
	staticFileService output.StaticFilePort
	allowedExtensions map[string]bool
}

// NewStaticFileHandler creates and returns a new instance of StaticFileHandler.

// It receives an implementation of the StaticFilePort interface to handle file operations.
// It also defines a set of allowed file extensions that can be served.
func NewStaticFileHandler(staticFileService output.StaticFilePort) *StaticFileHandler {
	// Define allowed file extensions for security
	allowedExtensions := map[string]bool{
		".css":   true,
		".js":    true,
		".jpg":   true,
		".jpeg":  true,
		".png":   true,
		".gif":   true,
		".svg":   true,
		".ico":   true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
	}

	return &StaticFileHandler{
		staticFileService: staticFileService,
		allowedExtensions: allowedExtensions,
	}
}

// RegisterRoutes configures the routes for serving static files.

// It registers specific directory routes for common static asset folders (e.g. "/css/", "/js/", "/assets/") and defines a route for serving individual static files.
func (h *StaticFileHandler) RegisterRoutes(router *mux.Router) {
	// Register static directories using defined prefixes.
	h.registerStaticDir(router, "/css/", "css")
	h.registerStaticDir(router, "/js/", "js")
	h.registerStaticDir(router, "/assets/", "assets")

	// Register route for individual static files.
	router.HandleFunc("/static/{*}", h.HandleStaticFile)
}

// registerStaticDir registers a route that serves files from a specified static directory.

// It creates a file handler using the static file service, and sets up a route with the given prefix that serves files from the specified directory.
func (h *StaticFileHandler) registerStaticDir(router *mux.Router, prefix, dir string) {
	handler := h.staticFileService.GetFileHandler(prefix, dir)
	router.PathPrefix(prefix).Handler(handler)
}

// HandleStaticFile handles HTTP requests for individual static files.

// It extracts the requested file path, validates the file extension against the allowed list, and checks if the file path is valid via the static file service. If the file passes validation, it sets the appropriate Content-Type header and serves the file.
// If the file is not allowed or not found, it responds with an error.
func (h *StaticFileHandler) HandleStaticFile(w http.ResponseWriter, r *http.Request) {
	// Extract the file path from the URL variables.
	vars := mux.Vars(r)
	path := vars["*"]

	// Validate the file extension.
	ext := strings.ToLower(filepath.Ext(path))
	if _, ok := h.allowedExtensions[ext]; !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Ensure the requested path is valid.
	if !h.staticFileService.IsValidPath(path) {
		http.NotFound(w, r)
		return
	}

	// Set the Content-Type header based on the file's MIME type.
	mimeType := h.staticFileService.GetMimeType(path)
	w.Header().Set("Content-Type", mimeType)

	// Retrieve the file handler and serve the file.
	handler := h.staticFileService.GetFileHandler("/static/", "")
	handler.ServeHTTP(w, r)
}
