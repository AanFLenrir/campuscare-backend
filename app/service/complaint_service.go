package service

import (
	"context"
	"errors"

	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/repository"
)

// ComplaintService provides business logic for complaints.
type ComplaintService struct {
	complaintRepo        repository.ComplaintRepository
	assignmentRepo       repository.AssignmentRepository
	complaintHistoryRepo repository.ComplaintHistoryRepository
	userRepo             repository.UserRepository
	roleRepo             repository.RoleRepository
}

func NewComplaintService(
	compRepo repository.ComplaintRepository,
	assignRepo repository.AssignmentRepository,
	histRepo repository.ComplaintHistoryRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) *ComplaintService {
	return &ComplaintService{
		complaintRepo:        compRepo,
		assignmentRepo:       assignRepo,
		complaintHistoryRepo: histRepo,
		userRepo:             userRepo,
		roleRepo:             roleRepo,
	}
}

// CreateComplaintInput holds validated input for creating a complaint.
type CreateComplaintInput struct {
	FacilityID  string
	CategoryID  string
	Title       string
	Description string
	Priority    string
}

// CreateComplaint creates a new complaint for a mahasiswa.
func (s *ComplaintService) CreateComplaint(ctx context.Context, userID string, input CreateComplaintInput) (*model.Complaint, error) {
	if input.Title == "" || input.FacilityID == "" || input.CategoryID == "" || input.Priority == "" {
		return nil, errors.New("missing required fields")
	}
	comp := &model.Complaint{
		UserID:      userID,
		FacilityID:  input.FacilityID,
		CategoryID:  input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      "OPEN",
	}
	if err := s.complaintRepo.Create(ctx, comp); err != nil {
		return nil, err
	}
	hist := &model.ComplaintHistory{
		ComplaintID: comp.ID,
		UserID:      userID,
		OldStatus:   "",
		NewStatus:   "OPEN",
		Notes:       "Created complaint",
	}
	_ = s.complaintHistoryRepo.Create(ctx, hist)
	return comp, nil
}

// ListMyComplaints returns complaints owned by userID.
func (s *ComplaintService) ListMyComplaints(ctx context.Context, userID string) ([]*model.Complaint, error) {
	all, err := s.complaintRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var list []*model.Complaint
	for _, cmt := range all {
		if cmt.UserID == userID {
			list = append(list, cmt)
		}
	}
	return list, nil
}

// GetMyComplaint returns a complaint if it belongs to userID.
func (s *ComplaintService) GetMyComplaint(ctx context.Context, userID, complaintID string) (*model.Complaint, error) {
	comp, err := s.complaintRepo.GetByID(ctx, complaintID)
	if err != nil {
		return nil, err
	}
	if comp.UserID != userID {
		return nil, errors.New("forbidden")
	}
	return comp, nil
}

// ListAllComplaints returns all complaints (admin).
func (s *ComplaintService) ListAllComplaints(ctx context.Context) ([]*model.Complaint, error) {
	return s.complaintRepo.GetAll(ctx)
}

// GetComplaint returns any complaint by ID (admin).
func (s *ComplaintService) GetComplaint(ctx context.Context, complaintID string) (*model.Complaint, error) {
	return s.complaintRepo.GetByID(ctx, complaintID)
}

// UpdateComplaintStatus changes a complaint status and records history.
func (s *ComplaintService) UpdateComplaintStatus(ctx context.Context, officerID, complaintID, newStatus string) (*model.Complaint, error) {
	comp, err := s.complaintRepo.GetByID(ctx, complaintID)
	if err != nil {
		return nil, err
	}
	oldStatus := comp.Status
	comp.Status = newStatus
	if err := s.complaintRepo.Update(ctx, comp); err != nil {
		return nil, err
	}
	hist := &model.ComplaintHistory{
		ComplaintID: comp.ID,
		UserID:      officerID,
		OldStatus:   oldStatus,
		NewStatus:   newStatus,
		Notes:       "Status updated",
	}
	_ = s.complaintHistoryRepo.Create(ctx, hist)
	return comp, nil
}

// AssignComplaint creates an assignment record for a complaint.
func (s *ComplaintService) AssignComplaint(ctx context.Context, adminID, complaintID, officerID, notes string) (*model.Assignment, error) {
	comp, err := s.complaintRepo.GetByID(ctx, complaintID)
	if err != nil {
		return nil, err
	}
	assign := &model.Assignment{
		ComplaintID: comp.ID,
		OfficerID:   officerID,
		AssignedBy:  adminID,
		Notes:       notes,
	}
	if err := s.assignmentRepo.Create(ctx, assign); err != nil {
		return nil, err
	}
	return assign, nil
}
