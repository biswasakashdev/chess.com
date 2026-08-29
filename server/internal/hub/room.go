package hub

import (
	"encoding/json"
	"sync"

	"github.com/biswasakashdev/chess.com/server/internal/chess"
)

type MoveAction struct {
	PlayerID string
	MoveUCI  string
}

type GameRoom struct {
	ID          string
	WhitePlayer *Client
	BlackPlayer *Client
	Game        chess.Chess
	MoveChan    chan MoveAction
	StopChan    chan struct{}
	mu          sync.Mutex
}

func NewGameRoom(id string, white, black *Client) *GameRoom {
	return &GameRoom{
		ID:          id,
		WhitePlayer: white,
		BlackPlayer: black,
		Game:        chess.NewGame(),
		MoveChan:    make(chan MoveAction),
		StopChan:    make(chan struct{}),
	}
}

func (r *GameRoom) Run() {
	for {
		select {
		case <-r.StopChan:
			return
		case action := <-r.MoveChan:
			r.mu.Lock()

			// Check player turn
			turn := r.Game.GetTurn()
			if (turn == chess.TurnWhite && action.PlayerID != r.WhitePlayer.UserID) ||
				(turn == chess.TurnBlack && action.PlayerID != r.BlackPlayer.UserID) {
				r.mu.Unlock()
				r.sendToPlayer(action.PlayerID, EventError, ErrorPayload{Message: "Not your turn"})
				continue
			}

			// Validate and execute move
			err := r.Game.MakeMove(action.MoveUCI)
			if err != nil {
				r.mu.Unlock()
				r.sendToPlayer(action.PlayerID, EventError, ErrorPayload{Message: "Invalid move: " + err.Error()})
				continue
			}

			// Return the state of the board
			fen := r.Game.GetBoard()
			outcome := r.Game.GetOutCome()
			r.mu.Unlock()

			// Broadcast move to both players
			r.broadcast(EventMoveMade, MoveMadePayload{
				GameID: r.ID,
				Move:   action.MoveUCI,
				FEN:    fen,
			})

			// Check for game completion
			if outcome != chess.NoOutCome {
				payload := GameOverPayload{
					GameID: r.ID,
					Reason: "Win",
				}
				if outcome == chess.BlackWin {
					payload.Reason = ""
				} else {
					payload.Reason = ""
				}
				r.broadcast(EventGameOver, payload)
				return
			}
		}
	}
}

func (r *GameRoom) broadcast(eventType EventType, payload any) {
	bytes, _ := json.Marshal(payload)
	msg, _ := json.Marshal(Event{Type: eventType, Payload: bytes})

	r.WhitePlayer.Send <- msg
	r.BlackPlayer.Send <- msg
}

func (r *GameRoom) sendToPlayer(userID string, eventType EventType, payload any) {
	bytes, _ := json.Marshal(payload)
	msg, _ := json.Marshal(Event{Type: eventType, Payload: bytes})

	if r.WhitePlayer.UserID == userID {
		r.WhitePlayer.Send <- msg
	} else if r.BlackPlayer.UserID == userID {
		r.BlackPlayer.Send <- msg
	}
}
