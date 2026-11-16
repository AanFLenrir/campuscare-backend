package service_test

import (
	"context"
	"testing"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
)

func TestService_CreateService_Success(t *testing.T) {
	serviceRepo := newMockServiceRepo()
	clinicRepo := newMockClinicRepo()

	clinicID := "clinic-test"
	clinicRepo.Create(context.Background(), &model.Clinic{
		ID:       clinicID,
		Name:     "Pet Clinic",
		IsActive: true,
	})

	svc := service.NewServiceService(serviceRepo, clinicRepo)

	res, err := svc.CreateService(context.Background(), model.CreateServiceInput{
		ClinicID:        clinicID,
		Name:            "Dental Scaling",
		DurationMinutes: 45,
		Price:           200000,
	})

	if err != nil {
		t.Fatalf("expected create service success, got: %v", err)
	}
	if res.Name != "Dental Scaling" {
		t.Errorf("expected name Dental Scaling, got %s", res.Name)
	}
	if res.Price != 200000 {
		t.Errorf("expected price 200000, got %f", res.Price)
	}
}

func TestService_CreateService_InactiveClinic(t *testing.T) {
	serviceRepo := newMockServiceRepo()
	clinicRepo := newMockClinicRepo()

	clinicRepo.Create(context.Background(), &model.Clinic{
		ID:       "clinic-inactive",
		Name:     "Inactive Clinic",
		IsActive: false,
	})

	svc := service.NewServiceService(serviceRepo, clinicRepo)

	_, err := svc.CreateService(context.Background(), model.CreateServiceInput{
		ClinicID: "clinic-inactive",
		Name:     "Vaksinasi",
		Price:    100000,
	})
	if err == nil {
		t.Fatalf("expected error for inactive clinic, got nil")
	}
}
