package mappers

import (
	"chat-service/internal/dao/datastruct"
	"chat-service/internal/domain/entities"
	chatpb "proto/chat"
)

func MapChatDataStructToEntity(chat *datastruct.Chat) *entities.Chat {
	if chat == nil {
		return nil
	}

	return &entities.Chat{
		Id:        chat.Id,
		TicketId:  chat.TicketId,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}
}

func MapChatEntityToProto(chat *entities.Chat) *chatpb.Chat {
	if chat == nil {
		return nil
	}

	return &chatpb.Chat{
		Id:        chat.Id,
		TicketId:  chat.TicketId,
		CreatedAt: formatTime(chat.CreatedAt),
		UpdatedAt: formatTime(chat.UpdatedAt),
	}
}
