package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
)

func TestVaccination_Create_Success(t *testing.T) {
	vacRepo := newMockVaccinationRepo()
	petRepo := newMockPetRepo()

	petRepo.Create(context.Background(), &model.Pet{
		ID:       "pet-vac",
		OwnerID:  "owner-1",
		IsActive: true,
	})

	svc := service.NewVaccinationService(vacRepo, petRepo)

	nextDate := "2027-11-16"
	batch := "BATCH-001"
	res, err := svc.CreateVaccination(context.Background(), "vet-1", model.CreateVaccinationInput{
		PetID:           "pet-vac",
		VaccineName:     "Rabisin",
		VaccineDate:     "2026-11-16",
		NextVaccineDate: &nextDate,
		BatchNumber:     &batch,
		Notes:           "Vaksin tahunan rabies",
	})

	if err != nil {
		t.Fatalf("expected vaccination creation success, got: %v", err)
	}
	if res.VaccineName != "Rabisin" {
		t.Errorf("expected vaccine name Rabisin, got %s", res.VaccineName)
	}
}

func TestVaccination_Create_InvalidPet(t *testing.T) {
	vacRepo := newMockVaccinationRepo()
	petRepo := newMockPetRepo()

	svc := service.NewVaccinationService(vacRepo, petRepo)

	_, err := svc.CreateVaccination(context.Background(), "vet-1", model.CreateVaccinationInput{
		PetID:       "pet-nonexistent",
		VaccineName: "Rabisin",
		VaccineDate: "2026-11-16",
	})
	if err == nil {
		t.Fatalf("expected error for nonexistent pet, got nil")
	}
}

func TestVaccination_Get_ForbiddenForOtherOwner(t *testing.T) {
	vacRepo := newMockVaccinationRepo()
	petRepo := newMockPetRepo()

	petRepo.Create(context.Background(), &model.Pet{
		ID:       "pet-owned-by-1",
		OwnerID:  "owner-1",
		IsActive: true,
	})

	vacRepo.Create(context.Background(), &model.Vaccination{
		ID:          "vac-1",
		PetID:       "pet-owned-by-1",
		VaccineName: "Parvo Vaccine",
	})

	svc := service.NewVaccinationService(vacRepo, petRepo)

	// Owner-2 tries to get vaccination of Owner-1's pet
	_, err := svc.GetVaccination(context.Background(), "owner-2", model.RolePetOwner, "vac-1")
	if !errors.Is(err, service.ErrVaccinationForbidden) {
		t.Fatalf("expected ErrVaccinationForbidden, got: %v", err)
	}
}
