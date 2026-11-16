package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/service"
)

func TestClinic_CreateClinic_Success(t *testing.T) {
	repo := newMockClinicRepo()
	svc := service.NewClinicService(repo)

	c, err := svc.CreateClinic(context.Background(), model.CreateClinicInput{
		Name:    "Klinik Hewan Sejahtera",
		Address: "Jl. Diponegoro No. 12",
		Phone:   "031-123456",
	})

	if err != nil {
		t.Fatalf("expected clinic create success, got: %v", err)
	}
	if c.Name != "Klinik Hewan Sejahtera" {
		t.Errorf("expected clinic name Klinik Hewan Sejahtera, got %s", c.Name)
	}
	if !c.IsActive {
		t.Errorf("expected clinic to be active")
	}
}

func TestClinic_CreateClinic_MissingFields(t *testing.T) {
	repo := newMockClinicRepo()
	svc := service.NewClinicService(repo)

	_, err := svc.CreateClinic(context.Background(), model.CreateClinicInput{
		Name: "", // missing name
	})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}
}

func TestClinic_DeleteClinic_Success(t *testing.T) {
	repo := newMockClinicRepo()
	repo.Create(context.Background(), &model.Clinic{
		ID:       "c-del",
		Name:     "Klinik Hapus",
		IsActive: true,
	})

	svc := service.NewClinicService(repo)
	if err := svc.DeleteClinic(context.Background(), "c-del"); err != nil {
		t.Fatalf("expected delete success, got: %v", err)
	}

	c, _ := repo.FindByID(context.Background(), "c-del")
	if c.IsActive {
		t.Errorf("expected clinic to be deactivated")
	}
}
