package services

import (
	"context"
	"fmt"
	"strings"

	"core-service/internal/domain/entities"
	"core-service/internal/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository       *repositories.UserRepository
	specialistRepository *repositories.SpecialistRepository
}

func NewUserService(
	userRepository *repositories.UserRepository,
	specialistRepository *repositories.SpecialistRepository,
) *UserService {
	return &UserService{
		userRepository:       userRepository,
		specialistRepository: specialistRepository,
	}
}

func (s *UserService) RegisterUser(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*entities.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		return nil, fmt.Errorf("name, email and password are required")
	}

	existingUser, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.Create(ctx, name, email, string(passwordHashBytes))
	if err != nil {
		return nil, err
	}

	user.IsSpecialist = false

	return user, nil
}

func (s *UserService) LoginUser(
	ctx context.Context,
	email string,
	password string,
) (*entities.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	isSpecialist, err := s.specialistRepository.IsSpecialist(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	user.IsSpecialist = isSpecialist

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*entities.User, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	isSpecialist, err := s.specialistRepository.IsSpecialist(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	user.IsSpecialist = isSpecialist

	return user, nil
}
