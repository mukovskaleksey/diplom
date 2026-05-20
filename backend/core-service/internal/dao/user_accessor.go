package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"core-service/internal/dao/datastruct"
	"core-service/internal/domain/entities"
	"core-service/internal/mappers"
)

type UserAccessor struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewUserAccessor(db *sql.DB) *UserAccessor {
	return &UserAccessor{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (a *UserAccessor) Create(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
) (*entities.User, error) {
	query, args, err := a.qb.
		Insert("users").
		Columns("name", "email", "password_hash").
		Values(name, email, passwordHash).
		Suffix("RETURNING id, name, email, password_hash, created_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var user datastruct.User
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return mappers.MapUserDataStructToEntity(&user), nil
}

func (a *UserAccessor) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query, args, err := a.qb.
		Select("id", "name", "email", "password_hash", "created_at").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user datastruct.User
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return mappers.MapUserDataStructToEntity(&user), nil
}

func (a *UserAccessor) GetByID(ctx context.Context, id int64) (*entities.User, error) {
	query, args, err := a.qb.
		Select("id", "name", "email", "password_hash", "created_at").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user datastruct.User
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return mappers.MapUserDataStructToEntity(&user), nil
}
