// Package static provides an adapter for serving static assets in the ecommerce-platform application.
// It wraps Go's standard http.FileServer to mount directories and paths for CSS, JavaScript, and other assets.
package static

import (
	"net/http"
	"path/filepath"
)

// FileServer serves static files from a configured directory.

// It centralizes the logic for mounting asset directories on an http.ServeMux.
// Supported directories include /css/, /js/, and /assets/ by default, and custom paths via ServePath.
// This adapter is part of the infrastructure layer in a hexagonal architecture.
type FileServer struct {
	staticDir string
}

// NewFileServer creates a new FileServer with the given root directory.

// staticDir: the filesystem path where static assets are located.
func NewFileServer(staticDir string) *FileServer {
	return &FileServer{
		staticDir: staticDir,
	}
}

// ServeStatic mounts the standard asset directories on the provided ServeMux.

// It registers handlers for:
//   /css/    -> <staticDir>/css
//   /js/     -> <staticDir>/js
//   /assets/ -> <staticDir>/assets

// Each handler uses http.StripPrefix to remove the URL prefix before serving files.
func (fs *FileServer) ServeStatic(router *http.ServeMux) {
	// Define paths for static files
	cssHandler := http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(fs.staticDir, "css"))))

	jsHandler := http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join(fs.staticDir, "js"))))

	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(fs.staticDir, "assets"))))

	// Register handlers
	router.Handle("/css/", cssHandler)
	router.Handle("/js/", jsHandler)
	router.Handle("/assets/", assetsHandler)
}

// ServePath mounts a custom URL path to serve files from a corresponding subdirectory.

// path: the URL prefix (e.g., "/static/")
// It serves files from <staticDir>/<trimmed path> and strips the prefix before serving.
func (fs *FileServer) ServePath(path string, router *http.ServeMux) {
	handler := http.StripPrefix(path, http.FileServer(http.Dir(filepath.Join(fs.staticDir, path))))
	router.Handle(path+"*", handler)
}

// GetStaticDir returns the configured root directory for static files.
func (fs *FileServer) GetStaticDir() string {
	return fs.staticDir
}
