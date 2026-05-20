package repositories

import (
	"context"

	"core-service/internal/dao"
	"core-service/internal/domain/entities"
	"core-service/internal/mappers/datastruct"
)

type SpecialistRepository struct {
	specialistAccessor *dao.SpecialistAccessor
}

func NewSpecialistRepository(specialistAccessor *dao.SpecialistAccessor) *SpecialistRepository {
	return &SpecialistRepository{
		specialistAccessor: specialistAccessor,
	}
}

func (r *SpecialistRepository) FindByCategory(
	ctx context.Context,
	category string,
) ([]*entities.Specialist, error) {
	list, err := r.specialistAccessor.FindByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Specialist, 0, len(list))
	for _, item := range list {
		result = append(result, datastruct.MapSpecialistDataStructToEntity(item))
	}

	return result, nil
}

func (r *SpecialistRepository) IncrementLoad(
	ctx context.Context,
	specialistId int64,
	category string,
) error {
	return r.specialistAccessor.IncrementLoad(ctx, specialistId, category)
}

func (r *SpecialistRepository) DecrementLoad(
	ctx context.Context,
	specialistId int64,
	category string,
) error {
	return r.specialistAccessor.DecrementLoad(ctx, specialistId, category)
}

func (r *SpecialistRepository) GetById(
	ctx context.Context,
	id int64,
	category string,
) (*entities.Specialist, error) {
	return r.specialistAccessor.GetById(ctx, id, category)
}

func (r *SpecialistRepository) ListByUserID(
	ctx context.Context,
	userID int64,
) ([]*entities.Specialist, error) {
	return r.specialistAccessor.ListByUserID(ctx, userID)
}

func (r *SpecialistRepository) IsSpecialist(
	ctx context.Context,
	userID int64,
) (bool, error) {
	return r.specialistAccessor.IsSpecialist(ctx, userID)
}
