package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"chat-service/internal/dao/datastruct"
	"chat-service/internal/domain/entities"
	"chat-service/internal/mappers"
)

type ChatAccessor struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewChatAccessor(db *sql.DB) *ChatAccessor {
	return &ChatAccessor{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (a *ChatAccessor) GetChatByTicketID(ctx context.Context, ticketID int64) (*entities.Chat, error) {
	query, args, err := a.qb.
		Select("id", "ticket_id", "created_at", "updated_at").
		From("chats").
		Where(sq.Eq{"ticket_id": ticketID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var chat datastruct.Chat
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&chat.Id,
		&chat.TicketId,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get chat by ticket id: %w", err)
	}

	return mappers.MapChatDataStructToEntity(&chat), nil
}

func (a *ChatAccessor) GetChatByID(ctx context.Context, chatID int64) (*entities.Chat, error) {
	query, args, err := a.qb.
		Select("id", "ticket_id", "created_at", "updated_at").
		From("chats").
		Where(sq.Eq{"id": chatID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var chat datastruct.Chat
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&chat.Id,
		&chat.TicketId,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get chat by id: %w", err)
	}

	return mappers.MapChatDataStructToEntity(&chat), nil
}

func (a *ChatAccessor) CreateChat(ctx context.Context, ticketId int64) (*entities.Chat, error) {
	query, args, err := a.qb.
		Insert("chats").
		Columns("ticket_id").
		Values(ticketId).
		Suffix("RETURNING id, ticket_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var chat datastruct.Chat
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&chat.Id,
		&chat.TicketId,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create chat: %w", err)
	}

	return mappers.MapChatDataStructToEntity(&chat), nil
}

func (a *ChatAccessor) UpdateChatUpdatedAt(ctx context.Context, chatID int64) error {
	query, args, err := a.qb.
		Update("chats").
		Set("updated_at", sq.Expr("now()")).
		Where(sq.Eq{"id": chatID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update chat: %w", err)
	}

	return nil
}
