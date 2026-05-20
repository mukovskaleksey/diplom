package grpc_client

import (
	"api-gateway-service/internal/domain/entities"
	"context"
	corepb "proto/core"
)

func (c *CoreClient) CreateTicket(
	ctx context.Context,
	userID int64,
	message string,
) (*corepb.Ticket, error) {
	resp, err := c.client.CreateTicket(ctx, &corepb.CreateTicketRequest{
		UserId:  userID,
		Message: message,
	})
	if err != nil {
		return nil, err
	}

	return resp.Ticket, nil
}

func (c *CoreClient) GetTicket(
	ctx context.Context,
	id int64,
) (*corepb.Ticket, error) {
	resp, err := c.client.GetTicket(ctx, &corepb.GetTicketRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Ticket, nil
}

func (c *CoreClient) ListTickets(
	ctx context.Context,
	filter *entities.Filter,
) ([]*corepb.Ticket, error) {
	req := &corepb.ListTicketsRequest{}

	if filter != nil {
		if filter.UserId != nil {
			req.UserId = filter.UserId
		}
		if filter.SpecialistId != nil {
			req.SpecialistId = filter.SpecialistId
		}
	}

	resp, err := c.client.ListTickets(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Tickets, nil
}

func (c *CoreClient) AssignTicket(
	ctx context.Context,
	ticketID int64,
	specialistID int64,
) (*corepb.Ticket, error) {
	resp, err := c.client.AssignTicket(ctx, &corepb.AssignTicketRequest{
		TicketId:     ticketID,
		SpecialistId: specialistID,
	})
	if err != nil {
		return nil, err
	}

	return resp.Ticket, nil
}

func (c *CoreClient) CloseTicket(ctx context.Context, ticketID int64) (*corepb.Ticket, error) {
	resp, err := c.client.CloseTicket(ctx, &corepb.CloseTicketRequest{
		TicketId: ticketID,
	})
	if err != nil {
		return nil, err
	}

	return resp.Ticket, nil
}
