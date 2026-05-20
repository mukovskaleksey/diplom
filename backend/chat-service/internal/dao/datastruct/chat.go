package datastruct

import "time"

type Chat struct {
	Id        int64     `db:"id"`
	TicketId  int64     `db:"ticket_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
