package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/domain/dto"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/mapper"
	"encoding/json"
	"net/http"
)

type RegisterHandler struct {
	coreClient *grpc_client.CoreClient
}

func NewRegisterHandler(coreClient *grpc_client.CoreClient) *RegisterHandler {
	return &RegisterHandler{coreClient: coreClient}
}

func (h *RegisterHandler) GetPath() string {
	return paths.RegisterPath
}

func (h *RegisterHandler) GetMethod() string {
	return http.MethodPost
}

func (h *RegisterHandler) GetRequest() any {
	return &dto.RegisterRequest{}
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	request, ok := req.(*dto.RegisterRequest)
	if !ok {
		http.Error(w, "invalid request type", http.StatusBadRequest)
		return
	}

	user, err := h.coreClient.RegisterUser(
		r.Context(),
		request.Name,
		request.Email,
		request.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user": mapper.MapUserProtoToResponse(user),
	})
}
