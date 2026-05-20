package dao

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"core-service/internal/dao/datastruct"
	"core-service/internal/domain/entities"
)

type TicketAccessor struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewTicketAccessor(db *sql.DB) *TicketAccessor {
	return &TicketAccessor{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (a *TicketAccessor) Create(
	ctx context.Context,
	ticket *datastruct.Ticket,
) error {
	query, args, err := a.qb.
		Insert("tickets").
		Columns(
			"user_id",
			"message",
			"category",
			"status",
			"specialist_id",
		).
		Values(
			ticket.UserId,
			ticket.Message,
			ticket.Category,
			ticket.Status,
			ticket.SpecialistId,
		).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return err
	}

	return a.db.QueryRowContext(ctx, query, args...).Scan(
		&ticket.Id,
		&ticket.CreatedAt,
	)
}

func (a *TicketAccessor) GetById(
	ctx context.Context,
	id int64,
) (*datastruct.Ticket, error) {
	query, args, err := a.qb.
		Select(
			"id",
			"user_id",
			"message",
			"category",
			"status",
			"specialist_id",
			"created_at",
		).
		From("tickets").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var item datastruct.Ticket

	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&item.Id,
		&item.UserId,
		&item.Message,
		&item.Category,
		&item.Status,
		&item.SpecialistId,
		&item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (a *TicketAccessor) List(
	ctx context.Context,
	filters *entities.Filter,
) ([]*datastruct.Ticket, error) {
	qb := a.qb.
		Select(
			"id",
			"user_id",
			"message",
			"category",
			"status",
			"specialist_id",
			"created_at",
		).
		From("tickets")

	if filters != nil {
		if filters.UserId != nil {
			qb = qb.Where(sq.Eq{"user_id": *filters.UserId})
		}

		if filters.SpecialistId != nil {
			qb = qb.Where(sq.Eq{"specialist_id": *filters.SpecialistId})
		}
	}

	qb = qb.OrderBy(
		"CASE WHEN status = 'CLOSED' THEN 1 ELSE 0 END",
		"created_at DESC",
	)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*datastruct.Ticket, 0)

	for rows.Next() {
		var item datastruct.Ticket

		if err := rows.Scan(
			&item.Id,
			&item.UserId,
			&item.Message,
			&item.Category,
			&item.Status,
			&item.SpecialistId,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *TicketAccessor) Update(
	ctx context.Context,
	ticket *datastruct.Ticket,
) error {
	query, args, err := a.qb.
		Update("tickets").
		Set("category", ticket.Category).
		Set("status", ticket.Status).
		Set("specialist_id", ticket.SpecialistId).
		Where(sq.Eq{"id": ticket.Id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, query, args...)
	return err
}
