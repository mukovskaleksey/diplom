package repositories

import (
	"context"

	"chat-service/internal/dao"
	"chat-service/internal/dao/datastruct"
	"chat-service/internal/domain/entities"
)

type MessageRepository struct {
	messageAccessor *dao.MessageAccessor
}

func NewMessageRepository(messageAccessor *dao.MessageAccessor) *MessageRepository {
	return &MessageRepository{
		messageAccessor: messageAccessor,
	}
}

func (r *MessageRepository) GetMessagesByChatID(ctx context.Context, chatID int64) ([]entities.Message, error) {
	return r.messageAccessor.GetMessagesByChatID(ctx, chatID)
}

func (r *MessageRepository) CreateMessage(
	ctx context.Context,
	chatID int64,
	senderType string,
	senderID int64,
	body string,
) (*entities.Message, error) {
	return r.messageAccessor.CreateMessage(ctx, datastruct.CreateMessage{
		ChatId:     chatID,
		SenderType: senderType,
		SenderId:   senderID,
		Body:       body,
	})
}
