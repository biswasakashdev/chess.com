package routers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/biswasakashdev/chess.com/internal/auth"
	"github.com/biswasakashdev/chess.com/internal/dtos"
	"github.com/go-chi/chi/v5"
)

type authHandler struct {
	authServ *auth.AuthService
}

func NewAuthRouter(authService *auth.AuthService) chi.Router {
	authHandl := authHandler{
		authServ: authService,
	}
	router := chi.NewRouter()

	router.HandleFunc("/", authHandl.login)
	router.HandleFunc("/register", authHandl.register)

	return router
}

func (ah *authHandler) register(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req dtos.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	user, err := ah.authServ.Register(ctx, req.Username, req.Password, req.FirstName, req.LastName)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (ah *authHandler) login(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req dtos.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	authToken, err := ah.authServ.Login(ctx, req.Username, req.Password)

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	response := dtos.LoginResponse{
		Token: authToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
