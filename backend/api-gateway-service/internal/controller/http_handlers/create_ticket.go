package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/domain/dto"
	appErrors "api-gateway-service/internal/domain/errors"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/mapper"
	"encoding/json"
	"net/http"
)

type CreateTicketHandler struct {
	coreClient *grpc_client.CoreClient
}

func NewCreateTicketHandler(coreClient *grpc_client.CoreClient) *CreateTicketHandler {
	return &CreateTicketHandler{coreClient: coreClient}
}

func (h *CreateTicketHandler) GetPath() string {
	return paths.CreateTicketPath
}

func (h *CreateTicketHandler) GetMethod() string {
	return http.MethodPost
}

func (h *CreateTicketHandler) GetRequest() any {
	return &dto.CreateTicketRequest{}
}

func (h *CreateTicketHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	request, ok := req.(*dto.CreateTicketRequest)
	if !ok {
		http.Error(w, appErrors.ErrInvalidRequestType.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := h.coreClient.CreateTicket(
		r.Context(),
		request.UserId,
		request.Message,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ticket": mapper.MapTicketProtoToResponse(ticket),
	})
}
