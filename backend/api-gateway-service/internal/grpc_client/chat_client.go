package grpc_client

import (
	"context"
	"fmt"
	"time"

	chatpb "proto/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatClient struct {
	conn   *grpc.ClientConn
	client chatpb.ChatServiceClient
}

func NewChatClient(host string, port int) (*ChatClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &ChatClient{
		conn:   conn,
		client: chatpb.NewChatServiceClient(conn),
	}, nil
}

func (c *ChatClient) Close() error {
	return c.conn.Close()
}

func (c *ChatClient) GetOrCreateChatByTicket(ctx context.Context, ticketID int64) (*chatpb.GetOrCreateChatByTicketResponse, error) {
	return c.client.GetOrCreateChatByTicket(ctx, &chatpb.GetOrCreateChatByTicketRequest{
		TicketId: ticketID,
	})
}

func (c *ChatClient) GetMessages(ctx context.Context, chatID int64) (*chatpb.GetMessagesResponse, error) {
	return c.client.GetMessages(ctx, &chatpb.GetMessagesRequest{
		ChatId: chatID,
	})
}

func (c *ChatClient) SendMessage(ctx context.Context, chatID int64, senderType string, senderID int64, body string) (*chatpb.SendMessageResponse, error) {
	return c.client.SendMessage(ctx, &chatpb.SendMessageRequest{
		ChatId:     chatID,
		SenderType: senderType,
		SenderId:   senderID,
		Body:       body,
	})
}
