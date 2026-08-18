package users

import (
	"context"
	"time"

	"github.com/biswasakashdev/chess.com/server/internal/util"
)

type UserService struct {
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
	user, err := us.userRepo.FindByUsername(username, newCtx)

	if err != nil {
		return nil, util.InternalError
	}

	return user, nil
}

func (us *UserService) GetUserById(id string, ctx context.Context) (*User, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	user, err := us.userRepo.FindById(id, newCtx)

	if err != nil {
		return nil, util.InternalError
	}

	return user, nil
}
