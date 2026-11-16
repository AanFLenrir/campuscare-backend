package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrVaccinationNotFound = errors.New("vaccination record not found")

type VaccinationRepository interface {
	Create(ctx context.Context, v *model.Vaccination) error
	Update(ctx context.Context, v *model.Vaccination) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Vaccination, error)
	ListByPet(ctx context.Context, petID string) ([]model.Vaccination, error)
	ListAll(ctx context.Context) ([]model.Vaccination, error)
	Count(ctx context.Context) (int64, error)
	CountByPetOwner(ctx context.Context, ownerID string) (int64, error)
}

type vaccinationRepository struct {
	db *pgxpool.Pool
}

func NewVaccinationRepository(db *pgxpool.Pool) VaccinationRepository {
	return &vaccinationRepository{db: db}
}

func (r *vaccinationRepository) Create(ctx context.Context, v *model.Vaccination) error {
	query := `
		INSERT INTO vaccinations (pet_id, veterinarian_id, vaccine_name, vaccine_date, next_vaccine_date, batch_number, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		v.PetID, v.VeterinarianID, v.VaccineName, v.VaccineDate, v.NextVaccineDate, v.BatchNumber, v.Notes,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
}

func (r *vaccinationRepository) Update(ctx context.Context, v *model.Vaccination) error {
	query := `
		UPDATE vaccinations
		SET vaccine_name = $1, vaccine_date = $2, next_vaccine_date = $3, batch_number = $4, notes = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`
	err := r.db.QueryRow(ctx, query,
		v.VaccineName, v.VaccineDate, v.NextVaccineDate, v.BatchNumber, v.Notes, v.ID,
	).Scan(&v.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrVaccinationNotFound
		}
		return err
	}
	return nil
}

func (r *vaccinationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM vaccinations WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrVaccinationNotFound
	}
	return nil
}

func (r *vaccinationRepository) FindByID(ctx context.Context, id string) (*model.Vaccination, error) {
	query := `
		SELECT vac.id, vac.pet_id, vac.veterinarian_id, vac.vaccine_name,
		       vac.vaccine_date::TEXT, vac.next_vaccine_date::TEXT,
		       vac.batch_number, COALESCE(vac.notes, ''), vac.created_at, vac.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name
		FROM vaccinations vac
		JOIN pets p ON vac.pet_id = p.id
		JOIN users u ON p.owner_id = u.id
		JOIN users v ON vac.veterinarian_id = v.id
		WHERE vac.id = $1 LIMIT 1`
	var v model.Vaccination
	err := r.db.QueryRow(ctx, query, id).Scan(
		&v.ID, &v.PetID, &v.VeterinarianID, &v.VaccineName,
		&v.VaccineDate, &v.NextVaccineDate,
		&v.BatchNumber, &v.Notes, &v.CreatedAt, &v.UpdatedAt,
		&v.PetName, &v.PetSpecies,
		&v.OwnerName,
		&v.VeterinarianName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVaccinationNotFound
		}
		return nil, fmt.Errorf("find vaccination: %w", err)
	}
	return &v, nil
}

func (r *vaccinationRepository) ListByPet(ctx context.Context, petID string) ([]model.Vaccination, error) {
	query := `
		SELECT vac.id, vac.pet_id, vac.veterinarian_id, vac.vaccine_name,
		       vac.vaccine_date::TEXT, vac.next_vaccine_date::TEXT,
		       vac.batch_number, COALESCE(vac.notes, ''), vac.created_at, vac.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name
		FROM vaccinations vac
		JOIN pets p ON vac.pet_id = p.id
		JOIN users u ON p.owner_id = u.id
		JOIN users v ON vac.veterinarian_id = v.id
		WHERE vac.pet_id = $1
		ORDER BY vac.vaccine_date DESC`
	rows, err := r.db.Query(ctx, query, petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Vaccination
	for rows.Next() {
		var v model.Vaccination
		if err := rows.Scan(
			&v.ID, &v.PetID, &v.VeterinarianID, &v.VaccineName,
			&v.VaccineDate, &v.NextVaccineDate,
			&v.BatchNumber, &v.Notes, &v.CreatedAt, &v.UpdatedAt,
			&v.PetName, &v.PetSpecies,
			&v.OwnerName,
			&v.VeterinarianName,
		); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

func (r *vaccinationRepository) ListAll(ctx context.Context) ([]model.Vaccination, error) {
	query := `
		SELECT vac.id, vac.pet_id, vac.veterinarian_id, vac.vaccine_name,
		       vac.vaccine_date::TEXT, vac.next_vaccine_date::TEXT,
		       vac.batch_number, COALESCE(vac.notes, ''), vac.created_at, vac.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name
		FROM vaccinations vac
		JOIN pets p ON vac.pet_id = p.id
		JOIN users u ON p.owner_id = u.id
		JOIN users v ON vac.veterinarian_id = v.id
		ORDER BY vac.vaccine_date DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Vaccination
	for rows.Next() {
		var v model.Vaccination
		if err := rows.Scan(
			&v.ID, &v.PetID, &v.VeterinarianID, &v.VaccineName,
			&v.VaccineDate, &v.NextVaccineDate,
			&v.BatchNumber, &v.Notes, &v.CreatedAt, &v.UpdatedAt,
			&v.PetName, &v.PetSpecies,
			&v.OwnerName,
			&v.VeterinarianName,
		); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

func (r *vaccinationRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM vaccinations`).Scan(&count)
	return count, err
}

func (r *vaccinationRepository) CountByPetOwner(ctx context.Context, ownerID string) (int64, error) {
	query := `SELECT COUNT(*) FROM vaccinations vac JOIN pets p ON vac.pet_id = p.id WHERE p.owner_id = $1`
	var count int64
	err := r.db.QueryRow(ctx, query, ownerID).Scan(&count)
	return count, err
}
