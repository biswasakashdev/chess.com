package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/biswasakashdev/chess.com/internal/auth"
	cfg "github.com/biswasakashdev/chess.com/internal/config"
	"github.com/biswasakashdev/chess.com/internal/database"
	"github.com/biswasakashdev/chess.com/internal/game"
	"github.com/biswasakashdev/chess.com/internal/middleware"
	"github.com/biswasakashdev/chess.com/internal/users"
	"github.com/go-chi/chi/v5"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
)

const (
	Day time.Duration = time.Hour * 2
)

func main() {

	config := cfg.Load()

	conn := database.Connect()

	router := chi.NewRouter()

	// Logger middleware.
	router.Use(chiMiddle.Logger)

	// Create Repositories
	sqLiteUserRepo := users.NewSQLiteUserRepository(conn)

	// Create Services

	userService := users.NewUserService(sqLiteUserRepo)
	authService := auth.NewAuthService(sqLiteUserRepo, config.JwtSecret, Day)

	// Initialise handlers
	authHandler := auth.NewAuthHandler(authService)
	userHandler := users.NewUserHandler(userService)

	/* Public routes */

	// Auth rotes
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/", authHandler.Login)
	})

	/* Protected routes */

	router.Group(func(r chi.Router) {

		r.Use(middleware.AuthMiddleware(authService))

		r.Route("/api/v1/users", func(r chi.Router) {
			r.Get("/", userHandler.GetUser)
		})

		// Ws Server
		gameServer := game.NewGameServer()

		r.Handle("/ws", gameServer)

	})

	fmt.Printf("%s Server listening on %s port.", config.AppName, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router))

}
