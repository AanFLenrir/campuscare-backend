package model

import "time"

type User struct {
	ID            string    `json:"id"`
	RoleID        string    `json:"role_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	Phone         string    `json:"phone"`
	StudentNumber *string   `json:"student_number,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
