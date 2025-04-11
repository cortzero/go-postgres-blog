package post

import (
	"context"

	"github.com/cortzero/go-postgres-blog/internal/service/errors"
)

type Service interface {
	CreatePost(ctx context.Context, post *Post) *errors.CustomError
	UpdatePost(ctx context.Context, id uint, post *Post) *errors.CustomError
	DeletePost(ctx context.Context, id uint) *errors.CustomError
	GetAllPosts(ctx context.Context) ([]Post, *errors.CustomError)
	GetPostById(ctx context.Context, id uint) (Post, *errors.CustomError)
	GetPostsByUserId(ctx context.Context, userId uint) ([]Post, *errors.CustomError)
}
