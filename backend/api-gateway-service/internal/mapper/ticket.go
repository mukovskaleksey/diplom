package mapper

import (
	"api-gateway-service/internal/domain/dto"
	corepb "proto/core"
)

func MapTicketProtoToResponse(ticket *corepb.Ticket) *dto.TicketResponse {
	if ticket == nil {
		return nil
	}

	return &dto.TicketResponse{
		Id:           ticket.Id,
		UserId:       ticket.UserId,
		Message:      ticket.Message,
		Category:     ticket.Category,
		Status:       ticket.Status,
		SpecialistId: ticket.SpecialistId,
		CreatedAt:    ticket.CreatedAt,
	}
}
