package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/domain/entities"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/mapper"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ListUserTicketsHandler struct {
	coreClient *grpc_client.CoreClient
}

func NewListUserTicketsHandler(coreClient *grpc_client.CoreClient) *ListUserTicketsHandler {
	return &ListUserTicketsHandler{coreClient: coreClient}
}

func (h *ListUserTicketsHandler) GetPath() string {
	return paths.ListUserTicketsPath
}

func (h *ListUserTicketsHandler) GetMethod() string {
	return http.MethodGet
}

func (h *ListUserTicketsHandler) GetRequest() any {
	return nil
}

func (h *ListUserTicketsHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	filter := &entities.Filter{
		UserId: &userID,
	}

	tickets, err := h.coreClient.ListTickets(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]any, 0, len(tickets))
	for _, t := range tickets {
		resp = append(resp, mapper.MapTicketProtoToResponse(t))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"tickets": resp,
	})
}
