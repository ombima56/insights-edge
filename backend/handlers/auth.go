package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ombima56/insights-edge/backend/models"
	"github.com/ombima56/insights-edge/backend/service"
	"github.com/ombima56/insights-edge/backend/utils"
)

type AuthHandler struct {
	userService *service.UserService
	store       *sessions.CookieStore
}

func NewAuthHandler(userService *service.UserService, store *sessions.CookieStore) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		store:       store,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.HandleError(w, err)
		return
	}

	modelUser := &models.User{
		Email:        user.Email,
		Password:     user.Password,
		WalletAddr:   user.WalletAddr,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AccountType:  user.AccountType,
		CompanyName:  user.CompanyName,
		Industry:     user.Industry,
	}

	if err := h.userService.RegisterUser(modelUser); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var login models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		utils.HandleError(w, err)
		return
	}

	token, err := h.userService.LoginUser(login.Email, login.Password)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			utils.HandleError(w, err)
			return
		}
		utils.HandleError(w, err)
		return
	}

	session, err := h.store.Get(r, "auth-session")
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	session.Values["token"] = token
	if err := session.Save(r, w); err != nil {
		utils.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := &PageData{
		Title:           "Register",
		IsAuthenticated: isAuthenticated(r),
	}
	renderTemplate(w, "register", data)
}
