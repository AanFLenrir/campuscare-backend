package service

import (
	"context"
	"errors"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/repository"
)

// FacilityService provides business logic for Facility entity.
type FacilityService interface {
	Create(ctx context.Context, f *model.Facility) error
	GetAll(ctx context.Context) ([]*model.Facility, error)
	GetByID(ctx context.Context, id string) (*model.Facility, error)
	Update(ctx context.Context, f *model.Facility) error
	Delete(ctx context.Context, id string) error
}

type facilityService struct {
	repo repository.FacilityRepository
}

func NewFacilityService(repo repository.FacilityRepository) FacilityService {
	return &facilityService{repo: repo}
}

func (s *facilityService) Create(ctx context.Context, f *model.Facility) error {
	if f.Name == "" {
		return errors.New("facility name required")
	}
	return s.repo.Create(ctx, f)
}

func (s *facilityService) GetAll(ctx context.Context) ([]*model.Facility, error) {
	return s.repo.GetAll(ctx)
}

func (s *facilityService) GetByID(ctx context.Context, id string) (*model.Facility, error) {
	if id == "" {
		return nil, errors.New("facility id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *facilityService) Update(ctx context.Context, f *model.Facility) error {
	if f.ID == "" {
		return errors.New("facility id required for update")
	}
	return s.repo.Update(ctx, f)
}

func (s *facilityService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("facility id required for delete")
	}
	return s.repo.Delete(ctx, id)
}
