package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrMedicalRecordNotFound = errors.New("medical record not found")
var ErrDuplicateMedicalRecord = errors.New("medical record for this appointment already exists")

type MedicalRecordRepository interface {
	Create(ctx context.Context, record *model.MedicalRecord) error
	Update(ctx context.Context, record *model.MedicalRecord) error
	FindByID(ctx context.Context, id string) (*model.MedicalRecord, error)
	FindByAppointmentID(ctx context.Context, apptID string) (*model.MedicalRecord, error)
	ListByPet(ctx context.Context, petID string) ([]model.MedicalRecord, error)
	ListByVet(ctx context.Context, vetID string) ([]model.MedicalRecord, error)
	ListAll(ctx context.Context) ([]model.MedicalRecord, error)
	Count(ctx context.Context) (int64, error)
	CountByPetOwner(ctx context.Context, ownerID string) (int64, error)
	CountByVet(ctx context.Context, vetID string) (int64, error)
}

type medicalRecordRepository struct {
	db *pgxpool.Pool
}

func NewMedicalRecordRepository(db *pgxpool.Pool) MedicalRecordRepository {
	return &medicalRecordRepository{db: db}
}

func (r *medicalRecordRepository) Create(ctx context.Context, rec *model.MedicalRecord) error {
	query := `
		INSERT INTO medical_records (appointment_id, pet_id, veterinarian_id, diagnosis, symptoms, treatment, prescription, notes, examination_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query,
		rec.AppointmentID, rec.PetID, rec.VeterinarianID, rec.Diagnosis, rec.Symptoms, rec.Treatment, rec.Prescription, rec.Notes, rec.ExaminationDate,
	).Scan(&rec.ID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		if stringsContains(err.Error(), "unique") || stringsContains(err.Error(), "duplicate") {
			return ErrDuplicateMedicalRecord
		}
		return fmt.Errorf("create medical record: %w", err)
	}
	return nil
}

func (r *medicalRecordRepository) Update(ctx context.Context, rec *model.MedicalRecord) error {
	query := `
		UPDATE medical_records
		SET diagnosis = $1, symptoms = $2, treatment = $3, prescription = $4, notes = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`
	err := r.db.QueryRow(ctx, query,
		rec.Diagnosis, rec.Symptoms, rec.Treatment, rec.Prescription, rec.Notes, rec.ID,
	).Scan(&rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMedicalRecordNotFound
		}
		return err
	}
	return nil
}

func (r *medicalRecordRepository) FindByID(ctx context.Context, id string) (*model.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.appointment_id, mr.pet_id, mr.veterinarian_id, mr.diagnosis,
		       COALESCE(mr.symptoms, ''), COALESCE(mr.treatment, ''), COALESCE(mr.prescription, ''), COALESCE(mr.notes, ''),
		       mr.examination_date, mr.created_at, mr.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name
		FROM medical_records mr
		JOIN pets p ON mr.pet_id = p.id
		JOIN users u ON p.owner_id = u.id
		JOIN users v ON mr.veterinarian_id = v.id
		JOIN appointments a ON mr.appointment_id = a.id
		JOIN clinics c ON a.clinic_id = c.id
		WHERE mr.id = $1 LIMIT 1`
	var mr model.MedicalRecord
	err := r.db.QueryRow(ctx, query, id).Scan(
		&mr.ID, &mr.AppointmentID, &mr.PetID, &mr.VeterinarianID, &mr.Diagnosis,
		&mr.Symptoms, &mr.Treatment, &mr.Prescription, &mr.Notes,
		&mr.ExaminationDate, &mr.CreatedAt, &mr.UpdatedAt,
		&mr.PetName, &mr.PetSpecies,
		&mr.OwnerName,
		&mr.VeterinarianName,
		&mr.ClinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMedicalRecordNotFound
		}
		return nil, fmt.Errorf("find medical record: %w", err)
	}
	return &mr, nil
}

