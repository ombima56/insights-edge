package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ombima56/insights-edge/backend/models"
)

type FlashMessage struct {
	Type    string `json:"type"`    // "success", "error", "info", "warning"`
	Message string `json:"message"` // The message to display
}

var SessionStore *sessions.CookieStore

func InitSessionStore(key []byte) {
	SessionStore = sessions.NewCookieStore(key)
	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   true,      // Set to true in production
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return SessionStore.Get(r, "auth-session")
}

// AddFlashMessage adds a flash message to the session
func AddFlashMessage(w http.ResponseWriter, r *http.Request, msgType, message string) {
	session, err := GetSession(r)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get existing messages or create a new slice
	var messages []FlashMessage
	if raw, ok := session.Values["flash_messages"]; ok {
		if err := json.Unmarshal([]byte(raw.(string)), &messages); err == nil {
			messages = append(messages, FlashMessage{Type: msgType, Message: message})
		}
	} else {
		messages = []FlashMessage{{Type: msgType, Message: message}}
	}

	// Store the messages in the session
	flashBytes, _ := json.Marshal(messages)
	session.Values["flash_messages"] = string(flashBytes)
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetFlashMessages retrieves and clears flash messages from the session
func GetFlashMessages(r *http.Request, w http.ResponseWriter) []FlashMessage {
	session, err := GetSession(r)
	if err != nil {
		return []FlashMessage{}
	}

	// Get and clear the messages
	var messages []FlashMessage
	if raw, ok := session.Values["flash_messages"]; ok {
		if err := json.Unmarshal([]byte(raw.(string)), &messages); err == nil {
			delete(session.Values, "flash_messages")
			if err := session.Save(r, w); err != nil {
				return []FlashMessage{}
			}
		}
	}

	return messages
}

func HandleError(w http.ResponseWriter, err error) {
	log.Printf("Error: %v", err)
	switch err {
	case models.ErrInvalidCredentials:
		w.WriteHeader(http.StatusUnauthorized)
	case models.ErrUserExists:
		w.WriteHeader(http.StatusConflict)
	case models.ErrUserNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
