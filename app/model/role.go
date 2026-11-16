package model

import "time"

const (
	RoleAdmin        = "ADMIN"
	RolePetOwner     = "PET_OWNER"
	RoleVeterinarian = "VETERINARIAN"
)

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
