package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/mapper"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GetTicketHandler struct {
	coreClient *grpc_client.CoreClient
}

func NewGetTicketHandler(coreClient *grpc_client.CoreClient) *GetTicketHandler {
	return &GetTicketHandler{coreClient: coreClient}
}

func (h *GetTicketHandler) GetPath() string {
	return paths.GetTicketPath
}

func (h *GetTicketHandler) GetMethod() string {
	return http.MethodGet
}

func (h *GetTicketHandler) GetRequest() any {
	return nil
}

func (h *GetTicketHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid ticket id", http.StatusBadRequest)
		return
	}

	ticket, err := h.coreClient.GetTicket(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ticket": mapper.MapTicketProtoToResponse(ticket),
	})
}
