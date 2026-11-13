# System Overview

**Nama Sistem**: CAMPUSCARE (Campus Care – Campus Facility Complaint and Management System)

**Latar Belakang**

Fasilitas kampus seperti gedung, laboratorium, ruang kelas, dan area umum sering mengalami kerusakan atau masalah yang memerlukan penanganan cepat. Proses pelaporan secara manual (melalui kertas atau email) tidak efisien, mengakibatkan keterlambatan penanganan, kurangnya transparansi, dan kesulitan dalam memantau status pengaduan.

**Deskripsi Sistem**

CAMPUSCARE adalah sistem informasi berbasis REST API yang memungkinkan mahasiswa melaporkan kerusakan atau kebutuhan perbaikan fasilitas kampus, petugas menindaklanjuti pengaduan, dan admin mengelola seluruh proses. Sistem ini mencatat setiap langkah penanganan, menyediakan riwayat lengkap, serta menampilkan status secara real‑time.

**Permasalahan yang Diselesaikan**
1. Pelaporan pengaduan yang masih dilakukan secara manual.
2. Mahasiswa tidak dapat mengetahui perkembangan pengaduan mereka.
3. Petugas kekurangan alat pencatatan pekerjaan yang terstruktur.
4. Admin kesulitan memantau seluruh pengaduan secara menyeluruh.
5. Riwayat penanganan belum terdokumentasi dengan baik.

**Tujuan Sistem**
1. Mempermudah mahasiswa melaporkan kerusakan fasilitas.
2. Mempermudah petugas dalam menangani dan melacak pengaduan.
3. Mempermudah admin mengelola, memonitor, dan melaporkan statistik pengaduan.
4. Menyediakan informasi status pengaduan yang transparan.
5. Menyediakan riwayat penanganan yang dapat ditelusuri.

**Manfaat**
- Peningkatan efisiensi penanganan fasilitas.
- Transparansi proses bagi semua pemangku kepentingan.
- Dokumentasi lengkap untuk audit dan perbaikan berkelanjutan.

**Ruang Lingkup**
- Modul **Auth**: Registrasi, login, manajemen profil.
- Modul **User Management**: CRUD pengguna, manajemen peran.
- Modul **Facility & Category Management**: Pengelolaan data fasilitas dan kategori pengaduan.
- Modul **Complaint Management**: Pembuatan, pembaruan, penugasan, dan pelacakan pengaduan.
- Modul **Reporting**: Statistik dan laporan pengaduan untuk admin.

**Teknologi yang Digunakan**
- Bahasa pemrograman: **Go 1.26**
- Framework HTTP: **Fiber**
- Database: **PostgreSQL** dengan driver **pgx**
- Authentication: **JWT**
- Password hashing: **bcrypt**
- Migration: **golang-migrate**
- Versi kontrol: **Git** pada **GitHub**

---

*Dokumen ini disusun dalam bahasa Indonesia formal dan mudah dipahami.*
