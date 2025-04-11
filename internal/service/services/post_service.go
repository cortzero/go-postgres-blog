package services

import (
	"context"
	"time"

	"github.com/cortzero/go-postgres-blog/internal/model/post"
	"github.com/cortzero/go-postgres-blog/internal/service/errors"
)

type PostService struct {
	Repository post.Repository
}

func NewPostService(repository post.Repository) *PostService {
	return &PostService{
		Repository: repository,
	}
}

func (service *PostService) CreatePost(ctx context.Context, post *post.Post) *errors.CustomError {
	// Set the timestamp of creation of the post
	post.CreatedAt = time.Now()

	// Save post
	err := service.Repository.Create(ctx, post)
	if err != nil {
		return errors.NewCustomError(
			"ERROR_CREATING_POST",
			"An error occurred while creating the post.",
			err.Error(),
			time.Now(),
		)
	}
	return nil
}

func (service *PostService) UpdatePost(ctx context.Context, id uint, post *post.Post) *errors.CustomError {
	// Check if the post exists
	existingPost, error_existing := service.GetPostById(ctx, id)
	if error_existing != nil {
		return error_existing
	}

	// Updating the existing post
	existingPost.Title = post.Title
	existingPost.Body = post.Body
	existingPost.UpdatedAt = time.Now()

	// Update post
	error_update := service.Repository.Update(ctx, id, existingPost)
	if error_update != nil {
		return errors.NewCustomError(
			"ERROR_UPDATING_POST",
			"An error occurred while updating the post.",
			error_update.Error(),
			time.Now(),
		)
	}

	return nil
}

func (service *PostService) DeletePost(ctx context.Context, id uint) *errors.CustomError {
	// Check if the post exists
	_, error_get := service.GetPostById(ctx, id)
	if error_get != nil {
		return error_get
	}

	// Deleting the post
	error_deleting := service.Repository.Delete(ctx, id)
	if error_deleting != nil {
		return errors.NewCustomError(
			"ERROR_DELETING_POST",
			"An error occurred while deleting the post.",
			error_deleting.Error(),
			time.Now(),
		)
	}
	return nil
}

func (service *PostService) GetAllPosts(ctx context.Context) ([]post.Post, *errors.CustomError) {
	posts, err := service.Repository.GetAll(ctx)
	if err != nil {
		return nil, errors.NewCustomError(
			"ERROR_GETTING_POSTS",
			"An error occurred while getting all posts.",
			err.Error(),
			time.Now(),
		)
	}
	return posts, nil
}

func (service *PostService) GetPostById(ctx context.Context, id uint) (post.Post, *errors.CustomError) {
	existingPost, err := service.Repository.GetById(ctx, id)
	if err != nil {
		return post.Post{}, errors.NewCustomError(
			"ERROR_GETTING_POST",
			"An error occurred while getting a post by its id.",
			err.Error(),
			time.Now(),
		)
	}
	return existingPost, nil
}

func (service *PostService) GetPostsByUserId(ctx context.Context, userId uint) ([]post.Post, *errors.CustomError) {
	return nil, nil
}
