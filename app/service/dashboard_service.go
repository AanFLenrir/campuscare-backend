package service

import (
	"context"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
)

type DashboardService interface {
	GetDashboardSummary(ctx context.Context, userID, role string) (*model.DashboardSummary, error)
}

type dashboardService struct {
	userRepo   repository.UserRepository
	clinicRepo repository.ClinicRepository
	petRepo    repository.PetRepository
	apptRepo   repository.AppointmentRepository
	recordRepo repository.MedicalRecordRepository
	vacRepo    repository.VaccinationRepository
}

func NewDashboardService(
	userRepo repository.UserRepository,
	clinicRepo repository.ClinicRepository,
	petRepo repository.PetRepository,
	apptRepo repository.AppointmentRepository,
	recordRepo repository.MedicalRecordRepository,
	vacRepo repository.VaccinationRepository,
) DashboardService {
	return &dashboardService{
		userRepo:   userRepo,
		clinicRepo: clinicRepo,
		petRepo:    petRepo,
		apptRepo:   apptRepo,
		recordRepo: recordRepo,
		vacRepo:    vacRepo,
	}
}

func (s *dashboardService) GetDashboardSummary(ctx context.Context, userID, role string) (*model.DashboardSummary, error) {
	nowStr := time.Now().UTC().Format(time.RFC3339)

	switch role {
	case model.RoleAdmin:
		totalUsers, _ := s.userRepo.Count(ctx)
		totalClinics, _ := s.clinicRepo.Count(ctx)
		totalPets, _ := s.petRepo.Count(ctx)
		totalAppts, _ := s.apptRepo.Count(ctx)
		totalRecords, _ := s.recordRepo.Count(ctx)
		totalVacs, _ := s.vacRepo.Count(ctx)
		statusMap, _ := s.apptRepo.CountByStatus(ctx)

		return &model.DashboardSummary{
			Role:        role,
			GeneratedAt: nowStr,
			Data: model.AdminDashboardSummary{
				TotalUsers:           totalUsers,
				TotalClinics:         totalClinics,
				TotalPets:            totalPets,
				TotalAppointments:    totalAppts,
				TotalRecords:         totalRecords,
				TotalVaccinations:    totalVacs,
				AppointmentsByStatus: statusMap,
			},
		}, nil

	case model.RolePetOwner:
		totalPets, _ := s.petRepo.CountByOwner(ctx, userID)
		totalAppts, _ := s.apptRepo.CountByOwner(ctx, userID)
		upcomingAppts, _ := s.apptRepo.CountUpcomingByOwner(ctx, userID)
		totalRecords, _ := s.recordRepo.CountByPetOwner(ctx, userID)
		totalVacs, _ := s.vacRepo.CountByPetOwner(ctx, userID)

		return &model.DashboardSummary{
			Role:        role,
			GeneratedAt: nowStr,
			Data: model.PetOwnerDashboardSummary{
				TotalMyPets:          totalPets,
				TotalAppointments:    totalAppts,
				UpcomingAppointments: upcomingAppts,
				TotalMedicalRecords:  totalRecords,
				TotalVaccinations:    totalVacs,
			},
		}, nil

	case model.RoleVeterinarian:
		assigned, _ := s.apptRepo.CountAssignedToVet(ctx, userID)
		today, _ := s.apptRepo.CountTodayByVet(ctx, userID)
		completed, _ := s.apptRepo.CountCompletedByVet(ctx, userID)
		records, _ := s.recordRepo.CountByVet(ctx, userID)

		return &model.DashboardSummary{
			Role:        role,
			GeneratedAt: nowStr,
			Data: model.VeterinarianDashboardSummary{
				TotalAssignedAppointments: assigned,
				TodayAppointments:         today,
				CompletedAppointments:     completed,
				TotalRecordsCreated:       records,
			},
		}, nil

	default:
		return &model.DashboardSummary{
			Role:        role,
			GeneratedAt: nowStr,
			Data:        map[string]interface{}{},
		}, nil
	}
}
