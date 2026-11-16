# PETRIVA — Git Workflow & Version Control Standard

## 1. Alur Kerja Git (Git Feature Branch Workflow)
Proyek PETRIVA menggunakan alur kerja berbasis *feature branch* terpadu dengan peran **Git Master** sebagai penanggung jawab integrasi ke branch utama (`main`):

```
       [feat/auth-jwt] ───────────────┐
                                      ▼
[main] ───────────────────────────> [main] (Release / Stable)
                                      ▲
       [feat/pet-appointment] ────────┘
```

---

## 2. Aturan Branching & Penamaan

- `main`: Branch produksi yang selalu dalam kondisi stabil, dapat dikompilasi (`go build`), dan seluruh pengujian lulus (`go test`).
- `feat/<nama-fitur>`: Branch untuk penambahan fitur baru (contoh: `feat/appointment-service`, `feat/swagger-ui`).
- `fix/<nama-bug>`: Branch untuk perbaikan issue atau bug (contoh: `fix/conflict-check`, `fix/env-parser`).
- `docs/<nama-dokumen>`: Branch khusus pembaruan dokumentasi atau panduan setup.

---

## 3. Standar Pesan Commit (Conventional Commits)

Format commit:
`<tipe>(<lingkup opsional>): <deskripsi singkat imperative>`

Contoh tipe yang digunakan:
- `feat`: Penambahan fitur baru (`feat: implement pet management and ownership validation`)
- `fix`: Perbaikan bug (`fix: resolve scheduling conflict in appointment service`)
- `test`: Penambahan atau pembaruan pengujian (`test: add comprehensive mock tests for medical record service`)
- `docs`: Pembaruan dokumentasi (`docs: add Swagger UI guide and ERD diagram`)
- `refactor`: Perubahan struktur kode tanpa mengubah fungsionalitas (`refactor: align service layer with clean architecture`)
- `chore`: Pemeliharaan dependensi atau konfigurasi (`chore: add swagger dependencies to go.mod`)

---

## 4. Peran Git Master
1. Memastikan tidak ada *force push* ke branch `main`.
2. Menjaga riwayat commit bersih dan informatif (*micro commits*).
3. Melakukan verifikasi build (`go build ./...`) dan pengujian (`go test ./...`) sebelum melakukan merge ke branch `main`.
