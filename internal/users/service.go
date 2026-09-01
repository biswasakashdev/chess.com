package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/biswasakashdev/chess.com/internal/hub"
	"github.com/biswasakashdev/chess.com/internal/util"
)

type UserService struct {
	hb       hub.Hub
	userRepo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (us *UserService) GetUserByUsername(username string, ctx context.Context) (*User, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	user, err := us.userRepo.FindByUsername(newCtx, username)

	if err != nil {
		return nil, util.InternalError
	}

	return user, nil
}

func (us *UserService) GetUserById(id string, ctx context.Context) (*User, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	user, err := us.userRepo.FindById(newCtx, id)

	if err != nil {
		return nil, util.InternalError
	}

	return user, nil
}

func (us *UserService) GetUserByUsernameNotFriendWith(ctx context.Context, userId, username string) ([]*UserPayload, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	users, err := us.userRepo.FindByUsernameNotFriendsWith(newCtx, userId, username)

	if err != nil {
		return nil, err
	}

	return users, nil
}

var (
	ErrInvalidFriendType = errors.New("invalid friend type: must be 'all', 'active', or 'offline'")
)

func (us *UserService) FetchAllFriends(ctx context.Context, userId, friendType string) ([]*UserPayload, error) {
	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	friends, err := us.userRepo.FindFriendsByUserId(newCtx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch friends for user %s: %w", userId, err)
	}

	if len(friends) == 0 {
		return []*UserPayload{}, nil
	}

	switch friendType {
	case "", "all":
		return friends, nil

	case "active", "online":
		activeClients := us.hb.GetActiveClients()
		filtered := make([]*UserPayload, 0, len(friends))

		for _, friend := range friends {
			if _, isActive := activeClients[friend.Id]; isActive {
				filtered = append(filtered, friend)
			}
		}
		return filtered, nil

	case "offline":
		activeClients := us.hb.GetActiveClients()
		filtered := make([]*UserPayload, 0, len(friends))

		for _, friend := range friends {
			if _, isActive := activeClients[friend.Id]; !isActive {
				filtered = append(filtered, friend)
			}
		}
		return filtered, nil

	default:
		return nil, ErrInvalidFriendType
	}
}

// SendFriendRequest validates input and creates a friend request.
func (us *UserService) SendFriendRequest(ctx context.Context, userId, targetUserId string) error {
	if userId == targetUserId {
		return ErrSelfAction
	}

	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := us.userRepo.SendFriendRequest(newCtx, userId, targetUserId); err != nil {
		return fmt.Errorf("service failed to send friend request: %w", err)
	}

	return nil
}

var (
	ErrInvalidRequestType = errors.New("invalid request type: must be 'sent' or 'received'")
)

func (us *UserService) FetchAllRequests(ctx context.Context, userId, requestType string) ([]*UserPayload, error) {
	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var (
		requests []*UserPayload
		err      error
	)

	switch requestType {
	case "sent", "outgoing":
		requests, err = us.userRepo.FindRequestsSentByUserId(newCtx, userId)
	case "received", "incoming":
		requests, err = us.userRepo.FindRequestsByUserId(newCtx, userId)
	default:
		return nil, ErrInvalidRequestType
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s friend requests for user %s: %w", requestType, userId, err)
	}

	// Return an empty slice instead of nil for consistent JSON serialization ([] instead of null)
	if requests == nil {
		requests = []*UserPayload{}
	}

	return requests, nil
}

var (
	ErrInvalidFriendshipType = errors.New("invalid friendship status type")
	ErrSelfAction            = errors.New("cannot perform friendship action on yourself")
)

func (us *UserService) UpdateFriendShipStatus(ctx context.Context, userId, targetUserId, action string) error {
	if userId == targetUserId {
		return ErrSelfAction
	}

	var dbStatus string
	switch action {
	case "accept", "accepted":
		dbStatus = "accepted" // matches database enum/text value
	case "block", "blocked":
		dbStatus = "blocked"
	default:
		return ErrInvalidFriendshipType
	}

	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := us.userRepo.UpdateRequestStatus(newCtx, userId, targetUserId, dbStatus); err != nil {
		return fmt.Errorf("service failed to update friendship status: %w", err)
	}

	return nil
}

func (us *UserService) DeleteRequest(ctx context.Context, userId, targetUserId string) error {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := us.userRepo.DeleteRequest(newCtx, userId, targetUserId)

	// Handle or log or return custom error
	if err != nil {
		return err
	}

	return nil
}
