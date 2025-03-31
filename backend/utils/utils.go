package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ombima56/insights-edge/backend/models"
)

func HandleError(w http.ResponseWriter, err error) {
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
