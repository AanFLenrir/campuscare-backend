package model

import "time"

const (
	SpeciesCat    = "CAT"
	SpeciesDog    = "DOG"
	SpeciesBird   = "BIRD"
	SpeciesRabbit = "RABBIT"
	SpeciesOther  = "OTHER"

	GenderMale    = "MALE"
	GenderFemale  = "FEMALE"
	GenderUnknown = "UNKNOWN"
)

type Pet struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Name      string    `json:"name"`
	Species   string    `json:"species"`
	Breed     string    `json:"breed,omitempty"`
	Gender    string    `json:"gender"`
	BirthDate *string   `json:"birth_date,omitempty"` // YYYY-MM-DD
	Weight    *float64  `json:"weight,omitempty"`
	Color     string    `json:"color,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relational
	OwnerName  string `json:"owner_name,omitempty"`
	OwnerPhone string `json:"owner_phone,omitempty"`
	OwnerEmail string `json:"owner_email,omitempty"`
}

type CreatePetInput struct {
	Name      string   `json:"name"`
	Species   string   `json:"species"`
	Breed     string   `json:"breed,omitempty"`
	Gender    string   `json:"gender"`
	BirthDate *string  `json:"birth_date,omitempty"`
	Weight    *float64 `json:"weight,omitempty"`
	Color     string   `json:"color,omitempty"`
	Notes     string   `json:"notes,omitempty"`
}

type UpdatePetInput struct {
	Name      string   `json:"name,omitempty"`
	Species   string   `json:"species,omitempty"`
	Breed     string   `json:"breed,omitempty"`
	Gender    string   `json:"gender,omitempty"`
	BirthDate *string  `json:"birth_date,omitempty"`
	Weight    *float64 `json:"weight,omitempty"`
	Color     string   `json:"color,omitempty"`
	Notes     string   `json:"notes,omitempty"`
}
