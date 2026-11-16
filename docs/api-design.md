# PETRIVA — API Design & Endpoint Specification

## 1. Konvensi Umum
- **Base URL:** `http://localhost:3000/api/v1`
- **Swagger Documentation:** `http://localhost:3000/swagger/`
- **Format Pertukaran Data:** JSON (`Content-Type: application/json`)
- **Autentikasi:** HTTP Header `Authorization: Bearer <JWT_TOKEN>`

---

## 2. Format Response Standar

### Response Sukses (200 OK / 201 Created)
```json
{
  "success": true,
  "message": "Success",
  "data": { ... }
}
```

### Response Error (400 / 401 / 403 / 404 / 500)
```json
{
  "success": false,
  "message": "Deskripsi pesan kesalahan spesifik"
}
```

---

## 3. Rincian Endpoint

### 3.1. Authentication & Profil

#### `POST /api/v1/auth/register`
- **Akses:** Publik (hanya mendaftarkan role `PET_OWNER`)
- **Body:**
  ```json
  {
    "name": "Budi Santoso",
    "email": "budi@example.com",
    "password": "Password123!",
    "phone": "081234567890",
    "address": "Jl. Mawar No. 10, Surabaya"
  }
  ```
- **Response 201 Created:**
  ```json
  {
    "success": true,
    "message": "Created",
    "data": {
      "token": "eyJhbGciOi...",
      "user": {
        "id": "uuid-...",
        "name": "Budi Santoso",
        "email": "budi@example.com",
        "role": "PET_OWNER",
        "is_active": true
      }
    }
  }
  ```

