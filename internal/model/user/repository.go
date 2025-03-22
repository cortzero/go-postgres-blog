package user

import "context"

// Repository handles the CRUD operations for User
type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context) (User, error)
	GetByUsername(ctx context.Context) (User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id int, user User) error
	Delete(ctx context.Context, id int) error
}
