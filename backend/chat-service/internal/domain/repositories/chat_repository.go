package repositories

import (
	"context"

	"chat-service/internal/dao"
	"chat-service/internal/domain/entities"
)

type ChatRepository struct {
	chatAccessor *dao.ChatAccessor
}

func NewChatRepository(chatAccessor *dao.ChatAccessor) *ChatRepository {
	return &ChatRepository{
		chatAccessor: chatAccessor,
	}
}

func (r *ChatRepository) GetByTicketID(ctx context.Context, ticketID int64) (*entities.Chat, error) {
	return r.chatAccessor.GetChatByTicketID(ctx, ticketID)
}

func (r *ChatRepository) GetByID(ctx context.Context, chatID int64) (*entities.Chat, error) {
	return r.chatAccessor.GetChatByID(ctx, chatID)
}

func (r *ChatRepository) Create(ctx context.Context, ticketID int64) (*entities.Chat, error) {
	return r.chatAccessor.CreateChat(ctx, ticketID)
}

func (r *ChatRepository) UpdateUpdatedAt(ctx context.Context, chatID int64) error {
	return r.chatAccessor.UpdateChatUpdatedAt(ctx, chatID)
}
