package hub

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

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

func (h *Hub) SendFriendNotification() {

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
			UserData: dtos.UserPayload{
				Id:        userData.Id.String(),
				Username:  userData.Username,
				FirstName: userData.FirstName,
				LastName:  userData.LastName,
				Rating:    userData.Rating,
			},
		})

		payload = userPayload
	} else {
		userPayload, _ := json.Marshal(PresencePayload{
			PresenceType: PresenceTypeRemoveUser,
			RemoveUserId: userId,
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
		h.handleChallenge(client.UserID, payload.ToUserID)

	case EventChallengeAccept:
		var payload ChallengePayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			client.SendError("Bad payload")
			return
		}
		h.handleAcceptChallenge(payload.FromUserID, client.UserID)

	case EventMakeMove:
		var payload MovePayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			client.SendError("Bad payload")
			return
		}
		h.handleMove(client.UserID, payload)
	case EventSelectPiece:
		var payload SelectPiecePayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			client.SendError("Bad payload")
		}
		h.handleSelectPiece(payload)
	}

}

func (h *Hub) handleChallenge(fromID, toID string) {
	h.mu.RLock()
	targetClient, exists := h.clients[toID]
	h.mu.RUnlock()

	if !exists {
		h.sendErrorTo(fromID, "User is offline")
		return
	}

	payload, _ := json.Marshal(ChallengePayload{FromUserID: fromID, ToUserID: toID})
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
	wPayload, _ := json.Marshal(GameStartPayload{
		GameID:      gameID,
		WhitePlayer: white.UserID,
		BlackPlayer: black.UserID,
	})
	white.Send <- formatEvent(EventGameStart, wPayload)

	// Notify Black player
	bPayload, _ := json.Marshal(GameStartPayload{
		GameID:      gameID,
		WhitePlayer: white.UserID,
		BlackPlayer: black.UserID,
	})

	black.Send <- formatEvent(EventGameStart, bPayload)
}

func (h *Hub) handleMove(playerID string, payload MovePayload) {
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

func (h *Hub) handleSelectPiece(payload SelectPiecePayload) {

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
