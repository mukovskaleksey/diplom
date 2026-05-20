package http_handlers

import (
	"api-gateway-service/internal/controller/paths"
	"api-gateway-service/internal/domain/dto"
	"api-gateway-service/internal/grpc_client"
	"api-gateway-service/internal/mapper"
	"encoding/json"
	"net/http"
)

type LoginHandler struct {
	coreClient *grpc_client.CoreClient
}

func NewLoginHandler(coreClient *grpc_client.CoreClient) *LoginHandler {
	return &LoginHandler{coreClient: coreClient}
}

func (h *LoginHandler) GetPath() string {
	return paths.LoginPath
}

func (h *LoginHandler) GetMethod() string {
	return http.MethodPost
}

func (h *LoginHandler) GetRequest() any {
	return &dto.LoginRequest{}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	request, ok := req.(*dto.LoginRequest)
	if !ok {
		http.Error(w, "invalid request type", http.StatusBadRequest)
		return
	}

	user, err := h.coreClient.LoginUser(
		r.Context(),
		request.Email,
		request.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user": mapper.MapUserProtoToResponse(user),
	})
}
