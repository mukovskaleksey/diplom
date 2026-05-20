package datastruct

import "time"

type Ticket struct {
	Id           int64     `db:"id"`
	UserId       int64     `db:"user_id"`
	Message      string    `db:"message"`
	Category     string    `db:"category"`
	Status       string    `db:"status"`
	SpecialistId *int64    `db:"specialist_id"`
	CreatedAt    time.Time `db:"created_at"`
}
