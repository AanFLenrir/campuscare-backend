package service

import (
	"context"

	"github.com/AanFLenrir/campuscare-backend/app/repository"
)

// ReportSummary holds statistics about complaints.
type ReportSummary struct {
	TotalComplaints  int            `json:"total_complaints"`
	ByStatus         map[string]int `json:"by_status"`
	TotalAssignments int            `json:"total_assignments"`
	TotalHistories   int            `json:"total_histories"`
}

// ReportService provides reporting / analytics logic.
type ReportService struct {
	complaintRepo        repository.ComplaintRepository
	assignmentRepo       repository.AssignmentRepository
	complaintHistoryRepo repository.ComplaintHistoryRepository
}

func NewReportService(
	complaintRepo repository.ComplaintRepository,
	assignmentRepo repository.AssignmentRepository,
	histRepo repository.ComplaintHistoryRepository,
) *ReportService {
	return &ReportService{
		complaintRepo:        complaintRepo,
		assignmentRepo:       assignmentRepo,
		complaintHistoryRepo: histRepo,
	}
}

// GetSummary returns aggregated complaint statistics.
func (s *ReportService) GetSummary(ctx context.Context) (*ReportSummary, error) {
	complaints, err := s.complaintRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	byStatus := map[string]int{}
	for _, c := range complaints {
		byStatus[c.Status]++
	}
	histories, err := s.complaintHistoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &ReportSummary{
		TotalComplaints: len(complaints),
		ByStatus:        byStatus,
		TotalHistories:  len(histories),
	}, nil
}
