package controller

import (
	"net/http"
)

type Handler interface {
	GetPath() string
	GetMethod() string
	GetRequest() any
	Handle(w http.ResponseWriter, r *http.Request, req any)
}
