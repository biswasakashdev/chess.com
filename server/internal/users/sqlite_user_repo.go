package users

import (
	"context"
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

func (ur *SQLiteUserRepository) IsUsernameExits(username string) (bool, error) {

	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);
	`

	var isExists bool

	err := ur.db.QueryRow(query, username).Scan(&isExists)

	if err != nil {
		return false, err
	}
	return isExists, nil
}

func (ur *SQLiteUserRepository) FindByUsername(username string, ctx context.Context) (*User, error) {

	query := `
		SELECT id, username, hashed_password, first_name, last_name, created_at from users where username = $1;
	`

	var user User

	err := ur.db.QueryRowContext(ctx, query, username).Scan(
		&user.Id,
		&user.Username,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *SQLiteUserRepository) CreateUser(username, hashedPassword, firstName, lastName string) (*User, error) {
	user := &User{
		Id:             uuid.New(),
		Username:       username,
		HashedPassword: hashedPassword,
		FirstName:      firstName,
		LastName:       lastName,
		CreatedAt:      time.Now(),
	}
	query := `
        INSERT INTO users (id, username, hashed_password, first_name, last_name, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := ur.db.Exec(query, user.Id, user.Username, user.HashedPassword, user.FirstName, user.LastName, user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *SQLiteUserRepository) FindById(id string, ctx context.Context) (*User, error) {

	query := `
			SELECT id, username, hashed_password, first_name, last_name, created_at from users where id=$1;
		`

	var user User

	err := ur.db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Username,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
