package proto

import (
	"time"

	"core-service/internal/domain/entities"
	pb "proto/core"
)

func MapTicketEntityToProto(ticket *entities.Ticket) *pb.Ticket {
	if ticket == nil {
		return nil
	}

	var specialistId int64
	if ticket.SpecialistId != nil {
		specialistId = *ticket.SpecialistId
	}

	return &pb.Ticket{
		Id:           ticket.Id,
		UserId:       ticket.UserId,
		Message:      ticket.Message,
		Category:     ticket.Category,
		Status:       string(ticket.Status),
		SpecialistId: specialistId,
		CreatedAt:    ticket.CreatedAt.Format(time.RFC3339),
	}
}
