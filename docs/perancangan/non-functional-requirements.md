## Persyaratan Non‑Functional

1. **Keamanan**
   - Password disimpan dalam bentuk hash menggunakan bcrypt dengan cost default (`golang.org/x/crypto/bcrypt`).
   - Semua endpoint yang membutuhkan otorisasi harus menyertakan header `Authorization: Bearer <token>`.
   - JWT ditandatangani menggunakan RSA atau HMAC‑SHA256 (disimpan dalam environment variable `JWT_SECRET`).
   - Transport Layer Security (HTTPS) wajib dipakai pada lingkungan produksi (tidak di‑handle di kode, tetapi tercantum dalam dokumentasi).
2. **Kinerja**
   - Koneksi ke PostgreSQL dikelola dengan connection pool (`pgxpool`).
   - Rute di‑optimalkan dengan Fiber yang bersifat non‑blocking.
   - Index utama dibuat pada kolom yang sering dipakai pencarian (`email`, `role_id`, `facility_id`, `category_id`, `status`).
3. **Maintainability**
   - Menggunakan Clean Architecture sehingga kode terpisah menjadi entitas, use‑case (service), dan adapter (repository, framework).
   - Penamaan variabel dan fungsi menggunakan bahasa Inggris konsisten.
   - Dokumentasi inline menggunakan komentar GoDoc.
4. **Scalability**
   - Aplikasi stateless; dapat dijalankan pada banyak instance di belakang load balancer.
   - Database dapat di‑scale secara vertikal atau horizontal (partitioning pada tabel `complaints` bila diperlukan).
5. **Reliability & Availability**
   - Koneksi database memiliki retry otomatis pada saat inisialisasi.
   - Health check endpoint (`/api/v1/health`) mengembalikan status `200 OK` bila aplikasi dapat terhubung ke database.
6. **Integritas Data**
   - Menggunakan transaksi untuk operasi yang melibatkan lebih dari satu tabel (mis. assignment + status update).
   - `ON DELETE RESTRICT` pada foreign key untuk mencegah kehilangan data penting.
7. **Kepatuhan**
   - Tidak ada data sensitif (password, JWT secret) yang disimpan dalam repository.
   - `.env` dan file credential di‑exclude oleh `.gitignore`.
8. **Code Readability**
   - `go fmt`, `go vet`, dan `golint` (jika tersedia) dijalankan pada setiap commit.
   - Struktur folder mengikuti konvensi Clean Architecture yang telah ditetapkan.

*Dokumen ini disusun dalam bahasa Indonesia formal.*