func (r *medicalRecordRepository) FindByAppointmentID(ctx context.Context, apptID string) (*model.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.appointment_id, mr.pet_id, mr.veterinarian_id, mr.diagnosis,
		       COALESCE(mr.symptoms, ''), COALESCE(mr.treatment, ''), COALESCE(mr.prescription, ''), COALESCE(mr.notes, ''),
		       mr.examination_date, mr.created_at, mr.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name
		FROM medical_records mr
		JOIN pets p ON mr.pet_id = p.id
		JOIN users u ON p.owner_id = u.id
		JOIN users v ON mr.veterinarian_id = v.id
		JOIN appointments a ON mr.appointment_id = a.id
		JOIN clinics c ON a.clinic_id = c.id
		WHERE mr.appointment_id = $1 LIMIT 1`
	var mr model.MedicalRecord
	err := r.db.QueryRow(ctx, query, apptID).Scan(
		&mr.ID, &mr.AppointmentID, &mr.PetID, &mr.VeterinarianID, &mr.Diagnosis,
		&mr.Symptoms, &mr.Treatment, &mr.Prescription, &mr.Notes,
		&mr.ExaminationDate, &mr.CreatedAt, &mr.UpdatedAt,
		&mr.PetName, &mr.PetSpecies,
		&mr.OwnerName,
		&mr.VeterinarianName,
		&mr.ClinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMedicalRecordNotFound
		}
		return nil, err
	}
	return &mr, nil
}

func (r *medicalRecordRepository) ListByPet(ctx context.Context, petID string) ([]model.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.appointment_id, mr.pet_id, mr.veterinarian_id, mr.diagnosis,
		       COALESCE(mr.symptoms, ''), COALESCE(mr.treatment, ''), COALESCE(mr.prescription, ''), COALESCE(mr.notes, ''),
		       mr.examination_date, mr.created_at, mr.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name
		FROM medical_records mr
		JOIN pets p ON mr.pet_id = p.id
		JOIN users u ON p.owner_id = u.id
		JOIN users v ON mr.veterinarian_id = v.id
		JOIN appointments a ON mr.appointment_id = a.id
		JOIN clinics c ON a.clinic_id = c.id
		WHERE mr.pet_id = $1
		ORDER BY mr.examination_date DESC`
	rows, err := r.db.Query(ctx, query, petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.MedicalRecord
	for rows.Next() {
		var mr model.MedicalRecord
		if err := rows.Scan(
			&mr.ID, &mr.AppointmentID, &mr.PetID, &mr.VeterinarianID, &mr.Diagnosis,
			&mr.Symptoms, &mr.Treatment, &mr.Prescription, &mr.Notes,
			&mr.ExaminationDate, &mr.CreatedAt, &mr.UpdatedAt,
			&mr.PetName, &mr.PetSpecies,
			&mr.OwnerName,
			&mr.VeterinarianName,
			&mr.ClinicName,
		); err != nil {
			return nil, err
		}
		list = append(list, mr)
	}
	return list, nil
}

func (r *medicalRecordRepository) ListByVet(ctx context.Context, vetID string) ([]model.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.appointment_id, mr.pet_id, mr.veterinarian_id, mr.diagnosis,
		       COALESCE(mr.symptoms, ''), COALESCE(mr.treatment, ''), COALESCE(mr.prescription, ''), COALESCE(mr.notes, ''),
		       mr.examination_date, mr.created_at, mr.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name
		FROM medical_records mr
		JOIN pets p ON mr.pet_id = p.id
		JOIN users u ON p.owner_id = u.id
		JOIN users v ON mr.veterinarian_id = v.id
		JOIN appointments a ON mr.appointment_id = a.id
		JOIN clinics c ON a.clinic_id = c.id
		WHERE mr.veterinarian_id = $1
		ORDER BY mr.examination_date DESC`
	rows, err := r.db.Query(ctx, query, vetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.MedicalRecord
	for rows.Next() {
		var mr model.MedicalRecord
		if err := rows.Scan(
			&mr.ID, &mr.AppointmentID, &mr.PetID, &mr.VeterinarianID, &mr.Diagnosis,
			&mr.Symptoms, &mr.Treatment, &mr.Prescription, &mr.Notes,
			&mr.ExaminationDate, &mr.CreatedAt, &mr.UpdatedAt,
			&mr.PetName, &mr.PetSpecies,
			&mr.OwnerName,
			&mr.VeterinarianName,
			&mr.ClinicName,
		); err != nil {
			return nil, err
		}
		list = append(list, mr)
	}
	return list, nil
}

func (r *medicalRecordRepository) ListAll(ctx context.Context) ([]model.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.appointment_id, mr.pet_id, mr.veterinarian_id, mr.diagnosis,
		       COALESCE(mr.symptoms, ''), COALESCE(mr.treatment, ''), COALESCE(mr.prescription, ''), COALESCE(mr.notes, ''),
		       mr.examination_date, mr.created_at, mr.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name
		FROM medical_records mr
		JOIN pets p ON mr.pet_id = p.id
		JOIN users u ON p.owner_id = u.id
		JOIN users v ON mr.veterinarian_id = v.id
		JOIN appointments a ON mr.appointment_id = a.id
		JOIN clinics c ON a.clinic_id = c.id
		ORDER BY mr.examination_date DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.MedicalRecord
	for rows.Next() {
		var mr model.MedicalRecord
		if err := rows.Scan(
			&mr.ID, &mr.AppointmentID, &mr.PetID, &mr.VeterinarianID, &mr.Diagnosis,
			&mr.Symptoms, &mr.Treatment, &mr.Prescription, &mr.Notes,
			&mr.ExaminationDate, &mr.CreatedAt, &mr.UpdatedAt,
			&mr.PetName, &mr.PetSpecies,
			&mr.OwnerName,
			&mr.VeterinarianName,
			&mr.ClinicName,
		); err != nil {
			return nil, err
		}
		list = append(list, mr)
	}
	return list, nil
}

func (r *medicalRecordRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM medical_records`).Scan(&count)
	return count, err
}

func (r *medicalRecordRepository) CountByPetOwner(ctx context.Context, ownerID string) (int64, error) {
	query := `SELECT COUNT(*) FROM medical_records mr JOIN pets p ON mr.pet_id = p.id WHERE p.owner_id = $1`
	var count int64
	err := r.db.QueryRow(ctx, query, ownerID).Scan(&count)
	return count, err
}

func (r *medicalRecordRepository) CountByVet(ctx context.Context, vetID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM medical_records WHERE veterinarian_id = $1`, vetID).Scan(&count)
	return count, err
}
