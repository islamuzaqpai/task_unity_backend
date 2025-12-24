package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DepartmentRepositoryInterface interface {
	AddDepartment(department models.Department) (*models.Department, error)
	GetDepartmentById(id int) (*models.Department, error)
	GetAllDepartments() ([]models.Department, error)
	UpdateDepartment(id int, newDepartment models.Department) (*models.Department, error)
	DeleteDepartment(id int) error
}

type DepartmentRepository struct {
	Pool *pgxpool.Pool
}

func (departmentRepo *DepartmentRepository) AddDepartment(department models.Department) (*models.Department, error) {
	row := departmentRepo.Pool.QueryRow(context.Background(),
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

func (departmentRepo *DepartmentRepository) GetDepartmentById(id int) (*models.Department, error) {
	row := departmentRepo.Pool.QueryRow(context.Background(),
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

func (departmentRepo *DepartmentRepository) GetAllDepartments() ([]models.Department, error) {
	rows, err := departmentRepo.Pool.Query(context.Background(),
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

func (departmentRepo *DepartmentRepository) UpdateDepartment(id int, newDepartment models.Department) (*models.Department, error) {
	row := departmentRepo.Pool.QueryRow(context.Background(),
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

func (departmentRepo *DepartmentRepository) DeleteDepartment(id int) error {
	_, err := departmentRepo.Pool.Exec(context.Background(),
		"DELETE FROM departments WHERE id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete a department: %w", err)
	}

	return nil
}
