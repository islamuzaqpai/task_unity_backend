package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DepartmentRepositoryInterface interface {
	AddDepartment(ctx context.Context, department *models.Department) error
	GetDepartmentById(ctx context.Context, id int) (*models.Department, error)
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	UpdateDepartment(ctx context.Context, id int, newDepartment models.Department) error
	DeleteDepartment(ctx context.Context, id int) error
}

type DepartmentRepository struct {
	Pool *pgxpool.Pool
}

func (departmentRepo *DepartmentRepository) AddDepartment(ctx context.Context, department *models.Department) error {
	row := departmentRepo.Pool.QueryRow(ctx,
		"INSERT INTO departments (name) VALUES ($1) RETURNING id, created_at, updated_at",
		department.Name,
	)

	err := row.Scan(
		&department.Id,
		&department.CreatedAt,
		&department.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (departmentRepo *DepartmentRepository) GetDepartmentById(ctx context.Context, id int) (*models.Department, error) {
	row := departmentRepo.Pool.QueryRow(ctx,
		"SELECT id, name, created_at, updated_at, deleted_at FROM departments WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	var department models.Department
	err := row.Scan(
		&department.Id,
		&department.Name,
		&department.CreatedAt,
		&department.UpdatedAt,
		&department.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &department, nil
}

func (departmentRepo *DepartmentRepository) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	rows, err := departmentRepo.Pool.Query(ctx,
		"SELECT id, name, created_at, updated_at FROM departments WHERE deleted_at is null ",
	)

	if err != nil {
		return nil, fmt.Errorf("failed to select all departments: %w", err)
	}
	defer rows.Close()

	var departments []models.Department

	for rows.Next() {
		var department models.Department

		err = rows.Scan(
			&department.Id,
			&department.Name,
			&department.CreatedAt,
			&department.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		departments = append(departments, department)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return departments, nil
}

func (departmentRepo *DepartmentRepository) UpdateDepartment(ctx context.Context, id int, newDepartment *models.Department) error {
	row := departmentRepo.Pool.QueryRow(ctx,
		"UPDATE departments SET name = $1, updated_at = now() WHERE id = $2 RETURNING id",
		newDepartment.Name,
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

func (departmentRepo *DepartmentRepository) DeleteDepartment(ctx context.Context, id int) error {
	_, err := departmentRepo.Pool.Exec(ctx,
		"UPDATE departments SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete a department: %w", err)
	}

	return nil
}

func (departmentRepo *DepartmentRepository) DepartmentExists(ctx context.Context, departmentName string) (bool, error) {
	row := departmentRepo.Pool.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM departments WHERE name = $1)",
		departmentName,
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
