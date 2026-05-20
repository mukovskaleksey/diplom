package dto

type OpenChatResponse struct {
	Chat    ChatResponse `json:"chat"`
	Created bool         `json:"created"`
}

type ChatResponse struct {
	ID        int64  `json:"id"`
	TicketID  int64  `json:"ticket_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MessageResponse struct {
	ID         int64  `json:"id"`
	ChatID     int64  `json:"chat_id"`
	SenderType string `json:"sender_type"`
	SenderID   int64  `json:"sender_id"`
	Body       string `json:"body"`
	CreatedAt  string `json:"created_at"`
}

type SendMessageRequest struct {
	SenderType string `json:"sender_type"`
	SenderID   int64  `json:"sender_id"`
	Body       string `json:"body"`
}
