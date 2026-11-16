# PETRIVA — Laporan Kontribusi Anggota Kelompok

## 1. Identitas Kelompok
- **Nama Tim:** Kelompok 2 — PETRIVA
- **Program Studi:** D4 Teknik Informatika
- **Fakultas:** Fakultas Vokasi
- **Perguruan Tinggi:** Universitas Airlangga
- **Mata Kuliah:** Praktikum Pemrograman Backend / Arsitektur Perangkat Lunak

---

## 2. Susunan Anggota Tim

| No | Nama Lengkap | NIM | Peran Utama |
|---|---|---|---|
| 1 | **MOHAMMAD ANHAZ ABDILAH** | **434241102** | Git Master, Software Architect, Core Backend Engineer |
| 2 | **OKTA ROMAN SANJAYA** | **434241108** | Database Engineer, Service & Repository Developer |
| 3 | **DANU ADJI FERYANTO** | **434241111** | QA Engineer, Middleware, Testing & Swagger UI Documentation |

---

## 3. Rincian Kontribusi Anggota

### 3.1. Mohammad Anhaz Abdilah (NIM: 434241102) — Git Master & Lead Architect
- Merancang arsitektur Clean Architecture tanpa menggunakan controller atau handler.
- Menyiapkan boilerplate sistem, Fiber v2 setup, serta konfigurasi environment dan graceful shutdown pada `main.go`.
- Mengimplementasikan alur otentikasi JWT, hashing bcrypt, serta adapter route HTTP pada `route/route.go`.
- Mengintegrasikan Swagger UI ke dalam backend Go Fiber v2.
- Mengelola Git repository, branching strategy, dan code review integrasi branch.

### 3.2. Okta Roman Sanjaya (NIM: 434241108) — Database & Domain Engineer
- Merancang Entity Relationship Diagram (ERD) dan skema database PostgreSQL untuk seluruh entitas PETRIVA.
- Menyusun DDL migration files (up/down) dari tabel 000001 hingga 000009 lengkap dengan foreign keys dan indexes.
- Mengimplementasikan layer data access (PostgreSQL repositories) menggunakan `pgxpool.Pool` dan parameterized queries.
- Membangun business logic pada `pet_service.go`, `clinic_service.go`, dan `service_service.go`.

### 3.3. Danu Adji Feryanto (NIM: 434241111) — QA & Documentation Engineer
- Mengimplementasikan middleware otentikasi JWT (`authentication.go`) dan otorisasi berbasis peran (`authorization.go`).
- Mengembangkan logika bisnis reservasi konsultasi (`appointment_service.go`), rekam medis (`medical_record_service.go`), dan vaksinasi (`vaccination_service.go`).
- Menyusun suite unit test komprehensif menggunakan mock repository untuk seluruh service, helper, dan middleware.
- Menyusun dokumentasi spesifikasi REST API, ERD, dan panduan instalasi lokal.
