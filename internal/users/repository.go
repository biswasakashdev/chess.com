package users

import "context"

type UserRepository interface {
	CreateUser(username, hashedPassword, firstName, lastName string) (*User, error)
	FindByUsername(username string, ctx context.Context) (*User, error)
	FindById(id string, ctx context.Context) (*User, error)
	IsUsernameExits(username string) (bool, error)
}
