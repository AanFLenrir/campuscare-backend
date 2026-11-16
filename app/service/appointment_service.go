package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid appointment status transition")
	ErrAppointmentForbidden    = errors.New("you do not have permission to access or modify this appointment")
	ErrPastAppointmentTime     = errors.New("appointment date must be in the future")
	ErrVetNotAvailable         = errors.New("veterinarian is not available at the requested time (schedule conflict)")
	ErrVetClinicMismatch       = errors.New("veterinarian does not practice at the selected clinic")
	ErrServiceClinicMismatch   = errors.New("service does not belong to the selected clinic")
	ErrCancellationReasonReq   = errors.New("cancellation reason is required")
)

type AppointmentService interface {
	CreateAppointment(ctx context.Context, ownerID string, input model.CreateAppointmentInput) (*model.Appointment, error)
	UpdateStatus(ctx context.Context, callerID, callerRole, apptID, newStatus string, reason *string) (*model.Appointment, error)
	CancelAppointment(ctx context.Context, callerID, callerRole, apptID, reason string) (*model.Appointment, error)
	GetAppointment(ctx context.Context, callerID, callerRole, apptID string) (*model.Appointment, error)
	ListAppointments(ctx context.Context, callerID, callerRole string, statusFilter *string) ([]model.Appointment, error)
}

type appointmentService struct {
	apptRepo    repository.AppointmentRepository
	petRepo     repository.PetRepository
	userRepo    repository.UserRepository
	clinicRepo  repository.ClinicRepository
	serviceRepo repository.ServiceRepository
}

func NewAppointmentService(
	apptRepo repository.AppointmentRepository,
	petRepo repository.PetRepository,
	userRepo repository.UserRepository,
	clinicRepo repository.ClinicRepository,
	serviceRepo repository.ServiceRepository,
) AppointmentService {
	return &appointmentService{
		apptRepo:    apptRepo,
		petRepo:     petRepo,
		userRepo:    userRepo,
		clinicRepo:  clinicRepo,
		serviceRepo: serviceRepo,
	}
}

func (s *appointmentService) CreateAppointment(ctx context.Context, ownerID string, input model.CreateAppointmentInput) (*model.Appointment, error) {
	if input.PetID == "" || input.VeterinarianID == "" || input.ClinicID == "" || input.ServiceID == "" || input.AppointmentDate == "" {
		return nil, fmt.Errorf("%w: pet_id, veterinarian_id, clinic_id, service_id, and appointment_date are required", ErrInvalidInput)
	}

	// 1. Verify Pet belongs to the authenticated owner
	pet, err := s.petRepo.FindByID(ctx, input.PetID)
	if err != nil {
		return nil, fmt.Errorf("pet error: %w", err)
	}
	if pet.OwnerID != ownerID {
		return nil, ErrPetForbidden
	}
	if !pet.IsActive {
		return nil, errors.New("cannot create appointment for inactive pet")
	}

	// 2. Parse and validate appointment date
	apptTime, err := time.Parse(time.RFC3339, input.AppointmentDate)
	if err != nil {
		// Try fallback format "2006-01-02 15:04:05" or "2006-01-02T15:04"
		apptTime, err = time.Parse("2006-01-02 15:04:05", input.AppointmentDate)
		if err != nil {
			apptTime, err = time.Parse("2006-01-02T15:04", input.AppointmentDate)
			if err != nil {
				return nil, fmt.Errorf("%w: invalid appointment_date format, expected RFC3339", ErrInvalidInput)
			}
		}
	}
	if apptTime.Before(time.Now()) {
		return nil, ErrPastAppointmentTime
	}

	// 3. Verify Clinic is active
	clinic, err := s.clinicRepo.FindByID(ctx, input.ClinicID)
	if err != nil || !clinic.IsActive {
		return nil, errors.New("selected clinic is not active or found")
	}

	// 4. Verify Veterinarian exists, active, and practices at clinic
	vet, err := s.userRepo.FindByID(ctx, input.VeterinarianID)
	if err != nil || !vet.IsActive || vet.RoleName != model.RoleVeterinarian {
		return nil, errors.New("veterinarian not found or inactive")
	}
	if vet.ClinicID == nil || *vet.ClinicID != input.ClinicID {
		return nil, ErrVetClinicMismatch
	}

	// 5. Verify Service belongs to clinic
	svc, err := s.serviceRepo.FindByID(ctx, input.ServiceID)
	if err != nil || !svc.IsActive {
		return nil, errors.New("service not found or inactive")
	}
	if svc.ClinicID != input.ClinicID {
		return nil, ErrServiceClinicMismatch
	}

	// 6. Check for conflict in veterinarian schedule
	conflict, err := s.apptRepo.CheckConflict(ctx, input.VeterinarianID, apptTime, svc.DurationMinutes, "")
	if err != nil {
		return nil, fmt.Errorf("conflict check error: %w", err)
	}
	if conflict {
		return nil, ErrVetNotAvailable
	}

	appt := &model.Appointment{
		PetID:           input.PetID,
		OwnerID:         ownerID,
		VeterinarianID:  input.VeterinarianID,
		ClinicID:        input.ClinicID,
		ServiceID:       input.ServiceID,
		AppointmentDate: apptTime,
		Status:          model.StatusPending,
		Complaint:       input.Complaint,
	}

	if err := s.apptRepo.Create(ctx, appt); err != nil {
		return nil, err
	}
	return s.apptRepo.FindByID(ctx, appt.ID)
}

