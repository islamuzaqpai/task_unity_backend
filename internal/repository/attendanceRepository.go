package repository

import (
	"context"
	"database/sql"
	"enactus/internal/apperrors"
	"enactus/internal/models"
	"enactus/internal/models/inputs"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceRepositoryInterface interface {
	AddAttendance(ctx context.Context, attendance *models.Attendance) error
	GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error)
	GetAllAttendances(ctx context.Context) ([]models.Attendance, error)
	UpdateAttendance(ctx context.Context, id int, in *inputs.UpdateAttendanceInput) error
	DeleteAttendance(ctx context.Context, id int) error
}

type AttendanceRepository struct {
	Pool *pgxpool.Pool
}

func NewAttendanceRepository(pool *pgxpool.Pool) *AttendanceRepository {
	return &AttendanceRepository{Pool: pool}
}

func (attendanceRepo *AttendanceRepository) AddAttendance(ctx context.Context, attendance *models.Attendance) error {
	query := `INSERT INTO attendance (user_id, date, department_id, status, comment, marked_by) VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id, created_at, updated_at, deleted_at`

	err := attendanceRepo.Pool.QueryRow(ctx, query,
		attendance.UserId,
		attendance.Date,
		attendance.DepartmentId,
		attendance.Status,
		attendance.Comment,
		attendance.MarkedBy,
	).Scan(
		&attendance.Id,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
		&attendance.DeletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (attendanceRepo *AttendanceRepository) GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error) {
	query := `SELECT id, date, user_id, department_id, status, comment, marked_by, created_at, updated_at FROM attendance WHERE id = $1 AND deleted_at IS NULL`

	var attendance models.Attendance
	err := attendanceRepo.Pool.QueryRow(ctx, query, id).
		Scan(
			&attendance.Id,
			&attendance.Date,
			&attendance.UserId,
			&attendance.DepartmentId,
			&attendance.Status,
			&attendance.Comment,
			&attendance.MarkedBy,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &attendance, nil
}

func (attendanceRepo *AttendanceRepository) GetAllAttendances(ctx context.Context) ([]models.Attendance, error) {
	query := `SELECT id, date, user_id, department_id, status, comment, marked_by, created_at, updated_at FROM attendance WHERE deleted_at is null`

	rows, err := attendanceRepo.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to get all attendances: %w", err)
	}
	defer rows.Close()

	var attendances []models.Attendance
	for rows.Next() {
		var attendance models.Attendance

		err = rows.Scan(
			&attendance.Id,
			&attendance.Date,
			&attendance.UserId,
			&attendance.DepartmentId,
			&attendance.Status,
			&attendance.Comment,
			&attendance.MarkedBy,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		attendances = append(attendances, attendance)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return attendances, nil
}

func (attendanceRepo *AttendanceRepository) UpdateAttendance(ctx context.Context, id int, in *inputs.UpdateAttendanceInput) error {
	query := `UPDATE attendance SET `
	var args []any
	i := 1

	if in.Status != nil {
		query += fmt.Sprintf(" status = $%d,", i)
		args = append(args, in.Status)
		i++

		if in.MarkedBy != nil {
			query += fmt.Sprintf(" marked_by = $%d,", i)
			args = append(args, in.MarkedBy)
			i++
		}
	}

	if in.Comment != nil {
		query += fmt.Sprintf(" comment = $%d,", i)
		args = append(args, in.Comment)
		i++
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	result, err := attendanceRepo.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update attendance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (attendanceRepo *AttendanceRepository) DeleteAttendance(ctx context.Context, id int) error {
	query := `UPDATE attendance SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`

	res, err := attendanceRepo.Pool.Exec(ctx, query, id)

	if res.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to delete an attendance: %w", err)
	}

	return nil
}
