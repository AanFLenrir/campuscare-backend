package model

import "time"

type ComplaintHistory struct {
	ID          string    `json:"id"`
	ComplaintID string    `json:"complaint_id"`
	UserID      string    `json:"user_id"`
	OldStatus   string    `json:"old_status"`
	NewStatus   string    `json:"new_status"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
