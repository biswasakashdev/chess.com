package hub

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/biswasakashdev/chess.com/internal/chess"
	"github.com/biswasakashdev/chess.com/internal/dtos"
	usersRepo "github.com/biswasakashdev/chess.com/internal/repository/users"
	"github.com/biswasakashdev/chess.com/internal/ticket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/* Upgrade the conn to an websocket conn */
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type Hub struct {
	userRepo usersRepo.UserRepository
	// Ws ticket authentication
	ticketService *ticket.TicketService
	// Online clients
	clients map[string]*Client
	// Activated games
	rooms map[string]*GameRoom
	// Register Chnnel
	Register chan *Client
	// If client close the conn
	Unregister chan *Client
	mu         sync.RWMutex
	rng        *rand.Rand
}

func NewHub(ticketService *ticket.TicketService, userRepo usersRepo.UserRepository) *Hub {
	return &Hub{
		userRepo:      userRepo,
		ticketService: ticketService,
		clients:       make(map[string]*Client),
		rooms:         make(map[string]*GameRoom),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			h.broadcastPresence(client.UserID, true)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			h.broadcastPresence(client.UserID, false)
		}
	}
}

func (h *Hub) GetActiveClients() map[string]struct{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	activeUserList := make(map[string]struct{})
	for ids := range h.clients {
		activeUserList[ids] = struct{}{}
	}

	return activeUserList
}

var (
	ErrorInvalidRoomDetails       error = errors.New("Invalid game room")
	ErrorInvalidUserDetails       error = errors.New("Invalid user id")
	ErrorFailedToFetchUserDetails error = errors.New("Failed to fetch the userdetails")
)

func (h *Hub) GetGameDetails(ctx context.Context, roomId, userId string) (*dtos.GameDetails, error) {
	h.mu.RLock()
	val, ok := h.rooms[roomId]
	h.mu.RUnlock()

	if !ok {
		return nil, ErrorInvalidRoomDetails
	}
	if val.WhitePlayer.UserID != userId && val.BlackPlayer.UserID != userId {
		return nil, ErrorInvalidRoomDetails
	}

	roomDetails := val.GetRoomDetails()

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	userList, err := h.userRepo.FindAllUsersByList(ctx, []string{roomDetails.WhitePlayer, roomDetails.BlackPlayer})

	if err != nil || len(userList) != 2 {
		return nil, ErrorFailedToFetchUserDetails
	}

	var whitePlayer usersRepo.UserDTO
	var blackPlayer usersRepo.UserDTO

	if userList[0].Id == roomDetails.WhitePlayer {
		whitePlayer = *userList[0]
		blackPlayer = *userList[1]
	} else {
		blackPlayer = *userList[0]
		whitePlayer = *userList[1]
	}

	var currTurn dtos.GameTurn

	if roomDetails.Turn == chess.TurnWhite {
		currTurn = dtos.GameTurnWhite
	} else {
		currTurn = dtos.GameTurnBlack
	}

	movesHistory := make([]dtos.GameMove, 0, len(roomDetails.History))

	for i, val := range roomDetails.History {
		movesHistory[i] = dtos.GameMove{
			UserId: val.PlayerID,
			Move:   val.MoveUCI,
		}
	}

	return &dtos.GameDetails{
		GameId: roomDetails.GameId,
		Turn:   currTurn,
		Board:  roomDetails.Board,
		WhitePlayer: dtos.UserResp{
			Id:        whitePlayer.Id,
			FirstName: whitePlayer.FirstName,
			LastName:  whitePlayer.LastName,
			Username:  whitePlayer.Username,
		},
		BlackPlayer: dtos.UserResp{
			Id:        blackPlayer.Id,
			FirstName: blackPlayer.FirstName,
			LastName:  blackPlayer.LastName,
			Username:  blackPlayer.Username,
		},
		History: movesHistory,
	}, nil

}

func (h *Hub) DeleteRoom(roomId, userId string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	val, ok := h.rooms[roomId]

	if !ok {
		return ErrorInvalidRoomDetails
	}
	var username string

	if val.WhitePlayer.UserID == userId {
		username = val.WhitePlayer.Username
	} else if val.BlackPlayer.UserID == userId {
		username = val.BlackPlayer.Username
	} else {
		return ErrorInvalidUserDetails
	}
	val.StopChan <- username
	delete(h.rooms, roomId)
	return nil
}

