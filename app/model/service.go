package model

import "time"

type Service struct {
	ID              string    `json:"id"`
	ClinicID        string    `json:"clinic_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	DurationMinutes int       `json:"duration_minutes"`
	Price           float64   `json:"price"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	ClinicName string `json:"clinic_name,omitempty"`
}

type CreateServiceInput struct {
	ClinicID        string  `json:"clinic_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description,omitempty"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
}

type UpdateServiceInput struct {
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	DurationMinutes *int     `json:"duration_minutes,omitempty"`
	Price           *float64 `json:"price,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}
