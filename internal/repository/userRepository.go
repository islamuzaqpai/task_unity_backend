package repository

import (
	"context"
	"enactus/internal/models"
	"enactus/internal/models/inputs"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	GetUserById(ctx context.Context, id int) (*models.User, error)
	GetAuthUserByEmail(ctx context.Context, email string) (*inputs.AuthUser, string, error)
	EmailExists(ctx context.Context, email *string) (bool, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	AddUser(ctx context.Context, user *models.User) error
	UpdateUserProfile(ctx context.Context, id int, in inputs.UpdateUserProfileInput) error
	UpdateUserPassword(ctx context.Context, id int, newPassword string) error
	DeleteUser(ctx context.Context, id int) error
}

type UserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{Pool: pool}
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

func (userRepo *UserRepository) GetAuthUserByEmail(ctx context.Context, email string) (*inputs.AuthUser, string, error) {
	query := `SELECT users.id, users.email, users.password, users.deleted_at, roles.name
				FROM users
				JOIN users_roles ON users.id = users_roles.user_id
				JOIN roles ON users_roles.role_id = roles.id
				WHERE users.email = $1 AND users.deleted_at IS NULL`

	var (
		authUser inputs.AuthUser
		role     string
	)
	err := userRepo.Pool.QueryRow(ctx,
		query,
		email,
	).Scan(
		&authUser.Id,
		&authUser.Email,
		&authUser.Password,
		&authUser.DeletedAt,
		&role,
	)

	if err != nil {
		return nil, "", fmt.Errorf("failed to scan a row: %w", err)
	}

	return &authUser, role, nil
}

func (userRepo *UserRepository) EmailExists(ctx context.Context, email *string) (bool, error) {
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
	tx, err := userRepo.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO users (full_name, email, password, department_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err = tx.QueryRow(ctx, query, user.FullName, user.Email, user.Password, user.DepartmentId).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to scan user: %w", err)
	}

	query = `INSERT INTO users_roles (user_id, role_id) VALUES ($1, (SELECT id FROM roles WHERE name = 'user'))`
	_, err = tx.Exec(ctx, query, user.Id)
	if err != nil {
		return fmt.Errorf("failed to add role: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transactions: %w", err)
	}

	return nil
}

func (userRepo *UserRepository) UpdateUserProfile(ctx context.Context, id int, in inputs.UpdateUserProfileInput) error {
	query := `UPDATE users SET `
	args := []any{}
	i := 1

	if in.FullName != nil {
		query += fmt.Sprintf("full_name = $%d,", i)
		args = append(args, *in.FullName)
		i++
	}

	if in.Email != nil {
		query += fmt.Sprintf("email = $%d,", i)
		args = append(args, *in.Email)
		i++
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id = $%d", i)

	args = append(args, id)
	_, err := userRepo.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
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
