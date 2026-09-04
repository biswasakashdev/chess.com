package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/biswasakashdev/chess.com/internal/auth"
	cfg "github.com/biswasakashdev/chess.com/internal/config"
	"github.com/biswasakashdev/chess.com/internal/database"
	"github.com/biswasakashdev/chess.com/internal/hub"
	"github.com/biswasakashdev/chess.com/internal/middleware"
	userRepo "github.com/biswasakashdev/chess.com/internal/repository/users"
	handlers "github.com/biswasakashdev/chess.com/internal/routers"
	"github.com/biswasakashdev/chess.com/internal/service"
	"github.com/biswasakashdev/chess.com/internal/ticket"
	"github.com/go-chi/chi/v5"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
)

const (
	Day time.Duration = time.Hour * 24
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
	sqLiteUserRepo := userRepo.NewSQLiteUserRepository(conn)

	/* Initilise the Hub */

	// Create Services

	userService := service.NewUserService(sqLiteUserRepo)
	authService := auth.NewAuthService(sqLiteUserRepo, config.JwtSecret, Day)
	ticketService := ticket.NewTicketService()
	gameHub := hub.NewHub(ticketService, sqLiteUserRepo)
	friendshipService := service.NewfriendshipService(sqLiteUserRepo, gameHub)

	// Initialise routers
	authRouter := handlers.NewAuthRouter(authService)
	userRouter := handlers.NewUserRouter(userService)
	friendshipRouter := handlers.NewFriendshipRouter(friendshipService)
	ticketRouter := handlers.NewTicketRouter(ticketService, sqLiteUserRepo)
	gameRouter := handlers.NewGameRouter(gameHub)

	// Start the background processing.
	go gameHub.Run()

	/* Public routes */

	// Auth rotes
	router.Mount("/api/v1/auth", authRouter)

	/* Protected routes */

	router.Group(func(r chi.Router) {

		r.Use(middleware.AuthMiddleware(authService))

		r.Mount("/api/v1/users", userRouter)
		r.Mount("/api/v1/tickets", ticketRouter)
		r.Mount("/api/v1/friends", friendshipRouter)
		r.Mount("/api/v1/games", gameRouter)

	})

	// Ws Server
	router.HandleFunc("/api/ws", gameHub.WsHandle)

	fmt.Printf("%s Server listening on %s port.", config.AppName, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router))

}
