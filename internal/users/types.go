package users

// Player status
type PlayerStatus string

var (
	PlayerStatusOffLine PlayerStatus = "off-line"
	PlayerStatusIdle    PlayerStatus = "idle"
)

type PlayerStatusPayload struct {
	UserId string       `json:"userID"`
	Status PlayerStatus `json:"status"`
}

type UserPayload struct {
	Id        string       `json:"userId"`
	Username  string       `json:"username"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	Rating    int          `json:"rating"`
	Status    PlayerStatus `json:"status"`
}
