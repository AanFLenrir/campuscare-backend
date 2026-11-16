package model

import "time"

const (
	StatusPending    = "PENDING"
	StatusConfirmed  = "CONFIRMED"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	StatusCancelled  = "CANCELLED"
	StatusRejected   = "REJECTED"
)

type Appointment struct {
	ID                 string    `json:"id"`
	PetID              string    `json:"pet_id"`
	OwnerID            string    `json:"owner_id"`
	VeterinarianID     string    `json:"veterinarian_id"`
	ClinicID           string    `json:"clinic_id"`
	ServiceID          string    `json:"service_id"`
	AppointmentDate    time.Time `json:"appointment_date"`
	Status             string    `json:"status"`
	Complaint          string    `json:"complaint,omitempty"`
	CancellationReason *string   `json:"cancellation_reason,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Relational details
	PetName          string  `json:"pet_name,omitempty"`
	PetSpecies       string  `json:"pet_species,omitempty"`
	OwnerName        string  `json:"owner_name,omitempty"`
	VeterinarianName string  `json:"veterinarian_name,omitempty"`
	ClinicName       string  `json:"clinic_name,omitempty"`
	ServiceName      string  `json:"service_name,omitempty"`
	ServicePrice     float64 `json:"service_price,omitempty"`
}

type CreateAppointmentInput struct {
	PetID           string `json:"pet_id"`
	VeterinarianID  string `json:"veterinarian_id"`
	ClinicID        string `json:"clinic_id"`
	ServiceID       string `json:"service_id"`
	AppointmentDate string `json:"appointment_date"` // RFC3339 format (e.g. 2026-11-20T10:00:00Z)
	Complaint       string `json:"complaint,omitempty"`
}

type UpdateAppointmentStatusInput struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type CancelAppointmentInput struct {
	Reason string `json:"reason"`
}
