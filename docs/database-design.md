# PETRIVA — Database Design & Entity Relationship Diagram (ERD)

## 1. Konsep Database
PETRIVA menggunakan **PostgreSQL** dengan tipe data primary key **UUID v4** untuk menjamin keunikan identitas data tanpa bergantung pada sequence auto-increment numerik. Seluruh foreign key memiliki relasi integritas referensial yang ketat (`ON UPDATE CASCADE ON DELETE RESTRICT` atau `SET NULL`).

---

## 2. Entity Relationship Diagram (ERD)

```mermaid
erDiagram
    ROLES ||--o{ USERS : "has role"
    CLINICS ||--o{ USERS : "employs veterinarian"
    CLINICS ||--o{ SERVICES : "offers"
    USERS ||--o{ PETS : "owns"
    USERS ||--o{ APPOINTMENTS : "books (owner)"
    USERS ||--o{ APPOINTMENTS : "attends (vet)"
    CLINICS ||--o{ APPOINTMENTS : "hosts"
    SERVICES ||--o{ APPOINTMENTS : "applied in"
    PETS ||--o{ APPOINTMENTS : "scheduled for"
    APPOINTMENTS ||--|| MEDICAL_RECORDS : "generates (1-to-1)"
    PETS ||--o{ MEDICAL_RECORDS : "has"
    USERS ||--o{ MEDICAL_RECORDS : "diagnoses (vet)"
    PETS ||--o{ VACCINATIONS : "receives"
    USERS ||--o{ VACCINATIONS : "administers (vet)"

    ROLES {
        uuid id PK
        varchar name UK
        text description
        timestamptz created_at
        timestamptz updated_at
    }

    CLINICS {
        uuid id PK
        varchar name
        text address
        varchar phone
        varchar email
        text description
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    USERS {
        uuid id PK
        uuid role_id FK
        uuid clinic_id FK
        varchar name
        varchar email UK
        varchar password
        varchar phone
        text address
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    PETS {
        uuid id PK
        uuid owner_id FK
        varchar name
        varchar species
        varchar breed
        varchar gender
        date birth_date
        numeric weight
        varchar color
        text notes
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    SERVICES {
        uuid id PK
        uuid clinic_id FK
        varchar name
        text description
        integer duration_minutes
        numeric price
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    APPOINTMENTS {
        uuid id PK
        uuid pet_id FK
        uuid owner_id FK
        uuid veterinarian_id FK
        uuid clinic_id FK
        uuid service_id FK
        timestamptz appointment_date
        varchar status
        text complaint
        text cancellation_reason
        timestamptz created_at
        timestamptz updated_at
    }

    MEDICAL_RECORDS {
        uuid id PK
        uuid appointment_id UK,FK
        uuid pet_id FK
        uuid veterinarian_id FK
        text diagnosis
        text symptoms
        text treatment
        text prescription
        text notes
        timestamptz examination_date
        timestamptz created_at
        timestamptz updated_at
    }

    VACCINATIONS {
        uuid id PK
        uuid pet_id FK
        uuid veterinarian_id FK
        varchar vaccine_name
        date vaccine_date
        date next_vaccine_date
        varchar batch_number
        text notes
        timestamptz created_at
        timestamptz updated_at
    }
```

---

## 3. Rincian Skema Tabel

### 3.1. `roles`
Menyimpan peran pengguna sistem: `ADMIN`, `PET_OWNER`, `VETERINARIAN`.
- `id` (UUID, PK)
- `name` (VARCHAR(50), UNIQUE, NOT NULL)
- `description` (TEXT)
- `created_at`, `updated_at` (TIMESTAMPTZ)

### 3.2. `clinics`
Menyimpan informasi klinik cabang atau pusat.
- `id` (UUID, PK)
- `name` (VARCHAR(150), NOT NULL)
- `address` (TEXT, NOT NULL)
- `phone` (VARCHAR(30))
- `email` (VARCHAR(150))
- `description` (TEXT)
- `is_active` (BOOLEAN, DEFAULT true)
- `created_at`, `updated_at` (TIMESTAMPTZ)

