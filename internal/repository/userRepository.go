package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	GetUserById(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	AddUser(user *models.User) (*models.User, error)
	UpdateUserProfile(id int, newUser models.User) (*models.User, error)
	UpdateUserRole(userId, roleId int) error
	UpdateUserPassword(id int, newPassword string) error
	DeleteUser(id int) error
}

type UserRepository struct {
	Pool *pgxpool.Pool
}

func (userRepo *UserRepository) GetUserById(id int) (*models.User, error) {
	row := userRepo.Pool.QueryRow(context.Background(),
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

func (userRepo *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User

	rows, err := userRepo.Pool.Query(context.Background(), "SELECT id, full_name, email, department_id, created_at, updated_at, deleted_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

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

	defer rows.Close()

	return users, nil
}

func (userRepo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	row := userRepo.Pool.QueryRow(context.Background(),
		"SELECT id, full_name, email, department_id, created_at, updated_at, deleted_at FROM users WHERE email = $1",
		email,
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

func (userRepo *UserRepository) AddUser(user *models.User) (*models.User, error) {
	row := userRepo.Pool.QueryRow(context.Background(),
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
		return nil, fmt.Errorf("failed to add a user: %w", err)
	}

	return user, nil
}

func (userRepo *UserRepository) UpdateUserProfile(id int, newUser models.User) (*models.User, error) {
	row := userRepo.Pool.QueryRow(context.Background(),
		"UPDATE users SET full_name = $1, email = $2, updated_at = now() WHERE id = $3 RETURNING id, full_name, email, department_id, created_at, updated_at",
		newUser.FullName,
		newUser.Email,
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
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &user, nil
}

func (userRepo *UserRepository) UpdateUserRole(userId, newRoleId int) error {
	row := userRepo.Pool.QueryRow(context.Background(),
		"UPDATE users_roles SET user_id = $1,role_id = $2 WHERE id = $3 RETURNING user_id, role_id")

	err := row.Scan(
		&userId, //затестить
		&newRoleId,
	)

	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (userRepo *UserRepository) UpdateUserPassword(id int, newPassword string) error {
	_, err := userRepo.Pool.Exec(context.Background(),
		"UPDATE users SET password = $1 WHERE id = $2",
		newPassword,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update a password: %w", err)
	}

	return nil
}

func (userRepo *UserRepository) DeleteUser(id int) error {
	_, err := userRepo.Pool.Exec(context.Background(),
		"DELETE FROM users WHERE id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete an user: %w", err)
	}

	return nil
}
