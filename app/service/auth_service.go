package service

import (
	"context"
	"errors"
	"time"

	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/repository"
	"github.com/AanFLenrir/campuscare-backend/config"
	"github.com/AanFLenrir/campuscare-backend/helper"
	"github.com/google/uuid"
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	cfg      *config.Config
}

func NewAuthService(repo repository.UserRepository, roleRepo repository.RoleRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: repo, roleRepo: roleRepo, cfg: cfg}
}

// RegisterInput holds registration data.
type RegisterInput struct {
	Name          string
	Email         string
	Password      string
	Phone         string
	StudentNumber string
}

// RegisterResult is returned after successful registration.
type RegisterResult struct {
	ID    string
	Name  string
	Email string
	Role  string
}

// LoginResult is returned after successful login.
type LoginResult struct {
	AccessToken string
}

// Register creates a new MAHASISWA account.
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*RegisterResult, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	if _, err := s.userRepo.GetByEmail(ctx, input.Email); err == nil {
		return nil, errors.New("email already registered")
	}
	hash, err := helper.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	role, err := s.roleRepo.GetByName(ctx, "mahasiswa")
	if err != nil {
		return nil, errors.New("default role not found")
	}
	user := &model.User{
		ID:            uuid.New().String(),
		RoleID:        role.ID,
		Name:          input.Name,
		Email:         input.Email,
		Password:      hash,
		Phone:         input.Phone,
		StudentNumber: &input.StudentNumber,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return &RegisterResult{ID: user.ID, Name: user.Name, Email: user.Email, Role: role.Name}, nil
}

// Login validates credentials and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsActive {
		return nil, errors.New("account disabled")
	}
	if err := helper.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get role")
	}
	token, err := helper.GenerateToken(s.cfg, user.ID, role.Name)
	if err != nil {
		return nil, err
	}
	return &LoginResult{AccessToken: token}, nil
}

// GetProfile returns public profile for the given userID.
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}
