package entities

import "time"

type Message struct {
	Id         int64
	ChatId     int64
	SenderType string
	SenderId   int64
	Body       string
	CreatedAt  time.Time
}
