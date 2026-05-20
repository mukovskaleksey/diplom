package domain

import (
	"core-service/internal/dao/datastruct"
	"core-service/internal/domain/entities"
)

func MapTicketEntityToDataStruct(ticket *entities.Ticket) *datastruct.Ticket {
	if ticket == nil {
		return nil
	}

	return &datastruct.Ticket{
		Id:           ticket.Id,
		UserId:       ticket.UserId,
		Message:      ticket.Message,
		Category:     ticket.Category,
		Status:       string(ticket.Status),
		SpecialistId: ticket.SpecialistId,
		CreatedAt:    ticket.CreatedAt,
	}
}

func MapTicketDataStructToEntity(ticket *datastruct.Ticket) *entities.Ticket {
	if ticket == nil {
		return nil
	}

	return &entities.Ticket{
		Id:           ticket.Id,
		UserId:       ticket.UserId,
		Message:      ticket.Message,
		Category:     ticket.Category,
		Status:       entities.TicketStatus(ticket.Status),
		SpecialistId: ticket.SpecialistId,
		CreatedAt:    ticket.CreatedAt,
	}
}
