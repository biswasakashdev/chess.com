package hub

import (
	"encoding/json"
)

type EventType string

const (

	// Client event

	EventMakeMove        EventType = "make_move"
	EventChallengeReq    EventType = "challenge_request"
	EventChallengeAccept EventType = "challenge_accept"
	EventSelectPiece     EventType = "select_piece"

	// Sever Events

	EventMoveMade       EventType = "move_made"
	EventGameStart      EventType = "game_start"
	EventPresence       EventType = "presence"
	EventPlayerStatus   EventType = "player_status"
	EventGameOver       EventType = "game_over"
	EventAvailableMoves EventType = "available_moves"
	EventError          EventType = "error"
)

// Event
type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type UserPayload struct {
	Id        string `json:"id"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// Client payloads
type ChallengePayload struct {
	ToUserID     string      `json:"to_user_id,omitempty"`
	FromUserData UserPayload `json:"from_user_data"`
}

type MakeMovePayload struct {
	GameID string `json:"game_id"`
	Move   string `json:"move"` // UCI format e.g., "e2e4"
}

type SelectPiecePayload struct {
	GameId string `json:"gameId"`
	UserID string `json:"userId"`
	Piece  string `json:"piece"`
}

/*
Server payloads
*/

type GameStartPayload struct {
	GameID string `json:"game_id"`
}

type MoveMadePayload struct {
	UserId string `json:"user_id"`
	Move   string `json:"move"`
	FEN    string `json:"fen"`
}

type GameOverPayload struct {
	Result string `json:"result"`
	Reason string `json:"reason"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type AvailableMovesPayload struct {
	Moves []string `json:"moves"`
}

type PresenceType string

var (
	PresenceTypeAddUser    PresenceType = "add_user"
	PresenceTypeRemoveUser PresenceType = "remove_user"
)

type PresencePayload struct {
	PresenceType PresenceType `json:"presence_type"`
	UserPayload  `json:"user_data"`
}
