package repository

import (
	"context"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ComplaintHistoryRepository defines data access for complaint history records.
type ComplaintHistoryRepository interface {
	Create(ctx context.Context, h *model.ComplaintHistory) error
	GetByComplaintID(ctx context.Context, complaintID string) ([]*model.ComplaintHistory, error)
	GetAll(ctx context.Context) ([]*model.ComplaintHistory, error)
}

type complaintHistoryRepo struct {
	db *pgxpool.Pool
}

func NewComplaintHistoryRepository(db *pgxpool.Pool) ComplaintHistoryRepository {
	return &complaintHistoryRepo{db: db}
}

func (r *complaintHistoryRepo) Create(ctx context.Context, h *model.ComplaintHistory) error {
	if h.ID == "" {
		h.ID = uuid.New().String()
	}
	query := `INSERT INTO complaint_histories (id, complaint_id, user_id, old_status, new_status, notes, created_at)
              VALUES ($1,$2,$3,$4,$5,$6,now())`
	_, err := r.db.Exec(ctx, query, h.ID, h.ComplaintID, h.UserID, h.OldStatus, h.NewStatus, h.Notes)
	return err
}

func (r *complaintHistoryRepo) GetByComplaintID(ctx context.Context, complaintID string) ([]*model.ComplaintHistory, error) {
	query := `SELECT id, complaint_id, user_id, old_status, new_status, notes, created_at FROM complaint_histories WHERE complaint_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, complaintID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.ComplaintHistory
	for rows.Next() {
		var h model.ComplaintHistory
		if err := rows.Scan(&h.ID, &h.ComplaintID, &h.UserID, &h.OldStatus, &h.NewStatus, &h.Notes, &h.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &h)
	}
	return list, nil
}

func (r *complaintHistoryRepo) GetAll(ctx context.Context) ([]*model.ComplaintHistory, error) {
	query := `SELECT id, complaint_id, user_id, old_status, new_status, notes, created_at FROM complaint_histories ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.ComplaintHistory
	for rows.Next() {
		var h model.ComplaintHistory
		if err := rows.Scan(&h.ID, &h.ComplaintID, &h.UserID, &h.OldStatus, &h.NewStatus, &h.Notes, &h.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &h)
	}
	return list, nil
}
