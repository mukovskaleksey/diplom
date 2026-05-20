package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/domain/dto"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/ws"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type ChatWSHandler struct {
	chatClient *grpc_client.ChatClient
	hub        *ws.Hub
	upgrader   websocket.Upgrader
}

func NewChatWSHandler(chatClient *grpc_client.ChatClient, hub *ws.Hub) *ChatWSHandler {
	return &ChatWSHandler{
		chatClient: chatClient,
		hub:        hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *ChatWSHandler) GetPath() string {
	return paths.WebSocketChatPath
}

func (h *ChatWSHandler) GetMethod() string {
	return http.MethodGet
}

func (h *ChatWSHandler) GetRequest() any {
	return nil
}

func (h *ChatWSHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat id", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		Conn:   conn,
		ChatID: chatID,
	}

	h.hub.AddClient(client)
	defer func() {
		h.hub.RemoveClient(client)
		_ = conn.Close()
	}()

	_ = conn.SetReadDeadline(time.Time{})

	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var incoming dto.WSIncomingMessage
		if err = json.Unmarshal(rawMessage, &incoming); err != nil {
			continue
		}

		if incoming.Type != "send_message" {
			continue
		}

		resp, err := h.chatClient.SendMessage(
			context.Background(),
			chatID,
			incoming.Data.SenderType,
			incoming.Data.SenderID,
			incoming.Data.Body,
		)
		if err != nil {
			continue
		}

		outgoing := dto.WSOutgoingMessage{
			Type: "message_created",
			Data: dto.WSMessagePayload{
				ID:         resp.Message.Id,
				ChatID:     resp.Message.ChatId,
				SenderType: resp.Message.SenderType,
				SenderID:   resp.Message.SenderId,
				Body:       resp.Message.Body,
				CreatedAt:  resp.Message.CreatedAt,
			},
		}

		payload, err := json.Marshal(outgoing)
		if err != nil {
			continue
		}

		h.hub.Broadcast(chatID, payload)
	}
}
