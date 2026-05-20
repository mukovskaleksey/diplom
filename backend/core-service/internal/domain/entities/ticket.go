package entities

import "time"

type TicketStatus string

const (
	StatusNew        TicketStatus = "new"
	StatusAssigned   TicketStatus = "assigned"
	StatusInProgress TicketStatus = "in_progress"
	StatusClosed     TicketStatus = "closed"
)

type Ticket struct {
	Id           int64
	UserId       int64
	Message      string
	Category     string
	Status       TicketStatus
	SpecialistId *int64
	CreatedAt    time.Time
}
