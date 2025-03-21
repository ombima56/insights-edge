package handlers

import (
	"net/http"
)

// HomeHandler serves the main page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Serve the index.html file from the templates directory
	http.ServeFile(w, r, "frontend/templates/index.html")
}

// AuthHandler serves the authentication page
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the auth.html file from the templates directory
	http.ServeFile(w, r, "frontend/templates/auth.html")
}

// SetupStaticFileServer configures the static file server for CSS, JS, etc.
func SetupStaticFileServer() {
	fs := http.FileServer(http.Dir("frontend/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
