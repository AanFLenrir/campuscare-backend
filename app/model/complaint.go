package model

import "time"

type Complaint struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	FacilityID  string    `json:"facility_id"`
	CategoryID  string    `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Priority    string    `json:"priority"` // LOW, MEDIUM, HIGH, URGENT
	Status      string    `json:"status"`   // ENUM list
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
