package hub

import (
	"encoding/json"

	"github.com/biswasakashdev/chess.com/internal/dtos"
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

// Client payloads
type ChallengePayload struct {
	FromUserID string `json:"from_user_id,omitempty"`
	ToUserID   string `json:"to_user_id"`
}

type MovePayload struct {
	GameID string `json:"game_id"`
	Move   string `json:"move"` // UCI format e.g., "e2e4"
}

/*
Server payloads
*/

type GameStartPayload struct {
	GameID      string `json:"game_id"`
	WhitePlayer string `json:"white_player"`
	BlackPlayer string `json:"black_player"`
}

type MoveMadePayload struct {
	Move string `json:"move"`
	FEN  string `json:"fen"`
}

type GameOverPayload struct {
	Result string `json:"result"`
	Reason string `json:"reason"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type SelectPiecePayload struct {
	GameId string `json:"gameId"`
	UserID string `json:"userId"`
	Piece  string `json:"piece"`
}

type AvailableMovesPayload struct {
	ForUserId string   `json:"forUserId"`
	Moves     []string `json:"moves"`
}

type PresenceType string

var (
	PresenceTypeAddUser    PresenceType = "add_user"
	PresenceTypeRemoveUser PresenceType = "remove_user"
)

type PresencePayload struct {
	PresenceType PresenceType     `json:"presence_type"`
	RemoveUserId string           `json:"remove_user_id,omitempty"`
	UserData     dtos.UserPayload `json:"user_data"`
}
