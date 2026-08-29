package ticket

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type TicketService struct {
	tickets map[string]string
	mut     sync.RWMutex
}

func NewTicketService() *TicketService {
	return &TicketService{
		tickets: make(map[string]string),
	}
}

func (ts *TicketService) InitTicket(userID string) (string, error) {
	ts.mut.Lock()
	defer ts.mut.Unlock()
	newTicket, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Failed to initialise a ticket")
		return "", err
	}
	ts.tickets[newTicket.String()] = userID
	return newTicket.String(), nil
}

func (ts *TicketService) GetTicket(ticket string) (string, error) {
	ts.mut.RLock()
	defer ts.mut.RUnlock()
	val, ok := ts.tickets[ticket]

	if !ok {
		fmt.Println("Failed to initialise a ticket")
		return "", errors.New("Invalid ticket")
	}
	return val, nil
}

func (ts *TicketService) DeleteTicket(ticket string) error {
	ts.mut.RLock()
	defer ts.mut.RUnlock()
	_, ok := ts.tickets[ticket]

	if !ok {
		fmt.Println("Failed to initialise a ticket")
		return errors.New("Invalid ticket")
	}
	delete(ts.tickets, ticket)
	return nil
}
