package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
)

// ── Tests for PetService ─────────────────────────────────────────────────────

func TestPet_CreatePet_Success(t *testing.T) {
	repo := newMockPetRepo()
	svc := service.NewPetService(repo)

	weight := 4.5
	res, err := svc.CreatePet(context.Background(), "owner-1", model.CreatePetInput{
		Name:    "Milo",
		Species: model.SpeciesCat,
		Breed:   "Persian",
		Gender:  model.GenderMale,
		Weight:  &weight,
		Color:   "White",
	})

	if err != nil {
		t.Fatalf("expected pet create success, got: %v", err)
	}
	if res.Name != "Milo" {
		t.Errorf("expected pet name Milo, got %s", res.Name)
	}
	if res.Species != model.SpeciesCat {
		t.Errorf("expected species CAT, got %s", res.Species)
	}
	if res.OwnerID != "owner-1" {
		t.Errorf("expected owner-1, got %s", res.OwnerID)
	}
}

func TestPet_CreatePet_InvalidSpecies(t *testing.T) {
	repo := newMockPetRepo()
	svc := service.NewPetService(repo)

	_, err := svc.CreatePet(context.Background(), "owner-1", model.CreatePetInput{
		Name:    "Dino",
		Species: "DINOSAUR",
	})
	if !errors.Is(err, service.ErrInvalidSpecies) {
		t.Fatalf("expected ErrInvalidSpecies, got: %v", err)
	}
}

func TestPet_UpdatePet_ForbiddenForOtherOwner(t *testing.T) {
	repo := newMockPetRepo()
	repo.Create(context.Background(), &model.Pet{
		ID:       "pet-1",
		OwnerID:  "owner-1",
		Name:     "Milo",
		Species:  model.SpeciesCat,
		IsActive: true,
	})

	svc := service.NewPetService(repo)

	// Owner-2 tries to update Owner-1's pet
	_, err := svc.UpdatePet(context.Background(), "owner-2", "pet-1", model.UpdatePetInput{
		Name: "Milo Renamed",
	})
	if !errors.Is(err, service.ErrPetForbidden) {
		t.Fatalf("expected ErrPetForbidden, got: %v", err)
	}
}

func TestPet_DeletePet_Success(t *testing.T) {
	repo := newMockPetRepo()
	repo.Create(context.Background(), &model.Pet{
		ID:       "pet-del",
		OwnerID:  "owner-1",
		Name:     "Luna",
		Species:  model.SpeciesDog,
		IsActive: true,
	})

	svc := service.NewPetService(repo)

	err := svc.DeletePet(context.Background(), "owner-1", "pet-del")
	if err != nil {
		t.Fatalf("expected delete success, got: %v", err)
	}

	p, _ := repo.FindByID(context.Background(), "pet-del")
	if p.IsActive {
		t.Errorf("expected pet to be deactivated")
	}
}
