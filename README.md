# Order Service API & Vue 3 Dashboard (Golang Echo + Vue 3 App.vue)

Backend dan Frontend terpadu untuk manajemen data pesanan (Order) yang dirancang siap terintegrasi dengan **Product API** dan **Payment Gateway API**.

> **Catatan Teknologi:** Seluruh antarmuka visual dibuat dalam format komponen resmi Vue 3: **[App.vue](file:///c:/xampp/htdocs/Order/public/App.vue)** *(Single File Component: `<template>`, `<script setup>`, `<style>`)*. Backend menggunakan **100% Go (Golang Echo)**. **Tidak ada kode PHP (0% PHP code)**.

---

## 🏛️ Struktur Arsitektur

```
c:\xampp\htdocs\Order/
├── public/
│   ├── App.vue              # Komponen visual utama (Vue 3 Single File Component)
│   └── docs.html            # Swagger UI Interactive API documentation
├── src/
│   ├── App.vue              # File SFC untuk workflow Vite/NPM
│   └── main.js              # Entrypoint Vue
├── package.json             # Konfigurasi Vue 3 & Vite
├── vite.config.js           # Konfigurasi proxy Vite ke backend Go Echo
├── openapi.json             # Spesifikasi OpenAPI 3.0 untuk tim frontend
├── order_api_postman_collection.json # Postman collection v2.1
├── internal/                # Layered Clean Architecture Backend (Go)
├── main.go                  # Server Go Echo (Mount App.vue langsung tanpa index.html)
├── go.mod
└── go.sum
```

---

## 🔑 Daftar Key Token untuk Setiap Tim

Semua endpoint kini telah diamankan dengan middleware autentikasi token. Konfigurasi tersimpan di file **[.env](file:///c:/xampp/htdocs/Order/.env)**:

| Penerima / Tim | Nama Header | Nilai Key Token | Fungsi |
| :--- | :--- | :--- | :--- |
| **🎨 Tim Frontend** | `X-API-Key` *(atau `Bearer`)* | `ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba` | Untuk memanggil API order (`/api/v1/orders`) dari web/mobile client. |
| **💳 Tim Payment Gateway** | `X-Webhook-Token` | `whsec_pay_8841a0e9b62c451fbcda332156890abc` | Untuk memvalidasi notifikasi webhook pembayaran (`/api/v1/payments/webhook`). |
| **📦 Tim Produk** | `X-Product-Key` | `prd_sec_int_771f2e9b0a3c4d5e1234567890abcdef` | Untuk autentikasi komunikasi internal service-to-service antara Order & Produk. |
| **👑 Master Internal** | `X-API-Key` | `ord_master_admin_secret_99887766554433221100` | Token admin dengan akses penuh ke seluruh endpoint. |

---

## 🚀 Cara Menjalankan

### Cara 1: Langsung via Server Go (Tanpa Build Step)
Server Go Echo secara otomatis melayani dan merender [App.vue](file:///c:/xampp/htdocs/Order/public/App.vue) langsung di browser via browser SFC engine:
```bash
go run main.go
```
Buka browser:
- **Tampilan Visual App.vue:** `http://localhost:1323/` atau `http://localhost:1323/dashboard`
- **File Mentah App.vue:** `http://localhost:1323/App.vue`
- **Interactive Swagger Docs:** `http://localhost:1323/docs`

### Cara 2: Workflow Standar Vue Vite (Opsional untuk Tim Frontend)
Jika tim frontend ingin menjalankan dengan hot-reload via Vite:
```bash
npm install
npm run dev
```
*(Aplikasi Vue akan berjalan di `http://localhost:3000` dengan proxy otomatis ke backend Go Echo di port 1323).*

---

## 🔑 Daftar Key Token untuk Setiap Tim

Semua endpoint kini telah diamankan dengan middleware autentikasi token. Konfigurasi tersimpan di file **[.env](file:///c:/xampp/htdocs/Order/.env)**:

| Penerima / Tim | Nama Header | Nilai Key Token | Fungsi |
| :--- | :--- | :--- | :--- |
| **🎨 Tim Frontend** | `X-API-Key` *(atau `Bearer`)* | `ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba` | Untuk memanggil API order (`/api/v1/orders`) dari web/mobile client. |
| **💳 Tim Payment Gateway** | `X-Webhook-Token` | `whsec_pay_8841a0e9b62c451fbcda332156890abc` | Untuk memvalidasi notifikasi webhook pembayaran (`/api/v1/payments/webhook`). |
| **📦 Tim Produk** | `X-Product-Key` | `prd_sec_int_771f2e9b0a3c4d5e1234567890abcdef` | Untuk autentikasi komunikasi internal service-to-service antara Order & Produk. |
| **👑 Master Internal** | `X-API-Key` | `ord_master_admin_secret_99887766554433221100` | Token admin dengan akses penuh ke seluruh endpoint. |

---

## 🚀 Cara Menjalankan Server & Membuka Tampilan Visual

Jalankan perintah berikut di terminal:

```bash
go run main.go
```

Buka browser Anda dan akses:
- **Tampilan Visual Web Dashboard:** `http://localhost:1323/` atau `http://localhost:1323/dashboard`
- **Interactive Swagger API Docs (Untuk Tim Frontend):** `http://localhost:1323/docs`
- **Spesifikasi OpenAPI JSON:** `http://localhost:1323/openapi.json`
- **Info API JSON:** `http://localhost:1323/api/info`

---

## 🗄️ Database MySQL (XAMPP)

Seluruh skema database dan seed data telah disiapkan di folder `database/` dan sudah di-import ke MySQL XAMPP:

| Nama Database | Skema File | Tabel | Deskripsi |
| :--- | :--- | :--- | :--- |
| **`order_db`** | `database/order_schema.sql` | `orders`, `order_items` | Database utama Order Service untuk menyimpan data pesanan & item transaksi. |
| **`preloved_product`** | `database/product_schema.sql` | `categories`, `products`, `product_images` | Database Product Microservice (port 3001) yang menyimpan katalog produk Pre-Loved. |

### Konfigurasi Koneksi Database (`.env`):
```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=order_db
```
*Sistem dilengkapi fitur **Auto-Detect**: jika MySQL aktif, data tersimpan permanen di MySQL; jika MySQL mati, sistem otomatis fallback ke In-Memory tanpa membuat server crash.*

---

## 📦 File Spesifikasi API untuk Diberikan ke Tim Lain

Di dalam root project ini telah disediakan 2 format standar industri:
1. **[openapi.json](file:///c:/xampp/htdocs/Order/openapi.json)** (OpenAPI / Swagger 3.0 Standard)
   - Format universal untuk generator kode (TypeScript, Axios, Flutter, Android/iOS) dan Swagger Editor.
2. **[order_api_postman_collection.json](file:///c:/xampp/htdocs/Order/order_api_postman_collection.json)** (Postman Collection v2.1)
   - Tim frontend / QA cukup klik **Import** di Postman dan semua request beserta contoh JSON body langsung siap diuji coba.

---

## 🖥️ Fitur Tampilan Visual (Dashboard)
1. **Statistik Real-time**: Menampilkan Total Order, Menunggu Bayar, Sudah Dibayar, dan Total Omset.
2. **Katalog & Pembuatan Order**: Memilih produk dengan kalkulasi harga otomatis secara dinamis.
3. **Simulasi Bayar Instan**: Tombol "Simulasi Bayar" yang langsung memicu webhook payment gateway untuk mengubah status ke `PAID`.
4. **Pembatalan Order**: Pembatalan pesanan yang masih berstatus `PENDING`.
5. **Modal Detail Order**: Menampilkan breakdown belanja, status pembayaran, dan URL sandbox payment gateway.

---

## 📡 Daftar Endpoint API

### 1. Root / Health Check
- **Method:** `GET`
- **URL:** `http://localhost:1323/`

---

### 2. Buat Order Baru (Create Order)
- **Method:** `POST`
- **URL:** `http://localhost:1323/api/v1/orders`
- **Header:** `Content-Type: application/json`, `X-API-Key: ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba`
- **Body:**
```json
{
  "customer_name": "Budi Santoso",
  "customer_email": "budi@example.com",
  "items": [
    {
      "product_id": "1",
      "quantity": 1
    }
  ]
}
```
*(Catatan: Produk Pre-Loved bawaan: `1` (Jaket Denim Levi's Vintage 90s, Rp450.000), `2` (Sepatu Sneakers Converse 70s, Rp350.000), `3` (Kamera Analog Vintage 35mm, Rp750.000))*

---

### 3. Ambil Semua Daftar Order
- **Method:** `GET`
- **URL:** `http://localhost:1323/api/v1/orders`

---

### 4. Ambil Detail Order berdasarkan ID
- **Method:** `GET`
- **URL:** `http://localhost:1323/api/v1/orders/{ORDER_ID}`
- **Contoh:** `http://localhost:1323/api/v1/orders/ORD-20261114-1A2B3C`

---

### 5. Batalkan Order (Cancel)
- **Method:** `PUT`
- **URL:** `http://localhost:1323/api/v1/orders/{ORDER_ID}/cancel`

---

### 6. Webhook Notifikasi Payment Gateway
- **Method:** `POST`
- **URL:** `http://localhost:1323/api/v1/payments/webhook`
- **Body:**
```json
{
  "order_id": "ORD-20261114-1A2B3C",
  "payment_id": "PAY-1731600000",
  "transaction_status": "settlement",
  "payment_method": "qris",
  "gross_amount": 195000
}
```

---

### 7. Simulasi Pembayaran Berhasil (Testing Shortcut)
- **Method:** `POST`
- **URL:** `http://localhost:1323/api/v1/orders/{ORDER_ID}/pay`

