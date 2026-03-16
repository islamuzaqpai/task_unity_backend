package inputs

type UsersRolesInput struct {
	Id     int
	UserId int
	RoleId int
}

type UpdateAttendanceInput struct {
	Status  string
	Comment string
}
