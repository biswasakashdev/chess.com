package hub

import "encoding/json"

type EventType string

const (
	EventPresence        EventType = "presence"
	EventChallengeReq    EventType = "challenge_request"
	EventChallengeAccept EventType = "challenge_accept"
	EventGameStart       EventType = "game_start"
	EventMakeMove        EventType = "make_move"
	EventMoveMade        EventType = "move_made"
	EventGameOver        EventType = "game_over"
	EventAvailableMoves  EventType = "available_moves"
	EventError           EventType = "error"
)

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Client payloads
type ChallengePayload struct {
	FromUserID string `json:"from_user_id,omitempty"`
	ToUserID   string `json:"to_user_id"`
}

// Server payloads
type GameStartPayload struct {
	GameID      string `json:"game_id"`
	WhitePlayer string `json:"white_player"`
	BlackPlayer string `json:"black_player"`
	YourColor   string `json:"your_color,omitempty"`
}

type MovePayload struct {
	GameID string `json:"game_id"`
	Move   string `json:"move"` // UCI format e.g., "e2e4"
}

type MoveMadePayload struct {
	GameID string `json:"game_id"`
	Move   string `json:"move"`
	FEN    string `json:"fen"`
}

type GameOverPayload struct {
	GameID string `json:"game_id"`
	Result string `json:"result"`
	Reason string `json:"reason"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}
