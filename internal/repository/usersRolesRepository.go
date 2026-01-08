package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRolesRepositoryInterface interface {
	SetUserRole(ctx context.Context, in models.UsersRolesInput) error
	DeleteUserRole(ctx context.Context, id int) error
}

type UsersRolesRepository struct {
	Pool *pgxpool.Pool
}

func (usersRolesRepo *UsersRolesRepository) SetUserRole(ctx context.Context, in models.UsersRolesInput) error {
	_, err := usersRolesRepo.Pool.Exec(ctx,
		"INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET role_id = EXCLUDED.role_id",
		in.UserId,
		in.RoleId,
	)

	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (usersRolesRepo *UsersRolesRepository) DeleteUserRole(ctx context.Context, userId int) error {
	_, err := usersRolesRepo.Pool.Exec(ctx,
		"DELETE FROM users_roles WHERE user_id = $1",
		userId,
	)

	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	return nil
}
