package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
)

var (
	ErrInvalidSpecies = errors.New("invalid species: must be CAT, DOG, BIRD, RABBIT, or OTHER")
	ErrInvalidGender  = errors.New("invalid gender: must be MALE, FEMALE, or UNKNOWN")
	ErrPetForbidden   = errors.New("you do not have permission to access or modify this pet")
)

type PetService interface {
	CreatePet(ctx context.Context, ownerID string, input model.CreatePetInput) (*model.Pet, error)
	UpdatePet(ctx context.Context, ownerID, petID string, input model.UpdatePetInput) (*model.Pet, error)
	DeletePet(ctx context.Context, ownerID, petID string) error
	GetPet(ctx context.Context, callerID, callerRole, petID string) (*model.Pet, error)
	ListPets(ctx context.Context, callerID, callerRole string, ownerFilter *string) ([]model.Pet, error)
}

type petService struct {
	petRepo repository.PetRepository
}

func NewPetService(petRepo repository.PetRepository) PetService {
	return &petService{petRepo: petRepo}
}

func isValidSpecies(s string) bool {
	switch s {
	case model.SpeciesCat, model.SpeciesDog, model.SpeciesBird, model.SpeciesRabbit, model.SpeciesOther:
		return true
	default:
		return false
	}
}

func isValidGender(g string) bool {
	switch g {
	case model.GenderMale, model.GenderFemale, model.GenderUnknown:
		return true
	default:
		return false
	}
}

func (s *petService) CreatePet(ctx context.Context, ownerID string, input model.CreatePetInput) (*model.Pet, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Species = strings.TrimSpace(strings.ToUpper(input.Species))
	input.Gender = strings.TrimSpace(strings.ToUpper(input.Gender))

	if input.Name == "" || input.Species == "" {
		return nil, fmt.Errorf("%w: name and species are required", ErrInvalidInput)
	}

	if !isValidSpecies(input.Species) {
		return nil, ErrInvalidSpecies
	}

	if input.Gender == "" {
		input.Gender = model.GenderUnknown
	} else if !isValidGender(input.Gender) {
		return nil, ErrInvalidGender
	}

	if input.Weight != nil && *input.Weight < 0 {
		return nil, fmt.Errorf("%w: weight cannot be negative", ErrInvalidInput)
	}

	pet := &model.Pet{
		OwnerID:   ownerID,
		Name:      input.Name,
		Species:   input.Species,
		Breed:     input.Breed,
		Gender:    input.Gender,
		BirthDate: input.BirthDate,
		Weight:    input.Weight,
		Color:     input.Color,
		Notes:     input.Notes,
		IsActive:  true,
	}

	if err := s.petRepo.Create(ctx, pet); err != nil {
		return nil, err
	}
	return s.petRepo.FindByID(ctx, pet.ID)
}

func (s *petService) UpdatePet(ctx context.Context, ownerID, petID string, input model.UpdatePetInput) (*model.Pet, error) {
	existing, err := s.petRepo.FindByID(ctx, petID)
	if err != nil {
		return nil, err
	}

	if existing.OwnerID != ownerID {
		return nil, ErrPetForbidden
	}

	if input.Name != "" {
		existing.Name = strings.TrimSpace(input.Name)
	}
	if input.Species != "" {
		species := strings.TrimSpace(strings.ToUpper(input.Species))
		if !isValidSpecies(species) {
			return nil, ErrInvalidSpecies
		}
		existing.Species = species
	}
	if input.Gender != "" {
		gender := strings.TrimSpace(strings.ToUpper(input.Gender))
		if !isValidGender(gender) {
			return nil, ErrInvalidGender
		}
		existing.Gender = gender
	}
	if input.Breed != "" {
		existing.Breed = input.Breed
	}
	if input.BirthDate != nil {
		existing.BirthDate = input.BirthDate
	}
	if input.Weight != nil {
		if *input.Weight < 0 {
			return nil, fmt.Errorf("%w: weight cannot be negative", ErrInvalidInput)
		}
		existing.Weight = input.Weight
	}
	if input.Color != "" {
		existing.Color = input.Color
	}
	if input.Notes != "" {
		existing.Notes = input.Notes
	}

	if err := s.petRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return s.petRepo.FindByID(ctx, petID)
}

func (s *petService) DeletePet(ctx context.Context, ownerID, petID string) error {
	existing, err := s.petRepo.FindByID(ctx, petID)
	if err != nil {
		return err
	}
	if existing.OwnerID != ownerID {
		return ErrPetForbidden
	}
	return s.petRepo.Delete(ctx, petID)
}

func (s *petService) GetPet(ctx context.Context, callerID, callerRole, petID string) (*model.Pet, error) {
	pet, err := s.petRepo.FindByID(ctx, petID)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if callerRole == model.RolePetOwner && pet.OwnerID != callerID {
		return nil, ErrPetForbidden
	}

	return pet, nil
}

func (s *petService) ListPets(ctx context.Context, callerID, callerRole string, ownerFilter *string) ([]model.Pet, error) {
	if callerRole == model.RolePetOwner {
		// Owner can only see own pets
		return s.petRepo.ListByOwner(ctx, callerID)
	}

	// Admin / Vet
	if ownerFilter != nil && *ownerFilter != "" {
		return s.petRepo.ListByOwner(ctx, *ownerFilter)
	}
	return s.petRepo.ListAll(ctx)
}