func (s *appointmentService) UpdateStatus(ctx context.Context, callerID, callerRole, apptID, newStatus string, reason *string) (*model.Appointment, error) {
	appt, err := s.apptRepo.FindByID(ctx, apptID)
	if err != nil {
		return nil, err
	}

	// State machine validation
	switch newStatus {
	case model.StatusConfirmed:
		// Only Admin or assigned Vet can confirm
		if callerRole != model.RoleAdmin && (callerRole != model.RoleVeterinarian || appt.VeterinarianID != callerID) {
			return nil, ErrAppointmentForbidden
		}
		if appt.Status != model.StatusPending {
			return nil, fmt.Errorf("%w: only PENDING appointment can be CONFIRMED", ErrInvalidStatusTransition)
		}

	case model.StatusInProgress:
		// Only assigned Vet can start progress
		if callerRole != model.RoleAdmin && (callerRole != model.RoleVeterinarian || appt.VeterinarianID != callerID) {
			return nil, ErrAppointmentForbidden
		}
		if appt.Status != model.StatusConfirmed {
			return nil, fmt.Errorf("%w: only CONFIRMED appointment can be moved to IN_PROGRESS", ErrInvalidStatusTransition)
		}

	case model.StatusCompleted:
		// Completed by Vet or Admin
		if callerRole != model.RoleAdmin && (callerRole != model.RoleVeterinarian || appt.VeterinarianID != callerID) {
			return nil, ErrAppointmentForbidden
		}
		if appt.Status != model.StatusInProgress && appt.Status != model.StatusConfirmed {
			return nil, fmt.Errorf("%w: appointment must be CONFIRMED or IN_PROGRESS before COMPLETED", ErrInvalidStatusTransition)
		}

	case model.StatusRejected:
		// Admin or assigned Vet rejects
		if callerRole != model.RoleAdmin && (callerRole != model.RoleVeterinarian || appt.VeterinarianID != callerID) {
			return nil, ErrAppointmentForbidden
		}
		if appt.Status != model.StatusPending {
			return nil, fmt.Errorf("%w: only PENDING appointment can be REJECTED", ErrInvalidStatusTransition)
		}

	default:
		return nil, ErrInvalidStatusTransition
	}

	return s.apptRepo.UpdateStatus(ctx, apptID, newStatus, reason)
}

func (s *appointmentService) CancelAppointment(ctx context.Context, callerID, callerRole, apptID, reason string) (*model.Appointment, error) {
	if reason == "" {
		return nil, ErrCancellationReasonReq
	}

	appt, err := s.apptRepo.FindByID(ctx, apptID)
	if err != nil {
		return nil, err
	}

	// Owner can cancel own appointment, Admin can cancel any
	if callerRole == model.RolePetOwner && appt.OwnerID != callerID {
		return nil, ErrAppointmentForbidden
	}

	// Only PENDING or CONFIRMED can be cancelled
	if appt.Status != model.StatusPending && appt.Status != model.StatusConfirmed {
		return nil, fmt.Errorf("%w: only PENDING or CONFIRMED appointments can be cancelled (current: %s)", ErrInvalidStatusTransition, appt.Status)
	}

	return s.apptRepo.UpdateStatus(ctx, apptID, model.StatusCancelled, &reason)
}

func (s *appointmentService) GetAppointment(ctx context.Context, callerID, callerRole, apptID string) (*model.Appointment, error) {
	appt, err := s.apptRepo.FindByID(ctx, apptID)
	if err != nil {
		return nil, err
	}

	if callerRole == model.RolePetOwner && appt.OwnerID != callerID {
		return nil, ErrAppointmentForbidden
	}
	if callerRole == model.RoleVeterinarian && appt.VeterinarianID != callerID {
		return nil, ErrAppointmentForbidden
	}

	return appt, nil
}

func (s *appointmentService) ListAppointments(ctx context.Context, callerID, callerRole string, statusFilter *string) ([]model.Appointment, error) {
	if callerRole == model.RolePetOwner {
		return s.apptRepo.ListByOwner(ctx, callerID)
	}
	if callerRole == model.RoleVeterinarian {
		return s.apptRepo.ListByVet(ctx, callerID)
	}
	// Admin
	status := ""
	if statusFilter != nil {
		status = *statusFilter
	}
	return s.apptRepo.ListAll(ctx, status)
}
