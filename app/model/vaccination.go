package model

import "time"

type Vaccination struct {
	ID              string    `json:"id"`
	PetID           string    `json:"pet_id"`
	VeterinarianID  string    `json:"veterinarian_id"`
	VaccineName     string    `json:"vaccine_name"`
	VaccineDate     string    `json:"vaccine_date"` // YYYY-MM-DD
	NextVaccineDate *string   `json:"next_vaccine_date,omitempty"`
	BatchNumber     *string   `json:"batch_number,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relational details
	PetName          string `json:"pet_name,omitempty"`
	PetSpecies       string `json:"pet_species,omitempty"`
	OwnerName        string `json:"owner_name,omitempty"`
	VeterinarianName string `json:"veterinarian_name,omitempty"`
}

type CreateVaccinationInput struct {
	PetID           string  `json:"pet_id"`
	VaccineName     string  `json:"vaccine_name"`
	VaccineDate     string  `json:"vaccine_date"` // YYYY-MM-DD
	NextVaccineDate *string `json:"next_vaccine_date,omitempty"`
	BatchNumber     *string `json:"batch_number,omitempty"`
	Notes           string  `json:"notes,omitempty"`
}

type UpdateVaccinationInput struct {
	VaccineName     string  `json:"vaccine_name,omitempty"`
	VaccineDate     string  `json:"vaccine_date,omitempty"`
	NextVaccineDate *string `json:"next_vaccine_date,omitempty"`
	BatchNumber     *string `json:"batch_number,omitempty"`
	Notes           string  `json:"notes,omitempty"`
}
