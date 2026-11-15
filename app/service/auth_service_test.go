package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/service"
	"github.com/AanFLenrir/campuscare-backend/config"
	"github.com/AanFLenrir/campuscare-backend/helper"
)

type mockAuthUserRepo struct {
	usersByID    map[string]*model.User
	usersByEmail map[string]*model.User
}

func newMockAuthUserRepo() *mockAuthUserRepo {
	return &mockAuthUserRepo{
		usersByID:    make(map[string]*model.User),
		usersByEmail: make(map[string]*model.User),
	}
}

func (m *mockAuthUserRepo) Create(_ context.Context, u *model.User) error {
	m.usersByID[u.ID] = u
	m.usersByEmail[u.Email] = u
	return nil
}

func (m *mockAuthUserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	u, ok := m.usersByEmail[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockAuthUserRepo) GetByID(_ context.Context, id string) (*model.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockAuthUserRepo) GetAll(_ context.Context) ([]*model.User, error) {
	var list []*model.User
	for _, u := range m.usersByID {
		list = append(list, u)
	}
	return list, nil
}

func (m *mockAuthUserRepo) Update(_ context.Context, u *model.User) error {
	m.usersByID[u.ID] = u
	m.usersByEmail[u.Email] = u
	return nil
}

func (m *mockAuthUserRepo) Delete(_ context.Context, id string) error {
	if u, ok := m.usersByID[id]; ok {
		delete(m.usersByEmail, u.Email)
		delete(m.usersByID, id)
	}
	return nil
}

type mockAuthRoleRepo struct {
	rolesByName map[string]*model.Role
	rolesByID   map[string]*model.Role
}

func newMockAuthRoleRepo() *mockAuthRoleRepo {
	r := &mockAuthRoleRepo{
		rolesByName: make(map[string]*model.Role),
		rolesByID:   make(map[string]*model.Role),
	}
	mhsRole := &model.Role{ID: "role-mhs", Name: "mahasiswa", Description: "Mahasiswa Role"}
	adminRole := &model.Role{ID: "role-admin", Name: "admin", Description: "Admin Role"}
	petugasRole := &model.Role{ID: "role-petugas", Name: "petugas", Description: "Petugas Role"}

	r.rolesByName["mahasiswa"] = mhsRole
	r.rolesByName["admin"] = adminRole
	r.rolesByName["petugas"] = petugasRole

	r.rolesByID["role-mhs"] = mhsRole
	r.rolesByID["role-admin"] = adminRole
	r.rolesByID["role-petugas"] = petugasRole
	return r
}

func (m *mockAuthRoleRepo) GetByName(_ context.Context, name string) (*model.Role, error) {
	r, ok := m.rolesByName[name]
	if !ok {
		return nil, errors.New("role not found")
	}
	return r, nil
}

func (m *mockAuthRoleRepo) GetByID(_ context.Context, id string) (*model.Role, error) {
	r, ok := m.rolesByID[id]
	if !ok {
		return nil, errors.New("role not found")
	}
	return r, nil
}

// ── Unit Tests for AuthService ───────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	cfg := &config.Config{JWTSecret: "testsecret", JWTExpireHour: 24}

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	input := service.RegisterInput{
		Name:          "Danu Adji Ferianto",
		Email:         "danu@campuscare.id",
		Password:      "securepassword123",
		Phone:         "081234567890",
		StudentNumber: "434241111",
	}

	result, err := svc.Register(context.Background(), input)
	if err != nil {
		t.Fatalf("expected successful registration, got error: %v", err)
	}

	if result.Email != input.Email {
		t.Errorf("expected email %s, got %s", input.Email, result.Email)
	}
	if result.Role != "mahasiswa" {
		t.Errorf("expected role mahasiswa, got %s", result.Role)
	}
}

func TestRegister_MissingFields(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	cfg := &config.Config{JWTSecret: "testsecret", JWTExpireHour: 24}

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	// Missing name
	_, err := svc.Register(context.Background(), service.RegisterInput{
		Email:    "test@test.com",
		Password: "password123",
	})
	if err == nil {
		t.Errorf("expected error for missing name, got nil")
	}

	// Missing email
	_, err = svc.Register(context.Background(), service.RegisterInput{
		Name:     "Test User",
		Password: "password123",
	})
	if err == nil {
		t.Errorf("expected error for missing email, got nil")
	}

	// Missing password
	_, err = svc.Register(context.Background(), service.RegisterInput{
		Name:  "Test User",
		Email: "test@test.com",
	})
	if err == nil {
		t.Errorf("expected error for missing password, got nil")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	cfg := &config.Config{JWTSecret: "testsecret", JWTExpireHour: 24}

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	input := service.RegisterInput{
		Name:     "User Pertama",
		Email:    "duplicate@campuscare.id",
		Password: "password123",
	}

	_, err := svc.Register(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error on first register: %v", err)
	}

	// Coba registrasi ulang dengan email sama
	_, err = svc.Register(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error for duplicate email, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	cfg := &config.Config{JWTSecret: "testsecret", JWTExpireHour: 24}

	hashedPassword, _ := helper.HashPassword("mypassword123")
	userRepo.Create(context.Background(), &model.User{
		ID:        "usr-001",
		RoleID:    "role-mhs",
		Name:      "Danu Adji",
		Email:     "danu.login@campuscare.id",
		Password:  hashedPassword,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	result, err := svc.Login(context.Background(), "danu.login@campuscare.id", "mypassword123")
	if err != nil {
		t.Fatalf("expected successful login, got error: %v", err)
	}

	if result.AccessToken == "" {
		t.Errorf("expected non-empty access token")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	cfg := &config.Config{JWTSecret: "testsecret", JWTExpireHour: 24}

	hashedPassword, _ := helper.HashPassword("mypassword123")
	userRepo.Create(context.Background(), &model.User{
		ID:        "usr-001",
		RoleID:    "role-mhs",
		Name:      "Danu Adji",
		Email:     "danu.login@campuscare.id",
		Password:  hashedPassword,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	// Wrong password
	_, err := svc.Login(context.Background(), "danu.login@campuscare.id", "wrongpassword")
	if err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}

	// Non-existent email
	_, err = svc.Login(context.Background(), "nonexistent@campuscare.id", "mypassword123")
	if err == nil {
		t.Errorf("expected error for non-existent email, got nil")
	}
}

func TestLogin_AccountDisabled(t *testing.T) {
	userRepo := newMockAuthUserRepo()
	roleRepo := newMockAuthRoleRepo()
	cfg := &config.Config{JWTSecret: "testsecret", JWTExpireHour: 24}

	hashedPassword, _ := helper.HashPassword("mypassword123")
	userRepo.Create(context.Background(), &model.User{
		ID:        "usr-002",
		RoleID:    "role-mhs",
		Name:      "Nonaktif User",
		Email:     "disabled@campuscare.id",
		Password:  hashedPassword,
		IsActive:  false, // Akun nonaktif
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	svc := service.NewAuthService(userRepo, roleRepo, cfg)

	_, err := svc.Login(context.Background(), "disabled@campuscare.id", "mypassword123")
	if err == nil {
		t.Errorf("expected error for disabled account, got nil")
	}
}
