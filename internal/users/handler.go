package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserHandler struct {
	userServ *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userServ: userService,
	}
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	user := User{
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
