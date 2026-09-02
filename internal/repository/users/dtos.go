package users

type FriendshipStatus string

var (
	FriendshipStatusBlocked FriendshipStatus = "blocked"
	FriendshipStatusPending FriendshipStatus = "pending"
	FriendshipStatusAccept  FriendshipStatus = "accepted"
)

type UserDTO struct {
	Id        string
	Username  string
	FirstName string
	LastName  string
	Rating    int
}
