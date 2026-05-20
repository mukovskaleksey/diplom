package entities

import "time"

type Chat struct {
	Id        int64
	TicketId  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
