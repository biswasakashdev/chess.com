package routers

import (
	"errors"
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/middleware"
	usersReposit "github.com/biswasakashdev/chess.com/internal/repository/users"
	"github.com/biswasakashdev/chess.com/internal/ticket"
	"github.com/biswasakashdev/chess.com/internal/util"
	"github.com/go-chi/chi/v5"
)

type ticketsHandler struct {
	ticketService *ticket.TicketService
	userRepo      usersReposit.UserRepository
}

func NewTicketRouter(ticketService *ticket.TicketService, userRepo usersReposit.UserRepository) chi.Router {

	ticketHandl := ticketsHandler{
		ticketService: ticketService,
		userRepo:      userRepo,
	}
	router := chi.NewRouter()

	router.Post("/", ticketHandl.CreateTicket)

	return router
}

func (th *ticketsHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := util.GetContextValue(ctx, middleware.UserIDKey)

	if err != nil {
		util.BuildErrResponse(w, errors.New("Invalid UserID"))
		return
	}

	user, err := th.userRepo.FindById(ctx, userID)
	if err != nil {
		util.BuildErrResponse(w, errors.New("Unavailable"))
	}

	tkt, err := th.ticketService.InitTicket(user.Id.String(), user.Username)
	if err != nil {
		util.BuildErrResponse(w, err)
		return
	}

	util.BuildResponseBodyWithCode(w, map[string]string{
		"ticket": tkt,
	}, http.StatusCreated)
}
