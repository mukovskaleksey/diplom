package paths

const (
	CreateTicketPath          = "/api/v1/tickets"
	ListSpecialistTicketsPath = "/api/v1/specialists/tickets"
	ListUserTicketsPath       = "/api/v1/users/{user_id}/tickets"
	GetTicketPath             = "/api/v1/tickets/{id}"
	AssignTicketPath          = "/api/v1/tickets/{id}/assign"
	CloseTicketPath           = "/api/v1/tickets/{id}/close"

	OpenTicketChatPath  = "/api/v1/tickets/{ticket_id}/chat/open"
	GetChatMessagesPath = "/api/v1/chats/{chat_id}/messages"
	SendMessagePath     = "/api/v1/chats/{chat_id}/messages"

	WebSocketChatPath = "/api/v1/ws/chats/{chat_id}"

	RegisterPath = "/api/v1/auth/register"
	LoginPath    = "/api/v1/auth/login"
	GetUserPath  = "/api/v1/users/{id}"
)
