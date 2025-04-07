package user

import "context"

// Repository handles the CRUD operations for User
type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id uint) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id uint, user User) error
	Delete(ctx context.Context, id uint) error
}
