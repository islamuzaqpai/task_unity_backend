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

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateUserProfileInput struct {
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
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

type UpdatePasswordInput struct {
	Password string `json:"password"`
}
