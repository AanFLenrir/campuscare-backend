package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
	"github.com/AanFLenrir/PetrivaBackend/helper"
)

var (
	ErrRoleInvalid       = errors.New("invalid role specified")
	ErrClinicRequired    = errors.New("clinic_id is required for veterinarian role")
	ErrInvalidVeterinary = errors.New("veterinarian must be assigned to an active clinic")
)

type UserService interface {
	ListUsers(ctx context.Context, roleFilter string) ([]model.UserResponse, error)
	GetUser(ctx context.Context, id string) (*model.UserResponse, error)
	CreateUser(ctx context.Context, input model.CreateUserInput) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	ListVeterinarians(ctx context.Context, clinicID *string) ([]model.UserResponse, error)
}

type userService struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	clinicRepo repository.ClinicRepository
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, clinicRepo repository.ClinicRepository) UserService {
	return &userService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		clinicRepo: clinicRepo,
	}
}

func (s *userService) ListUsers(ctx context.Context, roleFilter string) ([]model.UserResponse, error) {
	users, err := s.userRepo.List(ctx, roleFilter)
	if err != nil {
		return nil, err
	}
	res := make([]model.UserResponse, len(users))
	for i, u := range users {
		res[i] = model.UserResponse{
			ID:         u.ID,
			Name:       u.Name,
			Email:      u.Email,
			Phone:      u.Phone,
			Address:    u.Address,
			Role:       u.RoleName,
			ClinicID:   u.ClinicID,
			ClinicName: u.ClinicName,
			IsActive:   u.IsActive,
			CreatedAt:  u.CreatedAt,
		}
	}
	return res, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		Address:    u.Address,
		Role:       u.RoleName,
		ClinicID:   u.ClinicID,
		ClinicName: u.ClinicName,
		IsActive:   u.IsActive,
		CreatedAt:  u.CreatedAt,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.UserResponse, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	input.RoleName = strings.TrimSpace(strings.ToUpper(input.RoleName))

	if input.Email == "" || input.Password == "" || input.Name == "" || input.RoleName == "" {
		return nil, fmt.Errorf("%w: name, email, password, and role are required", ErrInvalidInput)
	}

	role, err := s.roleRepo.FindByName(ctx, input.RoleName)
	if err != nil {
		return nil, ErrRoleInvalid
	}

	// Validate VETERINARIAN role requirements
	if input.RoleName == model.RoleVeterinarian {
		if input.ClinicID == nil || *input.ClinicID == "" {
			return nil, ErrClinicRequired
		}
		// Check clinic exists and is active
		clinic, err := s.clinicRepo.FindByID(ctx, *input.ClinicID)
		if err != nil || !clinic.IsActive {
			return nil, ErrInvalidVeterinary
		}
	}

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		RoleID:   role.ID,
		RoleName: role.Name,
		ClinicID: input.ClinicID,
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

	return s.GetUser(ctx, user.ID)
}

func (s *userService) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (*model.UserResponse, error) {
	existing, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		existing.Name = strings.TrimSpace(input.Name)
	}
	if input.Phone != "" {
		existing.Phone = input.Phone
	}
	if input.Address != "" {
		existing.Address = input.Address
	}
	if input.IsActive != nil {
		existing.IsActive = *input.IsActive
	}
	if input.ClinicID != nil {
		if *input.ClinicID != "" {
			clinic, err := s.clinicRepo.FindByID(ctx, *input.ClinicID)
			if err != nil || !clinic.IsActive {
				return nil, ErrInvalidVeterinary
			}
			existing.ClinicID = input.ClinicID
		} else {
			existing.ClinicID = nil
		}
	}
	if input.Password != "" {
		hashed, err := helper.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}
		existing.Password = hashed
	} else {
		existing.Password = "" // will not overwrite
	}

	if err := s.userRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) ListVeterinarians(ctx context.Context, clinicID *string) ([]model.UserResponse, error) {
	vets, err := s.userRepo.ListVeterinarians(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	res := make([]model.UserResponse, len(vets))
	for i, v := range vets {
		res[i] = model.UserResponse{
			ID:         v.ID,
			Name:       v.Name,
			Email:      v.Email,
			Phone:      v.Phone,
			Address:    v.Address,
			Role:       v.RoleName,
			ClinicID:   v.ClinicID,
			ClinicName: v.ClinicName,
			IsActive:   v.IsActive,
			CreatedAt:  v.CreatedAt,
		}
	}
	return res, nil
}
