package services

import (
	"context"
	"fmt"
	"strings"

	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/repositories"
)

type ChatService struct {
	chatRepository    *repositories.ChatRepository
	messageRepository *repositories.MessageRepository
}

func NewChatService(
	chatRepository *repositories.ChatRepository,
	messageRepository *repositories.MessageRepository,
) *ChatService {
	return &ChatService{
		chatRepository:    chatRepository,
		messageRepository: messageRepository,
	}
}

func (s *ChatService) GetOrCreateChatByTicket(ctx context.Context, ticketID int64) (*entities.Chat, bool, error) {
	chat, err := s.chatRepository.GetByTicketID(ctx, ticketID)
	if err != nil {
		return nil, false, fmt.Errorf("get chat by ticket id: %w", err)
	}

	if chat != nil {
		return chat, false, nil
	}

	createdChat, err := s.chatRepository.Create(ctx, ticketID)
	if err != nil {
		return nil, false, fmt.Errorf("create chat: %w", err)
	}

	return createdChat, true, nil
}

func (s *ChatService) GetMessagesByChatID(ctx context.Context, chatID int64) ([]entities.Message, error) {
	chat, err := s.chatRepository.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat by id: %w", err)
	}
	if chat == nil {
		return nil, fmt.Errorf("chat not found")
	}

	messages, err := s.messageRepository.GetMessagesByChatID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get messages by chat id: %w", err)
	}

	return messages, nil
}

func (s *ChatService) SendMessage(
	ctx context.Context,
	chatID int64,
	senderType string,
	senderID int64,
	body string,
) (*entities.Message, error) {
	if senderType != "user" && senderType != "specialist" && senderType != "system" {
		return nil, fmt.Errorf("invalid sender type")
	}

	body = strings.TrimSpace(body)
	if body == "" {
		return nil, fmt.Errorf("message body is empty")
	}

	chat, err := s.chatRepository.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat by id: %w", err)
	}
	if chat == nil {
		return nil, fmt.Errorf("chat not found")
	}

	message, err := s.messageRepository.CreateMessage(ctx, chatID, senderType, senderID, body)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	err = s.chatRepository.UpdateUpdatedAt(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("update chat updated_at: %w", err)
	}

	return message, nil
}
