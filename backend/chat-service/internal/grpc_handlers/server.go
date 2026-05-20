package grpc_handlers

import (
	"chat-service/internal/mappers"
	"context"

	chatpb "proto/chat"

	"chat-service/internal/services"
)

type ChatServer struct {
	chatpb.UnimplementedChatServiceServer
	chatService *services.ChatService
}

func NewChatServer(chatService *services.ChatService) *ChatServer {
	return &ChatServer{
		chatService: chatService,
	}
}

func (s *ChatServer) GetOrCreateChatByTicket(
	ctx context.Context,
	req *chatpb.GetOrCreateChatByTicketRequest,
) (*chatpb.GetOrCreateChatByTicketResponse, error) {
	chat, created, err := s.chatService.GetOrCreateChatByTicket(ctx, req.TicketId)
	if err != nil {
		return nil, err
	}

	return &chatpb.GetOrCreateChatByTicketResponse{
		Chat:    mappers.MapChatEntityToProto(chat),
		Created: created,
	}, nil
}

func (s *ChatServer) GetChatByTicket(
	ctx context.Context,
	req *chatpb.GetChatByTicketRequest,
) (*chatpb.GetChatByTicketResponse, error) {
	chat, _, err := s.chatService.GetOrCreateChatByTicket(ctx, req.TicketId)
	if err != nil {
		return nil, err
	}

	return &chatpb.GetChatByTicketResponse{
		Chat: mappers.MapChatEntityToProto(chat),
	}, nil
}

func (s *ChatServer) GetMessages(
	ctx context.Context,
	req *chatpb.GetMessagesRequest,
) (*chatpb.GetMessagesResponse, error) {
	messages, err := s.chatService.GetMessagesByChatID(ctx, req.ChatId)
	if err != nil {
		return nil, err
	}

	return &chatpb.GetMessagesResponse{
		Messages: mappers.MapMessageEntitiesToProto(messages),
	}, nil
}

func (s *ChatServer) SendMessage(
	ctx context.Context,
	req *chatpb.SendMessageRequest,
) (*chatpb.SendMessageResponse, error) {
	message, err := s.chatService.SendMessage(
		ctx,
		req.ChatId,
		req.SenderType,
		req.SenderId,
		req.Body,
	)
	if err != nil {
		return nil, err
	}

	return &chatpb.SendMessageResponse{
		Message: mappers.MapMessageEntityToProto(message),
	}, nil
}
