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
	ErrVaccinationForbidden = errors.New("you do not have permission to access or modify this vaccination record")
)

type VaccinationService interface {
	CreateVaccination(ctx context.Context, vetID string, input model.CreateVaccinationInput) (*model.Vaccination, error)
	UpdateVaccination(ctx context.Context, callerID, callerRole, id string, input model.UpdateVaccinationInput) (*model.Vaccination, error)
	DeleteVaccination(ctx context.Context, callerRole, id string) error
	GetVaccination(ctx context.Context, callerID, callerRole, id string) (*model.Vaccination, error)
	ListVaccinations(ctx context.Context, callerID, callerRole string, petID *string) ([]model.Vaccination, error)
}

type vaccinationService struct {
	vacRepo repository.VaccinationRepository
	petRepo repository.PetRepository
}

func NewVaccinationService(vacRepo repository.VaccinationRepository, petRepo repository.PetRepository) VaccinationService {
	return &vaccinationService{
		vacRepo: vacRepo,
		petRepo: petRepo,
	}
}

func (s *vaccinationService) CreateVaccination(ctx context.Context, vetID string, input model.CreateVaccinationInput) (*model.Vaccination, error) {
	input.VaccineName = strings.TrimSpace(input.VaccineName)
	if input.PetID == "" || input.VaccineName == "" || input.VaccineDate == "" {
		return nil, fmt.Errorf("%w: pet_id, vaccine_name, and vaccine_date are required", ErrInvalidInput)
	}

	pet, err := s.petRepo.FindByID(ctx, input.PetID)
	if err != nil || !pet.IsActive {
		return nil, errors.New("active pet not found")
	}

	v := &model.Vaccination{
		PetID:           input.PetID,
		VeterinarianID:  vetID,
		VaccineName:     input.VaccineName,
		VaccineDate:     input.VaccineDate,
		NextVaccineDate: input.NextVaccineDate,
		BatchNumber:     input.BatchNumber,
		Notes:           input.Notes,
	}

	if err := s.vacRepo.Create(ctx, v); err != nil {
		return nil, err
	}
	return s.vacRepo.FindByID(ctx, v.ID)
}

func (s *vaccinationService) UpdateVaccination(ctx context.Context, callerID, callerRole, id string, input model.UpdateVaccinationInput) (*model.Vaccination, error) {
	existing, err := s.vacRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if callerRole != model.RoleAdmin && (callerRole != model.RoleVeterinarian || existing.VeterinarianID != callerID) {
		return nil, ErrVaccinationForbidden
	}

	if input.VaccineName != "" {
		existing.VaccineName = strings.TrimSpace(input.VaccineName)
	}
	if input.VaccineDate != "" {
		existing.VaccineDate = input.VaccineDate
	}
	if input.NextVaccineDate != nil {
		existing.NextVaccineDate = input.NextVaccineDate
	}
	if input.BatchNumber != nil {
		existing.BatchNumber = input.BatchNumber
	}
	if input.Notes != "" {
		existing.Notes = input.Notes
	}

	if err := s.vacRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return s.vacRepo.FindByID(ctx, id)
}

func (s *vaccinationService) DeleteVaccination(ctx context.Context, callerRole, id string) error {
	if callerRole != model.RoleAdmin && callerRole != model.RoleVeterinarian {
		return ErrVaccinationForbidden
	}
	return s.vacRepo.Delete(ctx, id)
}

func (s *vaccinationService) GetVaccination(ctx context.Context, callerID, callerRole, id string) (*model.Vaccination, error) {
	v, err := s.vacRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if callerRole == model.RolePetOwner {
		pet, err := s.petRepo.FindByID(ctx, v.PetID)
		if err != nil || pet.OwnerID != callerID {
			return nil, ErrVaccinationForbidden
		}
	}

	return v, nil
}

func (s *vaccinationService) ListVaccinations(ctx context.Context, callerID, callerRole string, petID *string) ([]model.Vaccination, error) {
	if callerRole == model.RolePetOwner {
		if petID != nil && *petID != "" {
			pet, err := s.petRepo.FindByID(ctx, *petID)
			if err != nil || pet.OwnerID != callerID {
				return nil, ErrVaccinationForbidden
			}
			return s.vacRepo.ListByPet(ctx, *petID)
		}
		// All pets of owner
		pets, err := s.petRepo.ListByOwner(ctx, callerID)
		if err != nil {
			return nil, err
		}
		var allList []model.Vaccination
		for _, p := range pets {
			list, err := s.vacRepo.ListByPet(ctx, p.ID)
			if err == nil {
				allList = append(allList, list...)
			}
		}
		return allList, nil
	}

	if petID != nil && *petID != "" {
		return s.vacRepo.ListByPet(ctx, *petID)
	}
	return s.vacRepo.ListAll(ctx)
}
