package data

import (
	"context"
	"log"
	"time"

	"github.com/cortzero/go-postgres-blog/internal/model/post"
)

type PostRepository struct {
	Data *Data
}

func NewPostRepository(connection *Data) *PostRepository {
	return &PostRepository{
		Data: connection,
	}
}

func (repository *PostRepository) GetAll(ctx context.Context) ([]post.Post, error) {
	query := `
	SELECT id, title, body, user_id, created_at, updated_at
	FROM posts;
	`
	rows, err := repository.Data.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []post.Post
	for rows.Next() {
		var p post.Post
		rows.Scan(&p.ID, &p.Title, &p.Body, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		posts = append(posts, p)
	}
	return posts, nil
}

func (repository *PostRepository) GetById(ctx context.Context, id uint) (post.Post, error) {
	query := `
	SELECT id, user_id, title, body, created_at, COALESCE(updated_at, '0001-01-01T00:00:00Z')
	FROM posts
	WHERE id = $1;
	`
	logger := log.Default()
	row := repository.Data.DB.QueryRowContext(ctx, query, id)
	var p post.Post
	err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Body, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		logger.Panicln(err.Error())
		return post.Post{}, err
	}
	return p, nil
}

func (repository *PostRepository) GetByUser(ctx context.Context, userId uint) ([]post.Post, error) {
	query := `
	SELECT id, title, body, user_id, created_at, updated_at
	FROM posts
	WHERE user_id=$1;
	`
	rows, err := repository.Data.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []post.Post
	for rows.Next() {
		var p post.Post
		rows.Scan(&p.ID, &p.Title, &p.Body, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		posts = append(posts, p)
	}
	return posts, nil
}

func (repository *PostRepository) Create(ctx context.Context, post *post.Post) error {
	insert := `
	INSERT INTO posts (user_id, title, body, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`
	stmt, err := repository.Data.DB.PrepareContext(ctx, insert)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, post.UserID, post.Title, post.Body, time.Now(), nil)
	err = row.Scan(&post.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PostRepository) Update(ctx context.Context, id uint, post post.Post) error {
	update := `
	UPDATE posts SET title=$1, body=$2, updated_at=$3
	WHERE id=$4;
	`
	stmt, err := repository.Data.DB.PrepareContext(ctx, update)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, post.Title, post.Body, post.UpdatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository PostRepository) Delete(ctx context.Context, id uint) error {
	delete := `
	DELETE FROM posts WHERE id=$1;
	`
	stmt, err := repository.Data.DB.PrepareContext(ctx, delete)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
