package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/domain/dto"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/mapper"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AssignTicketHandler struct {
	coreClient *grpc_client.CoreClient
}

func NewAssignTicketHandler(coreClient *grpc_client.CoreClient) *AssignTicketHandler {
	return &AssignTicketHandler{coreClient: coreClient}
}

func (h *AssignTicketHandler) GetPath() string {
	return paths.AssignTicketPath
}

func (h *AssignTicketHandler) GetMethod() string {
	return http.MethodPost
}

func (h *AssignTicketHandler) GetRequest() any {
	return &dto.AssignTicketRequest{}
}

func (h *AssignTicketHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	request, ok := req.(*dto.AssignTicketRequest)
	if !ok {
		http.Error(w, "invalid request type", http.StatusBadRequest)
		return
	}

	ticketIDStr := chi.URLParam(r, "id")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid ticket id", http.StatusBadRequest)
		return
	}

	if request.SpecialistID <= 0 {
		http.Error(w, "invalid specialist id", http.StatusBadRequest)
		return
	}

	ticket, err := h.coreClient.AssignTicket(r.Context(), ticketID, request.SpecialistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ticket": mapper.MapTicketProtoToResponse(ticket),
	})
}
