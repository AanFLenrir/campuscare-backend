package service

import (
	"context"

	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/repository"
)

// AssignmentService provides business logic for officer assignments.
type AssignmentService struct {
	repo repository.AssignmentRepository
}

func NewAssignmentService(repo repository.AssignmentRepository) *AssignmentService {
	return &AssignmentService{repo: repo}
}

// ListMyAssignments returns assignments for the given officerID.
func (s *AssignmentService) ListMyAssignments(ctx context.Context, officerID string) ([]*model.Assignment, error) {
	return s.repo.GetAllByOfficer(ctx, officerID)
}
