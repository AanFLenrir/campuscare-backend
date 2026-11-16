package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
)

var (
	ErrMedicalRecordForbidden = errors.New("you do not have permission to access or create this medical record")
	ErrApptNotReadyForRecord  = errors.New("appointment must be CONFIRMED or IN_PROGRESS before creating a medical record")
)

type MedicalRecordService interface {
	CreateMedicalRecord(ctx context.Context, vetID string, input model.CreateMedicalRecordInput) (*model.MedicalRecord, error)
	UpdateMedicalRecord(ctx context.Context, callerID, callerRole, recordID string, input model.UpdateMedicalRecordInput) (*model.MedicalRecord, error)
	GetMedicalRecord(ctx context.Context, callerID, callerRole, recordID string) (*model.MedicalRecord, error)
	ListMedicalRecords(ctx context.Context, callerID, callerRole string, petID *string) ([]model.MedicalRecord, error)
}

type medicalRecordService struct {
	recordRepo repository.MedicalRecordRepository
	apptRepo   repository.AppointmentRepository
	petRepo    repository.PetRepository
}

func NewMedicalRecordService(
	recordRepo repository.MedicalRecordRepository,
	apptRepo repository.AppointmentRepository,
	petRepo repository.PetRepository,
) MedicalRecordService {
	return &medicalRecordService{
		recordRepo: recordRepo,
		apptRepo:   apptRepo,
		petRepo:    petRepo,
	}
}

func (s *medicalRecordService) CreateMedicalRecord(ctx context.Context, vetID string, input model.CreateMedicalRecordInput) (*model.MedicalRecord, error) {
	input.Diagnosis = strings.TrimSpace(input.Diagnosis)
	if input.AppointmentID == "" || input.Diagnosis == "" {
		return nil, fmt.Errorf("%w: appointment_id and diagnosis are required", ErrInvalidInput)
	}

	// 1. Verify Appointment exists
	appt, err := s.apptRepo.FindByID(ctx, input.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("appointment error: %w", err)
	}

	// 2. Only the assigned veterinarian can create medical record
	if appt.VeterinarianID != vetID {
		return nil, ErrMedicalRecordForbidden
	}

	// 3. Status check
	if appt.Status != model.StatusConfirmed && appt.Status != model.StatusInProgress && appt.Status != model.StatusCompleted {
		return nil, ErrApptNotReadyForRecord
	}

	examDate := time.Now()
	if input.ExaminationDate != nil && *input.ExaminationDate != "" {
		if parsed, err := time.Parse(time.RFC3339, *input.ExaminationDate); err == nil {
			examDate = parsed
		}
	}

	rec := &model.MedicalRecord{
		AppointmentID:   input.AppointmentID,
		PetID:           appt.PetID,
		VeterinarianID:  vetID,
		Diagnosis:       input.Diagnosis,
		Symptoms:        input.Symptoms,
		Treatment:       input.Treatment,
		Prescription:    input.Prescription,
		Notes:           input.Notes,
		ExaminationDate: examDate,
	}

	if err := s.recordRepo.Create(ctx, rec); err != nil {
		return nil, err
	}

	// Automatically transition appointment to COMPLETED if not already
	if appt.Status != model.StatusCompleted {
		_, _ = s.apptRepo.UpdateStatus(ctx, appt.ID, model.StatusCompleted, nil)
	}

	return s.recordRepo.FindByID(ctx, rec.ID)
}

func (s *medicalRecordService) UpdateMedicalRecord(ctx context.Context, callerID, callerRole, recordID string, input model.UpdateMedicalRecordInput) (*model.MedicalRecord, error) {
	existing, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}

	// Only assigned Vet or Admin can update
	if callerRole != model.RoleAdmin && (callerRole != model.RoleVeterinarian || existing.VeterinarianID != callerID) {
		return nil, ErrMedicalRecordForbidden
	}

	if input.Diagnosis != "" {
		existing.Diagnosis = strings.TrimSpace(input.Diagnosis)
	}
	if input.Symptoms != "" {
		existing.Symptoms = input.Symptoms
	}
	if input.Treatment != "" {
		existing.Treatment = input.Treatment
	}
	if input.Prescription != "" {
		existing.Prescription = input.Prescription
	}
	if input.Notes != "" {
		existing.Notes = input.Notes
	}

	if err := s.recordRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return s.recordRepo.FindByID(ctx, recordID)
}

func (s *medicalRecordService) GetMedicalRecord(ctx context.Context, callerID, callerRole, recordID string) (*model.MedicalRecord, error) {
	rec, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}

	// If owner, check pet ownership
	if callerRole == model.RolePetOwner {
		pet, err := s.petRepo.FindByID(ctx, rec.PetID)
		if err != nil || pet.OwnerID != callerID {
			return nil, ErrMedicalRecordForbidden
		}
	}

	return rec, nil
}

func (s *medicalRecordService) ListMedicalRecords(ctx context.Context, callerID, callerRole string, petID *string) ([]model.MedicalRecord, error) {
	if callerRole == model.RolePetOwner {
		if petID != nil && *petID != "" {
			// Verify pet belongs to owner
			pet, err := s.petRepo.FindByID(ctx, *petID)
			if err != nil || pet.OwnerID != callerID {
				return nil, ErrMedicalRecordForbidden
			}
			return s.recordRepo.ListByPet(ctx, *petID)
		}
		// List all records across all pets of owner
		pets, err := s.petRepo.ListByOwner(ctx, callerID)
		if err != nil {
			return nil, err
		}
		var allRecords []model.MedicalRecord
		for _, p := range pets {
			recs, err := s.recordRepo.ListByPet(ctx, p.ID)
			if err == nil {
				allRecords = append(allRecords, recs...)
			}
		}
		return allRecords, nil
	}

	if callerRole == model.RoleVeterinarian {
		if petID != nil && *petID != "" {
			return s.recordRepo.ListByPet(ctx, *petID)
		}
		return s.recordRepo.ListByVet(ctx, callerID)
	}

	// Admin
	if petID != nil && *petID != "" {
		return s.recordRepo.ListByPet(ctx, *petID)
	}
	return s.recordRepo.ListAll(ctx)
}
