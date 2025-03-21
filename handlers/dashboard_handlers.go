package handlers

import (
	"html/template"
	"net/http"
)

// DashboardHandler serves the dashboard HTML and handles API requests
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/dashboard" {
		serveDashboardTemplate(w, r)
		return
	}

	http.NotFound(w, r)
}

// serveDashboardTemplate serves the dashboard HTML template
func serveDashboardTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("frontend/templates/dashboard.html")
	if err != nil {
		http.Error(w, "Failed to load dashboard template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Dashboard",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
