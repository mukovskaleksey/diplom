package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/domain/dto"
	appErrors "api-gateway-service/internal/domain/errors"
	"api-gateway-service/internal/grpc_client"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SendMessageHandler struct {
	chatClient *grpc_client.ChatClient
}

func NewSendMessageHandler(chatClient *grpc_client.ChatClient) *SendMessageHandler {
	return &SendMessageHandler{chatClient: chatClient}
}

func (h *SendMessageHandler) GetPath() string {
	return paths.SendMessagePath
}

func (h *SendMessageHandler) GetMethod() string {
	return http.MethodPost
}

func (h *SendMessageHandler) GetRequest() any {
	return &dto.SendMessageRequest{}
}

func (h *SendMessageHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	request, ok := req.(*dto.SendMessageRequest)
	if !ok {
		http.Error(w, appErrors.ErrInvalidRequestType.Error(), http.StatusBadRequest)
		return
	}

	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat id", http.StatusBadRequest)
		return
	}

	resp, err := h.chatClient.SendMessage(
		r.Context(),
		chatID,
		request.SenderType,
		request.SenderID,
		request.Body,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": map[string]any{
			"id":          resp.Message.Id,
			"chat_id":     resp.Message.ChatId,
			"sender_type": resp.Message.SenderType,
			"sender_id":   resp.Message.SenderId,
			"body":        resp.Message.Body,
			"created_at":  resp.Message.CreatedAt,
		},
	})
}
