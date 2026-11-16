-- Up migration for vaccinations table
CREATE TABLE IF NOT EXISTS vaccinations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    veterinarian_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    vaccine_name VARCHAR(150) NOT NULL,
    vaccine_date DATE NOT NULL,
    next_vaccine_date DATE,
    batch_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_vaccinations_pet_id ON vaccinations(pet_id);
CREATE INDEX IF NOT EXISTS idx_vaccinations_vet_id ON vaccinations(veterinarian_id);
CREATE INDEX IF NOT EXISTS idx_vaccinations_date ON vaccinations(vaccine_date);

