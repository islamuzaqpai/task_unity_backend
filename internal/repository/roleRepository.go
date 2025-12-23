package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepositoryInterface interface {
	AddRole(role models.Role) (*models.Role, error)
	GetRoleById(id int) (*models.Role, error)
	GetAllRoles() ([]models.Role, error)
	UpdateRole(id int, newRole models.Role) (*models.Role, error)
	DeleteRole(id int) error
}

type RoleRepository struct {
	Pool *pgxpool.Pool
}

func (roleRepo *RoleRepository) AddRole(role models.Role) (*models.Role, error) {
	row := roleRepo.Pool.QueryRow(context.Background(),
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
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &role, nil
}

func (roleRepo *RoleRepository) GetRoleById(id int) (*models.Role, error) {
	row := roleRepo.Pool.QueryRow(context.Background(),
		"SELECT id, name, department_id FROM roles WHERE id = $1",
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

func (roleRepo *RoleRepository) GetAllRoles() ([]models.Role, error) {
	rows, err := roleRepo.Pool.Query(context.Background(),
		"SELECT id, name, department_id FROM roles")

	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

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

	defer rows.Close()

	return roles, nil
}

func (roleRepo *RoleRepository) UpdateRole(id int, newRole models.Role) (*models.Role, error) {
	row := roleRepo.Pool.QueryRow(context.Background(),
		"UPDATE roles SET name = $1, department_id = $2 WHERE id = $3 RETURNING id, name, department_id",
		newRole.Name,
		newRole.DepartmentId,
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

func (roleRepo *RoleRepository) DeleteRole(id int) error {
	_, err := roleRepo.Pool.Exec(context.Background(),
		"DELETE FROM roles WHERE id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete a role: %w", err)
	}

	return nil
}
