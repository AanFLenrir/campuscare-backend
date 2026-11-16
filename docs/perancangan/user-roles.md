# PETRIVA — User Roles & Permissions Matrix

## 1. Definisi Tiga Role Utama

### 1. ADMIN
Administrator sistem yang mengelola entitas master dan memiliki akses supervisi global:
- Mengelola data master klinik dan layanan klinik.
- Mengelola akun pengguna (membuat dokter hewan dengan klinik terdaftar, mengelola admin, aktivasi/deaktivasi user).
- Memantau seluruh jadwal janji temu (appointment) dari semua klinik.
- Melihat statistik dashboard global.

### 2. PET_OWNER
Pemilik hewan peliharaan:
- Mendaftarkan akun secara mandiri dan melakukan login.
- Mengelola profil hewan peliharaan miliknya (tambah, perbarui, nonaktifkan).
- Melihat katalog klinik, dokter hewan, dan layanan yang tersedia.
- Membuat jadwal appointment untuk hewannya.
- Membatalkan appointment yang masih berstatus `PENDING` atau `CONFIRMED`.
- Melihat rekam medis dan riwayat vaksinasi hewannya.

### 3. VETERINARIAN
Dokter hewan terdaftar yang bertugas di klinik tertentu:
- Melihat jadwal appointment yang ditugaskan kepadanya.
- Memperbarui status penanganan (`CONFIRMED`, `IN_PROGRESS`, `COMPLETED`, `REJECTED`).
- Mengisi rekam medis hasil pemeriksaan (diagnosis, gejala, tindakan, resep).
- Mencatat riwayat vaksinasi hewan.
- Melihat riwayat penanganan pasien hewan sebelumnya.

---

## 2. Matriks Hak Akses (RBAC Matrix)

| Modul & Aksi | Endpoint | Publik | ADMIN | PET_OWNER | VETERINARIAN |
|---|---|:---:|:---:|:---:|:---:|
| **Auth: Register** | `POST /auth/register` | ✅ | ❌ | ❌ | ❌ |
| **Auth: Login** | `POST /auth/login` | ✅ | ✅ | ✅ | ✅ |
| **Auth: Profile** | `GET /auth/me` | ❌ | ✅ | ✅ | ✅ |
| **Dashboard** | `GET /dashboard/summary` | ❌ | ✅ (Global) | ✅ (Personal) | ✅ (Assigned) |
| **Klinik: Baca** | `GET /clinics`, `/:id` | ✅ | ✅ | ✅ | ✅ |
| **Klinik: Kelola** | `POST, PUT, DELETE /clinics` | ❌ | ✅ | ❌ | ❌ |
| **Layanan: Baca** | `GET /services`, `/:id` | ✅ | ✅ | ✅ | ✅ |
| **Layanan: Kelola** | `POST, PUT, DELETE /services` | ❌ | ✅ | ❌ | ❌ |
| **Dokter: Baca** | `GET /veterinarians`, `/:id` | ✅ | ✅ | ✅ | ✅ |
| **User: Kelola** | `GET, POST, PUT, DELETE /users` | ❌ | ✅ | ❌ | ❌ |
| **Pet: Tambah** | `POST /pets` | ❌ | ❌ | ✅ | ❌ |
| **Pet: Baca** | `GET /pets`, `/:id` | ❌ | ✅ (Semua) | ✅ (Milik Sendiri) | ✅ (Pasien) |
| **Pet: Update/Hapus**| `PUT, DELETE /pets/:id` | ❌ | ❌ | ✅ (Milik Sendiri) | ❌ |
| **Appt: Buat** | `POST /appointments` | ❌ | ❌ | ✅ | ❌ |
| **Appt: Lihat** | `GET /appointments`, `/:id` | ❌ | ✅ (Semua) | ✅ (Milik Sendiri) | ✅ (Ditugaskan) |
| **Appt: Batalkan** | `PUT /appointments/:id/cancel` | ❌ | ✅ | ✅ (Milik Sendiri) | ❌ |
| **Appt: Ubah Status**| `PUT /appointments/:id/status` | ❌ | ✅ | ❌ | ✅ (Ditugaskan) |
| **Medis: Catat** | `POST /medical-records` | ❌ | ❌ | ❌ | ✅ (Ditugaskan) |
| **Medis: Lihat** | `GET /medical-records`, `/:id`| ❌ | ✅ (Semua) | ✅ (Milik Sendiri) | ✅ (Ditangani) |
| **Vaksin: Catat** | `POST /vaccinations` | ❌ | ❌ | ❌ | ✅ |
| **Vaksin: Lihat** | `GET /vaccinations`, `/:id` | ❌ | ✅ (Semua) | ✅ (Milik Sendiri) | ✅ |
| **Vaksin: Hapus** | `DELETE /vaccinations/:id` | ❌ | ✅ | ❌ | ✅ |
