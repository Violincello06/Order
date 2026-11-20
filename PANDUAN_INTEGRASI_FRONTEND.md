# 📘 Panduan Integrasi Order Service API untuk Tim Frontend

Dokumen ini adalah panduan resmi integrasi antara aplikasi client (Web Frontend / Mobile App) dengan **Order Service API**.

---

## 1. 📁 Berkas / File Spesifikasi yang Disertakan

Bersama dengan dokumen ini, terdapat berkas spesifikasi teknis yang dapat diimpor langsung:

| Nama File | Format | Kegunaan |
| :--- | :--- | :--- |
| **`openapi.json`** | OpenAPI 3.0 / Swagger | Digunakan untuk Swagger UI, Swagger Editor, atau generator otomatis tipe data/SDK (misal: *openapi-typescript-codegen*, *RTK Query*, *Orval*). |
| **`order_api_postman_collection.json`** | Postman Collection v2.1 | Siap diimpor langsung ke Postman (klik tombol **Import** di Postman) untuk mencoba seluruh endpoint API secara instan. |
| **`src/App.vue`** | Vue 3 Component | Referensi implementasi UI reaktif (kalkulasi cart dinamis, thumbnail gambar produk, pemanggilan API, dan modal pembayaran). |

---

## 2. 🌐 Base URL & Kredensial Autentikasi

Semua request wajib menyertakan token autentikasi di dalam request header:

* **Base URL (Local Development):** `http://localhost:1323/api/v1`
* **Swagger Interactive Docs (Web UI):** `http://localhost:1323/docs`
* **Header Wajib:**
  ```http
  Content-Type: application/json
  X-API-Key: ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba
  ```
  *(Catatan: Anda juga bisa menggunakan format standard `Authorization: Bearer ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba`)*

---

## 3. 📋 Daftar Endpoint API

### A. Membuat Order Baru (Create Order)
* **Method:** `POST`
* **URL:** `/api/v1/orders`
* **Headers:**
  ```http
  Content-Type: application/json
  X-API-Key: ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba
  ```
* **Request Body:**
  ```json
  {
    "customer_name": "Ralip",
    "customer_email": "pranaja0852@gmail.com",
    "items": [
      {
        "product_id": "1",
        "quantity": 1
      }
    ]
  }
  ```
* **Response Sukses (`201 Created`):**
  ```json
  {
    "success": true,
    "message": "Order berhasil dibuat",
    "data": {
      "id": "ORD-20261119-50F98A44",
      "customer_name": "Ralip",
      "customer_email": "pranaja0852@gmail.com",
      "total_amount": 450000,
      "status": "PENDING",
      "payment_id": "PAY-1795072110",
      "payment_url": "http://localhost:1323/sandbox/pay/PAY-1795072110",
      "items": [
        {
          "product_id": "1",
          "product_name": "Jaket Denim Levi's Vintage 90s Size L",
          "price": 450000,
          "quantity": 1,
          "subtotal": 450000,
          "image_url": "https://images.unsplash.com/photo-1576995853123-5a10305d93c0?w=400&auto=format&fit=crop&q=80"
        }
      ],
      "created_at": "2026-11-19T14:08:30Z",
      "updated_at": "2026-11-19T14:08:30Z"
    }
  }
  ```

---

### B. Mengambil Riwayat Semua Order (Get Orders)
* **Method:** `GET`
* **URL:** `/api/v1/orders`
* **Headers:** `X-API-Key: ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba`
* **Response Sukses (`200 OK`):**
  ```json
  {
    "success": true,
    "message": "Daftar order berhasil diambil",
    "data": [
      {
        "id": "ORD-20261119-50F98A44",
        "customer_name": "Ralip",
        "customer_email": "pranaja0852@gmail.com",
        "total_amount": 450000,
        "status": "PAID",
        "payment_id": "PAY-1795072110",
        "payment_url": "http://localhost:1323/sandbox/pay/PAY-1795072110",
        "created_at": "2026-11-19T14:08:30Z"
      }
    ]
  }
  ```

---

