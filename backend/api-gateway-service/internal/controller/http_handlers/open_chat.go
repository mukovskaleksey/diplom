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

type OpenTicketChatHandler struct {
	chatClient *grpc_client.ChatClient
}

func NewOpenTicketChatHandler(chatClient *grpc_client.ChatClient) *OpenTicketChatHandler {
	return &OpenTicketChatHandler{chatClient: chatClient}
}

func (h *OpenTicketChatHandler) GetPath() string {
	return paths.OpenTicketChatPath
}

func (h *OpenTicketChatHandler) GetMethod() string {
	return http.MethodPost
}

func (h *OpenTicketChatHandler) GetRequest() any {
	return nil
}

func (h *OpenTicketChatHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	ticketIDStr := chi.URLParam(r, "ticket_id")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		http.Error(w, appErrors.ErrInvalidRequestType.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.chatClient.GetOrCreateChatByTicket(r.Context(), ticketID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"chat": map[string]any{
			"id":         resp.Chat.Id,
			"ticket_id":  resp.Chat.TicketId,
			"created_at": resp.Chat.CreatedAt,
			"updated_at": resp.Chat.UpdatedAt,
		},
		"created": resp.Created,
	})
}
