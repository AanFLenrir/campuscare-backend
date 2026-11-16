package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
)

func TestUser_CreateUser_Admin_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	clinicRepo := newMockClinicRepo()
	svc := service.NewUserService(userRepo, roleRepo, clinicRepo)

	res, err := svc.CreateUser(context.Background(), model.CreateUserInput{
		Name:     "Admin Baru",
		Email:    "newadmin@petriva.com",
		Password: "Password123!",
		Phone:    "081234567890",
		RoleName: model.RoleAdmin,
	})

	if err != nil {
		t.Fatalf("expected create user success, got: %v", err)
	}
	if res.Email != "newadmin@petriva.com" {
		t.Errorf("expected email newadmin@petriva.com, got %s", res.Email)
	}
	if res.Role != model.RoleAdmin {
		t.Errorf("expected role ADMIN, got %s", res.Role)
	}
}

func TestUser_CreateVeterinarian_WithClinic(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	clinicRepo := newMockClinicRepo()

	clinicID := "clinic-central"
	clinicRepo.Create(context.Background(), &model.Clinic{
		ID:       clinicID,
		Name:     "Petriva Clinic",
		Address:  "Jl. Ahmad Yani",
		IsActive: true,
	})

	svc := service.NewUserService(userRepo, roleRepo, clinicRepo)

	// Veterinarian requires clinic_id
	_, err := svc.CreateUser(context.Background(), model.CreateUserInput{
		Name:     "Dr. Siti",
		Email:    "drsiti@petriva.com",
		Password: "Password123!",
		RoleName: model.RoleVeterinarian,
		ClinicID: nil, // missing
	})
	if !errors.Is(err, service.ErrClinicRequired) {
		t.Fatalf("expected ErrClinicRequired, got: %v", err)
	}

	// With valid clinic_id
	res, err := svc.CreateUser(context.Background(), model.CreateUserInput{
		Name:     "Dr. Siti",
		Email:    "drsiti@petriva.com",
		Password: "Password123!",
		RoleName: model.RoleVeterinarian,
		ClinicID: &clinicID,
	})
	if err != nil {
		t.Fatalf("expected veterinarian creation success: %v", err)
	}
	if res.Role != model.RoleVeterinarian {
		t.Errorf("expected role VETERINARIAN, got %s", res.Role)
	}
}

func TestUser_CreateUser_InvalidRole(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	clinicRepo := newMockClinicRepo()
	svc := service.NewUserService(userRepo, roleRepo, clinicRepo)

	_, err := svc.CreateUser(context.Background(), model.CreateUserInput{
		Name:     "Invalid Role User",
		Email:    "test@petriva.com",
		Password: "Password123!",
		RoleName: "SUPERUSER",
	})
	if !errors.Is(err, service.ErrRoleInvalid) {
		t.Fatalf("expected ErrRoleInvalid, got: %v", err)
	}
}

func TestUser_GetUser_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	clinicRepo := newMockClinicRepo()

	userRepo.Create(context.Background(), &model.User{
		ID:        "usr-find-1",
		RoleName:  model.RolePetOwner,
		Name:      "Rina",
		Email:     "rina@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	svc := service.NewUserService(userRepo, roleRepo, clinicRepo)

	u, err := svc.GetUser(context.Background(), "usr-find-1")
	if err != nil {
		t.Fatalf("expected user found, got: %v", err)
	}
	if u.ID != "usr-find-1" {
		t.Errorf("expected ID usr-find-1, got %s", u.ID)
	}
}

func TestUser_DeleteUser_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	clinicRepo := newMockClinicRepo()

	userRepo.Create(context.Background(), &model.User{
		ID:       "usr-del",
		Email:    "del@example.com",
		IsActive: true,
	})

	svc := service.NewUserService(userRepo, roleRepo, clinicRepo)

	if err := svc.DeleteUser(context.Background(), "usr-del"); err != nil {
		t.Fatalf("expected delete user success, got: %v", err)
	}

	u, _ := userRepo.FindByID(context.Background(), "usr-del")
	if u.IsActive {
		t.Errorf("expected user is_active to be false after delete")
	}
}
