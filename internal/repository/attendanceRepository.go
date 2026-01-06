package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceRepositoryInterface interface {
	AddAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error)
	GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error)
	GetAllAttendance(ctx context.Context) ([]models.Attendance, error)
	UpdateAttendance(ctx context.Context, id int, newAttendance models.Attendance) (*models.Attendance, error)
	DeleteAttendance(ctx context.Context, id int) error
}

type AttendanceRepository struct {
	Pool *pgxpool.Pool
}

func (attendanceRepo *AttendanceRepository) AddAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error) {
	row := attendanceRepo.Pool.QueryRow(ctx,
		"INSERT INTO attendance (user_id, department_id, status, comment, marked_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, user_id, department_id, status, comment, marked_by, created_at, updated_at",
		attendance.UserId,
		attendance.DepartmentId,
		attendance.Status,
		attendance.Comment,
		attendance.MarkedBy,
	)

	err := row.Scan(
		&attendance.Id,
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

	return attendance, nil
}

func (attendanceRepo *AttendanceRepository) GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error) {
	row := attendanceRepo.Pool.QueryRow(ctx,
		"SELECT id, user_id, department_id, status, comment, marked_by, created_at, updated_at FROM attendance WHERE id = $1",
		id,
	)

	var attendance models.Attendance
	err := row.Scan(
		&attendance.Id,
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

	return &attendance, nil
}

func (attendanceRepo *AttendanceRepository) GetAllAttendance(ctx context.Context) ([]models.Attendance, error) {
	rows, err := attendanceRepo.Pool.Query(ctx,
		"SELECT id, user_id, department_id, status, comment, marked_by, created_at, updated_at FROM attendance WHERE deleted_at is null")

	if err != nil {
		return nil, fmt.Errorf("failed to select all attendances: %w", err)
	}

	var attendances []models.Attendance
	for rows.Next() {
		var attendance models.Attendance

		err = rows.Scan(
			&attendance.Id,
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

	defer rows.Close()

	return attendances, nil
}

func (attendanceRepo *AttendanceRepository) UpdateAttendance(ctx context.Context, id int, newAttendance models.Attendance) (*models.Attendance, error) {
	row := attendanceRepo.Pool.QueryRow(ctx,
		"UPDATE attendance SET user_id = $1, department_id = $2, status = $3, comment = $4, marked_by = $5  WHERE id = $6 AND attendance.deleted_at IS NULL RETURNING id, user_id, department_id, status, comment, marked_by, created_at, updated_at",
		newAttendance.UserId,
		newAttendance.DepartmentId,
		newAttendance.Status,
		newAttendance.Comment,
		newAttendance.MarkedBy,
		id,
	)

	var attendance models.Attendance
	err := row.Scan(
		&attendance.Id,
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

	return &attendance, nil

}
func (attendanceRepo *AttendanceRepository) DeleteAttendance(ctx context.Context, id int) error {
	_, err := attendanceRepo.Pool.Exec(ctx,
		"DELETE FROM attendance WHERE id = $1",
		id)

	if err != nil {
		return fmt.Errorf("failed to delete an attendance: %w", err)
	}

	return nil
}
