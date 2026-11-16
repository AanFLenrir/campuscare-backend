-- Up migration for roles table
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed basic roles
INSERT INTO roles (name, description) VALUES
    ('ADMIN', 'Administrator sistem dengan hak akses penuh'),
    ('PET_OWNER', 'Pemilik hewan peliharaan'),
    ('VETERINARIAN', 'Dokter hewan penanggung jawab pemeriksaan')
ON CONFLICT (name) DO NOTHING;
