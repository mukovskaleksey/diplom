package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	appErrors "api-gateway-service/internal/domain/errors"
	"api-gateway-service/internal/grpc_client"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GetChatMessagesHandler struct {
	chatClient *grpc_client.ChatClient
}

func NewGetChatMessagesHandler(chatClient *grpc_client.ChatClient) *GetChatMessagesHandler {
	return &GetChatMessagesHandler{chatClient: chatClient}
}

func (h *GetChatMessagesHandler) GetPath() string {
	return paths.GetChatMessagesPath
}

func (h *GetChatMessagesHandler) GetMethod() string {
	return http.MethodGet
}

func (h *GetChatMessagesHandler) GetRequest() any {
	return nil
}

func (h *GetChatMessagesHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, appErrors.ErrInvalidRequestType.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.chatClient.GetMessages(r.Context(), chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages := make([]map[string]any, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		messages = append(messages, map[string]any{
			"id":          msg.Id,
			"chat_id":     msg.ChatId,
			"sender_type": msg.SenderType,
			"sender_id":   msg.SenderId,
			"body":        msg.Body,
			"created_at":  msg.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"messages": messages,
	})
}