### C. Mengambil Detail Satu Order (Get Order Detail)
* **Method:** `GET`
* **URL:** `/api/v1/orders/{order_id}`
* **Contoh:** `/api/v1/orders/ORD-20261119-50F98A44`
* **Headers:** `X-API-Key: ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba`
* **Response Sukses (`200 OK`):**
  Mengembalikan detail lengkap order bersangkutan, rincian barang, serta status pembayaran.

---

### D. Membatalkan Order (Cancel Order)
* **Method:** `PUT`
* **URL:** `/api/v1/orders/{order_id}/cancel`
* **Catatan:** Hanya order berstatus `PENDING` yang dapat dibatalkan.
* **Headers:** `X-API-Key: ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba`
* **Response Sukses (`200 OK`):**
  ```json
  {
    "success": true,
    "message": "Order berhasil dibatalkan",
    "data": {
      "id": "ORD-20261119-50F98A44",
      "status": "CANCELLED"
    }
  }
  ```

---

### E. Simulasi Pembayaran Berhasil (Local Testing Only)
* **Method:** `POST`
* **URL:** `/api/v1/orders/{order_id}/pay`
* **Kegunaan:** Endpoint simulasi yang dapat dipanggil oleh tim frontend untuk mengetes perubahan status menjadi `PAID` tanpa harus menunggu webhook Payment Gateway pihak ketiga.
* **Headers:** `X-API-Key: ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba`

---

## 4. 🔄 Siklus Status Order (Order Status Lifecycle)

| Status | Penjelasan | Aksi yang Diperbolehkan |
| :--- | :--- | :--- |
| **`PENDING`** | Order baru dibuat, menunggu customer menyelesaikan pembayaran di Payment Gateway. | Buka `payment_url`, Batalkan (`PUT /cancel`), atau Bayar. |
| **`PAID`** | Pembayaran sukses diverifikasi (via Webhook Payment Gateway). | Tampilkan tanda lunas / struk pesanan. Tidak dapat dibatalkan. |
| **`CANCELLED`** | Pesanan dibatalkan oleh pengguna atau expired. | Read-only. Stok barang otomatis di-rollback ke Product Service. |

---

## 5. 💻 Contoh Kode Integrasi (JavaScript / Axios & Fetch)

### Contoh A: Menggunakan Fetch API
```javascript
const API_BASE_URL = 'http://localhost:1323/api/v1';
const API_KEY = 'ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba';

async function createOrder(customerName, customerEmail, cartItems) {
  try {
    const response = await fetch(`${API_BASE_URL}/orders`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': API_KEY
      },
      body: JSON.stringify({
        customer_name: customerName,
        customer_email: customerEmail,
        items: cartItems // Contoh: [{ product_id: "1", quantity: 1 }]
      })
    });

    const result = await response.json();
    if (!result.success) {
      throw new Error(result.message);
    }

    console.log('Order Berhasil Dibuat:', result.data);
    
    // Arahkan user ke halaman pembayaran Sandbox jika ada
    if (result.data.payment_url) {
      window.open(result.data.payment_url, '_blank');
    }
    
    return result.data;
  } catch (err) {
    console.error('Gagal membuat order:', err.message);
    alert('Error: ' + err.message);
  }
}
```

### Contoh B: Menggunakan Axios
```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:1323/api/v1',
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': 'ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba'
  }
});

// Ambil riwayat order
export const getOrders = async () => {
  const res = await api.get('/orders');
  return res.data.data;
};

// Buat order baru
export const checkoutOrder = async (orderPayload) => {
  const res = await api.post('/orders', orderPayload);
  return res.data;
};
```

---

## 6. ⚠️ Format Error Response

Jika terjadi kesalahan validasi (misal stok habis atau token salah), API akan mengembalikan HTTP status `4xx` dengan struktur konsisten:

```json
{
  "success": false,
  "message": "validasi produk gagal: stok produk tidak tersedia (sudah dibooking atau terjual)"
}
```

---

## 7. 📞 Kontak & Dokumentasi Lanjutan
- **Dokumentasi UI Swagger:** `http://localhost:1323/docs`
- **File Postman:** [order_api_postman_collection.json](order_api_postman_collection.json)
- **File OpenAPI:** [openapi.json](openapi.json)
