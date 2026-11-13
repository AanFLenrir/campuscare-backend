package model

import "time"

type Assignment struct {
	ID          string    `json:"id"`
	ComplaintID string    `json:"complaint_id"`
	OfficerID   string    `json:"officer_id"`
	AssignedBy  string    `json:"assigned_by"`
	AssignedAt  time.Time `json:"assigned_at"`
	Notes       string    `json:"notes,omitempty"`
}
