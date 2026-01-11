package service

import (
	"context"
	"enactus/internal/models"
	"enactus/internal/repository"
	"fmt"
)

type RoleServiceInterface interface {
	AddRole(ctx context.Context, role *models.Role) (*models.Role, error)
	GetAllRoles(ctx context.Context) ([]models.Role, error)
}

type RoleService struct {
	RoleRepo *repository.RoleRepository
}

func (roleS *RoleService) AddRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	checkRole, err := roleS.RoleRepo.RoleExists(ctx, role.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check a role: %w", err)
	}

	if checkRole {
		return nil, fmt.Errorf("role is already exists")
	}

	err = roleS.RoleRepo.AddRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to add a role: %w", err)
	}

	return role, nil
}

func (roleS *RoleService) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	roles, err := roleS.RoleRepo.GetAllRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

	return roles, nil
}
