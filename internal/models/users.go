package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID
	FirstName      string
	LastName       string
	Username       string
	Rating         int
	HashedPassword string
	CreatedAt      time.Time
}
