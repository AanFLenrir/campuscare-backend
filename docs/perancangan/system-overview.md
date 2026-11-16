# PETRIVA — System Overview

## 1. Identitas Sistem
- **Nama Sistem:** PETRIVA
- **Kepanjangan:** Pet Care Management System
- **Jenis Aplikasi:** Backend RESTful Web Service
- **Platform:** Go (Fiber v2) + PostgreSQL (pgx v5)
- **Lisensi:** Open-Source Educational Project

---

## 2. Latar Belakang & Masalah
Pemilik hewan peliharaan (*pet owner*) seringkali mengalami kesulitan dalam:
1. Memantau riwayat kesehatan, vaksinasi, dan rekam medis hewan peliharaan mereka secara terorganisir.
2. Melakukan penjadwalan pemeriksaan (*appointment*) ke klinik hewan tanpa kepastian waktu dan dokter penanggung jawab.
3. Terjadinya jadwal bentrok (*conflict schedule*) pada dokter hewan yang sama saat penanganan pasien hewan.
4. Hilangnya buku vaksinasi fisik konvensional yang menyulitkan dokter memeriksa riwayat imunisasi hewan.

Klinik dan dokter hewan juga memerlukan sistem pencatatan rekam medis digital yang terstruktur, aman, dan mudah diakses untuk memberikan diagnosa dan pengobatan yang berkesinambungan.

---

## 3. Tujuan Sistem
1. Menyediakan REST API yang andal untuk manajemen data profil hewan, reservasi konsultasi, rekam medis, dan vaksinasi.
2. Mencegah terjadinya konflik jadwal dokter hewan melalui validasi jadwal otomatis.
3. Menjamin keamanan data rekam medis dengan pembatasan hak akses (*ownership & role-based access control*).
4. Menyediakan dokumentasi API interaktif menggunakan Swagger UI untuk memudahkan integrasi client/frontend.

---

## 4. Ruang Lingkup & Batasan Sistem
### Yang Dicakup:
- Autentikasi JWT dan otorisasi tiga role: `ADMIN`, `PET_OWNER`, dan `VETERINARIAN`.
- Manajemen klinik dan katalog layanan medis/perawatan.
- Pengelolaan profil hewan peliharaan (spesies, ras, jenis kelamin, berat, tanggal lahir).
- Pembuatan dan pengelolaan status janji temu (appointment) dengan pencegahan jadwal bentrok.
- Pencatatan rekam medis digital (diagnosis, gejala, tindakan, dan resep).
- Pencatatan buku vaksinasi digital beserta tanggal vaksin berikutnya.
- Dashboard statistik ringkasan per role.

### Batasan Sistem (*Out of Scope*):
- Tidak mencakup marketplace jual-beli hewan.
- Tidak menyediakan e-commerce pakan atau aksesoris hewan.
- Tidak mengintegrasikan gateway pembayaran online.
- Tidak menggunakan kecerdasan buatan (AI) untuk diagnosis otomatis.
- Tidak terhubung dengan perangkat IoT medis.
