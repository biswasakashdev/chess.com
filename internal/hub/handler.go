package hub

import (
	"errors"
	"log"
	"net/http"

	"github.com/biswasakashdev/chess.com/internal/util"
)

func (h *Hub) WsHandle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	ticket := queryParams.Get("ticket")
	if ticket == "" {
		util.BuildErrResponse(w, errors.New("Invalid Token"))
		return
	}

	userDet, err := h.ticketService.GetTicket(ticket)
	if err != nil {
		util.BuildErrResponse(w, errors.New("Invalid Token"))
		return
	}

	h.ticketService.DeleteTicket(ticket)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	client := &Client{
		Hub:      h,
		Conn:     conn,
		UserID:   userDet.UserId,
		Username: userDet.Username,
		Send:     make(chan []byte, 256),
	}

	h.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
