package datastruct

import "time"

type Message struct {
	Id         int64     `db:"id"`
	ChatId     int64     `db:"chat_id"`
	SenderType string    `db:"sender_type"`
	SenderId   int64     `db:"sender_id"`
	Body       string    `db:"body"`
	CreatedAt  time.Time `db:"created_at"`
}
