-- Up migration for medical_records table
CREATE TABLE IF NOT EXISTS medical_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID UNIQUE NOT NULL REFERENCES appointments(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    pet_id UUID NOT NULL REFERENCES pets(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    veterinarian_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    diagnosis TEXT NOT NULL,
    symptoms TEXT,
    treatment TEXT,
    prescription TEXT,
    notes TEXT,
    examination_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_medical_records_pet_id ON medical_records(pet_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_vet_id ON medical_records(veterinarian_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_appointment_id ON medical_records(appointment_id);

