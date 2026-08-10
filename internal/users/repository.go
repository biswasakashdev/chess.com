package users

type UserRepository interface {
	Save() (*User, error)
	CreateUser(username, hashedPassword, firstName, lastName string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	IsUsernameExits(username string) (bool, error)
}
