package data

import (
	"context"

	"github.com/cortzero/go-postgres-blog/internal/model/user"
)

type UserRepositoy struct {
	Data *Data
}

func NewUserRepository() *UserRepositoy {
	return &UserRepositoy{
		Data: New(),
	}
}

func (repository *UserRepositoy) GetAll(ctx context.Context) ([]user.User, error) {
	query := `
	SELECT id, first_name, last_name, username, email, picture, created_at, updated_at
	FROM users;
	`
	rows, err := repository.Data.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []user.User
	for rows.Next() {
		var user user.User
		rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email,
			&user.Picture, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}
	return users, nil
}

func (repository *UserRepositoy) GetById(ctx context.Context, id uint) (user.User, error) {
	query := `
	SELECT id, first_name, last_name, username, email, picture, created_at, updated_at
	FROM users
	WHERE id = $1;
	`
	row := repository.Data.DB.QueryRowContext(ctx, query, id)
	var u user.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username, &u.Email,
		&u.Picture, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return user.User{}, err
	}
	return u, nil
}

func (repository *UserRepositoy) GetByUsername(ctx context.Context, username string) (user.User, error) {
	query := `
	SELECT id, first_name, last_name, username, password, email, picture, created_at, updated_at
	FROM users
	WHERE username = $1;
	`
	row := repository.Data.DB.QueryRowContext(ctx, query, username)
	var u user.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username, &u.PasswordHash, &u.Email,
		&u.Picture, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return user.User{}, nil
	}
	return u, nil
}

func (repository *UserRepositoy) Create(ctx context.Context, user *user.User) error {
	insert := `
	INSERT INTO users (first_name, last_name, username, password, email, picture, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id;
	`
	// Sets default photo
	// if user.Picture == "" {
	// 	user.Picture = "https://placekitten.com/g/300/300"
	// }

	// Hashes the password
	if err := user.HashPassword(); err != nil {
		return err
	}

	row := repository.Data.DB.QueryRowContext(ctx, insert,
		user.FirstName, user.LastName, user.Username, user.PasswordHash, user.Email, user.Picture, user.CreatedAt, user.UpdatedAt,
	)

	err := row.Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repository *UserRepositoy) Update(ctx context.Context, id uint, user user.User) error {
	update := `
	UPDATE users SET first_name=$1, last_name=$2, email=$3, picture=$4, updated_at=$5
	WHERE id=$6;
	`
	stmt, err := repository.Data.DB.PrepareContext(ctx, update)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Email, user.Picture, user.UpdatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *UserRepositoy) Delete(ctx context.Context, id uint) error {
	delete := `
	DELETE FROM users WHERE id=$1;
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
