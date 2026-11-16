package model

import "time"

type User struct {
	ID        string    `json:"id"`
	RoleID    string    `json:"role_id"`
	ClinicID  *string   `json:"clinic_id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password in JSON
	Phone     string    `json:"phone,omitempty"`
	Address   string    `json:"address,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relational fields populated when needed
	RoleName   string  `json:"role,omitempty"`
	ClinicName *string `json:"clinic_name,omitempty"`
}

type UserResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone,omitempty"`
	Address    string    `json:"address,omitempty"`
	Role       string    `json:"role"`
	ClinicID   *string   `json:"clinic_id,omitempty"`
	ClinicName *string   `json:"clinic_name,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type CreateUserInput struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    string  `json:"phone,omitempty"`
	Address  string  `json:"address,omitempty"`
	RoleName string  `json:"role"`
	ClinicID *string `json:"clinic_id,omitempty"`
}

type UpdateUserInput struct {
	Name     string  `json:"name,omitempty"`
	Phone    string  `json:"phone,omitempty"`
	Address  string  `json:"address,omitempty"`
	Password string  `json:"password,omitempty"`
	ClinicID *string `json:"clinic_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
