## Functional Requirements

### Authentication
- **Registrasi**: Mahasiswa dapat mendaftar dengan nama, email, password, nomor mahasiswa, dan nomor telepon. Password disimpan dalam bentuk hash (bcrypt).
- **Login**: Menghasilkan token JWT yang berisi `user_id` dan `role` dengan waktu kedaluwarsa (1 hari).
- **Logout**: Token dapat di‑invalidasi pada client (server stateless, tidak ada endpoint khusus).
- **Pengambilan profil**: Endpoint untuk mendapatkan data profil pengguna yang sedang terautentikasi.
- **Perubahan profil**: Pengguna dapat memperbarui nama, email, nomor telepon, dan foto profil.
- **Pengamanan password**: Password harus memenuhi minimal 8 karakter, mengandung huruf besar, huruf kecil, angka, dan simbol.
- **Autentikasi menggunakan JWT**: Setiap request ke endpoint yang dilindungi harus menyertakan header `Authorization: Bearer <token>`.

### User Management (Admin)
- **Melihat pengguna**: Daftar seluruh pengguna dengan filter role dan status.
- **Menambahkan pengguna**: Admin dapat menambah pengguna baru (misalnya petugas) tanpa proses registrasi mahasiswa.
- **Mengubah data pengguna**: Update nama, email, role, status aktif/non‑aktif.
- **Menghapus atau menonaktifkan pengguna**: Soft delete dengan flag `is_active`.
- **Pengelolaan role**: Assign role (admin, mahasiswa, petugas) ke pengguna.

### Complaint Management
- **Membuat pengaduan** (Mahasiswa): Mengirim judul, deskripsi, kategori, fasilitas, prioritas, dan optional foto.
- **Melihat pengaduan**: Semua pengguna dapat melihat daftar pengaduan; mahasiswa hanya melihat miliknya, petugas melihat yang ditugaskan, admin melihat semua.
- **Mengubah pengaduan**: Mahasiswa dapat mengubah pengaduan miliknya selama status belum `ASSIGNED`.
- **Menghapus pengaduan**: Mahasiswa dapat meng‑delete (soft) pengaduan sebelum diverifikasi.
- **Melihat detail pengaduan**: Menampilkan semua informasi termasuk gambar.
- **Mengunggah gambar pendukung**: Menyimpan file gambar pada server (path) dan menyimpan referensi di tabel `complaint_images`.
- **Melihat riwayat pengaduan**: Daftar perubahan status dan catatan.

### Category Management (Admin)
- **Menambahkan kategori**
- **Melihat kategori**
- **Mengubah kategori**
- **Menghapus atau menonaktifkan kategori**

### Facility Management (Admin)
- **Menambahkan fasilitas** (nama, lokasi, deskripsi, kondisi)
- **Melihat fasilitas**
- **Mengubah fasilitas**
- **Menghapus atau menonaktifkan fasilitas**

### Assignment Management (Admin)
- **Menugaskan petugas** ke suatu pengaduan.
- **Petugas melihat tugas** yang ditugaskan kepadanya.
- **Admin melihat daftar penugasan** seluruhnya.
- **Mencatat waktu penugasan**.

### Complaint Tracking
- **Mahasiswa melihat status pengaduan** (REAL‑TIME melalui endpoint).
- **Petugas memperbarui status** (`IN_PROGRESS`, `RESOLVED`, dll.).
- **Sistem mencatat setiap perubahan status** ke tabel `complaint_histories`.
- **Mahasiswa dapat melihat riwayat penanganan**.

### Report Management (Admin)
- **Melihat jumlah pengaduan** keseluruhan.
- **Melihat pengaduan berdasarkan status**.
- **Melihat pengaduan berdasarkan kategori**.
- **Melihat laporan penyelesaian** (jumlah selesai per periode).

### Non‑Functional Requirements (see separate document)
