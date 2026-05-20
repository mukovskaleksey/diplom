package dto

type CreateTicketRequest struct {
	UserId  int64  `json:"user_id"`
	Message string `json:"message"`
}

type TicketResponse struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	Message      string `json:"message"`
	Category     string `json:"category"`
	Status       string `json:"status"`
	SpecialistId int64  `json:"specialist_id,omitempty"`
	CreatedAt    string `json:"created_at"`
}

type AssignTicketRequest struct {
	SpecialistID int64 `json:"specialist_id"`
	TicketID     int64 `json:"ticket_id"`
}
