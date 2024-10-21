package services

import (
	"context"

	"shanraq.com/internal/models"
)

type UserService interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (*models.User, error)
	// Other metod
}

type UserService struct {
	repo models.UserRepository
}

func (s *UserService) Register(ctx context.Context, user *models.User) error {
	// bissines logic
	return s.repo.CreateUser(ctx, user)
}

// Implement other metods