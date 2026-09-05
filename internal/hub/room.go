package hub

import (
	"encoding/json"
	"errors"
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
	Moves       []MoveAction
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
		Moves:       make([]MoveAction, 0),
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

			r.Moves = append(r.Moves, action)

			// Return the state of the board
			outcome := r.Game.GetOutCome()
			r.mu.Unlock()

			currTurn := dtos.GameTurnBlack

			if r.Game.GetTurn() == chess.TurnWhite {
				currTurn = dtos.GameTurnWhite
			}

			// Broadcast move to both players
			r.broadcast(EventMoveMade, MoveMadePayload{
				UserId: action.PlayerID,
				Turn:   currTurn,
				Move:   action.MoveUCI,
			})

			// Check for game completion
			if outcome != chess.NoOutCome {
				payload := GameOverPayload{
					GameId: r.ID,
					Reason: "Win",
				}
				if outcome == chess.BlackWin {
					payload.Reason = "Black wins"
				} else {
					payload.Reason = "White wins"
				}
				r.broadcast(EventGameOver, payload)
				return
			}
		}
	}
}

type RoomDetails struct {
	GameId      string
	Board       string
	Turn        chess.Turn
	WhitePlayer string
	BlackPlayer string
	History     []MoveAction
}

func (r *GameRoom) GetRoomDetails() *RoomDetails {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return &RoomDetails{
		GameId:      r.ID,
		WhitePlayer: r.WhitePlayer.UserID,
		BlackPlayer: r.BlackPlayer.UserID,
		Turn:        r.Game.GetTurn(),
		Board:       r.Game.GetBoard(),
		History:     r.Moves,
	}

}

func (r *GameRoom) GetAllAvailableMoves(userId, moveUCI string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	turn := r.Game.GetTurn()

	if (turn == chess.TurnWhite && userId != r.WhitePlayer.UserID) ||
		(turn == chess.TurnBlack && userId != r.BlackPlayer.UserID) {
		return "", errors.New("Invalid user or gameId")
	}
	return r.Game.GetAvailableMoves(moveUCI), nil
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
