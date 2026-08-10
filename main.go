package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/auth"
	"github.com/biswasakashdev/chess.com/internal/database"
	"github.com/biswasakashdev/chess.com/internal/game"
	"github.com/biswasakashdev/chess.com/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	conn := database.Connect()

	router := chi.NewRouter()

	// Logger middleware.
	router.Use(middleware.Logger)

	// Create Repositories
	userRepository := users.NewSQLiteUserRepository(conn)

	// Initialise handlers
	authHandler := auth.NewAuthHandler(userRepository)
	userHandler := users.NewUserHandler(userRepository)

	/* Public routes */

	// Auth rotes
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/", authHandler.Login)
	})

	router.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", userHandler.GetUser)
	})

	// Ws Server
	gameServer := game.NewGameServer()

	router.Handle("/ws", gameServer)

	fmt.Println("Go Server listening on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", router))
}