#### `POST /api/v1/auth/login`
- **Akses:** Publik (seluruh role: `ADMIN`, `PET_OWNER`, `VETERINARIAN`)
- **Body:**
  ```json
  {
    "email": "budi@example.com",
    "password": "Password123!"
  }
  ```
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Success",
    "data": {
      "token": "eyJhbGciOi...",
      "user": {
        "id": "uuid-...",
        "name": "Budi Santoso",
        "email": "budi@example.com",
        "role": "PET_OWNER",
        "is_active": true
      }
    }
  }
  ```

#### `GET /api/v1/auth/me`
- **Akses:** Protected (Semua Role)
- **Header:** `Authorization: Bearer <token>`
- **Response 200 OK:** Data profil pengguna yang sedang terautentikasi.

---

### 3.2. Dashboard Summary

#### `GET /api/v1/dashboard/summary`
- **Akses:** Protected (Semua Role)
- **Keterangan:** Mengembalikan data agregasi statistik yang disesuaikan secara dinamis berdasarkan role token:
  - **ADMIN:** Total user, total klinik, total hewan, total appointment, total rekam medis, total vaksinasi, dan rincian status appointment.
  - **PET_OWNER:** Total hewan peliharaan milik sendiri, total appointment, jadwal mendatang (*upcoming*), riwayat rekam medis & vaksinasi.
  - **VETERINARIAN:** Total appointment ditugaskan, jadwal hari ini, pemeriksaan selesai, dan rekam medis yang telah dibuat.

---

### 3.3. Pets (Hewan Peliharaan)

#### `POST /api/v1/pets`
- **Akses:** `PET_OWNER`
- **Body:**
  ```json
  {
    "name": "Milo",
    "species": "CAT",
    "breed": "British Shorthair",
    "gender": "MALE",
    "birth_date": "2023-05-10",
    "weight": 4.25,
    "color": "Grey",
    "notes": "Suka makanan basah, alergi salmon"
  }
  ```

#### `GET /api/v1/pets`
- **Akses:** `PET_OWNER`, `ADMIN`
- **Keterangan:** `PET_OWNER` otomatis hanya melihat hewan miliknya sendiri (diambil dari token). `ADMIN` dapat melihat seluruh hewan atau filter `?owner_id=`.

#### `GET /api/v1/pets/:id`
- **Akses:** `PET_OWNER` (pemilik), `VETERINARIAN` (yang menangani appointment), `ADMIN`.

#### `PUT /api/v1/pets/:id`
- **Akses:** `PET_OWNER` (hanya pemilik sah).

#### `DELETE /api/v1/pets/:id`
- **Akses:** `PET_OWNER` (*soft delete* / nonaktifkan).

---

### 3.4. Appointments (Reservasi Konsultasi)

#### `POST /api/v1/appointments`
- **Akses:** `PET_OWNER`
- **Body:**
  ```json
  {
    "pet_id": "uuid-pet",
    "veterinarian_id": "uuid-vet",
    "clinic_id": "uuid-clinic",
    "service_id": "uuid-service",
    "appointment_date": "2026-11-20T10:00:00Z",
    "complaint": "Kucing lemas dan tidak nafsu makan selama 2 hari"
  }
  ```
- **Validasi Bisnis:**
  1. Hewan benar-benar milik pemilik yang login.
  2. Waktu appointment harus di masa depan.
  3. Dokter hewan aktif dan terdaftar pada klinik yang dipilih.
  4. Layanan aktif dan terdaftar pada klinik yang dipilih.
  5. Dokter hewan tidak memiliki jadwal bentrok pada jam tersebut.
  6. Status awal: `PENDING`.

#### `GET /api/v1/appointments`
- **Akses:** Semua Role (terfilter otomatis berdasarkan kepemilikan / penugasan).

#### `PUT /api/v1/appointments/:id/cancel`
- **Akses:** `PET_OWNER` (pemilik), `ADMIN`
- **Body:**
  ```json
  {
    "reason": "Ada keperluan mendadak di luar kota"
  }
  ```
- **Validasi:** Hanya appointment berstatus `PENDING` atau `CONFIRMED` yang dapat dibatalkan.

#### `PUT /api/v1/appointments/:id/status`
- **Akses:** `VETERINARIAN` (yang ditugaskan), `ADMIN`
- **Body:**
  ```json
  {
    "status": "CONFIRMED"
  }
  ```
- **Aturan Transisi Status:**
  - `PENDING` → `CONFIRMED`
  - `PENDING` → `REJECTED`
  - `CONFIRMED` → `IN_PROGRESS`
  - `IN_PROGRESS` → `COMPLETED`

---

### 3.5. Medical Records (Rekam Medis)

#### `POST /api/v1/medical-records`
- **Akses:** `VETERINARIAN`
- **Body:**
  ```json
  {
    "appointment_id": "uuid-appointment",
    "diagnosis": "Gastroenteritis akut",
    "symptoms": "Muntah berbusa, diare ringan, dehidrasi sedang",
    "treatment": "Injeksi antiemetik dan cairan infus rehidrasi 200ml",
    "prescription": "Antibiotik cair 2x1ml (5 hari), Probiotik 1x1 sachet",
    "notes": "Kontrol kembali jika dalam 3 hari muntah belum membaik"
  }
  ```
- **Validasi Bisnis:**
  1. Hanya dokter hewan penanggung jawab appointment yang boleh mengisi rekam medis.
  2. Appointment harus berstatus `CONFIRMED`, `IN_PROGRESS`, atau `COMPLETED`.
  3. Mencegah duplikasi rekam medis untuk appointment yang sama.
  4. Status appointment otomatis bertransisi menjadi `COMPLETED`.

#### `GET /api/v1/medical-records`
- **Akses:** Semua Role (terlindungi per hewan/pemilik).

---

### 3.6. Vaccinations (Buku Vaksinasi)

#### `POST /api/v1/vaccinations`
- **Akses:** `VETERINARIAN`
- **Body:**
  ```json
  {
    "pet_id": "uuid-pet",
    "vaccine_name": "Rabisin (Anti-Rabies)",
    "vaccine_date": "2026-11-16",
    "next_vaccine_date": "2027-11-16",
    "batch_number": "RAB-2026-X89",
    "notes": "Kondisi sehat saat vaksinasi, suhu tubuh 38.4C"
  }
  ```

---

### 3.7. Clinics & Services Catalog

- `GET /api/v1/clinics` — Publik
- `GET /api/v1/clinics/:id` — Publik
- `POST /api/v1/clinics` — `ADMIN`
- `PUT /api/v1/clinics/:id` — `ADMIN`
- `DELETE /api/v1/clinics/:id` — `ADMIN`
- `GET /api/v1/services` — Publik (filter `?clinic_id=`)
- `GET /api/v1/services/:id` — Publik
- `POST /api/v1/services` — `ADMIN`
- `PUT /api/v1/services/:id` — `ADMIN`
- `DELETE /api/v1/services/:id` — `ADMIN`

