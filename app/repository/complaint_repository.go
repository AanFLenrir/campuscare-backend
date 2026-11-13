package repository

import (
	"context"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ComplaintRepository interface {
	Create(ctx context.Context, c *model.Complaint) error
	GetByID(ctx context.Context, id string) (*model.Complaint, error)
	GetAll(ctx context.Context) ([]*model.Complaint, error)
	Update(ctx context.Context, c *model.Complaint) error
	Delete(ctx context.Context, id string) error
	// Additional query for student own complaints could be added later.
}

type complaintRepo struct {
	db *pgxpool.Pool
}

func NewComplaintRepository(db *pgxpool.Pool) ComplaintRepository {
	return &complaintRepo{db: db}
}

func (r *complaintRepo) Create(ctx context.Context, c *model.Complaint) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	query := `INSERT INTO complaints (id, user_id, facility_id, category_id, title, description, priority, status, created_at, updated_at)
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now(),now())`
	_, err := r.db.Exec(ctx, query, c.ID, c.UserID, c.FacilityID, c.CategoryID, c.Title, c.Description, c.Priority, c.Status)
	return err
}

func (r *complaintRepo) GetByID(ctx context.Context, id string) (*model.Complaint, error) {
	query := `SELECT id, user_id, facility_id, category_id, title, description, priority, status, created_at, updated_at FROM complaints WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	var c model.Complaint
	err := row.Scan(&c.ID, &c.UserID, &c.FacilityID, &c.CategoryID, &c.Title, &c.Description, &c.Priority, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *complaintRepo) GetAll(ctx context.Context) ([]*model.Complaint, error) {
	rows, err := r.db.Query(ctx, `SELECT id, user_id, facility_id, category_id, title, description, priority, status, created_at, updated_at FROM complaints`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Complaint
	for rows.Next() {
		var c model.Complaint
		if err := rows.Scan(&c.ID, &c.UserID, &c.FacilityID, &c.CategoryID, &c.Title, &c.Description, &c.Priority, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, nil
}

func (r *complaintRepo) Update(ctx context.Context, c *model.Complaint) error {
	query := `UPDATE complaints SET title=$1, description=$2, priority=$3, status=$4, updated_at=now() WHERE id=$5`
	_, err := r.db.Exec(ctx, query, c.Title, c.Description, c.Priority, c.Status, c.ID)
	return err
}

func (r *complaintRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM complaints WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
