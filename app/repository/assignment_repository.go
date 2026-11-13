package repository

import (
	"context"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AssignmentRepository defines the data access methods for Assignment entities.
type AssignmentRepository interface {
	Create(ctx context.Context, a *model.Assignment) error
	GetByID(ctx context.Context, id string) (*model.Assignment, error)
	GetAllByOfficer(ctx context.Context, officerID string) ([]*model.Assignment, error)
	Update(ctx context.Context, a *model.Assignment) error
	Delete(ctx context.Context, id string) error
}

type assignmentRepo struct {
	db *pgxpool.Pool
}

func NewAssignmentRepository(db *pgxpool.Pool) AssignmentRepository {
	return &assignmentRepo{db: db}
}

func (r *assignmentRepo) Create(ctx context.Context, a *model.Assignment) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	query := `INSERT INTO assignments (id, complaint_id, officer_id, assigned_by, assigned_at, notes)
              VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(ctx, query, a.ID, a.ComplaintID, a.OfficerID, a.AssignedBy, a.AssignedAt, a.Notes)
	return err
}

func (r *assignmentRepo) GetByID(ctx context.Context, id string) (*model.Assignment, error) {
	query := `SELECT id, complaint_id, officer_id, assigned_by, assigned_at, notes FROM assignments WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	var a model.Assignment
	err := row.Scan(&a.ID, &a.ComplaintID, &a.OfficerID, &a.AssignedBy, &a.AssignedAt, &a.Notes)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *assignmentRepo) GetAllByOfficer(ctx context.Context, officerID string) ([]*model.Assignment, error) {
	query := `SELECT id, complaint_id, officer_id, assigned_by, assigned_at, notes FROM assignments WHERE officer_id=$1`
	rows, err := r.db.Query(ctx, query, officerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Assignment
	for rows.Next() {
		var a model.Assignment
		if err := rows.Scan(&a.ID, &a.ComplaintID, &a.OfficerID, &a.AssignedBy, &a.AssignedAt, &a.Notes); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	return list, nil
}

func (r *assignmentRepo) Update(ctx context.Context, a *model.Assignment) error {
	query := `UPDATE assignments SET complaint_id=$1, officer_id=$2, assigned_by=$3, assigned_at=$4, notes=$5 WHERE id=$6`
	_, err := r.db.Exec(ctx, query, a.ComplaintID, a.OfficerID, a.AssignedBy, a.AssignedAt, a.Notes, a.ID)
	return err
}

func (r *assignmentRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM assignments WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
