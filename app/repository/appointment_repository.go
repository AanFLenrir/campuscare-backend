package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAppointmentNotFound = errors.New("appointment not found")
var ErrAppointmentConflict = errors.New("veterinarian already has an appointment at this time")

type AppointmentRepository interface {
	Create(ctx context.Context, appt *model.Appointment) error
	UpdateStatus(ctx context.Context, id, status string, cancellationReason *string) (*model.Appointment, error)
	FindByID(ctx context.Context, id string) (*model.Appointment, error)
	ListByOwner(ctx context.Context, ownerID string) ([]model.Appointment, error)
	ListByVet(ctx context.Context, vetID string) ([]model.Appointment, error)
	ListAll(ctx context.Context, statusFilter string) ([]model.Appointment, error)
	CheckConflict(ctx context.Context, vetID string, apptTime time.Time, durationMinutes int, excludeApptID string) (bool, error)
	CountByStatus(ctx context.Context) (map[string]int64, error)
	Count(ctx context.Context) (int64, error)
	CountByOwner(ctx context.Context, ownerID string) (int64, error)
	CountUpcomingByOwner(ctx context.Context, ownerID string) (int64, error)
	CountAssignedToVet(ctx context.Context, vetID string) (int64, error)
	CountTodayByVet(ctx context.Context, vetID string) (int64, error)
	CountCompletedByVet(ctx context.Context, vetID string) (int64, error)
}

type appointmentRepository struct {
	db *pgxpool.Pool
}

func NewAppointmentRepository(db *pgxpool.Pool) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) Create(ctx context.Context, appt *model.Appointment) error {
	query := `
		INSERT INTO appointments (pet_id, owner_id, veterinarian_id, clinic_id, service_id, appointment_date, status, complaint, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		appt.PetID, appt.OwnerID, appt.VeterinarianID, appt.ClinicID, appt.ServiceID, appt.AppointmentDate, appt.Status, appt.Complaint,
	).Scan(&appt.ID, &appt.CreatedAt, &appt.UpdatedAt)
}

func (r *appointmentRepository) UpdateStatus(ctx context.Context, id, status string, cancellationReason *string) (*model.Appointment, error) {
	query := `
		UPDATE appointments
		SET status = $1, cancellation_reason = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, pet_id, owner_id, veterinarian_id, clinic_id, service_id, appointment_date, status, COALESCE(complaint, ''), cancellation_reason, created_at, updated_at`
	var a model.Appointment
	err := r.db.QueryRow(ctx, query, status, cancellationReason, id).Scan(
		&a.ID, &a.PetID, &a.OwnerID, &a.VeterinarianID, &a.ClinicID, &a.ServiceID, &a.AppointmentDate, &a.Status, &a.Complaint, &a.CancellationReason, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAppointmentNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *appointmentRepository) FindByID(ctx context.Context, id string) (*model.Appointment, error) {
	query := `
		SELECT a.id, a.pet_id, a.owner_id, a.veterinarian_id, a.clinic_id, a.service_id,
		       a.appointment_date, a.status, COALESCE(a.complaint, ''), a.cancellation_reason, a.created_at, a.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name,
		       s.name as service_name, s.price as service_price
		FROM appointments a
		JOIN pets p ON a.pet_id = p.id
		JOIN users u ON a.owner_id = u.id
		JOIN users v ON a.veterinarian_id = v.id
		JOIN clinics c ON a.clinic_id = c.id
		JOIN services s ON a.service_id = s.id
		WHERE a.id = $1 LIMIT 1`
	var a model.Appointment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.PetID, &a.OwnerID, &a.VeterinarianID, &a.ClinicID, &a.ServiceID,
		&a.AppointmentDate, &a.Status, &a.Complaint, &a.CancellationReason, &a.CreatedAt, &a.UpdatedAt,
		&a.PetName, &a.PetSpecies,
		&a.OwnerName,
		&a.VeterinarianName,
		&a.ClinicName,
		&a.ServiceName, &a.ServicePrice,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAppointmentNotFound
		}
		return nil, fmt.Errorf("find appointment: %w", err)
	}
	return &a, nil
}

func (r *appointmentRepository) ListByOwner(ctx context.Context, ownerID string) ([]model.Appointment, error) {
	query := `
		SELECT a.id, a.pet_id, a.owner_id, a.veterinarian_id, a.clinic_id, a.service_id,
		       a.appointment_date, a.status, COALESCE(a.complaint, ''), a.cancellation_reason, a.created_at, a.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name,
		       s.name as service_name, s.price as service_price
		FROM appointments a
		JOIN pets p ON a.pet_id = p.id
		JOIN users u ON a.owner_id = u.id
		JOIN users v ON a.veterinarian_id = v.id
		JOIN clinics c ON a.clinic_id = c.id
		JOIN services s ON a.service_id = s.id
		WHERE a.owner_id = $1
		ORDER BY a.appointment_date DESC`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appts []model.Appointment
	for rows.Next() {
		var a model.Appointment
		if err := rows.Scan(
			&a.ID, &a.PetID, &a.OwnerID, &a.VeterinarianID, &a.ClinicID, &a.ServiceID,
			&a.AppointmentDate, &a.Status, &a.Complaint, &a.CancellationReason, &a.CreatedAt, &a.UpdatedAt,
			&a.PetName, &a.PetSpecies,
			&a.OwnerName,
			&a.VeterinarianName,
			&a.ClinicName,
			&a.ServiceName, &a.ServicePrice,
		); err != nil {
			return nil, err
		}
		appts = append(appts, a)
	}
	return appts, nil
}

func (r *appointmentRepository) ListByVet(ctx context.Context, vetID string) ([]model.Appointment, error) {
	query := `
		SELECT a.id, a.pet_id, a.owner_id, a.veterinarian_id, a.clinic_id, a.service_id,
		       a.appointment_date, a.status, COALESCE(a.complaint, ''), a.cancellation_reason, a.created_at, a.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name,
		       s.name as service_name, s.price as service_price
		FROM appointments a
		JOIN pets p ON a.pet_id = p.id
		JOIN users u ON a.owner_id = u.id
		JOIN users v ON a.veterinarian_id = v.id
		JOIN clinics c ON a.clinic_id = c.id
		JOIN services s ON a.service_id = s.id
		WHERE a.veterinarian_id = $1
		ORDER BY a.appointment_date DESC`
	rows, err := r.db.Query(ctx, query, vetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appts []model.Appointment
	for rows.Next() {
		var a model.Appointment
		if err := rows.Scan(
			&a.ID, &a.PetID, &a.OwnerID, &a.VeterinarianID, &a.ClinicID, &a.ServiceID,
			&a.AppointmentDate, &a.Status, &a.Complaint, &a.CancellationReason, &a.CreatedAt, &a.UpdatedAt,
			&a.PetName, &a.PetSpecies,
			&a.OwnerName,
			&a.VeterinarianName,
			&a.ClinicName,
			&a.ServiceName, &a.ServicePrice,
		); err != nil {
			return nil, err
		}
		appts = append(appts, a)
	}
	return appts, nil
}

func (r *appointmentRepository) ListAll(ctx context.Context, statusFilter string) ([]model.Appointment, error) {
	query := `
		SELECT a.id, a.pet_id, a.owner_id, a.veterinarian_id, a.clinic_id, a.service_id,
		       a.appointment_date, a.status, COALESCE(a.complaint, ''), a.cancellation_reason, a.created_at, a.updated_at,
		       p.name as pet_name, p.species as pet_species,
		       u.name as owner_name,
		       v.name as vet_name,
		       c.name as clinic_name,
		       s.name as service_name, s.price as service_price
		FROM appointments a
		JOIN pets p ON a.pet_id = p.id
		JOIN users u ON a.owner_id = u.id
		JOIN users v ON a.veterinarian_id = v.id
		JOIN clinics c ON a.clinic_id = c.id
		JOIN services s ON a.service_id = s.id`
	args := []interface{}{}
	if statusFilter != "" {
		query += ` WHERE a.status = $1`
		args = append(args, statusFilter)
	}
	query += ` ORDER BY a.appointment_date DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appts []model.Appointment
	for rows.Next() {
		var a model.Appointment
		if err := rows.Scan(
			&a.ID, &a.PetID, &a.OwnerID, &a.VeterinarianID, &a.ClinicID, &a.ServiceID,
			&a.AppointmentDate, &a.Status, &a.Complaint, &a.CancellationReason, &a.CreatedAt, &a.UpdatedAt,
			&a.PetName, &a.PetSpecies,
			&a.OwnerName,
			&a.VeterinarianName,
			&a.ClinicName,
			&a.ServiceName, &a.ServicePrice,
		); err != nil {
			return nil, err
		}
		appts = append(appts, a)
	}
	return appts, nil
}

