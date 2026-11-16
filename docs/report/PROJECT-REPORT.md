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
<<<<<<< Updated upstream
| 2 | OKTA ROMAN SANJAYA | 434241108 | Backend Developer (Routing & Authentication) |
| 3 | DANU ADJI FERIANTO | 434241111 | QA & Backend Testing Engineer |
=======
| 2 | *(Anggota 2)* | — | *(isi sesuai data kelompok)* |
| 3 | *(Anggota 3)* | — | *(isi sesuai data kelompok)* |
>>>>>>> Stashed changes

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
│       ├── auth_service_test.go
│       ├── category_service.go
│       ├── complaint_service.go
│       ├── complaint_service_test.go
│       ├── facility_service.go
│       ├── report_service.go
│       ├── user_service.go
│       └── user_service_test.go
├── config/
│   └── config.go
├── database/
│   └── postgres.go
├── docs/
│   ├── perancangan/
│   ├── report/
│   │   └── PROJECT-REPORT.md
│   ├── git-workflow.md
│   └── team-contribution.md
├── helper/
│   ├── jwt.go
│   ├── jwt_test.go
│   ├── password.go
│   ├── password_test.go
│   └── response.go
├── middleware/
│   ├── authentication.go
│   ├── authorization.go
│   └── middleware_test.go
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

## 13. Instalasi dan Menjalankan Project

Untuk menjalankan backend CAMPUSCARE, diperlukan **Go (versi 1.21 ke atas)** dan **PostgreSQL (versi 14 ke atas)** yang telah terpasang pada komputer pengembang.

### Tahapan Umum Menjalankan Project:
1. Melakukan clone repository GitHub ke lokal.
2. Masuk ke direktori project.
3. Menyiapkan file konfigurasi environment (`.env`) berdasarkan template `.env.example`.
4. Menyiapkan database PostgreSQL dan menjalankan script migrasi.
5. Mengunduh dependency Go menggunakan `go mod tidy`.
6. Menjalankan aplikasi backend dengan `go run main.go`.

### Perintah Dasar yang Digunakan:

```bash
# 1. Clone repository
git clone https://github.com/AanFLenrir/campuscare-backend.git

# 2. Masuk ke direktori
cd campuscare-backend

# 3. Salin environment konfigurasi
cp .env.example .env

# 4. Unduh dependency Go
go mod tidy

# 5. Jalankan backend service
go run main.go
```

> **Catatan:** Nama file entry point adalah `main.go`. Konfigurasi environment pada file `.env` disesuaikan dengan kredensial database lokal (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, dan JWT_SECRET).

---

## 14. Pengelolaan Project Menggunakan GitHub

Pengembangan CAMPUSCARE menggunakan Git dan GitHub sebagai sarana pengelolaan kode sumber, pelacakan revisi, serta dokumentasi proyek secara kolaboratif antaranggota kelompok.

### Repository Utama:
- **URL Repository:** https://github.com/AanFLenrir/campuscare-backend
- **Visibilitas:** Public

### Struktur Branch & Alur Kolaborasi (Git Workflow):
- `main`: Branch produksi yang selalu siap dieksekusi dan stabil. Seluruh integrasi akhir digabungkan (*merge*) ke branch ini.
- Konvensi pesan commit mengikuti standar **Conventional Commits**:
  - `feat(<scope>)`: Penambahan fitur baru.
  - `test(<scope>)`: Penambahan atau pembaruan unit/integration test.
  - `docs`: Pembaruan dokumentasi atau laporan.
  - `refactor`: Perbaikan struktur kode tanpa mengubah fungsionalitas.
  - `chore`: Konfigurasi dependensi atau konfigurasi *build*.

Setiap anggota kelompok berkontribusi menggunakan akun GitHub masing-masing sesuai pembagian peran kerja.

### Riwayat Commit GitHub

