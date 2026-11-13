package repository

import (
	"context"
	"database/sql"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines data access methods for User entity.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error // soft delete (set is_active = false)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	query := `INSERT INTO users (id, role_id, name, email, password, phone, student_number, is_active, created_at, updated_at)
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now(),now())`
	_, err := r.db.Exec(ctx, query, user.ID, user.RoleID, user.Name, user.Email, user.Password, user.Phone, user.StudentNumber, user.IsActive)
	return err
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, role_id, name, email, password, phone, student_number, is_active, created_at, updated_at FROM users WHERE email=$1`
	row := r.db.QueryRow(ctx, query, email)
	var u model.User
	var studentNumber sql.NullString
	err := row.Scan(&u.ID, &u.RoleID, &u.Name, &u.Email, &u.Password, &u.Phone, &studentNumber, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if studentNumber.Valid {
		u.StudentNumber = &studentNumber.String
	}
	return &u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, role_id, name, email, password, phone, student_number, is_active, created_at, updated_at FROM users WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	var u model.User
	var studentNumber sql.NullString
	err := row.Scan(&u.ID, &u.RoleID, &u.Name, &u.Email, &u.Password, &u.Phone, &studentNumber, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if studentNumber.Valid {
		u.StudentNumber = &studentNumber.String
	}
	return &u, nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]*model.User, error) {
	rows, err := r.db.Query(ctx, `SELECT id, role_id, name, email, password, phone, student_number, is_active, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.User
	for rows.Next() {
		var u model.User
		var studentNumber sql.NullString
		if err := rows.Scan(&u.ID, &u.RoleID, &u.Name, &u.Email, &u.Password, &u.Phone, &studentNumber, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		if studentNumber.Valid {
			u.StudentNumber = &studentNumber.String
		}
		list = append(list, &u)
	}
	return list, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET role_id=$1, name=$2, email=$3, password=$4, phone=$5, student_number=$6, is_active=$7, updated_at=now() WHERE id=$8`
	_, err := r.db.Exec(ctx, query, user.RoleID, user.Name, user.Email, user.Password, user.Phone, user.StudentNumber, user.IsActive, user.ID)
	return err
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	// Soft delete: set is_active to false
	query := `UPDATE users SET is_active=false, updated_at=now() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
