package post

import "context"

// Repository handles the CRUD operations for Post
type Repository interface {
	GetAll(ctx context.Context) ([]Post, error)
	GetById(ctx context.Context, id uint) (Post, error)
	GetByUser(ctx context.Context, userId uint) ([]Post, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, id uint, post Post) error
	Delete(ctx context.Context, id uint) error
}
