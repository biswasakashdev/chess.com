package hub

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/biswasakashdev/chess.com/internal/chess"
	"github.com/biswasakashdev/chess.com/internal/dtos"
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
	StopChan    chan string
	mu          sync.RWMutex
}

func NewGameRoom(id string, white, black *Client) *GameRoom {
	return &GameRoom{
		ID:          id,
		WhitePlayer: white,
		BlackPlayer: black,
		Game:        chess.NewGame(),
		MoveChan:    make(chan MoveAction),
		StopChan:    make(chan string),
	}
}

// Handle user connections
func (r *GameRoom) Run() {
	for {
		select {
		case useranme := <-r.StopChan:

			payload := GameOverPayload{
				Reason: "Draw",
				Result: fmt.Sprintf("Username %s cancelled the game", useranme),
			}
			r.broadcast(EventGameOver, payload)
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
				Move: action.MoveUCI,
				FEN:  fen,
			})

			// Check for game completion
			if outcome != chess.NoOutCome {
				payload := GameOverPayload{
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

func (r *GameRoom) GetRoomDetails() *dtos.GameDetails {
	r.mu.RLock()
	defer r.mu.RUnlock()
	turn := "turn_black"
	if r.Game.GetTurn() == chess.TurnWhite {
		turn = "turn_white"
	}
	return &dtos.GameDetails{
		GameId:      r.ID,
		WhitePlayer: r.WhitePlayer.UserID,
		BlackPlayer: r.BlackPlayer.UserID,
		Turn:        turn,
		Board:       r.Game.GetBoard(),
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
