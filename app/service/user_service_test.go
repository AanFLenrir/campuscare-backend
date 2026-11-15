package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/service"
)

func TestCreateUser_Success(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	svc := service.NewUserService(userRepo, roleRepo)

	input := service.CreateUserInput{
		Name:          "Danu Admin",
		Email:         "admin.danu@campuscare.id",
		Password:      "adminsecret123",
		Phone:         "08111222333",
		StudentNumber: "434241111",
		RoleName:      "admin",
	}

	user, err := svc.CreateUser(context.Background(), input)
	if err != nil {
		t.Fatalf("expected successful user creation, got error: %v", err)
	}

	if user.Name != input.Name {
		t.Errorf("expected name %s, got %s", input.Name, user.Name)
	}
	if user.Email != input.Email {
		t.Errorf("expected email %s, got %s", input.Email, user.Email)
	}
	if user.Password != "" {
		t.Errorf("expected empty password in output, got %s", user.Password)
	}
}

func TestCreateUser_MissingFields(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	svc := service.NewUserService(userRepo, roleRepo)

	// Missing name
	_, err := svc.CreateUser(context.Background(), service.CreateUserInput{
		Email:    "test@campuscare.id",
		Password: "password",
		RoleName: "mahasiswa",
	})
	if err == nil {
		t.Errorf("expected error for missing name, got nil")
	}

	// Invalid role
	_, err = svc.CreateUser(context.Background(), service.CreateUserInput{
		Name:     "Test",
		Email:    "test@campuscare.id",
		Password: "password",
		RoleName: "unknown_role",
	})
	if err == nil {
		t.Errorf("expected error for invalid role, got nil")
	}
}

func TestGetUser_Success(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()

	userRepo.Create(context.Background(), &model.User{
		ID:        "usr-find-1",
		RoleID:    "role-mhs",
		Name:      "User Dicari",
		Email:     "dicari@campuscare.id",
		Password:  "hash123",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	svc := service.NewUserService(userRepo, roleRepo)

	user, err := svc.GetUser(context.Background(), "usr-find-1")
	if err != nil {
		t.Fatalf("expected user found, got error: %v", err)
	}
	if user.ID != "usr-find-1" {
		t.Errorf("expected user ID usr-find-1, got %s", user.ID)
	}
	if user.Password != "" {
		t.Errorf("expected cleared password, got %s", user.Password)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	svc := service.NewUserService(userRepo, roleRepo)

	_, err := svc.GetUser(context.Background(), "nonexistent-id")
	if err == nil {
		t.Errorf("expected error for non-existent user, got nil")
	}
}

func TestListUsers_Success(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()

	userRepo.Create(context.Background(), &model.User{ID: "u1", Name: "User 1", Email: "u1@test.com", Password: "p1"})
	userRepo.Create(context.Background(), &model.User{ID: "u2", Name: "User 2", Email: "u2@test.com", Password: "p2"})

	svc := service.NewUserService(userRepo, roleRepo)

	users, err := svc.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error listing users: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
	for _, u := range users {
		if u.Password != "" {
			t.Errorf("expected cleared password for user %s, got %s", u.ID, u.Password)
		}
	}
}

func TestUpdateUser_Success(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()

	userRepo.Create(context.Background(), &model.User{
		ID:        "u-update",
		RoleID:    "role-mhs",
		Name:      "Nama Lama",
		Email:     "lama@test.com",
		Password:  "oldhash",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	svc := service.NewUserService(userRepo, roleRepo)

	newName := "Nama Baru"
	updated, err := svc.UpdateUser(context.Background(), "u-update", service.UpdateUserInput{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("unexpected error updating user: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("expected updated name %s, got %s", newName, updated.Name)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()

	userRepo.Create(context.Background(), &model.User{
		ID:       "u-del",
		Email:    "del@test.com",
		IsActive: true,
	})

	svc := service.NewUserService(userRepo, roleRepo)

	err := svc.DeleteUser(context.Background(), "u-del")
	if err != nil {
		t.Fatalf("unexpected error deleting user: %v", err)
	}
}
