package hub

import (
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[string]*Client
	rooms      map[string]*GameRoom
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
	rng        *rand.Rand
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		rooms:      make(map[string]*GameRoom),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			h.broadcastPresence()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			h.broadcastPresence()
		}
	}
}

func (h *Hub) broadcastPresence() {
	h.mu.RLock()
	var onlineUsers []string
	for id := range h.clients {
		onlineUsers = append(onlineUsers, id)
	}
	h.mu.RUnlock()

	payload, _ := json.Marshal(onlineUsers)
	msg, _ := json.Marshal(Event{Type: EventPresence, Payload: payload})

	h.mu.RLock()
	for _, client := range h.clients {
		select {
		case client.Send <- msg:
		default:
		}
	}
	h.mu.RUnlock()
}

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
		YourColor:   "white",
	})
	white.Send <- formatEvent(EventGameStart, wPayload)

	// Notify Black player
	bPayload, _ := json.Marshal(GameStartPayload{
		GameID:      gameID,
		WhitePlayer: white.UserID,
		BlackPlayer: black.UserID,
		YourColor:   "black",
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
