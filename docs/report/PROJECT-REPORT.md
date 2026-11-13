# LAPORAN RANCANGAN SISTEM DAN BOILERPLATE BACKEND
## CAMPUSCARE — Campus Facility Complaint and Management System

---

**Program Studi:** D4 Teknik Informatika  
**Fakultas:** Vokasi  
**Universitas:** Universitas Airlangga  
**Tahun:** 2026

---

## Identitas Kelompok

| No | Nama | NIM | Peran |
|---|---|---|---|
| 1 | MOHAMMAD ANHAZ ABDILAH | 434241102 | Developer & Git Master |
| 2 | *(Anggota 2)* | — | *(isi sesuai data kelompok)* |
| 3 | *(Anggota 3)* | — | *(isi sesuai data kelompok)* |
| 4 | *(Anggota 4)* | — | *(isi sesuai data kelompok)* |

**Link Repository GitHub:** https://github.com/AanFLenrir/campuscare-backend

---

## 1. Latar Belakang

Fasilitas kampus seperti ruang kelas, laboratorium, dan area umum sering mengalami kerusakan yang tidak segera dilaporkan maupun ditangani. Ketiadaan sistem pelaporan yang terstruktur menyebabkan kerusakan terakumulasi dan mengganggu proses belajar mengajar.

CAMPUSCARE hadir sebagai solusi digital untuk:
- Memudahkan mahasiswa melaporkan kerusakan fasilitas secara sistematis.
- Membantu admin dan petugas memantau serta menangani pengaduan secara efisien.
- Menyediakan riwayat dan laporan penanganan yang dapat dipertanggungjawabkan.

---

## 2. Deskripsi Sistem

**Nama Sistem:** CAMPUSCARE — Campus Facility Complaint and Management System

**Deskripsi:** Sistem backend REST API berbasis Go yang mengelola pengaduan fasilitas kampus. Mahasiswa dapat membuat pengaduan, petugas dapat menangani dan memperbarui status, sedangkan admin dapat mengelola seluruh data dan memantau laporan.

---

## 3. Tujuan Sistem

1. Menyediakan platform pelaporan kerusakan fasilitas yang mudah digunakan.
2. Mempercepat respons penanganan pengaduan oleh petugas.
3. Memberikan transparansi status penanganan kepada pelapor.
4. Menyediakan data statistik pengaduan untuk pengambilan keputusan admin.
5. Membangun fondasi sistem yang scalable dan mudah dikembangkan.

---

## 4. Fitur Sistem

| No | Fitur | Deskripsi |
|---|---|---|
| 1 | Autentikasi | Register, login JWT, profil pengguna |
| 2 | Manajemen Pengguna | Admin CRUD pengguna, role assignment |
| 3 | Manajemen Fasilitas | CRUD fasilitas kampus |
| 4 | Manajemen Kategori | CRUD kategori pengaduan |
| 5 | Pengaduan | Buat, lihat, update status pengaduan |
| 6 | Penugasan Petugas | Admin menugaskan petugas ke pengaduan |
| 7 | Riwayat Pengaduan | Rekam perubahan status pengaduan |
| 8 | Laporan | Statistik pengaduan berdasarkan status |

---

## 5. Role Pengguna dan Hak Akses

### ADMIN
- Mengelola pengguna (CRUD)
- Mengelola fasilitas dan kategori
- Melihat seluruh pengaduan
- Mengubah status pengaduan
- Menugaskan petugas
- Mengakses laporan dan statistik

### MAHASISWA
- Registrasi dan login
- Membuat pengaduan baru
- Melihat pengaduan milik sendiri
- Melihat riwayat status pengaduan

### PETUGAS
- Melihat daftar tugas yang diberikan
- Memperbarui status penanganan pengaduan
- Menambahkan catatan penanganan

---

## 6. Rancangan Database

### Tabel: roles

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | UUID | Primary Key |
| name | VARCHAR(50) | Nama role (ADMIN, MAHASISWA, PETUGAS) |
| description | TEXT | Deskripsi role |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diperbarui |

### Tabel: users

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | UUID | Primary Key |
| role_id | UUID | FK → roles.id |
| name | VARCHAR(255) | Nama lengkap |
| email | VARCHAR(255) | Email unik |
| password | TEXT | Password (bcrypt hash) |
| phone | VARCHAR(20) | Nomor telepon |
| student_number | VARCHAR(50) | NIM (nullable) |
| is_active | BOOLEAN | Status akun |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diperbarui |

