package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrClinicNotFound = errors.New("clinic not found")

type ClinicRepository interface {
	Create(ctx context.Context, clinic *model.Clinic) error
	Update(ctx context.Context, clinic *model.Clinic) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Clinic, error)
	List(ctx context.Context, activeOnly bool) ([]model.Clinic, error)
	Count(ctx context.Context) (int64, error)
}

type clinicRepository struct {
	db *pgxpool.Pool
}

func NewClinicRepository(db *pgxpool.Pool) ClinicRepository {
	return &clinicRepository{db: db}
}

func (r *clinicRepository) Create(ctx context.Context, clinic *model.Clinic) error {
	query := `
		INSERT INTO clinics (name, address, phone, email, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		clinic.Name, clinic.Address, clinic.Phone, clinic.Email, clinic.Description, clinic.IsActive,
	).Scan(&clinic.ID, &clinic.CreatedAt, &clinic.UpdatedAt)
}

func (r *clinicRepository) Update(ctx context.Context, clinic *model.Clinic) error {
	query := `
		UPDATE clinics
		SET name = $1, address = $2, phone = $3, email = $4, description = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at`
	err := r.db.QueryRow(ctx, query,
		clinic.Name, clinic.Address, clinic.Phone, clinic.Email, clinic.Description, clinic.IsActive, clinic.ID,
	).Scan(&clinic.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrClinicNotFound
		}
		return err
	}
	return nil
}

func (r *clinicRepository) Delete(ctx context.Context, id string) error {
	// Soft delete: set is_active = false
	query := `UPDATE clinics SET is_active = false, updated_at = NOW() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrClinicNotFound
	}
	return nil
}

func (r *clinicRepository) FindByID(ctx context.Context, id string) (*model.Clinic, error) {
	query := `
		SELECT id, name, address, COALESCE(phone, ''), COALESCE(email, ''), COALESCE(description, ''), is_active, created_at, updated_at
		FROM clinics WHERE id = $1 LIMIT 1`
	var c model.Clinic
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Address, &c.Phone, &c.Email, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClinicNotFound
		}
		return nil, fmt.Errorf("find clinic error: %w", err)
	}
	return &c, nil
}

func (r *clinicRepository) List(ctx context.Context, activeOnly bool) ([]model.Clinic, error) {
	query := `
		SELECT id, name, address, COALESCE(phone, ''), COALESCE(email, ''), COALESCE(description, ''), is_active, created_at, updated_at
		FROM clinics`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clinics []model.Clinic
	for rows.Next() {
		var c model.Clinic
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Phone, &c.Email, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		clinics = append(clinics, c)
	}
	return clinics, nil
}

func (r *clinicRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM clinics WHERE is_active = true`).Scan(&count)
	return count, err
}
