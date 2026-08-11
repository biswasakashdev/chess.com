package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

type AuthHandler struct {
	authServ *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authServ: authService,
	}
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type RegisterResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	ah.authServ.Register(req.Username, req.Password, req.FirstName, req.LastName)

}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	authToken, err := ah.authServ.Login(req.Username, req.Password)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: authToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
