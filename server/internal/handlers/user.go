package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/biswasakashdev/chess.com/server/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type userHandler struct {
	userServ *users.UserService
}

func NewUserRouter(userService *users.UserService) chi.Router {

	userHandl := userHandler{
		userServ: userService,
	}
	router := chi.NewRouter()

	router.HandleFunc("/", userHandl.getUser)

	return router
}

func (uh *userHandler) getUser(w http.ResponseWriter, r *http.Request) {

	// ctx := r.Context()

	// userId := ctx.Value(middleware.UserIDKey).(string)

	user := users.User{
		FirstName:      "Akash",
		LastName:       "Biswas",
		Id:             uuid.New(),
		Username:       "abc@email.com",
		HashedPassword: "password",
		CreatedAt:      time.Now(),
	}

	// 2. Set the content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// 3. Set the HTTP status code
	w.WriteHeader(http.StatusOK)

	// 4. Encode the struct directly into the response writer
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
