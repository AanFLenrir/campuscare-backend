## API Design

### Public Endpoints
| Metode | Path | Deskripsi | Auth |
|--------|------|-----------|------|
| GET | `/api/v1/health` | Memeriksa kesehatan aplikasi dan koneksi database. | Tidak |
| POST | `/api/v1/auth/register` | Registrasi mahasiswa baru. Body JSON: `{name,email,password,phone,student_number}` | Tidak |
| POST | `/api/v1/auth/login` | Login, menghasilkan JWT. Body JSON: `{email,password}` | Tidak |

### Authenticated Endpoints (Semua memerlukan header `Authorization: Bearer <token>`)
| Metode | Path | Deskripsi | Role |
|--------|------|-----------|------|
| GET | `/api/v1/auth/me` | Informasi profil pengguna yang sedang terautentikasi. | Semua |
| GET | `/api/v1/users` | Daftar semua pengguna (admin only). | ADMIN |
| GET | `/api/v1/facilities` | Daftar fasilitas (admin & mahasiswa dapat melihat). | Semua |
| GET | `/api/v1/categories` | Daftar kategori pengaduan. | Semua |
| GET | `/api/v1/complaints` | Daftar semua pengaduan (admin). | ADMIN |
| GET | `/api/v1/complaints/:id` | Detail pengaduan (admin, petugas yang ditugaskan, atau mahasiswa pemilik). | ADMIN/ PETUGAS/ MAHASISWA |
| POST | `/api/v1/complaints` | Membuat pengaduan baru (mahasiswa). | MAHASISWA |
| GET | `/api/v1/my-complaints` | Daftar pengaduan milik mahasiswa yang sedang login. | MAHASISWA |
| POST | `/api/v1/facilities` | Tambah fasilitas baru. | ADMIN |
| POST | `/api/v1/categories` | Tambah kategori baru. | ADMIN |
| POST | `/api/v1/complaints/:id/assign` | Menugaskan petugas ke pengaduan. Body: `{officer_id, notes}` | ADMIN |
| PATCH | `/api/v1/complaints/:id/status` | Update status pengaduan (petugas). Body: `{status, notes}` | PETUGAS |

**Catatan:**
- Semua request body dan response menggunakan JSON.
- Validasi sederhana dilakukan pada layer Service.
- Setiap perubahan status akan otomatis menambah entri pada tabel `complaint_histories`.
- Upload gambar belum diimplementasikan pada fase awal (bisa ditambahkan kemudian).

*Dokumen ini ditulis dalam bahasa Indonesia formal.*
