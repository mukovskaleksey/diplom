package dao

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"core-service/internal/dao/datastruct"
	"core-service/internal/domain/entities"
	mapper "core-service/internal/mappers/datastruct"
)

type SpecialistAccessor struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewSpecialistAccessor(db *sql.DB) *SpecialistAccessor {
	return &SpecialistAccessor{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (a *SpecialistAccessor) FindByCategory(
	ctx context.Context,
	category string,
) ([]*datastruct.Specialist, error) {
	query, args, err := a.qb.
		Select(
			"user_id",
			"category",
			"current_load",
		).
		From("specialists").
		Where(sq.Eq{"category": category}).
		OrderBy("current_load ASC", "user_id ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*datastruct.Specialist, 0)
	for rows.Next() {
		var item datastruct.Specialist

		if err := rows.Scan(
			&item.UserId,
			&item.Category,
			&item.CurrentLoad,
		); err != nil {
			return nil, err
		}

		result = append(result, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *SpecialistAccessor) IncrementLoad(
	ctx context.Context,
	specialistId int64,
	category string,
) error {
	query, args, err := a.qb.
		Update("specialists").
		Set("current_load", sq.Expr("current_load + 1")).
		Where(sq.Eq{
			"user_id":  specialistId,
			"category": category,
		}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, query, args...)
	return err
}

func (a *SpecialistAccessor) DecrementLoad(
	ctx context.Context,
	specialistId int64,
	category string,
) error {
	query, args, err := a.qb.
		Update("specialists").
		Set("current_load", sq.Expr("GREATEST(current_load - 1, 0)")).
		Where(sq.Eq{
			"user_id":  specialistId,
			"category": category,
		}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, query, args...)
	return err
}

func (a *SpecialistAccessor) GetById(
	ctx context.Context,
	id int64,
	category string,
) (*entities.Specialist, error) {
	query, args, err := a.qb.
		Select("user_id", "category", "current_load").
		From("specialists").
		Where(sq.Eq{
			"user_id":  id,
			"category": category,
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var specialist datastruct.Specialist
	err = a.db.QueryRowContext(ctx, query, args...).Scan(
		&specialist.UserId,
		&specialist.Category,
		&specialist.CurrentLoad,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.MapSpecialistDataStructToEntity(&specialist), nil
}

func (a *SpecialistAccessor) ListByUserID(
	ctx context.Context,
	userID int64,
) ([]*entities.Specialist, error) {
	query, args, err := a.qb.
		Select("user_id", "category", "current_load").
		From("specialists").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("category ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*entities.Specialist, 0)
	for rows.Next() {
		var item datastruct.Specialist

		if err := rows.Scan(
			&item.UserId,
			&item.Category,
			&item.CurrentLoad,
		); err != nil {
			return nil, err
		}

		result = append(result, mapper.MapSpecialistDataStructToEntity(&item))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *SpecialistAccessor) IsSpecialist(
	ctx context.Context,
	userID int64,
) (bool, error) {
	query, args, err := a.qb.
		Select("1").
		From("specialists").
		Where(sq.Eq{"user_id": userID}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	var one int
	err = a.db.QueryRowContext(ctx, query, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
