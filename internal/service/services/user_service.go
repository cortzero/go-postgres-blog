package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cortzero/go-postgres-blog/internal/model/user"
	"github.com/cortzero/go-postgres-blog/internal/service/errors"
)

// UserService is a service layer component that manages the CRUD operations for users
type UserService struct {
	Repository user.Repository
}

func NewUserService(repository user.Repository) *UserService {
	return &UserService{
		Repository: repository,
	}
}

func (service *UserService) CreateUser(ctx context.Context, user *user.User) *errors.CustomError {
	// Sets the timestamp at which the user is created
	user.CreatedAt = time.Now()

	// Hashes the password
	if err := user.HashPassword(); err != nil {
		return errors.NewCustomError(
			"ERROR_HASHING_PASSWORD",
			"An error occurred while hashing the user password.",
			err.Error(),
			time.Now(),
		)
	}

	// Creating the user
	err := service.Repository.Create(ctx, user)
	if err != nil {
		return errors.NewCustomError(
			"ERROR_CREATING_USER",
			"An error occurred while creating the user.",
			err.Error(),
			time.Now(),
		)
	}
	return nil
}

func (service *UserService) UpdateUser(ctx context.Context, id uint, user *user.User) *errors.CustomError {
	// Check if user exists
	existingUser, err := service.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	// Check that all fields are not empty
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		return errors.NewCustomError(
			"EMPTY_FIELDS",
			"You cannot leave user fields empty.",
			"Fill the corresponding fields to update the user.",
			time.Now(),
		)
	}

	// Check if there is another user with the new email
	_, err = service.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return errors.NewCustomError(
			"REPEATED_EMAIL",
			"There can't be two users with the same email.",
			"You need to choose another email for this user.",
			time.Now(),
		)
	}

	// Updating the existing user
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Email = user.Email
	existingUser.UpdatedAt = time.Now()
	err_update := service.Repository.Update(ctx, id, existingUser)
	if err_update != nil {
		return errors.NewCustomError(
			"ERROR_UPDATING_USER",
			"An error occurred while updating the user.",
			err_update.Error(),
			time.Now(),
		)
	}
	return nil
}

func (service *UserService) DeleteUser(ctx context.Context, id uint) *errors.CustomError {
	// Check if the user exists with the given id
	_, err := service.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	// Deleting the user
	err_delete := service.Repository.Delete(ctx, id)
	if err_delete != nil {
		return errors.NewCustomError(
			"ERROR_DELETING",
			"An error occurred while removing the user.",
			err_delete.Error(),
			time.Now(),
		)
	}
	return nil
}

func (service *UserService) GetAllUsers(ctx context.Context) ([]user.User, *errors.CustomError) {
	users, err := service.Repository.GetAll(ctx)
	if err != nil {
		return nil, errors.NewCustomError(
			"RESOURCE_NOT_FOUND",
			"An error occurred while looking for all users.",
			err.Error(),
			time.Now(),
		)
	}
	return users, nil
}

func (service *UserService) GetUserById(ctx context.Context, id uint) (user.User, *errors.CustomError) {
	u, err := service.Repository.GetById(ctx, id)
	if err != nil {
		return user.User{}, errors.NewCustomError(
			"RESOURCE_NOT_FOUND",
			fmt.Sprintf("There is not a user with id '%d'.", id),
			err.Error(),
			time.Now(),
		)
	}
	return u, nil
}

func (service *UserService) GetUserByUsername(ctx context.Context, username string) (user.User, *errors.CustomError) {
	u, err := service.Repository.GetByUsername(ctx, username)
	if err != nil {
		return user.User{}, errors.NewCustomError(
			"RESOURCE_NOT_FOUND",
			fmt.Sprintf("There is not a user with username '%s'.", username),
			err.Error(),
			time.Now(),
		)
	}
	return u, nil
}

func (service *UserService) GetUserByEmail(ctx context.Context, email string) (user.User, *errors.CustomError) {
	u, err := service.Repository.GetByEmail(ctx, email)
	if err != nil {
		return user.User{}, errors.NewCustomError(
			"RESOURCE_NOT_FOUND",
			fmt.Sprintf("There is not a user with email '%s'.", email),
			err.Error(),
			time.Now(),
		)
	}
	return u, nil
}
