package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AuthMiddlewareWithStore(store *sessions.CookieStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "auth-session")
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			token, ok := session.Values["token"].(string)
			if !ok || token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
