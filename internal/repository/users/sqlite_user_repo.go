package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	mdls "github.com/biswasakashdev/chess.com/internal/models"
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

func (ur *SQLiteUserRepository) IsUsernameExits(ctx context.Context, username string) (bool, error) {

	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);
	`

	var isExists bool

	err := ur.db.QueryRowContext(ctx, query, username).Scan(&isExists)

	if err != nil {
		return false, err
	}
	return isExists, nil
}

func (ur *SQLiteUserRepository) FindByUsername(ctx context.Context, username string) (*mdls.User, error) {

	query := `
		SELECT id, username, hashed_password, first_name, last_name, created_at from users where username = $1;
	`

	var user mdls.User

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

func (ur *SQLiteUserRepository) CreateUser(username, hashedPassword, firstName, lastName string) (*mdls.User, error) {
	user := &mdls.User{
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

func (ur *SQLiteUserRepository) FindById(ctx context.Context, id string) (*mdls.User, error) {

	query := `
			SELECT id, username, hashed_password, first_name, last_name, created_at from users where id=$1;
		`

	var user mdls.User

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

// GetFriendsByUserID retrieves all accepted friends for a given user ID.
func (ur *SQLiteUserRepository) FindFriendsByUserId(ctx context.Context, userId string) ([]*UserDTO, error) {
	// Query searches both directions (user_id -> friend_id OR friend_id -> user_id)
	query := `
		SELECT
			u.id,
			u.username,
			u.first_name,
			u.last_name,
			u.rating
		FROM users u
		INNER JOIN friendships f
			ON (f.friend_id = u.id AND f.user_id = ?)
			OR (f.user_id = u.id AND f.friend_id = ?)
		WHERE f.status = 'accepted'
	`

	rows, err := ur.db.QueryContext(ctx, query, userId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query friends for user %s: %w", userId, err)
	}
	defer rows.Close()

	var friends []*UserDTO
	for rows.Next() {
		var userPlayload UserDTO
		if err := rows.Scan(&userPlayload.Id, &userPlayload.Username, &userPlayload.FirstName, &userPlayload.LastName, &userPlayload.Rating); err != nil {

			return nil, fmt.Errorf("failed to scan friend row: %w", err)
		}
		friends = append(friends, &userPlayload)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return friends, nil
}

func (ur *SQLiteUserRepository) FindByUsernameNotFriendsWith(ctx context.Context, userId, username string) ([]*UserDTO, error) {
	query := `
			SELECT
				u.id,
				u.username,
				u.first_name,
				u.last_name,
				u.rating
			FROM users u
			WHERE u.id != ?
			  AND u.username LIKE ?
			  AND u.id NOT IN (
				  SELECT friend_id FROM friendships WHERE user_id = ?
				  UNION
				  SELECT user_id FROM friendships WHERE friend_id = ?
			  )
			ORDER BY u.username ASC
			LIMIT 20
		`

	// Using wildcards for a prefix/partial search
	searchPattern := "%" + username + "%"

	rows, err := ur.db.QueryContext(ctx, query, userId, searchPattern, userId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to search non-friend users: %w", err)
	}
	defer rows.Close()

	var users []*UserDTO
	for rows.Next() {
		user := &UserDTO{}
		if err := rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Rating); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

var (
	ErrFriendshipAlreadyExists = errors.New("a friend request or friendship already exists between these users")
)

// SendFriendRequest inserts a pending friendship record if no relationship exists in either direction.
func (ur *SQLiteUserRepository) SendFriendRequest(ctx context.Context, userId, targetUserId string) error {
	// 1. Check if any relationship already exists in either direction
	checkQuery := `
		SELECT status FROM friendships
		WHERE (user_id = ? AND friend_id = ?)
		   OR (user_id = ? AND friend_id = ?)
		LIMIT 1
	`
	var existingStatus string
	err := ur.db.QueryRowContext(ctx, checkQuery, userId, targetUserId, targetUserId, userId).Scan(&existingStatus)
	if err == nil {
		return ErrFriendshipAlreadyExists
	}

	// 2. Insert the pending request
	insertQuery := `
		INSERT INTO friendships (user_id, friend_id, status)
		VALUES (?, ?, 'pending')
	`
	_, err = ur.db.ExecContext(ctx, insertQuery, userId, targetUserId)
	if err != nil {
		return fmt.Errorf("failed to insert friend request: %w", err)
	}

	return nil
}

// FindRequestsByUserId retrieves all users who have sent a pending friend request TO userId (Incoming Requests).
func (ur *SQLiteUserRepository) FindRequestsByUserId(ctx context.Context, userId string) ([]*UserDTO, error) {
	query := `
		SELECT
			u.id,
			u.username,
			u.first_name,
			u.last_name,
			u.rating
		FROM users u
		INNER JOIN friendships f ON f.user_id = u.id
		WHERE f.friend_id = ?
		  AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`

	rows, err := ur.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query incoming friend requests for user %s: %w", userId, err)
	}
	defer rows.Close()

	var requests []*UserDTO
	for rows.Next() {
		user := &UserDTO{}
		if err := rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Rating); err != nil {
			return nil, fmt.Errorf("failed to scan incoming request user row: %w", err)
		}
		requests = append(requests, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("incoming requests rows iteration error: %w", err)
	}

	return requests, nil
}

// FindRequestsSentByUserId retrieves all users to whom userId has sent a request (Outgoing Requests).
func (ur *SQLiteUserRepository) FindRequestsSentByUserId(ctx context.Context, userId string) ([]*UserDTO, error) {
	query := `
		SELECT
			u.id,
			u.username,
			u.first_name,
			u.last_name,
			u.rating
		FROM users u
		INNER JOIN friendships f ON f.friend_id = u.id
		WHERE f.user_id = ?
		  AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`

	rows, err := ur.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query outgoing friend requests sent by user %s: %w", userId, err)
	}
	defer rows.Close()

	var sentRequests []*UserDTO
	for rows.Next() {
		user := &UserDTO{}
		if err := rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Rating); err != nil {
			return nil, fmt.Errorf("failed to scan outgoing request user row: %w", err)
		}
		sentRequests = append(sentRequests, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("outgoing requests rows iteration error: %w", err)
	}

	return sentRequests, nil
}

var (
	ErrRequestNotFound = errors.New("friend request not found")
)

// UpdateRequestStatus updates the friendship status (e.g., 'accepted', 'blocked').
// userId is the recipient acting on the request, targetUserId is the user who sent/holds the other end.
func (ur *SQLiteUserRepository) UpdateRequestStatus(ctx context.Context, userId, targetUserId string, status FriendshipStatus) error {
	query := `
		UPDATE friendships
		SET status = ?
		WHERE (user_id = ? AND friend_id = ?)
		   OR (user_id = ? AND friend_id = ?)
	`

	res, err := ur.db.ExecContext(ctx, query, status, targetUserId, userId, userId, targetUserId)
	if err != nil {
		return fmt.Errorf("failed to update request status: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrRequestNotFound
	}

	return nil
}

// DeleteRequest removes the friendship entry (used for rejecting requests, canceling sent requests, or unfriending).
func (ur *SQLiteUserRepository) DeleteFriendship(ctx context.Context, userId, targetUserId string) error {
	query := `
		DELETE FROM friendships
		WHERE (user_id = ? AND friend_id = ?)
		   OR (user_id = ? AND friend_id = ?)
	`

	res, err := ur.db.ExecContext(ctx, query, userId, targetUserId, targetUserId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete friendship record: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrRequestNotFound
	}

	return nil
}