```
[MASUKKAN SCREENSHOT RIWAYAT COMMIT GITHUB DI SINI]
```
*(Screenshot riwayat commit diambil dari halaman GitHub: https://github.com/AanFLenrir/campuscare-backend/commits/main yang menampilkan kontribusi nyata dari setiap anggota kelompok).*

---

## 15. Pengujian Sistem

Pengujian dilakukan untuk memastikan seluruh komponen backend dapat berjalan sesuai fungsinya masing-masing, bebas dari regresi, dan siap dihubungkan dengan frontend.

### 15.1 Unit Testing
Unit testing digunakan untuk menguji fungsionalitas murni secara terisolasi tanpa ketergantungan database eksternal menggunakan mock repository. Pengujian berfokus pada:
- **AuthService (`app/service/auth_service_test.go`):** Validasi registrasi akun mahasiswa, validasi input wajib kosong, proteksi email duplikat, verifikasi password hash saat login, penolakan kredensial salah, dan penolakan akun yang berstatus nonaktif.
- **UserService (`app/service/user_service_test.go`):** Pengujian CRUD pengguna, validasi field wajib, pencarian ID ditemukan/tidak ditemukan, penyembunyian password hash dari respons, pembaruan data user, dan penghapusan akun.
- **ComplaintService (`app/service/complaint_service_test.go`):** Validasi pembuatan pengaduan, verifikasi field mandatory, hak akses kepemilikan pengaduan, pembaruan status, dan penugasan petugas oleh admin.
- **Helper Testing (`helper/password_test.go` & `helper/jwt_test.go`):** Keandalan fungsi *bcrypt hashing* dan validasi token JWT.

### 15.2 Integration Testing
Integration testing digunakan untuk memastikan hubungan antarbagian aplikasi, terutama interaksi endpoint HTTP Fiber dengan lapisan service, serta komunikasi backend dengan database PostgreSQL menggunakan parameterized query.

### 15.3 Pengujian Middleware
Pengujian middleware (`middleware/middleware_test.go`) dilakukan dengan memanggil instance HTTP Fiber menggunakan `httptest` guna memastikan:
1. Header `Authorization` wajib disertakan (menghasilkan status HTTP 401 jika kosong).
2. Format header harus sesuai format `Bearer <token>` (menghasilkan status HTTP 401 jika format salah).
3. Token JWT yang kedaluwarsa atau memiliki *signature* tidak valid ditolak secara otomatis.
4. Token valid mengizinkan akses ke endpoint dan menyimpan konteks `userID` serta `role`.
5. Hak akses berbasis role (*Role-Based Authorization*) dibatasi secara ketat (misal: role *MAHASISWA* dilarang mengakses endpoint khusus *ADMIN* dengan status HTTP 403 Forbidden).

### 15.4 Pemeriksaan Kode
Pemeriksaan kode statis dan build dilakukan dengan urutan perintah berikut:

```bash
gofmt -w .
go mod tidy
go vet ./...
go test -v ./...
go build ./...
```

### Tabel Hasil Pengujian Aktual (31 Test Cases — 100% PASS)

| No | Modul / Package | File Pengujian | Nama Test Function | Hasil |
|---|---|---|---|---|
| 1 | `app/service` | `auth_service_test.go` | `TestRegister_Success` | **PASS** |
| 2 | `app/service` | `auth_service_test.go` | `TestRegister_MissingFields` | **PASS** |
| 3 | `app/service` | `auth_service_test.go` | `TestRegister_DuplicateEmail` | **PASS** |
| 4 | `app/service` | `auth_service_test.go` | `TestLogin_Success` | **PASS** |
| 5 | `app/service` | `auth_service_test.go` | `TestLogin_InvalidCredentials` | **PASS** |
| 6 | `app/service` | `auth_service_test.go` | `TestLogin_AccountDisabled` | **PASS** |
| 7 | `app/service` | `complaint_service_test.go` | `TestCreateComplaint_Success` | **PASS** |
| 8 | `app/service` | `complaint_service_test.go` | `TestCreateComplaint_MissingFields` | **PASS** |
| 9 | `app/service` | `complaint_service_test.go` | `TestListMyComplaints` | **PASS** |
| 10 | `app/service` | `complaint_service_test.go` | `TestGetMyComplaint_Forbidden` | **PASS** |
| 11 | `app/service` | `complaint_service_test.go` | `TestUpdateComplaintStatus` | **PASS** |
| 12 | `app/service` | `complaint_service_test.go` | `TestAssignComplaint` | **PASS** |
| 13 | `app/service` | `user_service_test.go` | `TestCreateUser_Success` | **PASS** |
| 14 | `app/service` | `user_service_test.go` | `TestCreateUser_MissingFields` | **PASS** |
| 15 | `app/service` | `user_service_test.go` | `TestGetUser_Success` | **PASS** |
| 16 | `app/service` | `user_service_test.go` | `TestGetUser_NotFound` | **PASS** |
| 17 | `app/service` | `user_service_test.go` | `TestListUsers_Success` | **PASS** |
| 18 | `app/service` | `user_service_test.go` | `TestUpdateUser_Success` | **PASS** |
| 19 | `app/service` | `user_service_test.go` | `TestDeleteUser_Success` | **PASS** |
| 20 | `helper` | `jwt_test.go` | `TestGenerateAndValidateToken_Success` | **PASS** |
| 21 | `helper` | `jwt_test.go` | `TestValidateToken_InvalidSecret` | **PASS** |
| 22 | `helper` | `jwt_test.go` | `TestValidateToken_Expired` | **PASS** |
| 23 | `helper` | `password_test.go` | `TestHashPassword_Success` | **PASS** |
| 24 | `helper` | `password_test.go` | `TestCheckPassword_Correct` | **PASS** |
| 25 | `helper` | `password_test.go` | `TestCheckPassword_Incorrect` | **PASS** |
| 26 | `middleware` | `middleware_test.go` | `TestAuthMiddleware_MissingHeader` | **PASS** |
| 27 | `middleware` | `middleware_test.go` | `TestAuthMiddleware_InvalidHeaderFormat` | **PASS** |
| 28 | `middleware` | `middleware_test.go` | `TestAuthMiddleware_InvalidToken` | **PASS** |
| 29 | `middleware` | `middleware_test.go` | `TestAuthMiddleware_Success` | **PASS** |
| 30 | `middleware` | `middleware_test.go` | `TestAuthorize_AdminAllowed` | **PASS** |
| 31 | `middleware` | `middleware_test.go` | `TestAuthorize_ForbiddenRole` | **PASS** |
| 32 | `middleware` | `middleware_test.go` | `TestAuthorize_MultipleAllowedRoles` | **PASS** |

### Output Eksekusi Terminal (`go test -v ./...`):

```text
=== RUN   TestRegister_Success
--- PASS: TestRegister_Success (0.12s)
=== RUN   TestRegister_MissingFields
--- PASS: TestRegister_MissingFields (0.00s)
=== RUN   TestRegister_DuplicateEmail
--- PASS: TestRegister_DuplicateEmail (0.12s)
=== RUN   TestLogin_Success
--- PASS: TestLogin_Success (0.23s)
=== RUN   TestLogin_InvalidCredentials
--- PASS: TestLogin_InvalidCredentials (0.32s)
=== RUN   TestLogin_AccountDisabled
--- PASS: TestLogin_AccountDisabled (0.15s)
=== RUN   TestCreateComplaint_Success
--- PASS: TestCreateComplaint_Success (0.00s)
=== RUN   TestCreateComplaint_MissingFields
--- PASS: TestCreateComplaint_MissingFields (0.00s)
=== RUN   TestListMyComplaints
--- PASS: TestListMyComplaints (0.00s)
=== RUN   TestGetMyComplaint_Forbidden
--- PASS: TestGetMyComplaint_Forbidden (0.00s)
=== RUN   TestUpdateComplaintStatus
--- PASS: TestUpdateComplaintStatus (0.00s)
=== RUN   TestAssignComplaint
--- PASS: TestAssignComplaint (0.00s)
=== RUN   TestCreateUser_Success
--- PASS: TestCreateUser_Success (0.16s)
=== RUN   TestCreateUser_MissingFields
--- PASS: TestCreateUser_MissingFields (0.00s)
=== RUN   TestGetUser_Success
--- PASS: TestGetUser_Success (0.00s)
=== RUN   TestGetUser_NotFound
--- PASS: TestGetUser_NotFound (0.00s)
=== RUN   TestListUsers_Success
--- PASS: TestListUsers_Success (0.00s)
=== RUN   TestUpdateUser_Success
--- PASS: TestUpdateUser_Success (0.00s)
=== RUN   TestDeleteUser_Success
--- PASS: TestDeleteUser_Success (0.00s)
PASS
ok      github.com/AanFLenrir/campuscare-backend/app/service    5.436s
=== RUN   TestGenerateAndValidateToken_Success
--- PASS: TestGenerateAndValidateToken_Success (0.00s)
=== RUN   TestValidateToken_InvalidSecret
--- PASS: TestValidateToken_InvalidSecret (0.00s)
=== RUN   TestValidateToken_Expired
--- PASS: TestValidateToken_Expired (0.00s)
=== RUN   TestHashPassword_Success
--- PASS: TestHashPassword_Success (0.12s)
=== RUN   TestCheckPassword_Correct
--- PASS: TestCheckPassword_Correct (0.29s)
=== RUN   TestCheckPassword_Incorrect
--- PASS: TestCheckPassword_Incorrect (0.24s)
PASS
ok      github.com/AanFLenrir/campuscare-backend/helper 4.606s
=== RUN   TestAuthMiddleware_MissingHeader
--- PASS: TestAuthMiddleware_MissingHeader (0.00s)
=== RUN   TestAuthMiddleware_InvalidHeaderFormat
--- PASS: TestAuthMiddleware_InvalidHeaderFormat (0.00s)
=== RUN   TestAuthMiddleware_InvalidToken
--- PASS: TestAuthMiddleware_InvalidToken (0.00s)
=== RUN   TestAuthMiddleware_Success
--- PASS: TestAuthMiddleware_Success (0.00s)
=== RUN   TestAuthorize_AdminAllowed
--- PASS: TestAuthorize_AdminAllowed (0.00s)
=== RUN   TestAuthorize_ForbiddenRole
--- PASS: TestAuthorize_ForbiddenRole (0.00s)
=== RUN   TestAuthorize_MultipleAllowedRoles
--- PASS: TestAuthorize_MultipleAllowedRoles (0.00s)
PASS
ok      github.com/AanFLenrir/campuscare-backend/middleware     4.382s
```

```
[MASUKKAN SCREENSHOT HASIL PENGUJIAN DI SINI]
```

---

## 16. Catatan Penting

- File `.env` tidak dimasukkan ke repository (terdapat di `.gitignore`).
- Gunakan `.env.example` sebagai template konfigurasi.
- Credential database dan JWT secret tidak boleh di-commit ke GitHub.
- Seluruh pengujian wajib diverifikasi ulang sebelum kode digabungkan ke branch utama (*main*).

---

*Laporan ini dibuat berdasarkan implementasi aktual pada repository CAMPUSCARE Backend.*
