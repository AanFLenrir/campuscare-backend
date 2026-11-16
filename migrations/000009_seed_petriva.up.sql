-- Up migration: Seed initial demonstration data for PETRIVA
-- Roles are already seeded in 000001_create_roles_table.up.sql

-- 1. Seed Demo Clinic
INSERT INTO clinics (id, name, address, phone, email, description, is_active)
SELECT 
    'a0000000-0000-0000-0000-000000000001',
    'PETRIVA Central Animal Care & Clinic',
    'Jl. Dharmawangsa Dalam No. 45, Surabaya',
    '031-5021234',
    'surabaya@petriva.com',
    'Klinik perawatan dan kesehatan hewan terpadu dengan fasilitas lengkap dan dokter profesional.',
    true
WHERE NOT EXISTS (
    SELECT 1 FROM clinics WHERE id = 'a0000000-0000-0000-0000-000000000001'
)
ON CONFLICT (id) DO NOTHING;

-- 2. Seed Default Clinic Services
INSERT INTO services (id, clinic_id, name, description, duration_minutes, price, is_active)
SELECT 
    'b0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'Pemeriksaan Umum (General Checkup)',
    'Pemeriksaan fisik komprehensif, suhu tubuh, detak jantung, mata, dan telinga hewan.',
    30,
    75000.00,
    true
WHERE NOT EXISTS (
    SELECT 1 FROM services WHERE id = 'b0000000-0000-0000-0000-000000000001'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO services (id, clinic_id, name, description, duration_minutes, price, is_active)
SELECT 
    'b0000000-0000-0000-0000-000000000002',
    'a0000000-0000-0000-0000-000000000001',
    'Vaksinasi Rabies & Polivalen',
    'Pemberian vaksin untuk pencegahan rabies dan penyakit menular utama hewan.',
    20,
    150000.00,
    true
WHERE NOT EXISTS (
    SELECT 1 FROM services WHERE id = 'b0000000-0000-0000-0000-000000000002'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO services (id, clinic_id, name, description, duration_minutes, price, is_active)
SELECT 
    'b0000000-0000-0000-0000-000000000003',
    'a0000000-0000-0000-0000-000000000001',
    'Grooming & Perawatan Higienis',
    'Mandi medis/parasit, potong kuku, pembersihan telinga, dan perapian bulu.',
    60,
    120000.00,
    true
WHERE NOT EXISTS (
    SELECT 1 FROM services WHERE id = 'b0000000-0000-0000-0000-000000000003'
)
ON CONFLICT (id) DO NOTHING;

-- 3. Seed Demo Users (Password: Password123!)
-- Hash cost 10: $2a$10$TYAtvd/KIVZ/emVbpCtPh.UPklMjLCHsHA3o2/pipVbVQVJBDWIzi
-- Note: WHERE NOT EXISTS checks BOTH (id OR email) to avoid PK/Unique collision if re-run

-- 3a. Admin User
INSERT INTO users (id, role_id, clinic_id, name, email, password, phone, address, is_active)
SELECT 
    'c0000000-0000-0000-0000-000000000001',
    r.id,
    NULL,
    'Administrator Petriva',
    'admin@petriva.com',
    '$2a$10$TYAtvd/KIVZ/emVbpCtPh.UPklMjLCHsHA3o2/pipVbVQVJBDWIzi',
    '081234567890',
    'Surabaya, Jawa Timur',
    true
FROM roles r 
WHERE r.name = 'ADMIN'
  AND NOT EXISTS (
      SELECT 1 FROM users WHERE id = 'c0000000-0000-0000-0000-000000000001' OR email = 'admin@petriva.com'
  )
ON CONFLICT (email) DO NOTHING;

-- 3b. Veterinarian User (Assigned to PETRIVA Central Clinic)
INSERT INTO users (id, role_id, clinic_id, name, email, password, phone, address, is_active)
SELECT 
    'c0000000-0000-0000-0000-000000000002',
    r.id,
    'a0000000-0000-0000-0000-000000000001',
    'drh. Sarah Wijaya',
    'dokter@petriva.com',
    '$2a$10$TYAtvd/KIVZ/emVbpCtPh.UPklMjLCHsHA3o2/pipVbVQVJBDWIzi',
    '081234567891',
    'Jl. Kertajaya Indah No. 12, Surabaya',
    true
FROM roles r 
WHERE r.name = 'VETERINARIAN'
  AND NOT EXISTS (
      SELECT 1 FROM users WHERE id = 'c0000000-0000-0000-0000-000000000002' OR email = 'dokter@petriva.com'
  )
ON CONFLICT (email) DO NOTHING;

-- 3c. Pet Owner User
INSERT INTO users (id, role_id, clinic_id, name, email, password, phone, address, is_active)
SELECT 
    'c0000000-0000-0000-0000-000000000003',
    r.id,
    NULL,
    'Budi Santoso',
    'owner@petriva.com',
    '$2a$10$TYAtvd/KIVZ/emVbpCtPh.UPklMjLCHsHA3o2/pipVbVQVJBDWIzi',
    '081234567892',
    'Jl. Gubeng Kertajaya No. 8, Surabaya',
    true
FROM roles r 
WHERE r.name = 'PET_OWNER'
  AND NOT EXISTS (
      SELECT 1 FROM users WHERE id = 'c0000000-0000-0000-0000-000000000003' OR email = 'owner@petriva.com'
  )
ON CONFLICT (email) DO NOTHING;

-- 4. Seed Demo Pet (Belongs to owner@petriva.com)
INSERT INTO pets (id, owner_id, name, species, breed, gender, birth_date, weight, color, notes, is_active)
SELECT 
    'd0000000-0000-0000-0000-000000000001',
    u.id,
    'Milo',
    'CAT',
    'Domestic Shorthair',
    'MALE',
    '2023-05-10',
    4.25,
    'Orange Tabby',
    'Kucing peliharaan sehat, aktif, dan ramah.',
    true
FROM users u 
WHERE u.email = 'owner@petriva.com'
  AND NOT EXISTS (
      SELECT 1 FROM pets WHERE id = 'd0000000-0000-0000-0000-000000000001' OR (owner_id = u.id AND name = 'Milo')
  )
ON CONFLICT (id) DO NOTHING;
