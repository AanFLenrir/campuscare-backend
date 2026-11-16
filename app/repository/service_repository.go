package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrServiceNotFound = errors.New("service not found")

type ServiceRepository interface {
	Create(ctx context.Context, service *model.Service) error
	Update(ctx context.Context, service *model.Service) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Service, error)
	ListByClinic(ctx context.Context, clinicID string, activeOnly bool) ([]model.Service, error)
	ListAll(ctx context.Context, activeOnly bool) ([]model.Service, error)
}

type serviceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, s *model.Service) error {
	query := `
		INSERT INTO services (clinic_id, name, description, duration_minutes, price, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		s.ClinicID, s.Name, s.Description, s.DurationMinutes, s.Price, s.IsActive,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *serviceRepository) Update(ctx context.Context, s *model.Service) error {
	query := `
		UPDATE services
		SET name = $1, description = $2, duration_minutes = $3, price = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`
	err := r.db.QueryRow(ctx, query,
		s.Name, s.Description, s.DurationMinutes, s.Price, s.IsActive, s.ID,
	).Scan(&s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrServiceNotFound
		}
		return err
	}
	return nil
}

func (r *serviceRepository) Delete(ctx context.Context, id string) error {
	// Soft delete
	query := `UPDATE services SET is_active = false, updated_at = NOW() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrServiceNotFound
	}
	return nil
}

func (r *serviceRepository) FindByID(ctx context.Context, id string) (*model.Service, error) {
	query := `
		SELECT s.id, s.clinic_id, s.name, COALESCE(s.description, ''), s.duration_minutes, s.price, s.is_active, s.created_at, s.updated_at,
		       c.name as clinic_name
		FROM services s
		JOIN clinics c ON s.clinic_id = c.id
		WHERE s.id = $1 LIMIT 1`
	var s model.Service
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.ClinicID, &s.Name, &s.Description, &s.DurationMinutes, &s.Price, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		&s.ClinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrServiceNotFound
		}
		return nil, fmt.Errorf("find service: %w", err)
	}
	return &s, nil
}

func (r *serviceRepository) ListByClinic(ctx context.Context, clinicID string, activeOnly bool) ([]model.Service, error) {
	query := `
		SELECT s.id, s.clinic_id, s.name, COALESCE(s.description, ''), s.duration_minutes, s.price, s.is_active, s.created_at, s.updated_at,
		       c.name as clinic_name
		FROM services s
		JOIN clinics c ON s.clinic_id = c.id
		WHERE s.clinic_id = $1`
	if activeOnly {
		query += ` AND s.is_active = true`
	}
	query += ` ORDER BY s.name ASC`

	rows, err := r.db.Query(ctx, query, clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(
			&s.ID, &s.ClinicID, &s.Name, &s.Description, &s.DurationMinutes, &s.Price, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
			&s.ClinicName,
		); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func (r *serviceRepository) ListAll(ctx context.Context, activeOnly bool) ([]model.Service, error) {
	query := `
		SELECT s.id, s.clinic_id, s.name, COALESCE(s.description, ''), s.duration_minutes, s.price, s.is_active, s.created_at, s.updated_at,
		       c.name as clinic_name
		FROM services s
		JOIN clinics c ON s.clinic_id = c.id`
	if activeOnly {
		query += ` WHERE s.is_active = true`
	}
	query += ` ORDER BY s.name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(
			&s.ID, &s.ClinicID, &s.Name, &s.Description, &s.DurationMinutes, &s.Price, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
			&s.ClinicName,
		); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}
