package dtos

// Player status
type PlayerStatus string

var (
	PlayerStatusOffLine PlayerStatus = "off-line"
	PlayerStatusIdle    PlayerStatus = "idle"
)

type PlayerStatusPayload struct {
	UserId string       `json:"id"`
	Status PlayerStatus `json:"status"`
}

type UserPayload struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Rating    int    `json:"rating,omitempty"`
	IsActive  bool   `json:"isActive,omitempty"`
}
