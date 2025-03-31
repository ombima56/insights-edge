package middleware

import (
	"net/http"

	"github.com/ombima56/insights-edge/internal/database"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var exists bool
		err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM sessions WHERE token = ? AND expires_at > CURRENT_TIMESTAMP)", cookie.Value).Scan(&exists)
		if err != nil || !exists {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}