func (h *Hub) broadcastPresence(userId string, presenseTypeRegister bool) {
	h.mu.RLock()
	onlineUsers := make(map[string]struct{})
	for id := range h.clients {
		onlineUsers[id] = struct{}{}
	}
	h.mu.RUnlock()

	ctx := context.Background()

	userFriendsList, err := h.userRepo.FindFriendsByUserId(ctx, userId)

	if err != nil {
		return
	}

	userOnlineFriends := make([]*usersRepo.UserDTO, 0, 20)
	for _, val := range userFriendsList {
		if _, ok := onlineUsers[val.Id]; ok {
			userOnlineFriends = append(userOnlineFriends, val)
		}
	}

	var payload []byte

	if presenseTypeRegister {
		ctx = context.Background()

		userData, err := h.userRepo.FindById(ctx, userId)
		if err != nil {
			return
		}

		userPayload, _ := json.Marshal(PresencePayload{
			PresenceType: PresenceTypeAddUser,
			UserPayload: UserPayload{
				Id:        userData.Id.String(),
				Username:  userData.Username,
				FirstName: userData.FirstName,
				LastName:  userData.LastName,
			},
		})

		payload = userPayload
	} else {
		userPayload, _ := json.Marshal(PresencePayload{
			PresenceType: PresenceTypeRemoveUser,
			UserPayload: UserPayload{
				Id: userId,
			},
		})

		payload = userPayload
	}

	msg, _ := json.Marshal(Event{Type: EventPresence, Payload: payload})

	h.mu.RLock()
	for _, client := range userOnlineFriends {
		clConn, ok := h.clients[client.Id]
		if !ok {
			continue
		}

		// Non-blocking send or slow-client eviction
		select {
		case clConn.Send <- msg:
		default:
			// Slow client: buffer full. Close connection or drop message
			// to avoid locking up other goroutines.
		}
	}
	h.mu.RUnlock()
}

/*
 ** Handle user created events
 * - Send a challenge
 * - Accept a challenge
 * - Made a move
 * - Select a piece
 * - Send a friend request
 */
func (h *Hub) RouteEvent(client *Client, event Event) {
	switch event.Type {
	case EventChallengeReq:
		var payload ChallengePayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			client.SendError("Bad payload")
			return
		}

		userData := payload.FromUserData

		h.handleChallenge(client.UserID, payload.ToUserID, userData.Username, userData.FirstName, userData.LastName)

	case EventChallengeAccept:
		var payload ChallengePayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			client.SendError("Bad payload")
			return
		}
		h.handleAcceptChallenge(payload.FromUserData.Id, client.UserID)

	case EventMakeMove:
		var payload MakeMovePayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			client.SendError("Bad payload")
			return
		}
		h.handleMove(client.UserID, payload)
	}

}

func (h *Hub) handleChallenge(fromID, toID, username, firstName, lastName string) {
	h.mu.RLock()
	targetClient, exists := h.clients[toID]
	h.mu.RUnlock()

	if !exists {
		h.sendErrorTo(fromID, "User is offline")
		return
	}

	payload, _ := json.Marshal(ChallengePayload{ToUserID: toID,
		FromUserData: UserPayload{
			Id:        fromID,
			Username:  username,
			FirstName: firstName,
			LastName:  lastName,
		}})
	msg, _ := json.Marshal(Event{Type: EventChallengeReq, Payload: payload})
	targetClient.Send <- msg
}

func (h *Hub) handleAcceptChallenge(fromID, toID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientA, existsA := h.clients[fromID]
	clientB, existsB := h.clients[toID]

	if !existsA || !existsB {
		h.sendErrorTo(toID, "Opponent disconnected")
		return
	}

	// Randomly assign colors
	var white, black *Client
	if h.rng.Intn(2) == 0 {
		white, black = clientA, clientB
	} else {
		white, black = clientB, clientA
	}

	gameID := uuid.New().String()
	room := NewGameRoom(gameID, white, black)
	h.rooms[gameID] = room

	go room.Run()

	// Notify White player
	//
	gameStartPayload := GameStartPayload{
		GameID: gameID,
	}
	payload, _ := json.Marshal(gameStartPayload)
	white.Send <- formatEvent(EventGameStart, payload)
	black.Send <- formatEvent(EventGameStart, payload)
}

func (h *Hub) handleMove(playerID string, payload MakeMovePayload) {
	h.mu.RLock()
	room, exists := h.rooms[payload.GameID]
	h.mu.RUnlock()

	if !exists {
		h.sendErrorTo(playerID, "Game not found")
		return
	}

	room.MoveChan <- MoveAction{
		PlayerID: playerID,
		MoveUCI:  payload.Move,
	}
}

func (h *Hub) GetAvailableMoves(roomId, userId, piece string) *dtos.AvailableMovesResp {
	h.mu.RLock()
	val, ok := h.rooms[roomId]
	h.mu.RUnlock()

	if !ok {
		return &dtos.AvailableMovesResp{}
	}

	availableMoves, err := val.GetAllAvailableMoves(userId, piece)

	if err != nil {
		return &dtos.AvailableMovesResp{}
	}

	return &dtos.AvailableMovesResp{
		Moves: availableMoves,
	}
}

func (h *Hub) sendErrorTo(userID, message string) {
	h.mu.RLock()
	client, exists := h.clients[userID]
	h.mu.RUnlock()

	if exists {
		client.SendError(message)
	}
}

func formatEvent(t EventType, payload []byte) []byte {
	bytes, _ := json.Marshal(Event{Type: t, Payload: payload})
	return bytes
}