### Tabel: facilities

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | UUID | Primary Key |
| name | VARCHAR(255) | Nama fasilitas |
| location | VARCHAR(255) | Lokasi |
| description | TEXT | Deskripsi |
| condition | VARCHAR(50) | Kondisi fasilitas |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diperbarui |

### Tabel: categories

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | UUID | Primary Key |
| name | VARCHAR(100) | Nama kategori |
| description | TEXT | Deskripsi |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diperbarui |

### Tabel: complaints

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | UUID | Primary Key |
| user_id | UUID | FK → users.id |
| facility_id | UUID | FK → facilities.id |
| category_id | UUID | FK → categories.id |
| title | VARCHAR(255) | Judul pengaduan |
| description | TEXT | Detail pengaduan |
| priority | VARCHAR(20) | LOW/MEDIUM/HIGH/URGENT |
| status | VARCHAR(30) | Status pengaduan |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diperbarui |

### Tabel: complaint_images

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | UUID | Primary Key |
| complaint_id | UUID | FK → complaints.id |
| image_url | TEXT | URL/path gambar |
| created_at | TIMESTAMP | Waktu diunggah |

### Tabel: assignments

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | UUID | Primary Key |
| complaint_id | UUID | FK → complaints.id |
| officer_id | UUID | FK → users.id (petugas) |
| assigned_by | UUID | FK → users.id (admin) |
| assigned_at | TIMESTAMP | Waktu penugasan |
| notes | TEXT | Catatan penugasan |

### Tabel: complaint_histories

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | UUID | Primary Key |
| complaint_id | UUID | FK → complaints.id |
| user_id | UUID | FK → users.id |
| old_status | VARCHAR(30) | Status sebelumnya |
| new_status | VARCHAR(30) | Status baru |
| notes | TEXT | Catatan perubahan |
| created_at | TIMESTAMP | Waktu perubahan |

---

## 7. ERD (Entity-Relationship Diagram)

```
roles ──< users >──< complaints >──< complaint_images
                         │
                         ├──< assignments
                         │
                         └──< complaint_histories
                         
facilities ──< complaints
categories ──< complaints
```

Relasi:
- `roles` 1:N `users`
- `users` 1:N `complaints` (sebagai pelapor)
- `users` 1:N `assignments` (sebagai petugas)
- `facilities` 1:N `complaints`
- `categories` 1:N `complaints`
- `complaints` 1:N `complaint_images`
- `complaints` 1:N `assignments`
- `complaints` 1:N `complaint_histories`

---

## 8. Arsitektur Clean Architecture

```
┌────────────────────────────────────┐
│          HTTP Layer (route/)        │
│  Fiber callbacks: parse request,   │
│  call service, return response     │
├────────────────────────────────────┤
│        Service Layer (service/)     │
│  Pure business logic               │
│  Input: context + typed params     │
│  Output: domain objects + error    │
├────────────────────────────────────┤
│      Repository Layer (repository/) │
│  Database access with pgx          │
│  Parameterized queries             │
├────────────────────────────────────┤
│        Model Layer (model/)         │
│  Go structs, no business logic     │
└────────────────────────────────────┘
         │            ▲
         ▼            │
┌────────────────────────────────────┐
│          PostgreSQL Database        │
└────────────────────────────────────┘
```

**Prinsip:**
- Service tidak bergantung pada framework HTTP (tidak ada `*fiber.Ctx`).
- Repository hanya berinteraksi dengan database.
- Route hanya bertugas sebagai adapter HTTP.

---

## 9. Teknologi yang Digunakan

| Kategori | Teknologi | Versi |
|---|---|---|
| Bahasa Pemrograman | Go | 1.21+ |
| HTTP Framework | Fiber | v2 |
| Database | PostgreSQL | 14+ |
| Driver DB | pgx | v5 |
| Autentikasi | golang-jwt | v5 |
| Hashing | bcrypt (golang.org/x/crypto) | — |
| UUID | google/uuid | v1 |
| Format | gofmt | bawaan Go |

---

## 10. Struktur Project

