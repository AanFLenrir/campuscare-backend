package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")
var ErrDuplicateEmail = errors.New("email already registered")

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, roleFilter string) ([]model.User, error)
	ListVeterinarians(ctx context.Context, clinicID *string) ([]model.User, error)
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (role_id, clinic_id, name, email, password, phone, address, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query,
		user.RoleID, user.ClinicID, user.Name, user.Email, user.Password, user.Phone, user.Address, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if stringsContains(err.Error(), "unique") || stringsContains(err.Error(), "duplicate") {
			return ErrDuplicateEmail
		}
		return fmt.Errorf("create user error: %w", err)
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET name = $1, phone = $2, address = $3, clinic_id = $4, is_active = $5, updated_at = NOW()`
	args := []interface{}{user.Name, user.Phone, user.Address, user.ClinicID, user.IsActive}
	argIdx := 6

	if user.Password != "" {
		query += fmt.Sprintf(", password = $%d", argIdx)
		args = append(args, user.Password)
		argIdx++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING updated_at", argIdx)
	args = append(args, user.ID)

	err := r.db.QueryRow(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	// Soft delete / deactivate
	query := `UPDATE users SET is_active = false, updated_at = NOW() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT u.id, u.role_id, u.clinic_id, u.name, u.email, u.password,
		       COALESCE(u.phone, ''), COALESCE(u.address, ''), u.is_active, u.created_at, u.updated_at,
		       r.name as role_name, c.name as clinic_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		LEFT JOIN clinics c ON u.clinic_id = c.id
		WHERE u.id = $1 LIMIT 1`
	var u model.User
	var clinicName *string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.RoleID, &u.ClinicID, &u.Name, &u.Email, &u.Password,
		&u.Phone, &u.Address, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		&u.RoleName, &clinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	u.ClinicName = clinicName
	return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT u.id, u.role_id, u.clinic_id, u.name, u.email, u.password,
		       COALESCE(u.phone, ''), COALESCE(u.address, ''), u.is_active, u.created_at, u.updated_at,
		       r.name as role_name, c.name as clinic_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		LEFT JOIN clinics c ON u.clinic_id = c.id
		WHERE u.email = $1 LIMIT 1`
	var u model.User
	var clinicName *string
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.RoleID, &u.ClinicID, &u.Name, &u.Email, &u.Password,
		&u.Phone, &u.Address, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		&u.RoleName, &clinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	u.ClinicName = clinicName
	return &u, nil
}

func (r *userRepository) List(ctx context.Context, roleFilter string) ([]model.User, error) {
	query := `
		SELECT u.id, u.role_id, u.clinic_id, u.name, u.email, u.password,
		       COALESCE(u.phone, ''), COALESCE(u.address, ''), u.is_active, u.created_at, u.updated_at,
		       r.name as role_name, c.name as clinic_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		LEFT JOIN clinics c ON u.clinic_id = c.id`
	args := []interface{}{}
	if roleFilter != "" {
		query += ` WHERE r.name = $1`
		args = append(args, roleFilter)
	}
	query += ` ORDER BY u.created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var clinicName *string
		if err := rows.Scan(
			&u.ID, &u.RoleID, &u.ClinicID, &u.Name, &u.Email, &u.Password,
			&u.Phone, &u.Address, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
			&u.RoleName, &clinicName,
		); err != nil {
			return nil, err
		}
		u.ClinicName = clinicName
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) ListVeterinarians(ctx context.Context, clinicID *string) ([]model.User, error) {
	query := `
		SELECT u.id, u.role_id, u.clinic_id, u.name, u.email, u.password,
		       COALESCE(u.phone, ''), COALESCE(u.address, ''), u.is_active, u.created_at, u.updated_at,
		       r.name as role_name, c.name as clinic_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		LEFT JOIN clinics c ON u.clinic_id = c.id
		WHERE r.name = 'VETERINARIAN' AND u.is_active = true`
	args := []interface{}{}
	if clinicID != nil && *clinicID != "" {
		query += ` AND u.clinic_id = $1`
		args = append(args, *clinicID)
	}
	query += ` ORDER BY u.name ASC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vets []model.User
	for rows.Next() {
		var u model.User
		var clinicName *string
		if err := rows.Scan(
			&u.ID, &u.RoleID, &u.ClinicID, &u.Name, &u.Email, &u.Password,
			&u.Phone, &u.Address, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
			&u.RoleName, &clinicName,
		); err != nil {
			return nil, err
		}
		u.ClinicName = clinicName
		vets = append(vets, u)
	}
	return vets, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE is_active = true`).Scan(&count)
	return count, err
}

func stringsContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && stringSearch(s, substr)))
}

func stringSearch(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
