package users

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{
		db: db,
	}
}

func (ur *SQLiteUserRepository) Save() (*User, error) {

	return &User{
		FirstName:      "Akash",
		LastName:       "Biswas",
		Id:             uuid.New(),
		Username:       "abc123",
		HashedPassword: "hashed",
		CreatedAt:      time.Now(),
	}, nil
}

func (ur *SQLiteUserRepository) IsUsernameExits(username string) (bool, error) {
	return false, nil
}

func (ur *SQLiteUserRepository) GetUserByUsername(username string) (*User, error) {

	return &User{
		FirstName:      "Akash",
		LastName:       "Biswas",
		Id:             uuid.New(),
		Username:       "abc123",
		HashedPassword: "hashed",
		CreatedAt:      time.Now(),
	}, nil
}

func (ur *SQLiteUserRepository) CreateUser(username, hashedPassword, firstName, lastName string) (*User, error) {

	return &User{
		FirstName:      "Akash",
		LastName:       "Biswas",
		Id:             uuid.New(),
		Username:       "abc123",
		HashedPassword: "hashed",
		CreatedAt:      time.Now(),
	}, nil
}
