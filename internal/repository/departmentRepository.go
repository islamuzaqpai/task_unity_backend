package repository

import (
	"context"
	"database/sql"
	"enactus/internal/apperrors"
	"enactus/internal/models"
	"errors"
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

func NewDepartmentRepository(pool *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{Pool: pool}
}

func (departmentRepo *DepartmentRepository) AddDepartment(ctx context.Context, department *models.Department) error {
	query := `INSERT INTO departments (name) VALUES ($1) RETURNING id, created_at, updated_at`
	err := departmentRepo.Pool.QueryRow(ctx, query, department.Name).Scan(
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
	query := `SELECT id, name, created_at, updated_at, deleted_at FROM departments WHERE id = $1 AND deleted_at IS NULL`
	var department models.Department

	err := departmentRepo.Pool.QueryRow(ctx, query, id).Scan(
		&department.Id,
		&department.Name,
		&department.CreatedAt,
		&department.UpdatedAt,
		&department.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &department, nil
}

func (departmentRepo *DepartmentRepository) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	query := `SELECT id, name, created_at, updated_at FROM departments WHERE deleted_at is null`

	rows, err := departmentRepo.Pool.Query(ctx, query)

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
	query := `UPDATE departments SET name = $1, updated_at = now() WHERE id = $2 RETURNING id`

	row := departmentRepo.Pool.QueryRow(ctx,
		query,
		newDepartment.Name,
		id,
	)

	var returnedId int
	err := row.Scan(
		&returnedId,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrNotFound
		}

		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (departmentRepo *DepartmentRepository) DeleteDepartment(ctx context.Context, id int) error {
	query := `UPDATE departments SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`

	res, err := departmentRepo.Pool.Exec(ctx, query, id)

	if err != nil {
		if res.RowsAffected() == 0 {
			return apperrors.ErrNotFound
		}

		return fmt.Errorf("failed to delete a department: %w", err)
	}

	return nil
}

func (departmentRepo *DepartmentRepository) DepartmentExists(ctx context.Context, departmentName string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM departments WHERE name = $1)`

	row := departmentRepo.Pool.QueryRow(ctx, query, departmentName)

	var exists bool
	err := row.Scan(
		&exists,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apperrors.ErrNotFound
		}

		return false, fmt.Errorf("failed to scan: %w", err)
	}

	return exists, nil
}
