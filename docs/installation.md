# PETRIVA — Installation & Local Setup Guide

Panduan langkah demi langkah untuk menjalankan backend PETRIVA di lingkungan lokal (Windows / Linux / macOS).

---

## 1. Prasyarat Lingkungan
Pastikan perangkat lunak berikut telah terpasang:
- **Go:** Versi 1.21 atau lebih baru (`go version`)
- **PostgreSQL:** Versi 14 atau lebih baru (tersedia pada paket Laragon / Docker / Standalone)
- **Git:** Versi 2.30+ (`git version`)

---

## 2. Langkah Setup Proyek

### Langkah 1: Kloning / Buka Workspace
Pastikan Anda berada di direktori proyek:
```powershell
cd C:\laragon\www\campuscare-backend
```

### Langkah 2: Konfigurasi File Lingkungan (.env)
Salin template konfigurasi:
```powershell
Copy-Item .env.example .env
```

Buka file `.env` dan sesuaikan parameter berikut:
```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password_postgres_anda
DB_NAME=petriva
DB_SSLMODE=disable

JWT_SECRET=kunci_rahasia_jwt_petriva_2026
JWT_EXPIRE_HOUR=24

APP_PORT=3000
APP_ENV=development
```

> **Catatan:** Backend PETRIVA dilengkapi modul parser `.env` otomatis pada `config/config.go` sehingga variabel akan otomatis dimuat tanpa perlu konfigurasi environment Windows manual.

---

## 3. Setup Database & Migrasi

### Langkah 1: Buat Database di PostgreSQL
Buka terminal PowerShell atau HeidiSQL:
```sql
CREATE DATABASE petriva;
```

### Langkah 2: Eksekusi Migrasi DDL & Seed Data
Jalankan file SQL migrasi dari direktori `migrations/` secara berurutan:
```powershell
$psql = "C:\laragon\bin\postgresql\pgsql\bin\psql.exe" # Path psql Laragon

& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000001_create_roles_table.up.sql
& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000002_create_clinics_table.up.sql
& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000003_create_users_table.up.sql
& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000004_create_pets_table.up.sql
& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000005_create_services_table.up.sql
& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000006_create_appointments_table.up.sql
& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000007_create_medical_records_table.up.sql
& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000008_create_vaccinations_table.up.sql
& $psql -h 127.0.0.1 -U postgres -d petriva -f migrations\000009_seed_petriva.up.sql
```

---

## 4. Menjalankan Server Backend

```powershell
go run main.go
```

Output yang diharapkan pada konsol:
```
2026/11/16 23:00:00 PETRIVA server listening on port :3000
 ┌───────────────────────────────────────────────────┐
 │                   Fiber v2.52.6                   │
 │               http://127.0.0.1:3000               │
 │       (bound on host 0.0.0.0 and port 3000)       │
 └───────────────────────────────────────────────────┘
```

---

## 5. Verifikasi Ketersediaan Layanan

### Uji Health Check via PowerShell:
```powershell
Invoke-RestMethod -Uri "http://localhost:3000/api/v1/health" -Method Get
```

Response JSON:
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "status": "ok",
    "system": "PETRIVA - Pet Care Management System"
  }
}
```

### Akses Dokumentasi Swagger UI:
Buka browser dan navigasikan ke:
`http://localhost:3000/swagger/`

