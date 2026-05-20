package datastruct

type CreateMessage struct {
	ChatId     int64  `db:"chat_id"`
	SenderType string `db:"sender_type"`
	SenderId   int64  `db:"sender_id"`
	Body       string `db:"body"`
}
