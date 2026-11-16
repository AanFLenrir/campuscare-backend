-- Down migration for PETRIVA demonstration seed data
-- Reverses inserts in strict reverse dependency order

-- 1. Remove demo pet
DELETE FROM pets 
WHERE id = 'd0000000-0000-0000-0000-000000000001'
   OR (name = 'Milo' AND owner_id IN (SELECT id FROM users WHERE email = 'owner@petriva.com'));

-- 2. Remove demo users
DELETE FROM users 
WHERE email IN ('admin@petriva.com', 'dokter@petriva.com', 'owner@petriva.com')
   OR id IN (
       'c0000000-0000-0000-0000-000000000001',
       'c0000000-0000-0000-0000-000000000002',
       'c0000000-0000-0000-0000-000000000003'
   );

-- 3. Remove demo services
DELETE FROM services 
WHERE clinic_id = 'a0000000-0000-0000-0000-000000000001'
   OR id IN (
       'b0000000-0000-0000-0000-000000000001',
       'b0000000-0000-0000-0000-000000000002',
       'b0000000-0000-0000-0000-000000000003'
   );

-- 4. Remove demo clinic
DELETE FROM clinics 
WHERE id = 'a0000000-0000-0000-0000-000000000001'
   OR name = 'PETRIVA Central Animal Care & Clinic';
