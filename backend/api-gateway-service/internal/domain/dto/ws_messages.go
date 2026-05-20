package dto

type WSIncomingMessage struct {
	Type string           `json:"type"`
	Data WSSendMessageDTO `json:"data"`
}

type WSSendMessageDTO struct {
	SenderType string `json:"sender_type"`
	SenderID   int64  `json:"sender_id"`
	Body       string `json:"body"`
}

type WSOutgoingMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type WSMessagePayload struct {
	ID         int64  `json:"id"`
	ChatID     int64  `json:"chat_id"`
	SenderType string `json:"sender_type"`
	SenderID   int64  `json:"sender_id"`
	Body       string `json:"body"`
	CreatedAt  string `json:"created_at"`
}
