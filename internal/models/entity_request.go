package models

type RegisterInput struct {
	FullName     string
	Email        string
	Password     string
	DepartmentId *int
}
