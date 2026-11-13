## Architecture Design

Proyek **CAMPUSCARE** menggunakan **Clean Architecture** yang memisahkan kode menjadi tiga lapisan utama:

1. **Entities (Domain Layer)** – berada di `app/model`.  Entitas merepresentasikan objek bisnis (User, Role, Facility, Category, Complaint, Assignment, ComplaintImage, ComplaintHistory) dan tidak bergantung pada framework atau driver database.
2. **Use Cases / Services (Business Logic Layer)** – berada di `app/service`.  Service berisi logika aplikasi (registrasi, login, manajemen pengaduan, penugasan, dll.) dan berinteraksi dengan **port** berupa interface repository.
3. **Interface Adapters (Framework & Driver Layer)** – berada di beberapa paket:
   - `repository` – implementasi akses data menggunakan **pgx** (PostgreSQL driver).
   - `config` – konfigurasi aplikasi dan environment.
   - `database` – inisialisasi koneksi database dan penyediaan *connection pool*.
   - `middleware` – autentikasi JWT dan otorisasi per‑role.
   - `route` – pendaftaran endpoint Fiber yang memanggil service.
   - `helper` – utilitas umum (hash password, generate/validate JWT, standar response).

### Dependency Rule
- **Frameworks & Drivers** (Fiber, pgx) → **Interface Adapters** → **Use Cases** → **Entities**.
- Entitas tidak meng‑import apapun selain paket standar Go.
- Service tidak meng‑import driver database; hanya menggunakan interface repository yang didefinisikan di `app/repository`.
- Repository meng‑import driver (pgx) dan menerima konteks `*pgxpool.Pool` dari package `database`.

### Flow Contoh Registrasi
```
HTTP Request (Fiber) → route/auth.go → AuthService.Register → UserRepository.Create (pgx) → DB
```
Setelah user dibuat, service menghasilkan response JSON standar.

### Keuntungan
- **Testability**: Service dapat di‑unit‑test dengan *mock* repository.
- **Maintainability**: Perubahan pada database atau framework tidak mempengaruhi bisnis logic.
- **Scalability**: Karena lapisan terpisah, mudah menambahkan API lain atau mengganti framework.

*Dokumen ini disusun dalam bahasa Indonesia formal.*
