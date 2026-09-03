package service

import (
	"context"
	"time"

	"github.com/biswasakashdev/chess.com/internal/dtos"
	userRepo "github.com/biswasakashdev/chess.com/internal/repository/users"
	"github.com/biswasakashdev/chess.com/internal/util"
)

type UserService struct {
	userRepo userRepo.UserRepository
}

func NewUserService(userRepo userRepo.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (us *UserService) GetUserByUsernameNotFriendWith(ctx context.Context, userId, searchQuery string) ([]*dtos.UserResp, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	users, err := us.userRepo.FindByUsernameNotFriendsWith(newCtx, userId, searchQuery)

	if err != nil {
		return nil, err
	}

	return getUserPayloadList(users), nil
}

func (us *UserService) GetUserByUsername(username string, ctx context.Context) (*dtos.UserResp, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	user, err := us.userRepo.FindByUsername(newCtx, username)

	if err != nil {
		return nil, util.InternalError
	}

	userResp := dtos.UserResp{
		Id:        user.Id.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}

	return &userResp, nil
}

func (us *UserService) GetUserById(id string, ctx context.Context) (*dtos.UserResp, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	user, err := us.userRepo.FindById(newCtx, id)

	if err != nil {
		return nil, util.InternalError
	}

	userResp := dtos.UserResp{
		Id:        user.Id.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}

	return &userResp, nil
}
