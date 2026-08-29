package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"password"`
	CreatedAt      time.Time `json:"createdAt"`
}
