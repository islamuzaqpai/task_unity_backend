package models

import (
	"time"
)

type RegisterInput struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	DepartmentId *int   `json:"department_id"`
}

type UpdateUserProfileInput struct {
	FullName string
	Email    string
}

type UsersRolesInput struct {
	Id     int
	UserId int
	RoleId int
}

type UpdateTaskInput struct {
	Title       string
	Description string
	Deadline    time.Time
	AssigneeId  int
	Status      string
}

type UpdateAttendanceInput struct {
	Status  string
	Comment string
}
