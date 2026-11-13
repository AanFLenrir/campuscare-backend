package repository

import (
	"context"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *model.Category) error
	GetAll(ctx context.Context) ([]*model.Category, error)
	GetByID(ctx context.Context, id string) (*model.Category, error)
	Update(ctx context.Context, c *model.Category) error
	Delete(ctx context.Context, id string) error
}

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) Create(ctx context.Context, c *model.Category) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	query := `INSERT INTO categories (id, name, description, created_at, updated_at) VALUES ($1,$2,$3,now(),now())`
	_, err := r.db.Exec(ctx, query, c.ID, c.Name, c.Description)
	return err
}

func (r *categoryRepo) GetAll(ctx context.Context) ([]*model.Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, description, created_at, updated_at FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, nil
}

func (r *categoryRepo) GetByID(ctx context.Context, id string) (*model.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	var c model.Category
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepo) Update(ctx context.Context, c *model.Category) error {
	query := `UPDATE categories SET name=$1, description=$2, updated_at=now() WHERE id=$3`
	_, err := r.db.Exec(ctx, query, c.Name, c.Description, c.ID)
	return err
}

func (r *categoryRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
