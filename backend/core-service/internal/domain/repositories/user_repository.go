package repositories

import (
	"context"

	"core-service/internal/dao"
	"core-service/internal/domain/entities"
)

type UserRepository struct {
	userAccessor *dao.UserAccessor
}

func NewUserRepository(userAccessor *dao.UserAccessor) *UserRepository {
	return &UserRepository{
		userAccessor: userAccessor,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
) (*entities.User, error) {
	return r.userAccessor.Create(ctx, name, email, passwordHash)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return r.userAccessor.GetByEmail(ctx, email)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*entities.User, error) {
	return r.userAccessor.GetByID(ctx, id)
}
