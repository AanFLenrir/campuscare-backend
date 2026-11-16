package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
)

// ── Tests for AppointmentService ─────────────────────────────────────────────

func TestAppointment_Create_Success(t *testing.T) {
	apptRepo := newMockApptRepo()
	petRepo := newMockPetRepo()
	userRepo := newMockUserRepo()
	clinicRepo := newMockClinicRepo()
	serviceRepo := newMockServiceRepo()

	clinicID := "clinic-1"
	clinicRepo.Create(context.Background(), &model.Clinic{ID: clinicID, Name: "Clinic A", IsActive: true})

	petRepo.Create(context.Background(), &model.Pet{
		ID:       "pet-1",
		OwnerID:  "owner-1",
		Name:     "Milo",
		Species:  model.SpeciesCat,
		IsActive: true,
	})

	userRepo.Create(context.Background(), &model.User{
		ID:        "vet-1",
		RoleName:  model.RoleVeterinarian,
		ClinicID:  &clinicID,
		Name:      "Dr. Andi",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	serviceRepo.Create(context.Background(), &model.Service{
		ID:              "svc-1",
		ClinicID:        clinicID,
		Name:            "Checkup",
		DurationMinutes: 30,
		Price:           50000,
		IsActive:        true,
	})

	svc := service.NewAppointmentService(apptRepo, petRepo, userRepo, clinicRepo, serviceRepo)

	futureTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	appt, err := svc.CreateAppointment(context.Background(), "owner-1", model.CreateAppointmentInput{
		PetID:           "pet-1",
		VeterinarianID:  "vet-1",
		ClinicID:        clinicID,
		ServiceID:       "svc-1",
		AppointmentDate: futureTime,
		Complaint:       "Demam dan lemas",
	})

	if err != nil {
		t.Fatalf("expected successful appointment creation, got: %v", err)
	}
	if appt.Status != model.StatusPending {
		t.Errorf("expected initial status PENDING, got %s", appt.Status)
	}
	if appt.OwnerID != "owner-1" {
		t.Errorf("expected owner-1, got %s", appt.OwnerID)
	}
}

func TestAppointment_Create_PetNotOwned(t *testing.T) {
	apptRepo := newMockApptRepo()
	petRepo := newMockPetRepo()
	userRepo := newMockUserRepo()
	clinicRepo := newMockClinicRepo()
	serviceRepo := newMockServiceRepo()

	petRepo.Create(context.Background(), &model.Pet{
		ID:       "pet-someone-else",
		OwnerID:  "owner-2",
		IsActive: true,
	})

	svc := service.NewAppointmentService(apptRepo, petRepo, userRepo, clinicRepo, serviceRepo)

	futureTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	_, err := svc.CreateAppointment(context.Background(), "owner-1", model.CreateAppointmentInput{
		PetID:           "pet-someone-else",
		VeterinarianID:  "vet-1",
		ClinicID:        "clinic-1",
		ServiceID:       "svc-1",
		AppointmentDate: futureTime,
	})

	if !errors.Is(err, service.ErrPetForbidden) {
		t.Fatalf("expected ErrPetForbidden, got: %v", err)
	}
}

func TestAppointment_StatusTransitions(t *testing.T) {
	apptRepo := newMockApptRepo()
	petRepo := newMockPetRepo()
	userRepo := newMockUserRepo()
	clinicRepo := newMockClinicRepo()
	serviceRepo := newMockServiceRepo()

	apptRepo.Create(context.Background(), &model.Appointment{
		ID:             "appt-flow",
		OwnerID:        "owner-1",
		VeterinarianID: "vet-1",
		Status:         model.StatusPending,
	})

	svc := service.NewAppointmentService(apptRepo, petRepo, userRepo, clinicRepo, serviceRepo)

	// 1. PENDING -> CONFIRMED (by Vet)
	updated, err := svc.UpdateStatus(context.Background(), "vet-1", model.RoleVeterinarian, "appt-flow", model.StatusConfirmed, nil)
	if err != nil {
		t.Fatalf("expected CONFIRMED success, got: %v", err)
	}
	if updated.Status != model.StatusConfirmed {
		t.Errorf("expected CONFIRMED, got %s", updated.Status)
	}

	// 2. CONFIRMED -> IN_PROGRESS (by Vet)
	updated, err = svc.UpdateStatus(context.Background(), "vet-1", model.RoleVeterinarian, "appt-flow", model.StatusInProgress, nil)
	if err != nil {
		t.Fatalf("expected IN_PROGRESS success, got: %v", err)
	}
	if updated.Status != model.StatusInProgress {
		t.Errorf("expected IN_PROGRESS, got %s", updated.Status)
	}

	// 3. IN_PROGRESS -> COMPLETED (by Vet)
	updated, err = svc.UpdateStatus(context.Background(), "vet-1", model.RoleVeterinarian, "appt-flow", model.StatusCompleted, nil)
	if err != nil {
		t.Fatalf("expected COMPLETED success, got: %v", err)
	}
	if updated.Status != model.StatusCompleted {
		t.Errorf("expected COMPLETED, got %s", updated.Status)
	}
}

func TestAppointment_Cancel_Success(t *testing.T) {
	apptRepo := newMockApptRepo()
	petRepo := newMockPetRepo()
	userRepo := newMockUserRepo()
	clinicRepo := newMockClinicRepo()
	serviceRepo := newMockServiceRepo()

	apptRepo.Create(context.Background(), &model.Appointment{
		ID:             "appt-cancel",
		OwnerID:        "owner-1",
		VeterinarianID: "vet-1",
		Status:         model.StatusPending,
	})

	svc := service.NewAppointmentService(apptRepo, petRepo, userRepo, clinicRepo, serviceRepo)

	res, err := svc.CancelAppointment(context.Background(), "owner-1", model.RolePetOwner, "appt-cancel", "Ada urusan mendadak")
	if err != nil {
		t.Fatalf("expected cancel success, got: %v", err)
	}
	if res.Status != model.StatusCancelled {
		t.Errorf("expected CANCELLED, got %s", res.Status)
	}
	if res.CancellationReason == nil || *res.CancellationReason != "Ada urusan mendadak" {
		t.Errorf("expected cancellation reason recorded")
	}
}
