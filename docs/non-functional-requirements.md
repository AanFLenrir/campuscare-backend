# PETRIVA — Non-Functional Requirements

## 1. Keamanan (Security)
- **NFR-SEC-01 (Password Hashing):** Seluruh password pengguna disimpan dalam bentuk hash satu arah menggunakan algoritma `bcrypt` dengan work factor (cost) minimal 10. Tidak ada password plaintext yang disimpan di database atau dikembalikan dalam respon API.
- **NFR-SEC-02 (Stateless Authentication):** Autentikasi mengimplementasikan standar JSON Web Token (JWT) dengan algoritma HMAC-SHA256 (HS256). Token disematkan dalam HTTP Authorization header dengan skema `Bearer`.
- **NFR-SEC-03 (Role-Based Access Control):** Setiap endpoint bisnis dilindungi oleh middleware yang memvalidasi token dan hak akses role (`ADMIN`, `PET_OWNER`, `VETERINARIAN`).
- **NFR-SEC-04 (SQL Injection Prevention):** Seluruh query database ke PostgreSQL menggunakan parameterized query (`$1, $2, ...`) yang didukung penuh oleh library `pgx/v5`.
- **NFR-SEC-05 (Ownership Verification):** Sistem tidak mempercayai `owner_id` dari client request body/query; identitas pemilik diekstrak langsung dari payload token JWT pengguna terautentikasi.
- **NFR-SEC-06 (Sensitive Data Concealment):** Response JSON menggunakan struct tag `json:"-"` pada kolom password sehingga hash tidak pernah bocor ke jaringan.

## 2. Kinerja & Efisiensi (Performance)
- **NFR-PERF-01 (High Throughput):** Menggunakan framework Go Fiber v2 yang dibangun di atas FastHTTP untuk memaksimalkan throughput dan meminimalkan alokasi memori (*zero memory allocation*).
- **NFR-PERF-02 (Database Connection Pooling):** Koneksi ke PostgreSQL dikelola melalui pool koneksi `pgxpool.Pool` dengan timeout koneksi 5 detik pada inisialisasi awal.
- **NFR-PERF-03 (Indexing Strategy):** Indeks dibuat pada kolom-kolom pencarian frekuensi tinggi: foreign keys, email, status, dan tanggal appointment.

## 3. Keandalan & Integritas Data (Reliability & Integrity)
- **NFR-REL-01 (Referential Integrity):** Seluruh relasi tabel PostgreSQL dilindungi oleh Foreign Key Constraints dengan aksi `ON UPDATE CASCADE` dan `ON DELETE RESTRICT` untuk mencegah data yatim (*orphaned data*).
- **NFR-REL-02 (Graceful Shutdown):** Aplikasi mengimplementasikan penanganan sinyal OS (`SIGINT`, `SIGTERM`) untuk menutup pool koneksi database dan menyelesaikan request yang sedang berjalan sebelum aplikasi berhenti.
- **NFR-REL-03 (Soft Deletion):** Penghapusan profil hewan peliharaan menerapkan status *soft delete* (`is_active = false`) sehingga riwayat rekam medis dan vaksinasi di masa lampau tetap utuh.

## 4. Keterpeliharaan & Skalabilitas (Maintainability & Scalability)
- **NFR-MAINT-01 (Clean Architecture):** Kode terstruktur rapi ke dalam layer yang independen tanpa adanya layer controller/handler yang dilarang.
- **NFR-MAINT-02 (Automated Testing):** Semua logika bisnis dilindungi oleh unit test berbasis mock repository yang dapat dijalankan secara cepat dan independen tanpa dependensi database eksternal.
- **NFR-MAINT-03 (Living Documentation):** Dokumentasi API disinkronisasi menggunakan Swag CLI dan dapat diuji langsung melalui Swagger UI.

