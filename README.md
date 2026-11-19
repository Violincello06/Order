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
- **Header:** `Content-Type: application/json`
- **Body:**
```json
{
  "customer_name": "Budi Santoso",
  "customer_email": "budi@example.com",
  "items": [
    {
      "product_id": "PROD-001",
      "quantity": 2
    },
    {
      "product_id": "PROD-002",
      "quantity": 1
    }
  ]
}
```
*(Catatan: Produk dummy bawaan: `PROD-001` (Kopi Arabika, Rp75.000), `PROD-002` (Syrup Vanilla, Rp45.000), `PROD-003` (Paper Cup, Rp35.000))*

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

