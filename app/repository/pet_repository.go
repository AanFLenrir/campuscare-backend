package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPetNotFound = errors.New("pet not found")

type PetRepository interface {
	Create(ctx context.Context, pet *model.Pet) error
	Update(ctx context.Context, pet *model.Pet) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Pet, error)
	ListByOwner(ctx context.Context, ownerID string) ([]model.Pet, error)
	ListAll(ctx context.Context) ([]model.Pet, error)
	Count(ctx context.Context) (int64, error)
	CountByOwner(ctx context.Context, ownerID string) (int64, error)
}

type petRepository struct {
	db *pgxpool.Pool
}

func NewPetRepository(db *pgxpool.Pool) PetRepository {
	return &petRepository{db: db}
}

func (r *petRepository) Create(ctx context.Context, pet *model.Pet) error {
	query := `
		INSERT INTO pets (owner_id, name, species, breed, gender, birth_date, weight, color, notes, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		pet.OwnerID, pet.Name, pet.Species, pet.Breed, pet.Gender, pet.BirthDate, pet.Weight, pet.Color, pet.Notes, pet.IsActive,
	).Scan(&pet.ID, &pet.CreatedAt, &pet.UpdatedAt)
}

func (r *petRepository) Update(ctx context.Context, pet *model.Pet) error {
	query := `
		UPDATE pets
		SET name = $1, species = $2, breed = $3, gender = $4, birth_date = $5, weight = $6, color = $7, notes = $8, updated_at = NOW()
		WHERE id = $9 AND is_active = true
		RETURNING updated_at`
	err := r.db.QueryRow(ctx, query,
		pet.Name, pet.Species, pet.Breed, pet.Gender, pet.BirthDate, pet.Weight, pet.Color, pet.Notes, pet.ID,
	).Scan(&pet.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPetNotFound
		}
		return err
	}
	return nil
}

func (r *petRepository) Delete(ctx context.Context, id string) error {
	// Soft delete / deactivate so medical records are preserved
	query := `UPDATE pets SET is_active = false, updated_at = NOW() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrPetNotFound
	}
	return nil
}

func (r *petRepository) FindByID(ctx context.Context, id string) (*model.Pet, error) {
	query := `
		SELECT p.id, p.owner_id, p.name, p.species, COALESCE(p.breed, ''), p.gender,
		       p.birth_date::TEXT, p.weight, COALESCE(p.color, ''), COALESCE(p.notes, ''),
		       p.is_active, p.created_at, p.updated_at,
		       u.name as owner_name, COALESCE(u.phone, '') as owner_phone, u.email as owner_email
		FROM pets p
		JOIN users u ON p.owner_id = u.id
		WHERE p.id = $1 LIMIT 1`
	var p model.Pet
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.OwnerID, &p.Name, &p.Species, &p.Breed, &p.Gender,
		&p.BirthDate, &p.Weight, &p.Color, &p.Notes,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		&p.OwnerName, &p.OwnerPhone, &p.OwnerEmail,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPetNotFound
		}
		return nil, fmt.Errorf("find pet by id: %w", err)
	}
	return &p, nil
}

func (r *petRepository) ListByOwner(ctx context.Context, ownerID string) ([]model.Pet, error) {
	query := `
		SELECT p.id, p.owner_id, p.name, p.species, COALESCE(p.breed, ''), p.gender,
		       p.birth_date::TEXT, p.weight, COALESCE(p.color, ''), COALESCE(p.notes, ''),
		       p.is_active, p.created_at, p.updated_at,
		       u.name as owner_name, COALESCE(u.phone, '') as owner_phone, u.email as owner_email
		FROM pets p
		JOIN users u ON p.owner_id = u.id
		WHERE p.owner_id = $1 AND p.is_active = true
		ORDER BY p.created_at DESC`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []model.Pet
	for rows.Next() {
		var p model.Pet
		if err := rows.Scan(
			&p.ID, &p.OwnerID, &p.Name, &p.Species, &p.Breed, &p.Gender,
			&p.BirthDate, &p.Weight, &p.Color, &p.Notes,
			&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
			&p.OwnerName, &p.OwnerPhone, &p.OwnerEmail,
		); err != nil {
			return nil, err
		}
		pets = append(pets, p)
	}
	return pets, nil
}

func (r *petRepository) ListAll(ctx context.Context) ([]model.Pet, error) {
	query := `
		SELECT p.id, p.owner_id, p.name, p.species, COALESCE(p.breed, ''), p.gender,
		       p.birth_date::TEXT, p.weight, COALESCE(p.color, ''), COALESCE(p.notes, ''),
		       p.is_active, p.created_at, p.updated_at,
		       u.name as owner_name, COALESCE(u.phone, '') as owner_phone, u.email as owner_email
		FROM pets p
		JOIN users u ON p.owner_id = u.id
		ORDER BY p.created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []model.Pet
	for rows.Next() {
		var p model.Pet
		if err := rows.Scan(
			&p.ID, &p.OwnerID, &p.Name, &p.Species, &p.Breed, &p.Gender,
			&p.BirthDate, &p.Weight, &p.Color, &p.Notes,
			&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
			&p.OwnerName, &p.OwnerPhone, &p.OwnerEmail,
		); err != nil {
			return nil, err
		}
		pets = append(pets, p)
	}
	return pets, nil
}

func (r *petRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pets WHERE is_active = true`).Scan(&count)
	return count, err
}

func (r *petRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pets WHERE owner_id = $1 AND is_active = true`, ownerID).Scan(&count)
	return count, err
}
