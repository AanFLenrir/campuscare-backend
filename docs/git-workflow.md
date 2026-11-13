# Git Workflow — CAMPUSCARE Backend

## Struktur Branch

```
main          ← branch utama (production-ready)
feature/*     ← pengembangan fitur baru
fix/*         ← perbaikan bug
```

## Peran Git Master

Git Master bertugas:
- Mereview Pull Request dari anggota.
- Melakukan merge ke branch `main`.
- Memastikan `main` selalu dalam kondisi buildable.

## Panduan Micro Commits

Setiap commit harus kecil, fokus, dan deskriptif:

```bash
git add app/model/complaint.go
git commit -m "feat(model): add Complaint entity with status and priority fields"

git add app/repository/complaint_repository.go
git commit -m "feat(repo): implement ComplaintRepository with CRUD operations"

git add app/service/complaint_service.go
git commit -m "feat(service): implement ComplaintService pure business logic"

git add route/route.go
git commit -m "feat(route): wire complaint endpoints to service layer"
```

## Konvensi Pesan Commit

Format: `<type>(<scope>): <deskripsi>`

| Type | Keterangan |
|---|---|
| `feat` | Fitur baru |
| `fix` | Perbaikan bug |
| `refactor` | Refactoring tanpa perubahan fungsional |
| `docs` | Perubahan dokumentasi |
| `test` | Penambahan atau perbaikan test |
| `chore` | Perubahan konfigurasi, dependency, dsb |

## Larangan

- ❌ Jangan gunakan `git push --force` ke `main`.
- ❌ Jangan commit credential atau file `.env`.
- ❌ Jangan commit atas nama anggota lain.
- ❌ Jangan manipulasi tanggal commit.

## Push ke Remote

```bash
git push origin main
```

Pastikan seluruh pemeriksaan (`go build`, `go vet`, `go test`) berhasil sebelum push.

