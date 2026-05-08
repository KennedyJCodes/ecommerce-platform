// Package static provides an implementation of the StaticFilePort interface for serving static files from a configured directory in the ecommerce-platform application.
package static

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// StaticFileAdapter implements the output.StaticFilePort interface.

// It serves files from a base directory, validates requested paths to prevent directory traversal, and determines the correct MIME type for each file.
type StaticFileAdapter struct {
	staticDir string
}

// NewStaticFileAdapter creates a new StaticFileAdapter with the given directory.

// staticDir: the filesystem path where static assets are stored.
// Returns an initialized StaticFileAdapter.
func NewStaticFileAdapter(staticDir string) *StaticFileAdapter {
	return &StaticFileAdapter{
		staticDir: staticDir,
	}
}

// GetStaticDir returns the root directory configured for static file serving.
func (s *StaticFileAdapter) GetStaticDir() string {
	return s.staticDir
}

// GetFileHandler returns an HTTP handler that serves files from a specific subdirectory.

// prefix: the URL path prefix to strip (e.g., "/static/")
// subPath: a subdirectory under staticDir (e.g., "images")
// The returned handler serves files from filepath.Join(staticDir, subPath), optionally stripping the given URL prefix before filesystem lookup
func (s *StaticFileAdapter) GetFileHandler(prefix, subPath string) http.Handler {
	// Build the full filesystem path
	fullPath := filepath.Join(s.staticDir, subPath)

	// Create the standard file server
	fileServer := http.FileServer(http.Dir(fullPath))

	// If a non-root prefix is provided, strip it before serving
	if prefix != "" && prefix != "/" {
		return http.StripPrefix(prefix, fileServer)
	}

	return fileServer
}

// IsValidPath checks whether the given URL path corresponds to an existing file under the configured static directory, preventing directory traversal.

// path: the request URL path relative to the static root (e.g., "css/style.css").
// Returns true if the file exists and is not a directory.
func (s *StaticFileAdapter) IsValidPath(path string) bool {
	// Normalize and clean the path
	cleanPath := filepath.Clean(path)

	// Reject paths containing parent directory references
	if strings.Contains(cleanPath, "..") {
		return false
	}

	// Build the absolute filesystem path
	fullPath := filepath.Join(s.staticDir, cleanPath)

	// Check existence and type
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}

	// Ensure it's not a directory
	return !info.IsDir()
}

// GetMimeType returns the MIME type for the given filename based on its extension.

// filename: the name or path of the file (e.g., "style.css").
// Returns the detected MIME type, or "application/octet-stream" if unknown.
func (s *StaticFileAdapter) GetMimeType(filename string) string {
	// Extract the extension, e.g., ".css"
	ext := filepath.Ext(filename)

	// Lookup the MIME type
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// Default to binary stream if unknown
		return "application/octet-stream"
	}

	return mimeType
}
