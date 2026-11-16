package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
)

type ServiceService interface {
	CreateService(ctx context.Context, input model.CreateServiceInput) (*model.Service, error)
	UpdateService(ctx context.Context, id string, input model.UpdateServiceInput) (*model.Service, error)
	DeleteService(ctx context.Context, id string) error
	GetService(ctx context.Context, id string) (*model.Service, error)
	ListServices(ctx context.Context, clinicID *string, activeOnly bool) ([]model.Service, error)
}

type serviceService struct {
	serviceRepo repository.ServiceRepository
	clinicRepo  repository.ClinicRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository, clinicRepo repository.ClinicRepository) ServiceService {
	return &serviceService{
		serviceRepo: serviceRepo,
		clinicRepo:  clinicRepo,
	}
}

func (s *serviceService) CreateService(ctx context.Context, input model.CreateServiceInput) (*model.Service, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" || input.ClinicID == "" {
		return nil, fmt.Errorf("%w: name and clinic_id are required", ErrInvalidInput)
	}

	// Validate clinic exists and is active
	clinic, err := s.clinicRepo.FindByID(ctx, input.ClinicID)
	if err != nil || !clinic.IsActive {
		return nil, fmt.Errorf("active clinic not found: %w", err)
	}

	if input.DurationMinutes <= 0 {
		input.DurationMinutes = 30
	}
	if input.Price < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}

	svc := &model.Service{
		ClinicID:        input.ClinicID,
		Name:            input.Name,
		Description:     input.Description,
		DurationMinutes: input.DurationMinutes,
		Price:           input.Price,
		IsActive:        true,
	}

	if err := s.serviceRepo.Create(ctx, svc); err != nil {
		return nil, err
	}
	svc.ClinicName = clinic.Name
	return svc, nil
}

func (s *serviceService) UpdateService(ctx context.Context, id string, input model.UpdateServiceInput) (*model.Service, error) {
	existing, err := s.serviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		existing.Name = strings.TrimSpace(input.Name)
	}
	if input.Description != "" {
		existing.Description = input.Description
	}
	if input.DurationMinutes != nil && *input.DurationMinutes > 0 {
		existing.DurationMinutes = *input.DurationMinutes
	}
	if input.Price != nil && *input.Price >= 0 {
		existing.Price = *input.Price
	}
	if input.IsActive != nil {
		existing.IsActive = *input.IsActive
	}

	if err := s.serviceRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *serviceService) DeleteService(ctx context.Context, id string) error {
	return s.serviceRepo.Delete(ctx, id)
}

func (s *serviceService) GetService(ctx context.Context, id string) (*model.Service, error) {
	return s.serviceRepo.FindByID(ctx, id)
}

func (s *serviceService) ListServices(ctx context.Context, clinicID *string, activeOnly bool) ([]model.Service, error) {
	if clinicID != nil && *clinicID != "" {
		return s.serviceRepo.ListByClinic(ctx, *clinicID, activeOnly)
	}
	return s.serviceRepo.ListAll(ctx, activeOnly)
}
