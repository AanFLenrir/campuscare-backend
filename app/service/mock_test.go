package service_test

import (
	"context"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/AanFLenrir/PetrivaBackend/app/repository"
)

// ── Mock UserRepository ───────────────────────────────────────────────────────
type mockUserRepo struct {
	usersByID    map[string]*model.User
	usersByEmail map[string]*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		usersByID:    make(map[string]*model.User),
		usersByEmail: make(map[string]*model.User),
	}
}

func (m *mockUserRepo) Create(_ context.Context, u *model.User) error {
	if _, exists := m.usersByEmail[u.Email]; exists {
		return repository.ErrDuplicateEmail
	}
	if u.ID == "" {
		u.ID = "usr-" + u.Email
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	m.usersByID[u.ID] = u
	m.usersByEmail[u.Email] = u
	return nil
}

func (m *mockUserRepo) Update(_ context.Context, u *model.User) error {
	m.usersByID[u.ID] = u
	m.usersByEmail[u.Email] = u
	return nil
}

func (m *mockUserRepo) Delete(_ context.Context, id string) error {
	if u, ok := m.usersByID[id]; ok {
		u.IsActive = false
		return nil
	}
	return repository.ErrUserNotFound
}

func (m *mockUserRepo) FindByID(_ context.Context, id string) (*model.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	u, ok := m.usersByEmail[email]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) List(_ context.Context, roleFilter string) ([]model.User, error) {
	var list []model.User
	for _, u := range m.usersByID {
		if roleFilter == "" || u.RoleName == roleFilter {
			list = append(list, *u)
		}
	}
	return list, nil
}

func (m *mockUserRepo) ListVeterinarians(_ context.Context, clinicID *string) ([]model.User, error) {
	var list []model.User
	for _, u := range m.usersByID {
		if u.RoleName == model.RoleVeterinarian && u.IsActive {
			if clinicID == nil || (u.ClinicID != nil && *u.ClinicID == *clinicID) {
				list = append(list, *u)
			}
		}
	}
	return list, nil
}

func (m *mockUserRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.usersByID)), nil
}

// ── Mock RoleRepository ───────────────────────────────────────────────────────
type mockRoleRepo struct {
	rolesByName map[string]*model.Role
	rolesByID   map[string]*model.Role
}

func newMockRoleRepo() *mockRoleRepo {
	r := &mockRoleRepo{
		rolesByName: make(map[string]*model.Role),
		rolesByID:   make(map[string]*model.Role),
	}
	admin := &model.Role{ID: "role-admin", Name: model.RoleAdmin, Description: "Admin"}
	petOwner := &model.Role{ID: "role-owner", Name: model.RolePetOwner, Description: "Pet Owner"}
	vet := &model.Role{ID: "role-vet", Name: model.RoleVeterinarian, Description: "Veterinarian"}

	r.rolesByName[model.RoleAdmin] = admin
	r.rolesByName[model.RolePetOwner] = petOwner
	r.rolesByName[model.RoleVeterinarian] = vet

	r.rolesByID["role-admin"] = admin
	r.rolesByID["role-owner"] = petOwner
	r.rolesByID["role-vet"] = vet
	return r
}

func (m *mockRoleRepo) FindByName(_ context.Context, name string) (*model.Role, error) {
	r, ok := m.rolesByName[name]
	if !ok {
		return nil, repository.ErrRoleNotFound
	}
	return r, nil
}

func (m *mockRoleRepo) FindByID(_ context.Context, id string) (*model.Role, error) {
	r, ok := m.rolesByID[id]
	if !ok {
		return nil, repository.ErrRoleNotFound
	}
	return r, nil
}

func (m *mockRoleRepo) List(_ context.Context) ([]model.Role, error) {
	var list []model.Role
	for _, r := range m.rolesByID {
		list = append(list, *r)
	}
	return list, nil
}

// ── Mock ClinicRepository ─────────────────────────────────────────────────────
type mockClinicRepo struct {
	clinics map[string]*model.Clinic
}

func newMockClinicRepo() *mockClinicRepo {
	return &mockClinicRepo{clinics: make(map[string]*model.Clinic)}
}

