package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AanFLenrir/PetrivaBackend/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRoleNotFound = errors.New("role not found")

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*model.Role, error)
	FindByID(ctx context.Context, id string) (*model.Role, error)
	List(ctx context.Context) ([]model.Role, error)
}

type roleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	query := `SELECT id, name, COALESCE(description, ''), created_at, updated_at FROM roles WHERE name = $1 LIMIT 1`
	var role model.Role
	err := r.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("role query error: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	query := `SELECT id, name, COALESCE(description, ''), created_at, updated_at FROM roles WHERE id = $1 LIMIT 1`
	var role model.Role
	err := r.db.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("role query error: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]model.Role, error) {
	query := `SELECT id, name, COALESCE(description, ''), created_at, updated_at FROM roles ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list roles error: %w", err)
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
