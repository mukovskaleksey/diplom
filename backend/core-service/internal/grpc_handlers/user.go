package grpc_handlers

import (
	"context"
	"core-service/internal/mappers"

	pb "proto/core"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) RegisterUser(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {
	user, err := s.userService.RegisterUser(
		ctx,
		req.GetName(),
		req.GetEmail(),
		req.GetPassword(),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.RegisterUserResponse{
		User: mappers.MapUserEntityToProto(user),
	}, nil
}

func (s *Server) LoginUser(
	ctx context.Context,
	req *pb.LoginUserRequest,
) (*pb.LoginUserResponse, error) {
	user, err := s.userService.LoginUser(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
	)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.LoginUserResponse{
		User: mappers.MapUserEntityToProto(user),
	}, nil
}

func (s *Server) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
	user, err := s.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetUserResponse{
		User: mappers.MapUserEntityToProto(user),
	}, nil
}
