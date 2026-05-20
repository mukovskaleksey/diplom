package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/domain/entities"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/mapper"
	"encoding/json"
	"net/http"
	"strconv"
)

type ListSpecialistTicketsHandler struct {
	coreClient *grpc_client.CoreClient
}

func NewSpecialistListTicketsHandler(coreClient *grpc_client.CoreClient) *ListSpecialistTicketsHandler {
	return &ListSpecialistTicketsHandler{coreClient: coreClient}
}

func (h *ListSpecialistTicketsHandler) GetPath() string {
	return paths.ListSpecialistTicketsPath
}

func (h *ListSpecialistTicketsHandler) GetMethod() string {
	return http.MethodGet
}

func (h *ListSpecialistTicketsHandler) GetRequest() any {
	return nil
}

func (h *ListSpecialistTicketsHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	var filter *entities.Filter
	specialistIDStr := r.URL.Query().Get("specialist_id")
	if specialistIDStr != "" {
		specialistID, err := strconv.ParseInt(specialistIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid specialist_id", http.StatusBadRequest)
			return
		}

		filter = &entities.Filter{
			SpecialistId: &specialistID,
		}
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
