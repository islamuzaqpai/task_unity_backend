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
