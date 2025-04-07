package user

import (
	"context"

	"github.com/cortzero/go-postgres-blog/internal/service/errors"
)

// user.Service is the interface that a service layer component must fullfil to manage CRUD operations for users.
type Service interface {
	CreateUser(ctx context.Context, user *User) *errors.CustomError
	UpdateUser(ctx context.Context, id uint, user *User) *errors.CustomError
	DeleteUser(ctx context.Context, id uint) *errors.CustomError
	GetAllUsers(ctx context.Context) ([]User, *errors.CustomError)
	GetUserById(ctx context.Context, id uint) (User, *errors.CustomError)
	GetUserByUsername(ctx context.Context, username string) (User, *errors.CustomError)
	GetUserByEmail(ctx context.Context, email string) (User, *errors.CustomError)
}
