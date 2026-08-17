package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/biswasakashdev/chess.com/server/internal/auth"
	cfg "github.com/biswasakashdev/chess.com/server/internal/config"
	"github.com/biswasakashdev/chess.com/server/internal/database"
	"github.com/biswasakashdev/chess.com/server/internal/game"
	"github.com/biswasakashdev/chess.com/server/internal/handlers"
	"github.com/biswasakashdev/chess.com/server/internal/middleware"
	"github.com/biswasakashdev/chess.com/server/internal/users"
	"github.com/go-chi/chi/v5"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
)

const (
	Day time.Duration = time.Hour * 2
)

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

	/* Public routes */

	// Auth rotes
	router.Mount("/api/v1/auth", authRouter)

	/* Protected routes */

	router.Group(func(r chi.Router) {

		r.Use(middleware.AuthMiddleware(authService))

		r.Mount("/api/v1/users", userRouter)

		// Ws Server
		gameServer := game.NewGameServer()

		r.Handle("/ws", gameServer)

	})

	fmt.Printf("%s Server listening on %s port.", config.AppName, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router))

}
