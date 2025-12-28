package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceRepositoryInterface interface {
	AddAttendance(attendance *models.Attendance) (*models.Attendance, error)
	GetAttendanceById(id int) (*models.Attendance, error)
	GetAllAttendance() ([]models.Attendance, error)
	UpdateAttendance(id int, newAttendance models.Attendance) (*models.Attendance, error)
	DeleteAttendance(id int) error
}

type AttendanceRepository struct {
	Pool *pgxpool.Pool
}

func (attendanceRepo *AttendanceRepository) AddAttendance(attendance *models.Attendance) (*models.Attendance, error) {
	row := attendanceRepo.Pool.QueryRow(context.Background(),
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

func (attendanceRepo *AttendanceRepository) GetAttendanceById(id int) (*models.Attendance, error) {
	row := attendanceRepo.Pool.QueryRow(context.Background(),
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

func (attendanceRepo *AttendanceRepository) GetAllAttendance() ([]models.Attendance, error) {
	rows, err := attendanceRepo.Pool.Query(context.Background(),
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

func (attendanceRepo *AttendanceRepository) UpdateAttendance(id int, newAttendance models.Attendance) (*models.Attendance, error) {
	return nil, nil
}
func (attendanceRepo *AttendanceRepository) DeleteAttendance(id int) error {
	return nil
}
