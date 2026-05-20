package mappers

import (
	"chat-service/internal/dao/datastruct"
	"chat-service/internal/domain/entities"
	chatpb "proto/chat"
	"time"
)

func MapMessageDataStructToEntity(message *datastruct.Message) *entities.Message {
	if message == nil {
		return nil
	}

	return &entities.Message{
		Id:         message.Id,
		ChatId:     message.ChatId,
		SenderType: message.SenderType,
		SenderId:   message.SenderId,
		Body:       message.Body,
		CreatedAt:  message.CreatedAt,
	}
}

func MapMessagesDataStructToEntities(messages []datastruct.Message) []entities.Message {
	result := make([]entities.Message, 0, len(messages))

	for _, message := range messages {
		result = append(result, entities.Message{
			Id:         message.Id,
			ChatId:     message.ChatId,
			SenderType: message.SenderType,
			SenderId:   message.SenderId,
			Body:       message.Body,
			CreatedAt:  message.CreatedAt,
		})
	}

	return result
}

func MapMessageEntityToProto(message *entities.Message) *chatpb.Message {
	if message == nil {
		return nil
	}

	return &chatpb.Message{
		Id:         message.Id,
		ChatId:     message.ChatId,
		SenderType: message.SenderType,
		SenderId:   message.SenderId,
		Body:       message.Body,
		CreatedAt:  formatTime(message.CreatedAt),
	}
}

func MapMessageEntitiesToProto(messages []entities.Message) []*chatpb.Message {
	result := make([]*chatpb.Message, 0, len(messages))

	for i := range messages {
		result = append(result, MapMessageEntityToProto(&messages[i]))
	}

	return result
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
