package mappers

import (
	"core-service/internal/dao/datastruct"
	"core-service/internal/domain/entities"
	pb "proto/core"
	"time"
)

func MapUserDataStructToEntity(user *datastruct.User) *entities.User {
	if user == nil {
		return nil
	}

	return &entities.User{
		Id:           user.Id,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}
}

func MapUserEntityToProto(user *entities.User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Id:           user.Id,
		Name:         user.Name,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		IsSpecialist: user.IsSpecialist,
	}
}
