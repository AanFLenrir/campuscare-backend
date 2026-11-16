# PETRIVA — Testing Documentation

## 1. Strategi Pengujian
Pengujian pada sistem PETRIVA berfokus pada:
1. **Unit Testing:** Pengujian isolasi logika bisnis pada setiap service menggunakan mock repository in-memory tanpa ketergantungan pada koneksi PostgreSQL aktif.
2. **Security & Cryptography Testing:** Pengujian pembangkitan token JWT, validasi kedaluwarsa token, signing key, serta keandalan hashing password bcrypt.
3. **Middleware Testing:** Pengujian otentikasi header HTTP Authorization dan penegakan izin akses per role (RBAC).

---

## 2. Cara Menjalankan Pengujian

### Menjalankan Seluruh Test Suite:
```powershell
go test -v ./...
```

### Menjalankan Pengujian per Package Tertentu:
```powershell
# Uji seluruh business logic service
go test -v ./app/service/...

# Uji fungsi kriptografi & utilitas
go test -v ./helper/...

# Uji middleware autentikasi & otorisasi
go test -v ./middleware/...
```

### Menjalankan Analisis Statis Kode (Go Vet):
```powershell
go vet ./...
```

---

## 3. Daftar Skenario Pengujian

### 3.1. Autentikasi (`app/service/auth_service_test.go`)
- `TestAuth_Register_Success`: Memastikan pendaftaran pemilik hewan berhasil dan menghasilkan token JWT serta role `PET_OWNER`.
- `TestAuth_Register_InvalidInput`: Memvalidasi penolakan pendaftaran jika nama, email, atau password tidak memenuhi syarat.
- `TestAuth_Register_DuplicateEmail`: Memastikan sistem menolak pendaftaran dengan email yang sudah terdaftar.
- `TestAuth_Login_Success`: Memastikan login dengan kredensial valid berhasil menerbitkan JWT token.
- `TestAuth_Login_WrongPassword`: Memastikan sistem menolak login jika password salah.
- `TestAuth_Login_InactiveUser`: Memastikan akun yang dinonaktifkan ditolak saat login.

### 3.2. Manajemen Pengguna (`app/service/user_service_test.go`)
- `TestUser_CreateUser_Admin_Success`: Memastikan admin dapat membuat akun dengan role apapun.
- `TestUser_CreateVeterinarian_WithClinic`: Memastikan dokter hewan wajib diasosiasikan dengan klinik yang aktif.
- `TestUser_CreateUser_InvalidRole`: Memastikan sistem menolak role yang tidak dikenal.
- `TestUser_GetUser_Success`: Memastikan pencarian profil user berhasil tanpa mengekspos hash password.
- `TestUser_DeleteUser_Success`: Memastikan penghapusan user menjalankan mekanisme deaktivasi (*soft delete*).

### 3.3. Hewan Peliharaan (`app/service/pet_service_test.go`)
- `TestPet_CreatePet_Success`: Memastikan pembuatan data hewan berhasil.
- `TestPet_CreatePet_InvalidSpecies`: Memastikan spesies di luar daftar valid ditolak.
- `TestPet_UpdatePet_ForbiddenForOtherOwner`: Memastikan pemilik hewan tidak dapat mengubah hewan milik user lain (*ownership boundary*).
- `TestPet_DeletePet_Success`: Memastikan penghapusan hewan berhasil menonaktifkan data.

### 3.4. Reservasi Konsultasi (`app/service/appointment_service_test.go`)
- `TestAppointment_Create_Success`: Memastikan pembuatan appointment berhasil dengan status awal `PENDING`.
- `TestAppointment_Create_PetNotOwned`: Memastikan user tidak dapat membuat appointment untuk hewan milik orang lain.
- `TestAppointment_StatusTransitions`: Memastikan alur status berjalan benar (`PENDING` → `CONFIRMED` → `IN_PROGRESS` → `COMPLETED`).
- `TestAppointment_Cancel_Success`: Memastikan pembatalan appointment berhasil dengan alasan yang tersimpan.

### 3.5. Rekam Medis (`app/service/medical_record_service_test.go`)
- `TestMedicalRecord_Create_Success`: Memastikan dokter hewan yang ditugaskan dapat mencatat rekam medis dan appointment otomatis bertransisi menjadi `COMPLETED`.
- `TestMedicalRecord_Create_ForbiddenForOtherVet`: Memastikan dokter hewan lain dilarang mencatat rekam medis pada janji temu yang bukan penugasannya.

### 3.6. Kriptografi & Token (`helper/jwt_test.go`, `helper/password_test.go`)
- `TestGenerateAndValidateToken_Success`: Pembangkitan dan validasi token JWT.
- `TestValidateToken_InvalidSecret`: Penolakan token yang di-sign dengan secret key berbeda.
- `TestValidateToken_Expired`: Penolakan token yang telah melewati waktu kedaluwarsa.
- `TestHashPassword_Success`: Keberhasilan pembangkitan hash bcrypt.
- `TestCheckPassword_Correct` & `TestCheckPassword_Incorrect`: Verifikasi kesesuaian password plaintext dengan hash.

### 3.7. Middleware HTTP (`middleware/middleware_test.go`)
- `TestAuthMiddleware_MissingHeader`: Penolakan request tanpa header Authorization (401).
- `TestAuthMiddleware_InvalidHeaderFormat`: Penolakan skema non-Bearer (401).
- `TestAuthMiddleware_InvalidToken`: Penolakan token JWT palsu (401).
- `TestAuthMiddleware_Success`: Keberhasilan ekstraksi claims `userID` dan `role` ke dalam locals Fiber (200).
- `TestAuthorize_AdminAllowed`: Otorisasi akses role `ADMIN` (200).
- `TestAuthorize_ForbiddenRole`: Penolakan akses role terlarang (403).
- `TestAuthorize_MultipleAllowedRoles`: Pengujian akses multi-role (contoh: `ADMIN` dan `VETERINARIAN`).

