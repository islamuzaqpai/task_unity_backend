package service

import (
	"enactus/internal/models"
	"enactus/internal/repository"
	"enactus/internal/utils"
	"fmt"
	"unicode/utf8"
)

type UserServiceInterface interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (*models.User, error)
}

type UserService struct {
	UserRepo *repository.UserRepository
}

type RegisterInput struct {
	FullName     string
	Email        string
	Password     string
	DepartmentId *int
}

func (userS *UserService) Register(input RegisterInput) (*models.User, error) {
	checkEmail, err := userS.UserRepo.EmailExists(input.Email)
	if err != nil || checkEmail != false {
		return nil, fmt.Errorf("email already exist")
	}

	if utf8.RuneCountInString(input.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters long")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash a password: %w", err)
	}

	var user models.User
	user.FullName = input.FullName
	user.Email = input.Email
	user.Password = hashedPassword
	user.DepartmentId = input.DepartmentId

	added, err := userS.UserRepo.AddUser(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to add an user: %w", err)
	}

	return added, nil
}
