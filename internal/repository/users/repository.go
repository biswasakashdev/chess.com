package users

import (
	"context"

	mdls "github.com/biswasakashdev/chess.com/internal/models"
)

type UserRepository interface {
	CreateUser(username, hashedPassword, firstName, lastName string) (*mdls.User, error)

	FindByUsername(ctx context.Context, username string) (*mdls.User, error)

	FindById(ctx context.Context, id string) (*mdls.User, error)

	IsUsernameExits(ctx context.Context, username string) (bool, error)

	SendFriendRequest(ctx context.Context, userId, targetUserId string) error

	FindRequestsByUserId(ctx context.Context, userId string) ([]*UserDTO, error)

	FindRequestsSentByUserId(ctx context.Context, userId string) ([]*UserDTO, error)

	// Accept, Reject
	UpdateRequestStatus(ctx context.Context, userId, targetUserId string, status FriendshipStatus) error

	DeleteFriendship(ctx context.Context, userId, targetUserId string) error

	FindFriendsByUserId(ctx context.Context, userId string) ([]*UserDTO, error)

	FindByUsernameNotFriendsWith(ctx context.Context, userId, username string) ([]*UserDTO, error)

	FindAllUsersByList(ctx context.Context, userIdList []string) ([]*UserDTO, error)
}
