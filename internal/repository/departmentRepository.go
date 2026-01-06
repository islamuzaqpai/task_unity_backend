package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DepartmentRepositoryInterface interface {
	AddDepartment(ctx context.Context, department models.Department) (*models.Department, error)
	GetDepartmentById(ctx context.Context, id int) (*models.Department, error)
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	UpdateDepartment(ctx context.Context, id int, newDepartment models.Department) (*models.Department, error)
	DeleteDepartment(ctx context.Context, id int) error
}

type DepartmentRepository struct {
	Pool *pgxpool.Pool
}

func (departmentRepo *DepartmentRepository) AddDepartment(ctx context.Context, department models.Department) (*models.Department, error) {
	row := departmentRepo.Pool.QueryRow(ctx,
		"INSERT INTO departments (name) VALUES ($1) RETURNING id, name",
		department.Name,
	)

	err := row.Scan(
		&department.Id,
		&department.Name,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &department, nil
}

func (departmentRepo *DepartmentRepository) GetDepartmentById(ctx context.Context, id int) (*models.Department, error) {
	row := departmentRepo.Pool.QueryRow(ctx,
		"SELECT id, name, created_at, updated_at, deleted_at FROM departments WHERE id = $1",
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

	defer rows.Close()

	return departments, nil
}

func (departmentRepo *DepartmentRepository) UpdateDepartment(ctx context.Context, id int, newDepartment models.Department) (*models.Department, error) {
	row := departmentRepo.Pool.QueryRow(ctx,
		"UPDATE departments SET name = $1 WHERE id = $2 RETURNING id, name, created_at, updated_at",
		newDepartment.Name,
		id,
	)

	var department models.Department
	err := row.Scan(
		&department.Id,
		&department.Name,
		&department.CreatedAt,
		&department.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &department, nil
}

func (departmentRepo *DepartmentRepository) DeleteDepartment(ctx context.Context, id int) error {
	_, err := departmentRepo.Pool.Exec(ctx,
		"DELETE FROM departments WHERE id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete a department: %w", err)
	}

	return nil
}
