package model

import "time"

type ComplaintImage struct {
	ID          string    `json:"id"`
	ComplaintID string    `json:"complaint_id"`
	ImagePath   string    `json:"image_path"`
	CreatedAt   time.Time `json:"created_at"`
}
