package grpc_handlers

import (
	"context"
	"core-service/internal/mappers/domain"

	protoMapper "core-service/internal/mappers/proto"
	pb "proto/core"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateTicket(
	ctx context.Context,
	req *pb.CreateTicketRequest,
) (*pb.CreateTicketResponse, error) {
	ticket, err := s.ticketService.CreateTicket(
		ctx,
		req.GetUserId(),
		req.GetMessage(),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateTicketResponse{
		Ticket: protoMapper.MapTicketEntityToProto(ticket),
	}, nil
}

func (s *Server) GetTicket(
	ctx context.Context,
	req *pb.GetTicketRequest,
) (*pb.GetTicketResponse, error) {
	ticket, err := s.ticketService.GetTicket(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetTicketResponse{
		Ticket: protoMapper.MapTicketEntityToProto(ticket),
	}, nil
}

func (s *Server) ListTickets(
	ctx context.Context,
	req *pb.ListTicketsRequest,
) (*pb.ListTicketsResponse, error) {
	filters := domain.Filter(req)
	tickets, err := s.ticketService.ListTickets(ctx, filters)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Ticket, 0, len(tickets))
	for _, ticket := range tickets {
		result = append(result, protoMapper.MapTicketEntityToProto(ticket))
	}

	return &pb.ListTicketsResponse{
		Tickets: result,
	}, nil
}

func (s *Server) CloseTicket(
	ctx context.Context,
	req *pb.CloseTicketRequest,
) (*pb.CloseTicketResponse, error) {
	ticket, err := s.ticketService.CloseTicket(ctx, req.GetTicketId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CloseTicketResponse{
		Ticket: protoMapper.MapTicketEntityToProto(ticket),
	}, nil
}

func (s *Server) AssignTicket(
	ctx context.Context,
	req *pb.AssignTicketRequest,
) (*pb.AssignTicketResponse, error) {
	ticket, err := s.ticketService.AssignTicket(
		ctx,
		req.GetTicketId(),
		req.GetSpecialistId(),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.AssignTicketResponse{
		Ticket: protoMapper.MapTicketEntityToProto(ticket),
	}, nil
}
