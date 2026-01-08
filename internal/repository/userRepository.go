package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	GetUserById(ctx context.Context, id int) (*models.User, error)
	GetAuthUserByEmail(ctx context.Context, email string) (*models.AuthUser, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	AddUser(ctx context.Context, user *models.User) error
	UpdateUserProfile(ctx context.Context, id int, newUser models.User) error
	UpdateUserPassword(ctx context.Context, id int, newPassword string) error
	DeleteUser(ctx context.Context, id int) error
}

type UserRepository struct {
	Pool *pgxpool.Pool
}

func (userRepo *UserRepository) GetUserById(ctx context.Context, id int) (*models.User, error) {
	row := userRepo.Pool.QueryRow(ctx,
		"SELECT id, full_name, email, department_id, created_at, updated_at, deleted_at FROM users WHERE id = $1",
		id,
	)

	var user models.User
	err := row.Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.DepartmentId,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan a row: %w", err)
	}

	return &user, nil
}

func (userRepo *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	rows, err := userRepo.Pool.Query(ctx, "SELECT id, full_name, email, department_id, created_at, updated_at, deleted_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err = rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Email,
			&user.DepartmentId,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan a row: %w", err)
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

func (userRepo *UserRepository) GetAuthUserByEmail(ctx context.Context, email string) (*models.AuthUser, error) {
	row := userRepo.Pool.QueryRow(ctx,
		"SELECT id, email, password, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL",
		email,
	)

	var authUser models.AuthUser
	err := row.Scan(
		&authUser.Id,
		&authUser.Email,
		&authUser.Password,
		&authUser.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan a row: %w", err)
	}

	return &authUser, nil
}

func (userRepo *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	row := userRepo.Pool.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)",
		email,
	)

	var exists bool
	err := row.Scan(
		&exists,
	)
	if err != nil {
		return false, fmt.Errorf("failed to scan: %w", err)
	}

	return exists, nil
}

func (userRepo *UserRepository) AddUser(ctx context.Context, user *models.User) error {
	row := userRepo.Pool.QueryRow(ctx,
		"INSERT INTO users (full_name, email, password, department_id) VALUES ($1, $2, $3, $4) RETURNING id, full_name, email, department_id",
		user.FullName,
		user.Email,
		user.Password,
		user.DepartmentId,
	)

	err := row.Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.DepartmentId,
	)

	if err != nil {
		return fmt.Errorf("failed to add a user: %w", err)
	}

	return nil
}

func (userRepo *UserRepository) UpdateUserProfile(ctx context.Context, id int, in models.UpdateUserProfileInput) error {
	row := userRepo.Pool.QueryRow(ctx,
		"UPDATE users SET full_name = $1, email = $2, updated_at = now() WHERE id = $3 RETURNING id",
		in.FullName,
		in.Email,
		id,
	)

	var returnedId int

	err := row.Scan(
		&returnedId,
	)

	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (userRepo *UserRepository) UpdateUserPassword(ctx context.Context, id int, newPassword string) error {
	_, err := userRepo.Pool.Exec(ctx,
		"UPDATE users SET password = $1 WHERE id = $2",
		newPassword,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update a password: %w", err)
	}

	return nil
}

func (userRepo *UserRepository) DeleteUser(ctx context.Context, id int) error {
	_, err := userRepo.Pool.Exec(ctx,
		"UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete an user: %w", err)
	}

	return nil
}