### 3.3. `users`
Menyimpan akun pengguna seluruh role.
- `id` (UUID, PK)
- `role_id` (UUID, FK -> roles.id)
- `clinic_id` (UUID, NULLable, FK -> clinics.id) -> Wajib diisi jika role `VETERINARIAN`
- `name` (VARCHAR(150), NOT NULL)
- `email` (VARCHAR(150), UNIQUE, NOT NULL)
- `password` (VARCHAR(255), NOT NULL, bcrypt hash)
- `phone` (VARCHAR(30))
- `address` (TEXT)
- `is_active` (BOOLEAN, DEFAULT true)
- `created_at`, `updated_at` (TIMESTAMPTZ)

### 3.4. `pets`
Menyimpan profil hewan peliharaan milik pemilik.
- `id` (UUID, PK)
- `owner_id` (UUID, FK -> users.id)
- `name` (VARCHAR(100), NOT NULL)
- `species` (VARCHAR(50), NOT NULL) -> `CAT`, `DOG`, `BIRD`, `RABBIT`, `OTHER`
- `breed` (VARCHAR(100))
- `gender` (VARCHAR(20), NOT NULL) -> `MALE`, `FEMALE`, `UNKNOWN`
- `birth_date` (DATE)
- `weight` (NUMERIC(6,2))
- `color` (VARCHAR(50))
- `notes` (TEXT)
- `is_active` (BOOLEAN, DEFAULT true) -> *Soft delete support*
- `created_at`, `updated_at` (TIMESTAMPTZ)

### 3.5. `services`
Katalog tindakan atau perawatan yang disediakan klinik.
- `id` (UUID, PK)
- `clinic_id` (UUID, FK -> clinics.id)
- `name` (VARCHAR(150), NOT NULL)
- `description` (TEXT)
- `duration_minutes` (INTEGER, DEFAULT 30)
- `price` (NUMERIC(12,2), DEFAULT 0)
- `is_active` (BOOLEAN, DEFAULT true)
- `created_at`, `updated_at` (TIMESTAMPTZ)

### 3.6. `appointments`
Reservasi janji temu konsultasi / tindakan medis hewan.
- `id` (UUID, PK)
- `pet_id` (UUID, FK -> pets.id)
- `owner_id` (UUID, FK -> users.id)
- `veterinarian_id` (UUID, FK -> users.id)
- `clinic_id` (UUID, FK -> clinics.id)
- `service_id` (UUID, FK -> services.id)
- `appointment_date` (TIMESTAMPTZ, NOT NULL)
- `status` (VARCHAR(30), CHECK: `PENDING`, `CONFIRMED`, `IN_PROGRESS`, `COMPLETED`, `CANCELLED`, `REJECTED`)
- `complaint` (TEXT)
- `cancellation_reason` (TEXT, NULLable)
- `created_at`, `updated_at` (TIMESTAMPTZ)

### 3.7. `medical_records`
Catatan rekam medis pemeriksaan dokter hewan (relasi 1-to-1 dengan appointment).
- `id` (UUID, PK)
- `appointment_id` (UUID, UNIQUE, FK -> appointments.id)
- `pet_id` (UUID, FK -> pets.id)
- `veterinarian_id` (UUID, FK -> users.id)
- `diagnosis` (TEXT, NOT NULL)
- `symptoms` (TEXT)
- `treatment` (TEXT)
- `prescription` (TEXT)
- `notes` (TEXT)
- `examination_date` (TIMESTAMPTZ, DEFAULT now())
- `created_at`, `updated_at` (TIMESTAMPTZ)

### 3.8. `vaccinations`
Catatan buku vaksinasi hewan peliharaan.
- `id` (UUID, PK)
- `pet_id` (UUID, FK -> pets.id)
- `veterinarian_id` (UUID, FK -> users.id)
- `vaccine_name` (VARCHAR(150), NOT NULL)
- `vaccine_date` (DATE, NOT NULL)
- `next_vaccine_date` (DATE, NULLable)
- `batch_number` (VARCHAR(100), NULLable)
- `notes` (TEXT)
- `created_at`, `updated_at` (TIMESTAMPTZ)

