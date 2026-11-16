package model

import "time"

type MedicalRecord struct {
	ID              string    `json:"id"`
	AppointmentID   string    `json:"appointment_id"`
	PetID           string    `json:"pet_id"`
	VeterinarianID  string    `json:"veterinarian_id"`
	Diagnosis       string    `json:"diagnosis"`
	Symptoms        string    `json:"symptoms,omitempty"`
	Treatment       string    `json:"treatment,omitempty"`
	Prescription    string    `json:"prescription,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	ExaminationDate time.Time `json:"examination_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relational details
	PetName          string `json:"pet_name,omitempty"`
	PetSpecies       string `json:"pet_species,omitempty"`
	OwnerName        string `json:"owner_name,omitempty"`
	VeterinarianName string `json:"veterinarian_name,omitempty"`
	ClinicName       string `json:"clinic_name,omitempty"`
}

type CreateMedicalRecordInput struct {
	AppointmentID   string  `json:"appointment_id"`
	Diagnosis       string  `json:"diagnosis"`
	Symptoms        string  `json:"symptoms,omitempty"`
	Treatment       string  `json:"treatment,omitempty"`
	Prescription    string  `json:"prescription,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	ExaminationDate *string `json:"examination_date,omitempty"` // RFC3339, default now()
}

type UpdateMedicalRecordInput struct {
	Diagnosis    string `json:"diagnosis,omitempty"`
	Symptoms     string `json:"symptoms,omitempty"`
	Treatment    string `json:"treatment,omitempty"`
	Prescription string `json:"prescription,omitempty"`
	Notes        string `json:"notes,omitempty"`
}
