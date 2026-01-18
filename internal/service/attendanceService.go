package service

import (
	"context"
	"enactus/internal/models"
	"enactus/internal/repository"
	"fmt"
)

type AttendanceServiceInterface interface {
	AddAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error)
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
