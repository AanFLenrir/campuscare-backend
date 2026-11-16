# PETRIVA — Functional Requirements

## 1. Modul Autentikasi & Akun
- **FR-01:** Sistem dapat memproses pendaftaran akun baru untuk role `PET_OWNER`.
- **FR-02:** Sistem dapat melakukan verifikasi login menggunakan email dan password yang telah di-hash dengan `bcrypt`.
- **FR-03:** Sistem menerbitkan JSON Web Token (JWT) dengan masa kedaluwarsa yang dapat dikonfigurasi.
- **FR-04:** Sistem menyediakan endpoint profil pengguna yang sedang login (`/auth/me`).

## 2. Modul Master Data (Klinik & Layanan)
- **FR-05:** Admin dapat membuat, membaca, memperbarui, dan menonaktifkan data klinik.
- **FR-06:** Admin dapat mengelola katalog layanan (nama, deskripsi, durasi menit, tarif) untuk setiap klinik.
- **FR-07:** Pengguna umum (publik) dapat melihat daftar klinik dan layanan yang sedang aktif.

## 3. Modul Hewan Peliharaan (Pets)
- **FR-08:** Pemilik hewan dapat menambahkan hewan peliharaan miliknya dengan validasi spesies (`CAT`, `DOG`, `BIRD`, `RABBIT`, `OTHER`) dan gender (`MALE`, `FEMALE`, `UNKNOWN`).
- **FR-09:** Pemilik hewan hanya dapat melihat, mengubah, dan menonaktifkan hewan peliharaannya sendiri.
- **FR-10:** Sistem menerapkan *soft delete* pada data hewan untuk menjaga integritas riwayat rekam medis dan vaksinasi.

## 4. Modul Reservasi Janji Temu (Appointments)
- **FR-11:** Pemilik hewan dapat membuat janji temu baru dengan memilih hewan, dokter hewan, klinik, layanan, dan waktu appointment.
- **FR-12:** Sistem memvalidasi bahwa waktu appointment berada di masa mendatang dan dokter tidak memiliki jadwal bentrok.
- **FR-13:** Sistem memvalidasi bahwa dokter dan layanan terdaftar pada klinik yang dipilih.
- **FR-14:** Sistem mengelola transisi status janji temu: `PENDING` → `CONFIRMED` → `IN_PROGRESS` → `COMPLETED` / `CANCELLED` / `REJECTED`.
- **FR-15:** Pemilik dapat membatalkan janji temu miliknya yang masih berstatus `PENDING` atau `CONFIRMED` dengan menyertakan alasan pembatalan.

## 5. Modul Rekam Medis (Medical Records)
- **FR-16:** Dokter hewan yang ditugaskan dapat membuat rekam medis untuk janji temu yang ditanganinya.
- **FR-17:** Rekam medis mencakup diagnosis (wajib), gejala, tindakan medis, dan resep obat.
- **FR-18:** Pembuatan rekam medis otomatis mengubah status janji temu menjadi `COMPLETED`.
- **FR-19:** Sistem mencegah duplikasi pembuatan rekam medis untuk janji temu yang sama.
- **FR-20:** Pemilik hewan dapat membaca seluruh riwayat rekam medis hewan peliharaannya.

## 6. Modul Buku Vaksinasi Digital (Vaccinations)
- **FR-21:** Dokter hewan dapat mencatat riwayat vaksinasi hewan meliputi nama vaksin, tanggal vaksinasi, tanggal vaksin berikutnya, nomor batch, dan catatan medis.
- **FR-22:** Pemilik hewan dapat melihat riwayat vaksinasi hewannya.

## 7. Modul Dashboard & Dokumentasi
- **FR-23:** Sistem menyediakan ringkasan metrik statistik real-time sesuai role yang mengakses.
- **FR-24:** Sistem menyediakan dokumentasi API interaktif berbasis Swagger UI di endpoint `/swagger/`.
