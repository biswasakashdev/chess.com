package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/biswasakashdev/chess.com/internal/service/game"
	"github.com/gorilla/websocket"
)

// Configure Gorilla Upgrader with CORS enabled for Vite (localhost:5173)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow React app to connect
	},
}

// JSON Payload structures
type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type MoveMessage struct {
	From Position `json:"from"`
	To   Position `json:"to"`
}

type BoardResponse struct {
	Board [8][8]string `json:"board"`
	Turn  string       `json:"turn"`
}

// Simple Hub to manage connected WebSockets
type Server struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
	// Embed your existing Game logic here!
	board [8][8]string
	turn  string
}

func NewServer() *Server {
	s := &Server{
		clients: make(map[*websocket.Conn]bool),
		turn:    "WHITE",
	}
	s.initSampleBoard()
	return s
}

// Initial dummy setup using unicode symbols for quick visual testing
func (s *Server) initSampleBoard() {
	s.board = [8][8]string{
		{"r", "n", "b", "q", "k", "b", "n", "r"},
		{"p", "p", "p", "p", "p", "p", "p", "p"},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{"P", "P", "P", "P", "P", "P", "P", "P"},
		{"R", "N", "B", "Q", "K", "B", "N", "R"},
	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
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

	// Listen for incoming move commands from React
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.mu.Lock()
			delete(s.clients, conn)
			s.mu.Unlock()
			break
		}

		var move MoveMessage
		if err := json.Unmarshal(msg, &move); err == nil {
			s.mu.Lock()
			// Execute simple move (Swap empty square)
			piece := s.board[move.From.Row][move.From.Col]
			s.board[move.From.Row][move.From.Col] = "."
			s.board[move.To.Row][move.To.Col] = piece

			// Switch turns
			if s.turn == "WHITE" {
				s.turn = "BLACK"
			} else {
				s.turn = "WHITE"
			}
			s.mu.Unlock()

			// Send updated board to all clients
			s.broadcastState()
		}
	}
}

func (s *Server) broadcastState() {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := BoardResponse{Board: s.board, Turn: s.turn}
	data, _ := json.Marshal(resp)

	for client := range s.clients {
		client.WriteMessage(websocket.TextMessage, data)
	}
}

func main() {

	newGame := game.NewGame()

	game.PrintGame(newGame)

	newGame.MakeMove(game.Move{
		From: game.Position{
			Row: 6,
			Col: 2,
		},
		To: game.Position{
			Row: 4,
			Col: 2,
		},
	})

	game.PrintGame(newGame)

	// Server

	server := NewServer()
	http.HandleFunc("/ws", server.handleWS)

	fmt.Println("Go Server listening on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
