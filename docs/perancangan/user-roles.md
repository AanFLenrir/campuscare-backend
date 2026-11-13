## User Roles

| Role      | Deskripsi                                                                 | Hak Akses Utama                                                                                                                                                     |
|-----------|---------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **ADMIN** | Pengelola sistem dengan hak penuh.                                        | - Mengelola pengguna, peran, fasilitas, kategori, dan seluruh pengaduan.<br>- Menugaskan petugas.<br>- Mengakses laporan statistik dan riwayat semua pengaduan. |
| **MAHASISWA** | Pengguna mahasiswa yang dapat melaporkan masalah fasilitas.               | - Registrasi, login, melihat/menambah/menampilkan pengaduannya sendiri.<br>- Mengunggah foto pendukung.<br>- Melihat status dan riwayat penanganan pengaduannya. |
| **PETUGAS** | Petugas operasional yang menangani pengaduan yang ditugaskan kepadanya. | - Melihat pengaduan yang ditugaskan.<br>- Memperbarui status pengaduan.<br>- Menambahkan catatan penanganan.<br>- Menandai pengaduan selesai. |

**Catatan:**
- Hak akses di‑implementasikan pada middleware Authorization yang memeriksa nilai `role` pada token JWT.
- Setiap endpoint menegaskan role yang diperbolehkan di dalam kode.

