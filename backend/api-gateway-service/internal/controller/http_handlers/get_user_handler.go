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

type GetUserHandler struct {
	coreClient *grpc_client.CoreClient
}

func NewGetUserHandler(coreClient *grpc_client.CoreClient) *GetUserHandler {
	return &GetUserHandler{coreClient: coreClient}
}

func (h *GetUserHandler) GetPath() string {
	return paths.GetUserPath
}

func (h *GetUserHandler) GetMethod() string {
	return http.MethodGet
}

func (h *GetUserHandler) GetRequest() any {
	return nil
}

func (h *GetUserHandler) Handle(w http.ResponseWriter, r *http.Request, req any) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.coreClient.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user": mapper.MapUserProtoToResponse(user),
	})
}
