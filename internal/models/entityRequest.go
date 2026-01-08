package models

type RegisterInput struct {
	FullName     string
	Email        string
	Password     string
	DepartmentId *int
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
