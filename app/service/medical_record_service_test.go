package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
)

// ── Tests for MedicalRecordService ───────────────────────────────────────────

func TestMedicalRecord_Create_Success(t *testing.T) {
	recRepo := newMockMedicalRecordRepo()
	apptRepo := newMockApptRepo()
	petRepo := newMockPetRepo()

	petRepo.Create(context.Background(), &model.Pet{
		ID:       "pet-rec",
		OwnerID:  "owner-1",
		IsActive: true,
	})

	apptRepo.Create(context.Background(), &model.Appointment{
		ID:             "appt-rec",
		PetID:          "pet-rec",
		VeterinarianID: "vet-1",
		Status:         model.StatusConfirmed,
	})

	svc := service.NewMedicalRecordService(recRepo, apptRepo, petRepo)

	rec, err := svc.CreateMedicalRecord(context.Background(), "vet-1", model.CreateMedicalRecordInput{
		AppointmentID: "appt-rec",
		Diagnosis:     "Infeksi telinga luar (Otitis Externa)",
		Treatment:     "Pembersihan telinga dan tetes antiseptik",
		Prescription:  "Otopain drops 2x sehari",
	})

	if err != nil {
		t.Fatalf("expected create record success: %v", err)
	}
	if rec.Diagnosis != "Infeksi telinga luar (Otitis Externa)" {
		t.Errorf("unexpected diagnosis: %s", rec.Diagnosis)
	}

	// Verify appointment is marked COMPLETED
	appt, _ := apptRepo.FindByID(context.Background(), "appt-rec")
	if appt.Status != model.StatusCompleted {
		t.Errorf("expected appointment status to transition to COMPLETED, got %s", appt.Status)
	}
}

func TestMedicalRecord_Create_ForbiddenForOtherVet(t *testing.T) {
	recRepo := newMockMedicalRecordRepo()
	apptRepo := newMockApptRepo()
	petRepo := newMockPetRepo()

	apptRepo.Create(context.Background(), &model.Appointment{
		ID:             "appt-vet1",
		PetID:          "pet-1",
		VeterinarianID: "vet-1",
		Status:         model.StatusConfirmed,
	})

	svc := service.NewMedicalRecordService(recRepo, apptRepo, petRepo)

	// Vet-2 attempts to create record for Vet-1's appointment
	_, err := svc.CreateMedicalRecord(context.Background(), "vet-2", model.CreateMedicalRecordInput{
		AppointmentID: "appt-vet1",
		Diagnosis:     "Flu",
	})
	if !errors.Is(err, service.ErrMedicalRecordForbidden) {
		t.Fatalf("expected ErrMedicalRecordForbidden, got: %v", err)
	}
}
