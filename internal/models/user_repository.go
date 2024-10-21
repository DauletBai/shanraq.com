package models

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserEmail(ctx context.Context, email string) (*User, error)
	// Other metods
}