package service

import (
	"context"
	"enactus/internal/auth"
	"enactus/internal/models"
	"enactus/internal/repository"
	"enactus/internal/utils"
	"fmt"
	"unicode/utf8"
)

type UserServiceInterface interface {
	Register(ctx context.Context, input models.RegisterInput) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type UserService struct {
	UserRepo  *repository.UserRepository
	JwtSecret *auth.JWTSecret
}

func (userS *UserService) Register(ctx context.Context, input models.RegisterInput) (*models.User, error) {
	checkEmail, err := userS.UserRepo.EmailExists(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("email already exist")
	}

	if checkEmail {
		return nil, fmt.Errorf("email already exists")
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

	err = userS.UserRepo.AddUser(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to add an user: %w", err)
	}

	user.Password = ""
	return &user, nil
}

func (userS *UserService) Login(ctx context.Context, email, password string) (string, error) {
	authUser, err := userS.UserRepo.GetAuthUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to find user with this email: %w", err)
	}

	isValid, err := utils.ValidatePassword(password, authUser.Password)
	if err != nil {
		return "", fmt.Errorf("failed to validate a password: %w", err)
	}

	if !isValid {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	tokenStr, err := userS.JwtSecret.GenerateToken(authUser)
	if err != nil {
		return "", fmt.Errorf("failed to generate a token: %w", err)
	}

	return tokenStr, nil
}

func (userS *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := userS.UserRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}
