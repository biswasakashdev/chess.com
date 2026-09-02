package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/biswasakashdev/chess.com/internal/dtos"
	"github.com/biswasakashdev/chess.com/internal/hub"
	userRepo "github.com/biswasakashdev/chess.com/internal/repository/users"
)

type FriendShipService struct {
	userRepo userRepo.UserRepository
	hb       *hub.Hub
}

func NewfriendshipService(userRepository userRepo.UserRepository, hb *hub.Hub) *FriendShipService {
	return &FriendShipService{
		userRepo: userRepository,
		hb:       hb,
	}
}

var (
	ErrInvalidFriendType = errors.New("invalid friend type: must be 'all', 'active', or 'offline'")
)

func (fs *FriendShipService) FetchAllFriends(ctx context.Context, userId, activeStatus string) ([]*dtos.UserPayload, error) {
	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	friends, err := fs.userRepo.FindFriendsByUserId(newCtx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch friends for user %s: %w", userId, err)
	}
	if len(friends) == 0 {
		return []*dtos.UserPayload{}, nil
	}

	activeClients := fs.hb.GetActiveClients()

	switch activeStatus {
	case "":

		result := make([]*dtos.UserPayload, 0, len(friends))

		for _, user := range friends {
			_, ok := activeClients[user.Id]
			userPayload := &dtos.UserPayload{
				Id:        user.Id,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Username:  user.Username,
				Rating:    user.Rating,
				IsActive:  ok,
			}
			result = append(result, userPayload)

		}
		return result, nil

	case "online":
		filtered := make([]*dtos.UserPayload, 0, len(friends))

		for _, user := range friends {
			if _, isActive := activeClients[user.Id]; isActive {
				filtered = append(filtered, &dtos.UserPayload{
					Id:        user.Id,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Username:  user.Username,
					Rating:    user.Rating,
					IsActive:  true,
				})
			}
		}
		return filtered, nil

	case "offline":
		activeClients := fs.hb.GetActiveClients()
		filtered := make([]*dtos.UserPayload, 0, len(friends))

		for _, user := range friends {
			if _, isActive := activeClients[user.Id]; !isActive {
				filtered = append(filtered, &dtos.UserPayload{
					Id:        user.Id,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Username:  user.Username,
					Rating:    user.Rating,
					IsActive:  false,
				})
			}
		}
		return filtered, nil

	default:
		return nil, ErrInvalidFriendType
	}
}

// SendFriendRequest validates input and creates a friend request.
func (fs *FriendShipService) SendFriendRequest(ctx context.Context, userId, targetUserId string) error {
	if userId == targetUserId {
		return ErrSelfAction
	}

	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := fs.userRepo.SendFriendRequest(newCtx, userId, targetUserId); err != nil {
		return fmt.Errorf("service failed to send friend request: %w", err)
	}

	return nil
}

var (
	ErrInvalidRequestType = errors.New("invalid request type: must be 'sent' or 'received'")
)

func (fs *FriendShipService) FetchAllRequests(ctx context.Context, userId, requestType string) ([]*dtos.UserPayload, error) {
	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var (
		requests []*userRepo.UserDTO
		err      error
	)

	switch requestType {
	case "pending":
		requests, err = fs.userRepo.FindRequestsByUserId(newCtx, userId)
	case "sent":
		requests, err = fs.userRepo.FindRequestsSentByUserId(newCtx, userId)
	default:
		return nil, ErrInvalidRequestType
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s friend requests for user %s: %w", requestType, userId, err)
	}

	// Return an empty slice instead of nil for consistent JSON serialization ([] instead of null)
	if requests == nil {
		return []*dtos.UserPayload{}, nil
	}

	return getUserPayloadList(requests), nil
}

var (
	ErrInvalidFriendshipType = errors.New("invalid friendship status type")
	ErrSelfAction            = errors.New("cannot perform friendship action on yourself")
)

func (fs *FriendShipService) UpdateFriendShipStatus(ctx context.Context, userId, targetUserId, action string) error {
	if userId == targetUserId {
		return ErrSelfAction
	}

	var err error

	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	switch action {
	case "accept":
		err = fs.userRepo.UpdateRequestStatus(newCtx, userId, targetUserId, userRepo.FriendshipStatusAccept)
	case "cancel":
		err = fs.userRepo.DeleteFriendship(newCtx, userId, targetUserId)
	default:
		return ErrInvalidFriendshipType
	}
	if err != nil {
		return fmt.Errorf("service failed to update friendship status: %w", err)
	}

	return nil
}

func (fs *FriendShipService) DeleteRequest(ctx context.Context, userId, targetUserId string) error {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := fs.userRepo.DeleteFriendship(newCtx, userId, targetUserId)

	// Handle or log or return custom error
	if err != nil {
		return err
	}

	return nil
}

func getUserPayloadList(userList []*userRepo.UserDTO) []*dtos.UserPayload {
	var userListResp []*dtos.UserPayload = make([]*dtos.UserPayload, 0, len(userList))

	for _, user := range userList {
		userListResp = append(userListResp, &dtos.UserPayload{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Rating:    user.Rating,
		})
	}

	return userListResp
}
