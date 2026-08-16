package game

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Configure Gorilla Upgrader with CORS enabled for Vite (localhost:5173)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow React app to connect
	},
}

// Simple Hub to manage connected WebSockets
type Server struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
	game    *Game
}

// TODO: Update the server configuration.
func NewGameServer() *Server {
	s := &Server{
		clients: make(map[*websocket.Conn]bool),
		game:    newGame(),
	}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	// Send current state immediately on connect
	s.broadcastState()

	// Listen for incoming move commands from Clients
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.mu.Lock()
			delete(s.clients, conn)
			s.mu.Unlock()
			break
		}

		var move Move
		if err := json.Unmarshal(msg, &move); err == nil {
			s.mu.Lock()
			s.game.MakeMove(move)
			s.mu.Unlock()

			// Send updated board to all clients
			s.broadcastState()
		}
	}
}

func (s *Server) broadcastState() {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := BoardResponse{Board: s.game.Board.ToBoardRes(), Turn: s.game.Turn.String()}
	data, _ := json.Marshal(resp)

	for client := range s.clients {
		client.WriteMessage(websocket.TextMessage, data)
	}
}
