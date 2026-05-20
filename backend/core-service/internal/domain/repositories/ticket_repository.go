package repositories

import (
	"context"
	"core-service/internal/mappers/datastruct"

	"core-service/internal/dao"
	"core-service/internal/domain/entities"
)

type TicketRepository struct {
	ticketAccessor *dao.TicketAccessor
}

func NewTicketRepository(ticketAccessor *dao.TicketAccessor) *TicketRepository {
	return &TicketRepository{
		ticketAccessor: ticketAccessor,
	}
}

func (r *TicketRepository) Create(
	ctx context.Context,
	ticket *entities.Ticket,
) error {
	ds := datastruct.MapTicketEntityToDataStruct(ticket)

	if err := r.ticketAccessor.Create(ctx, ds); err != nil {
		return err
	}

	ticket.Id = ds.Id
	ticket.CreatedAt = ds.CreatedAt

	return nil
}

func (r *TicketRepository) GetById(
	ctx context.Context,
	id int64,
) (*entities.Ticket, error) {
	ds, err := r.ticketAccessor.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return datastruct.MapTicketDataStructToEntity(ds), nil
}

func (r *TicketRepository) List(
	ctx context.Context,
	filters *entities.Filter,
) ([]*entities.Ticket, error) {
	list, err := r.ticketAccessor.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Ticket, 0, len(list))
	for _, item := range list {
		result = append(result, datastruct.MapTicketDataStructToEntity(item))
	}

	return result, nil
}

func (r *TicketRepository) Update(
	ctx context.Context,
	ticket *entities.Ticket,
) error {
	ds := datastruct.MapTicketEntityToDataStruct(ticket)
	return r.ticketAccessor.Update(ctx, ds)
}
