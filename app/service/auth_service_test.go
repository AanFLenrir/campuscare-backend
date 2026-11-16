package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
	"github.com/AanFLenrir/PetrivaBackend/config"
	"github.com/AanFLenrir/PetrivaBackend/helper"
)

// ── Tests for AuthService ────────────────────────────────────────────────────

func TestAuth_Register_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	cfg := &config.Config{JWTSecret: "supersecret_petriva", JWTExpireHour: 24}

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	res, err := svc.Register(context.Background(), model.RegisterInput{
		Name:     "Budi Santoso",
		Email:    "budi@example.com",
		Password: "Password123!",
		Phone:    "08123456789",
		Address:  "Jl. Kebon Jeruk No. 1",
	})

	if err != nil {
		t.Fatalf("expected successful register, got error: %v", err)
	}

	if res.User.Email != "budi@example.com" {
		t.Errorf("expected email budi@example.com, got %s", res.User.Email)
	}
	if res.User.Role != model.RolePetOwner {
		t.Errorf("expected role PET_OWNER, got %s", res.User.Role)
	}
	if res.Token == "" {
		t.Errorf("expected non-empty JWT token")
	}
}

func TestAuth_Register_InvalidInput(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	cfg := &config.Config{JWTSecret: "supersecret_petriva", JWTExpireHour: 24}

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	// Missing name
	_, err := svc.Register(context.Background(), model.RegisterInput{
		Email:    "budi@example.com",
		Password: "Password123!",
	})
	if err == nil {
		t.Errorf("expected error for missing name, got nil")
	}

	// Password too short (< 6 chars)
	_, err = svc.Register(context.Background(), model.RegisterInput{
		Name:     "Budi",
		Email:    "budi@example.com",
		Password: "123",
	})
	if err == nil {
		t.Errorf("expected error for short password, got nil")
	}
}

func TestAuth_Register_DuplicateEmail(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	cfg := &config.Config{JWTSecret: "supersecret_petriva", JWTExpireHour: 24}

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	input := model.RegisterInput{
		Name:     "Budi",
		Email:    "duplicate@example.com",
		Password: "Password123!",
	}

	_, err := svc.Register(context.Background(), input)
	if err != nil {
		t.Fatalf("first registration should succeed: %v", err)
	}

	_, err = svc.Register(context.Background(), input)
	if err == nil {
		t.Fatalf("expected duplicate email error, got nil")
	}
}

func TestAuth_Login_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	cfg := &config.Config{JWTSecret: "supersecret_petriva", JWTExpireHour: 24}

	hashed, _ := helper.HashPassword("MySecretPass123!")
	userRepo.Create(context.Background(), &model.User{
		ID:        "usr-1",
		RoleID:    "role-owner",
		RoleName:  model.RolePetOwner,
		Name:      "Dewi Lestari",
		Email:     "dewi@example.com",
		Password:  hashed,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	res, err := svc.Login(context.Background(), model.LoginInput{
		Email:    "dewi@example.com",
		Password: "MySecretPass123!",
	})

	if err != nil {
		t.Fatalf("expected successful login, got: %v", err)
	}
	if res.Token == "" {
		t.Errorf("expected JWT token, got empty")
	}
	if res.User.Name != "Dewi Lestari" {
		t.Errorf("expected name Dewi Lestari, got %s", res.User.Name)
	}
}

func TestAuth_Login_WrongPassword(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	cfg := &config.Config{JWTSecret: "supersecret_petriva", JWTExpireHour: 24}

	hashed, _ := helper.HashPassword("CorrectPassword")
	userRepo.Create(context.Background(), &model.User{
		ID:       "usr-1",
		RoleName: model.RolePetOwner,
		Email:    "dewi@example.com",
		Password: hashed,
		IsActive: true,
	})

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	_, err := svc.Login(context.Background(), model.LoginInput{
		Email:    "dewi@example.com",
		Password: "WrongPassword",
	})
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestAuth_Login_InactiveUser(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	cfg := &config.Config{JWTSecret: "supersecret_petriva", JWTExpireHour: 24}

	hashed, _ := helper.HashPassword("Password123!")
	userRepo.Create(context.Background(), &model.User{
		ID:       "usr-inactive",
		Email:    "inactive@example.com",
		Password: hashed,
		IsActive: false,
	})

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	_, err := svc.Login(context.Background(), model.LoginInput{
		Email:    "inactive@example.com",
		Password: "Password123!",
	})
	if !errors.Is(err, service.ErrUserInactive) {
		t.Fatalf("expected ErrUserInactive, got: %v", err)
	}
}
