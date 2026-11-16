package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
	"github.com/AanFLenrir/PetrivaBackend/config"
	"github.com/AanFLenrir/PetrivaBackend/helper"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidInput       = errors.New("invalid input data")
)

type AuthService interface {
	Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error)
	Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error)
	GetProfile(ctx context.Context, userID string) (*model.UserResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(ctx context.Context, input model.RegisterInput) (*model.AuthResponse, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Name = strings.TrimSpace(input.Name)

	if input.Email == "" || input.Password == "" || input.Name == "" {
		return nil, fmt.Errorf("%w: name, email, and password are required", ErrInvalidInput)
	}

	if len(input.Password) < 6 {
		return nil, fmt.Errorf("%w: password must be at least 6 characters", ErrInvalidInput)
	}

	// Only PET_OWNER can register publicly
	role, err := s.roleRepo.FindByName(ctx, model.RolePetOwner)
	if err != nil {
		return nil, fmt.Errorf("role %s not configured: %w", model.RolePetOwner, err)
	}

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		RoleID:   role.ID,
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
		Address:  input.Address,
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := helper.GenerateToken(s.cfg, user.ID, role.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthResponse{
		Token: token,
		User: model.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Address:   user.Address,
			Role:      role.Name,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		return nil, fmt.Errorf("%w: email and password are required", ErrInvalidInput)
	}

	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err := helper.CheckPassword(user.Password, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := helper.GenerateToken(s.cfg, user.ID, user.RoleName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthResponse{
		Token: token,
		User: model.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Phone:      user.Phone,
			Address:    user.Address,
			Role:       user.RoleName,
			ClinicID:   user.ClinicID,
			ClinicName: user.ClinicName,
			IsActive:   user.IsActive,
			CreatedAt:  user.CreatedAt,
		},
	}, nil
}

func (s *authService) GetProfile(ctx context.Context, userID string) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		Address:    user.Address,
		Role:       user.RoleName,
		ClinicID:   user.ClinicID,
		ClinicName: user.ClinicName,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
	}, nil
}
