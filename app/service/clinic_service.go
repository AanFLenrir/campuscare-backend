package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
)

type ClinicService interface {
	CreateClinic(ctx context.Context, input model.CreateClinicInput) (*model.Clinic, error)
	UpdateClinic(ctx context.Context, id string, input model.UpdateClinicInput) (*model.Clinic, error)
	DeleteClinic(ctx context.Context, id string) error
	GetClinic(ctx context.Context, id string) (*model.Clinic, error)
	ListClinics(ctx context.Context, activeOnly bool) ([]model.Clinic, error)
}

type clinicService struct {
	clinicRepo repository.ClinicRepository
}

func NewClinicService(clinicRepo repository.ClinicRepository) ClinicService {
	return &clinicService{clinicRepo: clinicRepo}
}

func (s *clinicService) CreateClinic(ctx context.Context, input model.CreateClinicInput) (*model.Clinic, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)

	if input.Name == "" || input.Address == "" {
		return nil, fmt.Errorf("%w: clinic name and address are required", ErrInvalidInput)
	}

	clinic := &model.Clinic{
		Name:        input.Name,
		Address:     input.Address,
		Phone:       input.Phone,
		Email:       input.Email,
		Description: input.Description,
		IsActive:    true,
	}

	if err := s.clinicRepo.Create(ctx, clinic); err != nil {
		return nil, err
	}
	return clinic, nil
}

func (s *clinicService) UpdateClinic(ctx context.Context, id string, input model.UpdateClinicInput) (*model.Clinic, error) {
	existing, err := s.clinicRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		existing.Name = strings.TrimSpace(input.Name)
	}
	if input.Address != "" {
		existing.Address = strings.TrimSpace(input.Address)
	}
	if input.Phone != "" {
		existing.Phone = input.Phone
	}
	if input.Email != "" {
		existing.Email = input.Email
	}
	if input.Description != "" {
		existing.Description = input.Description
	}
	if input.IsActive != nil {
		existing.IsActive = *input.IsActive
	}

	if err := s.clinicRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *clinicService) DeleteClinic(ctx context.Context, id string) error {
	return s.clinicRepo.Delete(ctx, id)
}

func (s *clinicService) GetClinic(ctx context.Context, id string) (*model.Clinic, error) {
	return s.clinicRepo.FindByID(ctx, id)
}

func (s *clinicService) ListClinics(ctx context.Context, activeOnly bool) ([]model.Clinic, error) {
	return s.clinicRepo.List(ctx, activeOnly)
}
