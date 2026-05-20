package domain

import (
	"core-service/internal/domain/entities"
	corepb "proto/core"
)

func Filter(req *corepb.ListTicketsRequest) *entities.Filter {
	if req == nil {
		return nil
	}
	return &entities.Filter{
		UserId:       req.UserId,
		SpecialistId: req.SpecialistId,
	}
}
