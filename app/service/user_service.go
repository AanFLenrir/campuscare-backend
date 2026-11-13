package service

import (
	"context"
	"errors"

	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/repository"
	"github.com/AanFLenrir/campuscare-backend/helper"
)

// UserService provides admin CRUD for users.
type UserService struct {
	repo     repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(repo repository.UserRepository, roleRepo repository.RoleRepository) *UserService {
	return &UserService{repo: repo, roleRepo: roleRepo}
}

// UpdateUserInput holds optional fields for update.
type UpdateUserInput struct {
	Name          *string
	Email         *string
	Password      *string
	Phone         *string
	StudentNumber *string
	RoleName      *string
	IsActive      *bool
}

// CreateUserInput holds fields for creating a user.
type CreateUserInput struct {
	Name          string
	Email         string
	Password      string
	Phone         string
	StudentNumber string
	RoleName      string
}

// ListUsers returns all users (passwords cleared).
func (s *UserService) ListUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		u.Password = ""
	}
	return users, nil
}

// GetUser returns a user by ID (password cleared).
func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*model.User, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" || input.RoleName == "" {
		return nil, errors.New("missing required fields")
	}
	roleObj, err := s.roleRepo.GetByName(ctx, input.RoleName)
	if err != nil {
		return nil, errors.New("invalid role")
	}
	hash, err := helper.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		RoleID:        roleObj.ID,
		Name:          input.Name,
		Email:         input.Email,
		Password:      hash,
		Phone:         input.Phone,
		StudentNumber: &input.StudentNumber,
		IsActive:      true,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

// UpdateUser updates an existing user.
func (s *UserService) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*model.User, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Email != nil {
		existing.Email = *input.Email
	}
	if input.Password != nil {
		hash, err := helper.HashPassword(*input.Password)
		if err != nil {
			return nil, err
		}
		existing.Password = hash
	}
	if input.Phone != nil {
		existing.Phone = *input.Phone
	}
	if input.StudentNumber != nil {
		existing.StudentNumber = input.StudentNumber
	}
	if input.RoleName != nil {
		roleObj, err := s.roleRepo.GetByName(ctx, *input.RoleName)
		if err != nil {
			return nil, errors.New("invalid role")
		}
		existing.RoleID = roleObj.ID
	}
	if input.IsActive != nil {
		existing.IsActive = *input.IsActive
	}
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	existing.Password = ""
	return existing, nil
}

// DeleteUser soft-deletes a user.
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
