package service

import (
	"context"
	"enactus/internal/models"
	"enactus/internal/repository"
	"fmt"
)

type AttendanceServiceInterface interface {
	AddAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error)
	GetAllAttendances(ctx context.Context) ([]models.Attendance, error)
	GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error)
	UpdateAttendance(ctx context.Context, id int, in *models.UpdateAttendanceInput) error
	DeleteAttendance(ctx context.Context, id int) error
}

type AttendanceService struct {
	AttendanceRepo *repository.AttendanceRepository
}

func (attendanceS *AttendanceService) AddAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error) {
	err := attendanceS.AttendanceRepo.AddAttendance(ctx, attendance)
	if err != nil {
		return nil, fmt.Errorf("failed to add an attendance: %w", err)
	}

	return attendance, nil
}

func (attendanceS *AttendanceService) GetAllAttendances(ctx context.Context) ([]models.Attendance, error) {
	attendances, err := attendanceS.AttendanceRepo.GetAllAttendances(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all attendances: %w", err)
	}

	return attendances, nil
}

func (attendanceS *AttendanceService) GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error) {
	attendance, err := attendanceS.AttendanceRepo.GetAttendanceById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get an attendance: %w", err)
	}

	return attendance, nil
}

func (attendanceS *AttendanceService) UpdateAttendance(ctx context.Context, id int, in *models.UpdateAttendanceInput) error {
	err := attendanceS.AttendanceRepo.UpdateAttendance(ctx, id, in)
	if err != nil {
		return fmt.Errorf("failed to update an attendance: %w", err)
	}

	return nil
}

func (attendanceS *AttendanceService) DeleteAttendance(ctx context.Context, id int) error {
	err := attendanceS.AttendanceRepo.DeleteAttendance(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete an attendance: %w", err)
	}

	return nil
}
