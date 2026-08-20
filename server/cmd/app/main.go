package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/biswasakashdev/chess.com/server/internal/auth"
	cfg "github.com/biswasakashdev/chess.com/server/internal/config"
	"github.com/biswasakashdev/chess.com/server/internal/database"
	"github.com/biswasakashdev/chess.com/server/internal/handlers"
	"github.com/biswasakashdev/chess.com/server/internal/hub"
	"github.com/biswasakashdev/chess.com/server/internal/middleware"
	"github.com/biswasakashdev/chess.com/server/internal/users"
	"github.com/go-chi/chi/v5"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

const (
	Day time.Duration = time.Hour * 24
)

/* Upgrade the conn to an websocket conn */
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func main() {

	config := cfg.Load()

	/* Initialise the database connection */
	conn := database.Connect()

	/* Create schema */
	database.InitSchema(conn)

	/* Create the router */
	router := chi.NewRouter()

	// Logger middleware.
	router.Use(chiMiddle.Logger)

	// Create Repositories
	sqLiteUserRepo := users.NewSQLiteUserRepository(conn)

	// Create Services

	userService := users.NewUserService(sqLiteUserRepo)
	authService := auth.NewAuthService(sqLiteUserRepo, config.JwtSecret, Day)

	// Initialise handlers
	authRouter := handlers.NewAuthRouter(authService)
	userRouter := handlers.NewUserRouter(userService)

	/* Initilise the Hub */

	gameHub := hub.NewHub()

	// Start the background processing.
	go gameHub.Run()

	/* Public routes */

	// Auth rotes
	router.Mount("/api/v1/auth", authRouter)

	/* Protected routes */

	router.Group(func(r chi.Router) {

		r.Use(middleware.AuthMiddleware(authService))

		r.Mount("/api/v1/users", userRouter)

		// Ws Server

		r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			userID := r.URL.Query().Get("userId")
			if userID == "" {
				http.Error(w, "userId query parameter is required", http.StatusBadRequest)
				return
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println("Upgrade failed:", err)
				return
			}

			client := &hub.Client{
				Hub:    gameHub,
				Conn:   conn,
				UserID: userID,
				Send:   make(chan []byte, 256),
			}

			gameHub.Register <- client

			go client.WritePump()
			go client.ReadPump()
		})

	})

	fmt.Printf("%s Server listening on %s port.", config.AppName, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router))

}
