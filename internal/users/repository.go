package users

import "context"

type UserRepository interface {
	CreateUser(username, hashedPassword, firstName, lastName string) (*User, error)

	FindByUsername(ctx context.Context, username string) (*User, error)

	FindById(ctx context.Context, id string) (*User, error)

	IsUsernameExits(ctx context.Context, username string) (bool, error)

	SendFriendRequest(ctx context.Context, userId, targetUserId string) error

	FindRequestsByUserId(ctx context.Context, userId string) ([]*UserPayload, error)

	FindRequestsSentByUserId(ctx context.Context, userId string) ([]*UserPayload, error)

	// Accept, Reject
	UpdateRequestStatus(ctx context.Context, userId, targetUserId, reqStatus string) error

	DeleteRequest(ctx context.Context, userId, targetUserId string) error

	FindFriendsByUserId(ctx context.Context, userId string) ([]*UserPayload, error)

	FindByUsernameNotFriendsWith(ctx context.Context, userId, username string) ([]*UserPayload, error)
}