```
campuscare-backend/
├── app/
│   ├── model/
│   │   ├── assignment.go
│   │   ├── category.go
│   │   ├── complaint.go
│   │   ├── complaint_history.go
│   │   ├── complaint_image.go
│   │   ├── facility.go
│   │   ├── role.go
│   │   └── user.go
│   ├── repository/
│   │   ├── assignment_repository.go
│   │   ├── category_repository.go
│   │   ├── complaint_history_repository.go
│   │   ├── complaint_repository.go
│   │   ├── facility_repository.go
│   │   ├── role_repository.go
│   │   └── user_repository.go
│   └── service/
│       ├── assignment_service.go
│       ├── auth_service.go
│       ├── category_service.go
│       ├── complaint_service.go
│       ├── complaint_service_test.go
│       ├── facility_service.go
│       ├── report_service.go
│       └── user_service.go
├── config/
├── database/
├── docs/
│   ├── perancangan/
│   ├── report/
│   │   └── PROJECT-REPORT.md
│   ├── git-workflow.md
│   └── team-contribution.md
├── helper/
│   ├── jwt.go
│   ├── password.go
│   └── response.go
├── middleware/
│   ├── authentication.go
│   └── authorization.go
├── migrations/
├── route/
│   └── route.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

---

## 11. Rancangan Endpoint API

### Auth

| Method | Endpoint | Auth | Body | Response |
|---|---|---|---|---|
| POST | `/api/v1/auth/register` | — | `{name, email, password, phone, student_number}` | `{id, name, email, role}` |
| POST | `/api/v1/auth/login` | — | `{email, password}` | `{access_token}` |
| GET | `/api/v1/auth/me` | JWT | — | User profile |

### Complaints

| Method | Endpoint | Auth | Body | Response |
|---|---|---|---|---|
| POST | `/api/v1/my-complaints` | MAHASISWA | `{facility_id, category_id, title, description, priority}` | Complaint |
| GET | `/api/v1/my-complaints` | MAHASISWA | — | List complaint |
| GET | `/api/v1/my-complaints/:id` | MAHASISWA | — | Complaint detail |
| GET | `/api/v1/complaints` | ADMIN | — | All complaints |
| PUT | `/api/v1/complaints/:id/status` | ADMIN | `{status}` | Updated complaint |
| POST | `/api/v1/complaints/:id/assign` | ADMIN | `{officer_id, notes}` | Assignment |

### Assignments

| Method | Endpoint | Auth | Response |
|---|---|---|---|
| GET | `/api/v1/assignments` | PETUGAS | Officer's assignments |
| PUT | `/api/v1/assignments/:id/status` | PETUGAS | Updated complaint |

### Reports

| Method | Endpoint | Auth | Response |
|---|---|---|---|
| GET | `/api/v1/reports/summary` | ADMIN | Summary statistics |

---

## 12. Implementasi Boilerplate

Boilerplate telah diimplementasikan dengan:

1. **Model layer** — 8 entitas: `Role`, `User`, `Facility`, `Category`, `Complaint`, `ComplaintImage`, `Assignment`, `ComplaintHistory`.
2. **Repository layer** — 7 repository dengan interface dan implementasi PostgreSQL menggunakan parameterized query.
3. **Service layer** — 7 service dengan pure business logic (tanpa Fiber).
4. **Route layer** — HTTP callbacks inline yang menghubungkan request ke service.
5. **Middleware** — JWT Authentication dan Role-based Authorization.
6. **Helper** — `GenerateToken`, `HashPassword`, `CheckPassword`, `Success`, `Error`, `Created`.

---

## 13. Pengujian

### Unit Test

File: `app/service/complaint_service_test.go`

| Test | Hasil |
|---|---|
| `TestCreateComplaint_Success` | PASS |
| `TestCreateComplaint_MissingFields` | PASS |
| `TestListMyComplaints` | PASS |
| `TestGetMyComplaint_Forbidden` | PASS |
| `TestUpdateComplaintStatus` | PASS |
| `TestAssignComplaint` | PASS |

### Build & Static Analysis

```
gofmt -w .          → OK (no changes)
go mod tidy         → OK
go vet ./...        → OK (no issues)
go test ./...       → OK (6 tests PASS)
go build ./...      → OK (binary compiled)
```

### Integration Test

Memerlukan koneksi PostgreSQL aktif. Jalankan setelah konfigurasi database:

```bash
go test ./... -tags integration
```

---

## 14. Link Repository GitHub

**URL:** https://github.com/AanFLenrir/campuscare-backend

Repository bersifat **public**. Pastikan pengaturan visibilitas sudah dikonfirmasi melalui GitHub Settings.

---

## 15. Catatan Penting

- File `.env` tidak dimasukkan ke repository (terdapat di `.gitignore`).
- Gunakan `.env.example` sebagai template konfigurasi.
- Credential database dan JWT secret tidak boleh di-commit ke GitHub.
- Kolom anggota kelompok 2–4 harus diisi sesuai data asli sebelum laporan dikumpulkan.

---

*Laporan ini dibuat berdasarkan implementasi aktual pada repository CAMPUSCARE Backend.*
