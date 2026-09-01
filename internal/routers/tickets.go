package handlers

import (
	"errors"
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/middleware"
	"github.com/biswasakashdev/chess.com/internal/ticket"
	"github.com/biswasakashdev/chess.com/internal/util"
	"github.com/go-chi/chi/v5"
)

type ticketsHandler struct {
	ticketService *ticket.TicketService
}

func NewTicketRouter(ticketService *ticket.TicketService) chi.Router {

	ticketHandl := ticketsHandler{
		ticketService: ticketService,
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
	tkt, err := th.ticketService.InitTicket(userID)
	if err != nil {
		util.BuildErrResponse(w, err)
		return
	}

	util.BuildResponseWithBodyAndCode(w, map[string]string{
		"ticket": tkt,
	}, http.StatusCreated)
}
