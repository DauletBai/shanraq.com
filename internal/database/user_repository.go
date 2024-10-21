package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"shanraq.com/internal/models"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) models.UserRepository {
	return &userRepository {db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (gender, birthday, first_name, last_name, phone, email, password, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query,
		user.Gender, user.Birthday, user.FirstName, user.LastName, user.Phone, user.Email, user.Password, user.Role,
	)
	return err
}

// Implements other methods