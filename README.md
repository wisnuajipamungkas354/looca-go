# Looca Backend API (Go + Gin + GORM + MySQL)

Backend REST API untuk mengelola data website Looca (`looca-react` dan website publik).

---

## Persiapan & Instalasi

### 1. Buat Database MySQL
Pastikan MySQL sudah berjalan dan buat database bernama `looca`:
```sql
CREATE DATABASE looca CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```
*(Atau jalankan `docker compose up -d` jika ingin menggunakan Docker)*

### 2. Konfigurasi Environment
Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
Sesuaikan konfigurasi database (`DB_USER`, `DB_PASS`, `DB_NAME`, dll.) dan `JWT_SECRET`.

### 3. Akun Admin Default
Secara default pada `.env`:
- **Email:** `admin@looca.id`
- **Password:** `admin123`

---

## Seeding Data (Manual)

Untuk mengisi data awal (admin default, 5 kategori, pengaturan situs, 9 produk, 5 workshop, dan 4 event dengan total counter bernilai 0):

```bash
# Mode Development
go run . seed

# Mode Production (Binary)
./looca-go seed
```

> **PENTING (Deploy ke Server):**
> Folder `database/seed_assets/` berisi aset gambar dan video awal yang akan disalin otomatis ke folder `uploads/` saat proses seeding berjalan. Pastikan folder `database/seed_assets/` ikut di-upload ke server saat deployment.

---

## Menjalankan Server

```bash
# Menjalankan langsung
go run .

# Atau build binary
go build -o looca-go.exe .
./looca-go.exe
```

Server akan aktif di `http://localhost:8080`.

---

## Daftar Endpoint API

### Public (Tanpa Auth)
| Method | Endpoint | Keterangan |
|---|---|---|
| `GET` | `/health` | Health check server |
| `POST` | `/api/auth/login` | Login admin (return JWT token & user) |
| `GET` | `/api/public/categories` | Daftar semua kategori |
| `GET` | `/api/public/products` | Daftar produk (opsional filter `?category_id=`) |
| `GET` | `/api/public/workshops` | Daftar workshop |
| `GET` | `/api/public/events` | Daftar event yang aktif (`is_active = true`) |
| `GET` | `/api/public/settings` | Pengaturan situs (WhatsApp, medsos, Shopee) |
| `POST` | `/api/public/products/:id/click` | Tambah 1 klik pada produk |
| `POST` | `/api/public/workshops/:id/play` | Tambah 1 play pada workshop |
| `POST` | `/api/public/visit` | Catat 1 kunjungan website hari ini |

### Admin (Memerlukan Header `Authorization: Bearer <token>`)
| Method | Endpoint | Keterangan |
|---|---|---|
| `GET` | `/api/admin/me` | Info profil admin yang sedang login |
| `GET` | `/api/admin/stats` | Total produk, workshop, dan event |
| `GET` | `/api/admin/analytics` | Data analitik pengunjung (daily, weekly, monthly) |
| `GET` | `/api/admin/settings` | Ambil pengaturan situs |
| `PUT` | `/api/admin/settings` | Simpan perubahan pengaturan situs |
| `POST` | `/api/admin/uploads/image` | Upload gambar (max 1MB, format jpg/png/webp) |
| `POST` | `/api/admin/uploads/video` | Upload video (max 50MB, format mp4/webm) |
| `GET` | `/api/admin/categories` | List kategori |
| `POST` | `/api/admin/categories` | Tambah kategori baru |
| `PUT` | `/api/admin/categories/:id` | Update kategori |
| `DELETE` | `/api/admin/categories/:id` | Hapus kategori (gagal 409 jika dipakai produk) |
| `GET` | `/api/admin/products` | List produk |
| `POST` | `/api/admin/products` | Tambah produk baru |
| `GET` | `/api/admin/products/:id` | Detail produk |
| `PUT` | `/api/admin/products/:id` | Update produk |
| `DELETE` | `/api/admin/products/:id` | Hapus produk |
| `GET` | `/api/admin/workshops` | List workshop |
| `POST` | `/api/admin/workshops` | Tambah workshop baru |
| `GET` | `/api/admin/workshops/:id` | Detail workshop |
| `PUT` | `/api/admin/workshops/:id` | Update workshop |
| `DELETE` | `/api/admin/workshops/:id` | Hapus workshop |
| `GET` | `/api/admin/events` | List event |
| `POST` | `/api/admin/events` | Tambah event baru |
| `GET` | `/api/admin/events/:id` | Detail event |
| `PUT` | `/api/admin/events/:id` | Update event |
| `DELETE` | `/api/admin/events/:id` | Hapus event |
