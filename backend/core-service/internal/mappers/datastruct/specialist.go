package datastruct

import (
	"core-service/internal/dao/datastruct"
	"core-service/internal/domain/entities"
)

func MapSpecialistDataStructToEntity(specialist *datastruct.Specialist) *entities.Specialist {
	if specialist == nil {
		return nil
	}

	return &entities.Specialist{
		UserId:      specialist.UserId,
		Category:    specialist.Category,
		CurrentLoad: specialist.CurrentLoad,
	}
}
