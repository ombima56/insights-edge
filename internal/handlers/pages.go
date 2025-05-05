package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ombima56/insights-edge/internal/database"
)

type MarketInsight struct {
	ID          int     `json:"id"`
	Industry    string  `json:"industry"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Content     string  `json:"content"`
	UserID      int     `json:"user_id"`
	TrendValue  float64 `json:"trend_value"`
	CreatedAt   string  `json:"created_at"`
}

type PageData struct {
	Title           string
	IsAuthenticated bool
	User            *User
	Insights        []MarketInsight
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	if data == nil {
		data = &PageData{}
	}

	files := []string{
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "pages", tmpl+".html"),
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := &PageData{
		Title:           "Home",
		IsAuthenticated: isAuthenticated(r),
	}
	renderTemplate(w, "home", data)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var insights []MarketInsight
	rows, err := database.DB.Query(`
		SELECT id, industry, title, description, type, content, user_id, trend_value, created_at 
		FROM market_insights 
		ORDER BY created_at DESC
	`)
	if err != nil {
		// Log the specific error
		log.Printf("Database query error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var insight MarketInsight
		err := rows.Scan(&insight.ID, &insight.Industry, &insight.Title, &insight.Description, 
			&insight.Type, &insight.Content, &insight.UserID, &insight.TrendValue, &insight.CreatedAt)
		if err != nil {
			// Log the specific error
			log.Printf("Row scan error: %v", err)
			continue
		}
		insights = append(insights, insight)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := &PageData{
		Title:           "Dashboard",
		IsAuthenticated: true,
		User:            user,
		Insights:        insights,
	}
	renderTemplate(w, "dashboard", data)
}

func getCurrentUser(r *http.Request) *User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}

	var userID int
	err = database.DB.QueryRow("SELECT user_id FROM sessions WHERE token = ? AND expires_at > CURRENT_TIMESTAMP", cookie.Value).Scan(&userID)
	if err != nil {
		return nil
	}

	var user User
	err = database.DB.QueryRow(`
		SELECT id, email, first_name, last_name, account_type, company_name, industry 
		FROM users 
		WHERE id = ?`, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.AccountType,
		&user.CompanyName,
		&user.Industry)
	if err != nil {
		return nil
	}

	return &user
}

func isAuthenticated(r *http.Request) bool {
	return getCurrentUser(r) != nil
}
