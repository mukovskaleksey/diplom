package grpc_handlers

import (
	"core-service/internal/services"
	pb "proto/core"
)

type Server struct {
	pb.UnimplementedCoreServiceServer
	ticketService *services.TicketService
	userService   *services.UserService
}

func NewServer(
	ticketService *services.TicketService,
	userService *services.UserService,
) *Server {
	return &Server{
		ticketService: ticketService,
		userService:   userService,
	}
}
