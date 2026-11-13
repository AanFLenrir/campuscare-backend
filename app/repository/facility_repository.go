package repository

import (
	"context"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FacilityRepository interface {
	Create(ctx context.Context, f *model.Facility) error
	GetAll(ctx context.Context) ([]*model.Facility, error)
	GetByID(ctx context.Context, id string) (*model.Facility, error)
	Update(ctx context.Context, f *model.Facility) error
	Delete(ctx context.Context, id string) error
}

type facilityRepo struct {
	db *pgxpool.Pool
}

func NewFacilityRepository(db *pgxpool.Pool) FacilityRepository {
	return &facilityRepo{db: db}
}

func (r *facilityRepo) Create(ctx context.Context, f *model.Facility) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	query := `INSERT INTO facilities (id, name, location, description, condition, created_at, updated_at)
              VALUES ($1,$2,$3,$4,$5,now(),now())`
	_, err := r.db.Exec(ctx, query, f.ID, f.Name, f.Location, f.Description, f.Condition)
	return err
}

func (r *facilityRepo) GetAll(ctx context.Context) ([]*model.Facility, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, location, description, condition, created_at, updated_at FROM facilities`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Facility
	for rows.Next() {
		var f model.Facility
		if err := rows.Scan(&f.ID, &f.Name, &f.Location, &f.Description, &f.Condition, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &f)
	}
	return list, nil
}

func (r *facilityRepo) GetByID(ctx context.Context, id string) (*model.Facility, error) {
	query := `SELECT id, name, location, description, condition, created_at, updated_at FROM facilities WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	var f model.Facility
	err := row.Scan(&f.ID, &f.Name, &f.Location, &f.Description, &f.Condition, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *facilityRepo) Update(ctx context.Context, f *model.Facility) error {
	query := `UPDATE facilities SET name=$1, location=$2, description=$3, condition=$4, updated_at=now() WHERE id=$5`
	_, err := r.db.Exec(ctx, query, f.Name, f.Location, f.Description, f.Condition, f.ID)
	return err
}

func (r *facilityRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM facilities WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
