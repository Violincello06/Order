# Order Service API (Golang + Echo Framework)

Backend tahap awal untuk manajemen data pesanan (Order) yang dirancang siap terintegrasi dengan **Product API** dan **Payment Gateway API**.

---

## 🏛️ Struktur Arsitektur (Layered Architecture)

```
c:\xampp\htdocs\Order/
├── internal/
│   ├── models/              # Struktur data Order, OrderItem, & Request/Response DTO
│   │   └── order.go
│   ├── clients/             # Interface & Client integrasi Produk & Payment
│   │   ├── product_client.go# Kontrak & simulasi integrasi Product Service (stok, harga)
│   │   └── payment_client.go# Kontrak & simulasi integrasi Payment Gateway (Midtrans/Xendit)
│   ├── repository/          # Layer penyimpanan data (In-Memory, siap diganti MySQL)
│   │   └── order_repository.go
│   ├── service/             # Business logic pemesanan, kalkulasi, & validasi
│   │   └── order_service.go
│   └── handler/             # HTTP Controllers (Echo Framework)
│       └── order_handler.go
├── main.go                  # Entry point & Dependency Injection
├── go.mod
└── go.sum
```

---

## 🚀 Cara Menjalankan Server & Membuka Tampilan Visual

Jalankan perintah berikut di terminal:

```bash
go run main.go
```

Buka browser Anda dan akses:
- **Tampilan Visual Web Dashboard:** `http://localhost:1323/` atau `http://localhost:1323/dashboard`
- **Info API JSON:** `http://localhost:1323/api/info`

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

