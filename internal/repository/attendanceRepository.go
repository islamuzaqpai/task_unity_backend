package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceRepositoryInterface interface {
	AddAttendance(ctx context.Context, attendance *models.Attendance) error
	GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error)
	GetAllAttendance(ctx context.Context) ([]models.Attendance, error)
	UpdateAttendance(ctx context.Context, id int, newAttendance models.Attendance) error
	DeleteAttendance(ctx context.Context, id int) error
}

type AttendanceRepository struct {
	Pool *pgxpool.Pool
}

func (attendanceRepo *AttendanceRepository) AddAttendance(ctx context.Context, attendance *models.Attendance) error {
	row := attendanceRepo.Pool.QueryRow(ctx,
		"INSERT INTO attendance (user_id, attendance_date, department_id, status, comment, marked_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at, deleted_at",
		attendance.UserId,
		attendance.Date,
		attendance.DepartmentId,
		attendance.Status,
		attendance.Comment,
		attendance.MarkedBy,
	)

	err := row.Scan(
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
	row := attendanceRepo.Pool.QueryRow(ctx,
		"SELECT id, attendance_date, user_id, department_id, status, comment, marked_by, created_at, updated_at FROM attendance WHERE id = $1 AND deleted_at IS NULL",
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
		"SELECT id, attendance_date, user_id, department_id, status, comment, marked_by, created_at, updated_at FROM attendance WHERE deleted_at is null")

	if err != nil {
		return nil, fmt.Errorf("failed to select all attendances: %w", err)
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

func (attendanceRepo *AttendanceRepository) UpdateAttendance(ctx context.Context, id int, newAttendance models.Attendance) error {
	row := attendanceRepo.Pool.QueryRow(ctx,
		"UPDATE attendance SET user_id = $1, attendance_date = $2, department_id = $3, status = $4, comment = $5, marked_by = $6, updated_at = now()  WHERE id = $7 AND deleted_at IS NULL RETURNING id",
		newAttendance.UserId,
		newAttendance.Date,
		newAttendance.DepartmentId,
		newAttendance.Status,
		newAttendance.Comment,
		newAttendance.MarkedBy,
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
func (attendanceRepo *AttendanceRepository) DeleteAttendance(ctx context.Context, id int) error {
	_, err := attendanceRepo.Pool.Exec(ctx,
		"UPDATE attendance SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL",
		id)

	if err != nil {
		return fmt.Errorf("failed to delete an attendance: %w", err)
	}

	return nil
}
