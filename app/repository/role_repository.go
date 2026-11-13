package repository

import (
	"context"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*model.Role, error)
	GetByID(ctx context.Context, id string) (*model.Role, error)
}

type roleRepo struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &roleRepo{db: db}
}

func (r *roleRepo) GetByName(ctx context.Context, name string) (*model.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name=$1`
	row := r.db.QueryRow(ctx, query, name)
	var role model.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) GetByID(ctx context.Context, id string) (*model.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	var role model.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
