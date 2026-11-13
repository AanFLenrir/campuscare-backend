package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/service"
)

// ── Mock Repositories ──────────────────────────────────────────────────────────

type mockComplaintRepo struct {
	complaints map[string]*model.Complaint
}

func newMockComplaintRepo() *mockComplaintRepo {
	return &mockComplaintRepo{complaints: make(map[string]*model.Complaint)}
}

func (m *mockComplaintRepo) Create(_ context.Context, c *model.Complaint) error {
	c.ID = "comp-001"
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	m.complaints[c.ID] = c
	return nil
}

func (m *mockComplaintRepo) GetByID(_ context.Context, id string) (*model.Complaint, error) {
	c, ok := m.complaints[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (m *mockComplaintRepo) GetAll(_ context.Context) ([]*model.Complaint, error) {
	var list []*model.Complaint
	for _, c := range m.complaints {
		list = append(list, c)
	}
	return list, nil
}

func (m *mockComplaintRepo) Update(_ context.Context, c *model.Complaint) error {
	m.complaints[c.ID] = c
	return nil
}

func (m *mockComplaintRepo) Delete(_ context.Context, id string) error {
	delete(m.complaints, id)
	return nil
}

type mockHistoryRepo struct{}

func (m *mockHistoryRepo) Create(_ context.Context, _ *model.ComplaintHistory) error { return nil }
func (m *mockHistoryRepo) GetByComplaintID(_ context.Context, _ string) ([]*model.ComplaintHistory, error) {
	return nil, nil
}
func (m *mockHistoryRepo) GetAll(_ context.Context) ([]*model.ComplaintHistory, error) {
	return nil, nil
}

type mockAssignmentRepo struct {
	assignments map[string]*model.Assignment
}

func newMockAssignmentRepo() *mockAssignmentRepo {
	return &mockAssignmentRepo{assignments: make(map[string]*model.Assignment)}
}

func (m *mockAssignmentRepo) Create(_ context.Context, a *model.Assignment) error {
	a.ID = "assign-001"
	m.assignments[a.ID] = a
	return nil
}

func (m *mockAssignmentRepo) GetByID(_ context.Context, id string) (*model.Assignment, error) {
	a, ok := m.assignments[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return a, nil
}

func (m *mockAssignmentRepo) GetAllByOfficer(_ context.Context, officerID string) ([]*model.Assignment, error) {
	var list []*model.Assignment
	for _, a := range m.assignments {
		if a.OfficerID == officerID {
			list = append(list, a)
		}
	}
	return list, nil
}

func (m *mockAssignmentRepo) Update(_ context.Context, a *model.Assignment) error {
	m.assignments[a.ID] = a
	return nil
}

func (m *mockAssignmentRepo) Delete(_ context.Context, id string) error {
	delete(m.assignments, id)
	return nil
}

type mockRoleRepo struct{}
type mockUserRepo struct{}

func (m *mockRoleRepo) GetByName(_ context.Context, _ string) (*model.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) GetByID(_ context.Context, _ string) (*model.Role, error) {
	return nil, nil
}

func (m *mockUserRepo) Create(_ context.Context, _ *model.User) error { return nil }
func (m *mockUserRepo) GetByID(_ context.Context, _ string) (*model.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepo) GetByEmail(_ context.Context, _ string) (*model.User, error) {
	return nil, errors.New("not found")
}
func (m *mockUserRepo) GetAll(_ context.Context) ([]*model.User, error) { return nil, nil }
func (m *mockUserRepo) Update(_ context.Context, _ *model.User) error   { return nil }
func (m *mockUserRepo) Delete(_ context.Context, _ string) error        { return nil }

// ── Tests ──────────────────────────────────────────────────────────────────────

func TestCreateComplaint_Success(t *testing.T) {
	compRepo := newMockComplaintRepo()
	assignRepo := newMockAssignmentRepo()
	histRepo := &mockHistoryRepo{}
	svc := service.NewComplaintService(compRepo, assignRepo, histRepo, &mockUserRepo{}, &mockRoleRepo{})

	comp, err := svc.CreateComplaint(context.Background(), "user-001", service.CreateComplaintInput{
		FacilityID:  "fac-001",
		CategoryID:  "cat-001",
		Title:       "AC rusak",
		Description: "AC di ruang kelas tidak berfungsi",
		Priority:    "HIGH",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if comp == nil {
		t.Fatal("expected complaint, got nil")
	}
	if comp.Status != "OPEN" {
		t.Errorf("expected status OPEN, got %s", comp.Status)
	}
	if comp.UserID != "user-001" {
		t.Errorf("expected userID user-001, got %s", comp.UserID)
	}
}

func TestCreateComplaint_MissingFields(t *testing.T) {
	compRepo := newMockComplaintRepo()
	assignRepo := newMockAssignmentRepo()
	histRepo := &mockHistoryRepo{}
	svc := service.NewComplaintService(compRepo, assignRepo, histRepo, &mockUserRepo{}, &mockRoleRepo{})

	_, err := svc.CreateComplaint(context.Background(), "user-001", service.CreateComplaintInput{
		Title: "AC rusak", // missing facility, category, priority
	})

	if err == nil {
		t.Fatal("expected error for missing fields, got nil")
	}
}

func TestListMyComplaints(t *testing.T) {
	compRepo := newMockComplaintRepo()
	assignRepo := newMockAssignmentRepo()
	histRepo := &mockHistoryRepo{}
	svc := service.NewComplaintService(compRepo, assignRepo, histRepo, &mockUserRepo{}, &mockRoleRepo{})

	// Create complaint owned by user-001
	_, _ = svc.CreateComplaint(context.Background(), "user-001", service.CreateComplaintInput{
		FacilityID: "fac-001", CategoryID: "cat-001", Title: "Test", Priority: "LOW",
	})

	list, err := svc.ListMyComplaints(context.Background(), "user-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 complaint, got %d", len(list))
	}

	// Different user should get empty list
	list2, _ := svc.ListMyComplaints(context.Background(), "user-999")
	if len(list2) != 0 {
		t.Errorf("expected 0 complaints for other user, got %d", len(list2))
	}
}

func TestGetMyComplaint_Forbidden(t *testing.T) {
	compRepo := newMockComplaintRepo()
	assignRepo := newMockAssignmentRepo()
	histRepo := &mockHistoryRepo{}
	svc := service.NewComplaintService(compRepo, assignRepo, histRepo, &mockUserRepo{}, &mockRoleRepo{})

	// Create complaint owned by user-001
	comp, _ := svc.CreateComplaint(context.Background(), "user-001", service.CreateComplaintInput{
		FacilityID: "fac-001", CategoryID: "cat-001", Title: "Test", Priority: "LOW",
	})

	// Different user tries to access it
	_, err := svc.GetMyComplaint(context.Background(), "user-999", comp.ID)
	if err == nil || err.Error() != "forbidden" {
		t.Errorf("expected forbidden error, got: %v", err)
	}
}

func TestUpdateComplaintStatus(t *testing.T) {
	compRepo := newMockComplaintRepo()
	assignRepo := newMockAssignmentRepo()
	histRepo := &mockHistoryRepo{}
	svc := service.NewComplaintService(compRepo, assignRepo, histRepo, &mockUserRepo{}, &mockRoleRepo{})

	comp, _ := svc.CreateComplaint(context.Background(), "user-001", service.CreateComplaintInput{
		FacilityID: "fac-001", CategoryID: "cat-001", Title: "Test", Priority: "LOW",
	})

	updated, err := svc.UpdateComplaintStatus(context.Background(), "officer-001", comp.ID, "IN_PROGRESS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != "IN_PROGRESS" {
		t.Errorf("expected status IN_PROGRESS, got %s", updated.Status)
	}
}

func TestAssignComplaint(t *testing.T) {
	compRepo := newMockComplaintRepo()
	assignRepo := newMockAssignmentRepo()
	histRepo := &mockHistoryRepo{}
	svc := service.NewComplaintService(compRepo, assignRepo, histRepo, &mockUserRepo{}, &mockRoleRepo{})

	comp, _ := svc.CreateComplaint(context.Background(), "user-001", service.CreateComplaintInput{
		FacilityID: "fac-001", CategoryID: "cat-001", Title: "Test", Priority: "LOW",
	})

	assign, err := svc.AssignComplaint(context.Background(), "admin-001", comp.ID, "officer-001", "Silakan tangani")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if assign.OfficerID != "officer-001" {
		t.Errorf("expected officerID officer-001, got %s", assign.OfficerID)
	}
}
