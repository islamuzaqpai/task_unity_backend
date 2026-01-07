package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepositoryInterface interface {
	AddRole(ctx context.Context, role *models.Role) error
	GetRoleById(ctx context.Context, id int) (*models.Role, error)
	GetAllRoles(ctx context.Context) ([]models.Role, error)
	UpdateRole(ctx context.Context, id int, newRole models.Role) error
	DeleteRole(ctx context.Context, id int) error
}

type RoleRepository struct {
	Pool *pgxpool.Pool
}

func (roleRepo *RoleRepository) AddRole(ctx context.Context, role *models.Role) error {
	row := roleRepo.Pool.QueryRow(ctx,
		"INSERT INTO roles (name, department_id) VALUES ($1, $2) RETURNING id, name, department_id",
		role.Name,
		role.DepartmentId,
	)

	err := row.Scan(
		&role.Id,
		&role.Name,
		&role.DepartmentId,
	)

	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (roleRepo *RoleRepository) GetRoleById(ctx context.Context, id int) (*models.Role, error) {
	row := roleRepo.Pool.QueryRow(ctx,
		"SELECT id, name, department_id FROM roles WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	var role models.Role

	err := row.Scan(
		&role.Id,
		&role.Name,
		&role.DepartmentId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &role, nil
}

func (roleRepo *RoleRepository) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	rows, err := roleRepo.Pool.Query(ctx,
		"SELECT id, name, department_id FROM roles WHERE deleted_at IS NULL")

	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var role models.Role

		err = rows.Scan(
			&role.Id,
			&role.Name,
			&role.DepartmentId,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		roles = append(roles, role)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return roles, nil
}

func (roleRepo *RoleRepository) UpdateRole(ctx context.Context, id int, newRole models.Role) error {
	row := roleRepo.Pool.QueryRow(ctx,
		"UPDATE roles SET name = $1, department_id = $2, updated_at = now() WHERE id = $3 RETURNING id",
		newRole.Name,
		newRole.DepartmentId,
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

func (roleRepo *RoleRepository) DeleteRole(ctx context.Context, id int) error {
	_, err := roleRepo.Pool.Exec(ctx,
		"UPDATE roles SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete a role: %w", err)
	}

	return nil
}
