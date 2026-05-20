package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(router chi.Router, handlers ...Handler) {
	for _, h := range handlers {
		handler := h

		router.MethodFunc(handler.GetMethod(), handler.GetPath(), func(w http.ResponseWriter, r *http.Request) {
			req := handler.GetRequest()

			if req != nil && r.Body != nil && r.Method != http.MethodGet {
				if err := json.NewDecoder(r.Body).Decode(req); err != nil {
					http.Error(w, "invalid request body", http.StatusBadRequest)
					return
				}
			}

			handler.Handle(w, r, req)
		})
	}
}