func (r *appointmentRepository) CheckConflict(ctx context.Context, vetID string, apptTime time.Time, durationMinutes int, excludeApptID string) (bool, error) {
	// Look for existing active appointment within the window
	startWindow := apptTime.Add(-time.Duration(durationMinutes) * time.Minute)
	endWindow := apptTime.Add(time.Duration(durationMinutes) * time.Minute)

	query := `
		SELECT COUNT(*) FROM appointments
		WHERE veterinarian_id = $1
		  AND status IN ('PENDING', 'CONFIRMED', 'IN_PROGRESS')
		  AND appointment_date > $2 AND appointment_date < $3`
	args := []interface{}{vetID, startWindow, endWindow}

	if excludeApptID != "" {
		query += ` AND id != $4`
		args = append(args, excludeApptID)
	}

	var count int
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *appointmentRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	query := `SELECT status, COUNT(*) FROM appointments GROUP BY status`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		res[status] = count
	}
	return res, nil
}

func (r *appointmentRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM appointments`).Scan(&count)
	return count, err
}

func (r *appointmentRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM appointments WHERE owner_id = $1`, ownerID).Scan(&count)
	return count, err
}

func (r *appointmentRepository) CountUpcomingByOwner(ctx context.Context, ownerID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM appointments WHERE owner_id = $1 AND status IN ('PENDING', 'CONFIRMED') AND appointment_date >= NOW()`, ownerID).Scan(&count)
	return count, err
}

func (r *appointmentRepository) CountAssignedToVet(ctx context.Context, vetID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM appointments WHERE veterinarian_id = $1`, vetID).Scan(&count)
	return count, err
}

func (r *appointmentRepository) CountTodayByVet(ctx context.Context, vetID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM appointments WHERE veterinarian_id = $1 AND DATE(appointment_date) = CURRENT_DATE AND status IN ('PENDING', 'CONFIRMED', 'IN_PROGRESS')`, vetID).Scan(&count)
	return count, err
}

func (r *appointmentRepository) CountCompletedByVet(ctx context.Context, vetID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM appointments WHERE veterinarian_id = $1 AND status = 'COMPLETED'`, vetID).Scan(&count)
	return count, err
}
