package dao

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"chat-service/internal/dao/datastruct"
	"chat-service/internal/domain/entities"
	"chat-service/internal/mappers"
)

type MessageAccessor struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewMessageAccessor(db *sql.DB) *MessageAccessor {
	return &MessageAccessor{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (a *MessageAccessor) GetMessagesByChatID(ctx context.Context, chatID int64) ([]entities.Message, error) {
	query, args, err := a.qb.
		Select("id", "chat_id", "sender_type", "sender_id", "body", "created_at").
		From("messages").
		Where(sq.Eq{"chat_id": chatID}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get messages by chat id query: %w", err)
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}
	defer rows.Close()

	messages := make([]datastruct.Message, 0)

	for rows.Next() {
		var message datastruct.Message

		err = rows.Scan(
			&message.Id,
			&message.ChatId,
			&message.SenderType,
			&message.SenderId,
			&message.Body,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate messages: %w", err)
	}

	return mappers.MapMessagesDataStructToEntities(messages), nil
}

func (a *MessageAccessor) CreateMessage(ctx context.Context, create datastruct.CreateMessage) (*entities.Message, error) {
	query, args, err := a.qb.
		Insert("messages").
		Columns("chat_id", "sender_type", "sender_id", "body").
		Values(create.ChatId, create.SenderType, create.SenderId, create.Body).
		Suffix("RETURNING id, chat_id, sender_type, sender_id, body, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create message query: %w", err)
	}

	var message datastruct.Message
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&message.Id,
		&message.ChatId,
		&message.SenderType,
		&message.SenderId,
		&message.Body,
		&message.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return mappers.MapMessageDataStructToEntity(&message), nil
}