func (m *mockClinicRepo) Create(_ context.Context, c *model.Clinic) error {
	if c.ID == "" {
		c.ID = "clinic-1"
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	m.clinics[c.ID] = c
	return nil
}

func (m *mockClinicRepo) Update(_ context.Context, c *model.Clinic) error {
	c.UpdatedAt = time.Now()
	m.clinics[c.ID] = c
	return nil
}

func (m *mockClinicRepo) Delete(_ context.Context, id string) error {
	if c, ok := m.clinics[id]; ok {
		c.IsActive = false
		return nil
	}
	return repository.ErrClinicNotFound
}

func (m *mockClinicRepo) FindByID(_ context.Context, id string) (*model.Clinic, error) {
	c, ok := m.clinics[id]
	if !ok {
		return nil, repository.ErrClinicNotFound
	}
	return c, nil
}

func (m *mockClinicRepo) List(_ context.Context, activeOnly bool) ([]model.Clinic, error) {
	var list []model.Clinic
	for _, c := range m.clinics {
		if !activeOnly || c.IsActive {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (m *mockClinicRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.clinics)), nil
}

// ── Mock PetRepository ────────────────────────────────────────────────────────
type mockPetRepo struct {
	pets map[string]*model.Pet
}

func newMockPetRepo() *mockPetRepo {
	return &mockPetRepo{pets: make(map[string]*model.Pet)}
}

func (m *mockPetRepo) Create(_ context.Context, p *model.Pet) error {
	if p.ID == "" {
		p.ID = "pet-" + p.Name
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	m.pets[p.ID] = p
	return nil
}

func (m *mockPetRepo) Update(_ context.Context, p *model.Pet) error {
	p.UpdatedAt = time.Now()
	m.pets[p.ID] = p
	return nil
}

func (m *mockPetRepo) Delete(_ context.Context, id string) error {
	if p, ok := m.pets[id]; ok {
		p.IsActive = false
		return nil
	}
	return repository.ErrPetNotFound
}

func (m *mockPetRepo) FindByID(_ context.Context, id string) (*model.Pet, error) {
	p, ok := m.pets[id]
	if !ok {
		return nil, repository.ErrPetNotFound
	}
	return p, nil
}

func (m *mockPetRepo) ListByOwner(_ context.Context, ownerID string) ([]model.Pet, error) {
	var list []model.Pet
	for _, p := range m.pets {
		if p.OwnerID == ownerID && p.IsActive {
			list = append(list, *p)
		}
	}
	return list, nil
}

func (m *mockPetRepo) ListAll(_ context.Context) ([]model.Pet, error) {
	var list []model.Pet
	for _, p := range m.pets {
		list = append(list, *p)
	}
	return list, nil
}

func (m *mockPetRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.pets)), nil
}

func (m *mockPetRepo) CountByOwner(_ context.Context, ownerID string) (int64, error) {
	var count int64
	for _, p := range m.pets {
		if p.OwnerID == ownerID && p.IsActive {
			count++
		}
	}
	return count, nil
}

// ── Mock ServiceRepository ────────────────────────────────────────────────────
type mockServiceRepo struct {
	services map[string]*model.Service
}

func newMockServiceRepo() *mockServiceRepo {
	return &mockServiceRepo{services: make(map[string]*model.Service)}
}

func (m *mockServiceRepo) Create(_ context.Context, s *model.Service) error {
	if s.ID == "" {
		s.ID = "svc-" + s.Name
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	m.services[s.ID] = s
	return nil
}

func (m *mockServiceRepo) Update(_ context.Context, s *model.Service) error {
	s.UpdatedAt = time.Now()
	m.services[s.ID] = s
	return nil
}

func (m *mockServiceRepo) Delete(_ context.Context, id string) error {
	if s, ok := m.services[id]; ok {
		s.IsActive = false
		return nil
	}
	return repository.ErrServiceNotFound
}

func (m *mockServiceRepo) FindByID(_ context.Context, id string) (*model.Service, error) {
	s, ok := m.services[id]
	if !ok {
		return nil, repository.ErrServiceNotFound
	}
	return s, nil
}

func (m *mockServiceRepo) ListByClinic(_ context.Context, clinicID string, activeOnly bool) ([]model.Service, error) {
	var list []model.Service
	for _, s := range m.services {
		if s.ClinicID == clinicID && (!activeOnly || s.IsActive) {
			list = append(list, *s)
		}
	}
	return list, nil
}

func (m *mockServiceRepo) ListAll(_ context.Context, activeOnly bool) ([]model.Service, error) {
	var list []model.Service
	for _, s := range m.services {
		if !activeOnly || s.IsActive {
			list = append(list, *s)
		}
	}
	return list, nil
}

// ── Mock AppointmentRepository ────────────────────────────────────────────────
type mockApptRepo struct {
	appts map[string]*model.Appointment
}

func newMockApptRepo() *mockApptRepo {
	return &mockApptRepo{appts: make(map[string]*model.Appointment)}
}

func (m *mockApptRepo) Create(_ context.Context, a *model.Appointment) error {
	if a.ID == "" {
		a.ID = "appt-1"
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	m.appts[a.ID] = a
	return nil
}

func (m *mockApptRepo) UpdateStatus(_ context.Context, id, status string, reason *string) (*model.Appointment, error) {
	a, ok := m.appts[id]
	if !ok {
		return nil, repository.ErrAppointmentNotFound
	}
	a.Status = status
	a.CancellationReason = reason
	a.UpdatedAt = time.Now()
	return a, nil
}

func (m *mockApptRepo) FindByID(_ context.Context, id string) (*model.Appointment, error) {
	a, ok := m.appts[id]
	if !ok {
		return nil, repository.ErrAppointmentNotFound
	}
	return a, nil
}

func (m *mockApptRepo) ListByOwner(_ context.Context, ownerID string) ([]model.Appointment, error) {
	var list []model.Appointment
	for _, a := range m.appts {
		if a.OwnerID == ownerID {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (m *mockApptRepo) ListByVet(_ context.Context, vetID string) ([]model.Appointment, error) {
	var list []model.Appointment
	for _, a := range m.appts {
		if a.VeterinarianID == vetID {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (m *mockApptRepo) ListAll(_ context.Context, statusFilter string) ([]model.Appointment, error) {
	var list []model.Appointment
	for _, a := range m.appts {
		if statusFilter == "" || a.Status == statusFilter {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (m *mockApptRepo) CheckConflict(_ context.Context, vetID string, apptTime time.Time, durationMinutes int, excludeID string) (bool, error) {
	startWindow := apptTime.Add(-time.Duration(durationMinutes) * time.Minute)
	endWindow := apptTime.Add(time.Duration(durationMinutes) * time.Minute)

	for _, a := range m.appts {
		if a.ID == excludeID {
			continue
		}
		if a.VeterinarianID == vetID && (a.Status == model.StatusPending || a.Status == model.StatusConfirmed || a.Status == model.StatusInProgress) {
			if a.AppointmentDate.After(startWindow) && a.AppointmentDate.Before(endWindow) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *mockApptRepo) CountByStatus(_ context.Context) (map[string]int64, error) {
	res := make(map[string]int64)
	for _, a := range m.appts {
		res[a.Status]++
	}
	return res, nil
}

func (m *mockApptRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.appts)), nil
}

func (m *mockApptRepo) CountByOwner(_ context.Context, ownerID string) (int64, error) {
	var count int64
	for _, a := range m.appts {
		if a.OwnerID == ownerID {
			count++
		}
	}
	return count, nil
}

func (m *mockApptRepo) CountUpcomingByOwner(_ context.Context, ownerID string) (int64, error) {
	var count int64
	for _, a := range m.appts {
		if a.OwnerID == ownerID && (a.Status == model.StatusPending || a.Status == model.StatusConfirmed) {
			count++
		}
	}
	return count, nil
}

func (m *mockApptRepo) CountAssignedToVet(_ context.Context, vetID string) (int64, error) {
	var count int64
	for _, a := range m.appts {
		if a.VeterinarianID == vetID {
			count++
		}
	}
	return count, nil
}

func (m *mockApptRepo) CountTodayByVet(_ context.Context, _ string) (int64, error) {
	return 0, nil
}

func (m *mockApptRepo) CountCompletedByVet(_ context.Context, vetID string) (int64, error) {
	var count int64
	for _, a := range m.appts {
		if a.VeterinarianID == vetID && a.Status == model.StatusCompleted {
			count++
		}
	}
	return count, nil
}

// ── Mock MedicalRecordRepository ──────────────────────────────────────────────
type mockMedicalRecordRepo struct {
	records map[string]*model.MedicalRecord
}

func newMockMedicalRecordRepo() *mockMedicalRecordRepo {
	return &mockMedicalRecordRepo{records: make(map[string]*model.MedicalRecord)}
}

func (m *mockMedicalRecordRepo) Create(_ context.Context, r *model.MedicalRecord) error {
	for _, existing := range m.records {
		if existing.AppointmentID == r.AppointmentID {
			return repository.ErrDuplicateMedicalRecord
		}
	}
	if r.ID == "" {
		r.ID = "rec-1"
	}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	m.records[r.ID] = r
	return nil
}

func (m *mockMedicalRecordRepo) Update(_ context.Context, r *model.MedicalRecord) error {
	r.UpdatedAt = time.Now()
	m.records[r.ID] = r
	return nil
}

func (m *mockMedicalRecordRepo) FindByID(_ context.Context, id string) (*model.MedicalRecord, error) {
	r, ok := m.records[id]
	if !ok {
		return nil, repository.ErrMedicalRecordNotFound
	}
	return r, nil
}

func (m *mockMedicalRecordRepo) FindByAppointmentID(_ context.Context, apptID string) (*model.MedicalRecord, error) {
	for _, r := range m.records {
		if r.AppointmentID == apptID {
			return r, nil
		}
	}
	return nil, repository.ErrMedicalRecordNotFound
}

func (m *mockMedicalRecordRepo) ListByPet(_ context.Context, petID string) ([]model.MedicalRecord, error) {
	var list []model.MedicalRecord
	for _, r := range m.records {
		if r.PetID == petID {
			list = append(list, *r)
		}
	}
	return list, nil
}

func (m *mockMedicalRecordRepo) ListByVet(_ context.Context, vetID string) ([]model.MedicalRecord, error) {
	var list []model.MedicalRecord
	for _, r := range m.records {
		if r.VeterinarianID == vetID {
			list = append(list, *r)
		}
	}
	return list, nil
}

func (m *mockMedicalRecordRepo) ListAll(_ context.Context) ([]model.MedicalRecord, error) {
	var list []model.MedicalRecord
	for _, r := range m.records {
		list = append(list, *r)
	}
	return list, nil
}

func (m *mockMedicalRecordRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.records)), nil
}

func (m *mockMedicalRecordRepo) CountByPetOwner(_ context.Context, _ string) (int64, error) {
	return int64(len(m.records)), nil
}

func (m *mockMedicalRecordRepo) CountByVet(_ context.Context, vetID string) (int64, error) {
	var count int64
	for _, r := range m.records {
		if r.VeterinarianID == vetID {
			count++
		}
	}
	return count, nil
}

// ── Mock VaccinationRepository ────────────────────────────────────────────────
type mockVaccinationRepo struct {
	vacs map[string]*model.Vaccination
}

func newMockVaccinationRepo() *mockVaccinationRepo {
	return &mockVaccinationRepo{vacs: make(map[string]*model.Vaccination)}
}

func (m *mockVaccinationRepo) Create(_ context.Context, v *model.Vaccination) error {
	if v.ID == "" {
		v.ID = "vac-1"
	}
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	m.vacs[v.ID] = v
	return nil
}

func (m *mockVaccinationRepo) Update(_ context.Context, v *model.Vaccination) error {
	v.UpdatedAt = time.Now()
	m.vacs[v.ID] = v
	return nil
}

func (m *mockVaccinationRepo) Delete(_ context.Context, id string) error {
	if _, ok := m.vacs[id]; ok {
		delete(m.vacs, id)
		return nil
	}
	return repository.ErrVaccinationNotFound
}

func (m *mockVaccinationRepo) FindByID(_ context.Context, id string) (*model.Vaccination, error) {
	v, ok := m.vacs[id]
	if !ok {
		return nil, repository.ErrVaccinationNotFound
	}
	return v, nil
}

func (m *mockVaccinationRepo) ListByPet(_ context.Context, petID string) ([]model.Vaccination, error) {
	var list []model.Vaccination
	for _, v := range m.vacs {
		if v.PetID == petID {
			list = append(list, *v)
		}
	}
	return list, nil
}

func (m *mockVaccinationRepo) ListAll(_ context.Context) ([]model.Vaccination, error) {
	var list []model.Vaccination
	for _, v := range m.vacs {
		list = append(list, *v)
	}
	return list, nil
}

func (m *mockVaccinationRepo) Count(_ context.Context) (int64, error) {
	return int64(len(m.vacs)), nil
}

func (m *mockVaccinationRepo) CountByPetOwner(_ context.Context, _ string) (int64, error) {
	return int64(len(m.vacs)), nil
}
