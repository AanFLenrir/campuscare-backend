# PETRIVA — Architecture Design

## 1. Prinsip Desain: Clean Architecture
PETRIVA dirancang dengan mematuhi prinsip **Clean Architecture** (Robert C. Martin). Tujuan utama arsitektur ini adalah memisahkan antara detail teknis (framework, database, protokol HTTP) dengan inti logika bisnis aplikasi (*business rules*).

### Ketentuan Mutlak Proyek:
- **DILARANG** menggunakan layer `controller/`, `controllers/`, `handler/`, maupun `handlers/`.
- **DILARANG** mengimpor atau menggunakan `*fiber.Ctx` di dalam layer `app/service/`.
- Komunikasi antar layer diatur melalui antarmuka (*interfaces*) dan injeksi dependensi (*dependency injection*).

---

## 2. Diagram Aliran Dependensi

```
┌─────────────────────────────────────────────────────────────┐
│                    FRAMEWORKS & DRIVERS                     │
│  - Fiber v2 HTTP Server (main.go)                          │
│  - PostgreSQL Database via pgxpool (database/postgres.go)   │
│  - Swagger UI (github.com/gofiber/swagger)                 │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                   INTERFACE ADAPTERS                        │
│  - route/route.go (Fiber Callbacks langsung)                │
│  - middleware/ (JWT Auth & Role Authorization)              │
│  - helper/ (bcrypt hashing, JWT token, API response)        │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                   APPLICATION BUSINESS RULES                │
│  - app/service/ (Pure Business Logic & Ownership Checks)    │
│  - Menggunakan context.Context & Interface Repository       │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                  ENTERPRISE BUSINESS RULES                  │
│  - app/model/ (Entities, DTOs, Enums, Status Constants)     │
│  - app/repository/ (Interface & SQL Implementations)        │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Pembagian Tanggung Jawab Layer

### Layer 1: Entities & Models (`app/model/`)
- Murni berisi definisi struct data entitas domain.
- Request DTO dan Response DTO.
- Konstanta enum status dan role.
- Tidak memiliki dependensi ke database, Fiber, atau framework eksternal.

### Layer 2: Repositories (`app/repository/`)
- Menentukan interface akses data ke PostgreSQL.
- Implementasi konkret menggunakan koneksi pool `pgxpool.Pool`.
- Seluruh query menggunakan parameterized query (`$1, $2, ...`) untuk mencegah SQL Injection.
- Bertanggung jawab memetakan hasil query baris SQL menjadi struct entitas.

### Layer 3: Services (`app/service/`)
- Pusat seluruh logika bisnis PETRIVA.
- Validasi data masukan (misal: format tanggal, kelayakan spesies/gender, nominal tarif).
- Validasi kepemilikan data (*ownership*): memastikan user hanya dapat mengakses data yang berhak.
- State machine transisi status appointment.
- Tidak mengetahui bahwa aplikasi berjalan di atas HTTP atau Fiber.

### Layer 4: Routes & HTTP Adapters (`route/route.go`)
- Adapter HTTP yang mendaftarkan endpoint pada Fiber.
- Callback Fiber langsung mengekstrak parameter, query string, dan payload JSON.
- Mengambil identitas user ID dan role dari context locals yang disediakan middleware.
- Memanggil method service yang sesuai.
- Memformat dan mengirimkan respon HTTP standar JSON (`helper.Success`, `helper.Created`, `helper.Error`).

