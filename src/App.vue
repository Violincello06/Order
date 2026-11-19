<template>
  <div class="min-h-screen bg-slate-50 text-slate-800 flex flex-col font-sans">
    
    <!-- Navbar Header -->
    <header class="bg-white border-b border-slate-200 sticky top-0 z-30 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-emerald-500 to-teal-400 flex items-center justify-center text-white shadow-emerald-100 shadow-md">
            <i class="fa-brands fa-vuejs text-2xl"></i>
          </div>
          <div>
            <div class="flex items-center space-x-2">
              <h1 class="text-base font-bold tracking-tight text-slate-900 leading-none">Order Management</h1>
              <span class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-800">App.vue (SFC)</span>
              <span class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-blue-100 text-blue-800">Go Echo (Non-PHP)</span>
            </div>
            <span class="text-xs text-slate-500 font-medium">Terkoneksi ke Product API & Payment Gateway</span>
          </div>
        </div>

        <div class="flex items-center space-x-3">
          <!-- Backend Connection Badge -->
          <span :class="[
            'inline-flex items-center px-2.5 py-1 rounded-full text-xs font-semibold border transition',
            apiConnected ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-rose-50 text-rose-700 border-rose-200'
          ]">
            <span :class="['w-2 h-2 rounded-full mr-1.5', apiConnected ? 'bg-emerald-500' : 'bg-rose-500 animate-pulse']"></span>
            {{ apiConnected ? 'Echo Backend Online' : 'Echo Backend Offline' }}
          </span>

          <a href="/docs" target="_blank" class="px-2.5 py-1.5 text-xs font-medium text-slate-600 hover:text-indigo-600 hover:bg-slate-100 rounded-lg transition flex items-center gap-1 border border-slate-200">
            <i class="fa-solid fa-book-open"></i> Swagger Docs
          </a>

          <button @click="fetchOrders" :disabled="loading" class="p-2 text-slate-500 hover:text-emerald-600 hover:bg-slate-100 rounded-lg transition" title="Refresh Data">
            <i :class="['fa-solid fa-arrows-rotate', { 'fa-spin': loading }]"></i>
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8 flex-1 w-full">
      
      <!-- Top Metrics Overview (Computed in Vue) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p class="text-xs font-semibold uppercase text-slate-400 tracking-wider">Total Order</p>
            <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ stats.total }}</h3>
          </div>
          <div class="w-12 h-12 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center text-xl">
            <i class="fa-solid fa-receipt"></i>
          </div>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p class="text-xs font-semibold uppercase text-slate-400 tracking-wider">Menunggu Bayar</p>
            <h3 class="text-2xl font-bold text-amber-600 mt-1">{{ stats.pending }}</h3>
          </div>
          <div class="w-12 h-12 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center text-xl">
            <i class="fa-regular fa-clock"></i>
          </div>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p class="text-xs font-semibold uppercase text-slate-400 tracking-wider">Sudah Dibayar</p>
            <h3 class="text-2xl font-bold text-emerald-600 mt-1">{{ stats.paid }}</h3>
          </div>
          <div class="w-12 h-12 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center text-xl">
            <i class="fa-solid fa-circle-check"></i>
          </div>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p class="text-xs font-semibold uppercase text-slate-400 tracking-wider">Total Pendapatan</p>
            <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ formatRupiah(stats.revenue) }}</h3>
          </div>
          <div class="w-12 h-12 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center text-xl">
            <i class="fa-solid fa-wallet"></i>
          </div>
        </div>

      </div>

      <!-- Main Layout: Form (Left) & Order List (Right) -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        
        <!-- Left: Form Buat Order Baru (Vue Two-Way Binding) -->
        <div class="lg:col-span-5 bg-white p-6 rounded-2xl border border-slate-200 shadow-sm space-y-6">
          <div class="flex items-center justify-between border-b border-slate-100 pb-4">
            <div>
              <h2 class="text-base font-bold text-slate-900">Buat Orderan Baru</h2>
              <p class="text-xs text-slate-500">Kalkulasi dinamis via App.vue & validasi ke Go Echo</p>
            </div>
            <span class="text-emerald-700 text-xs font-bold bg-emerald-50 px-2 py-0.5 rounded-md flex items-center gap-1">
              <i class="fa-solid fa-bolt"></i> Vue Reactive
            </span>
          </div>

          <form @submit.prevent="handleCreateOrder" class="space-y-4">
            <div>
              <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">Nama Customer</label>
              <input type="text" v-model="form.customer_name" required placeholder="Contoh: Budi Santoso"
                class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500">
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">Email Customer</label>
              <input type="email" v-model="form.customer_email" required placeholder="budi@example.com"
                class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500">
            </div>

            <!-- List Produk (v-for) -->
            <div>
              <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-2">Pilih Produk (Katalog Produk)</label>
              <div class="space-y-2.5">
                <div v-for="product in products" :key="product.id"
                  class="p-3 rounded-xl border border-slate-200 bg-slate-50/70 flex items-center justify-between gap-3">
                  <div class="flex-1">
                    <span class="text-xs font-bold text-slate-800 block">{{ product.name }}</span>
                    <span class="text-xs text-slate-500">
                      {{ formatRupiah(product.price) }} &bull; 
                      <span class="text-emerald-600 font-medium">Stok: {{ product.stock }}</span>
                    </span>
                  </div>
                  <div class="flex items-center space-x-2">
                    <label class="text-xs text-slate-400">Qty:</label>
                    <input type="number" min="0" :max="product.stock" v-model.number="form.items[product.id]"
                      class="w-16 px-2 py-1.5 rounded-lg border border-slate-200 text-sm text-center font-bold bg-white focus:outline-none focus:border-emerald-500">
                  </div>
                </div>
              </div>
            </div>

            <!-- Computed Total Preview -->
            <div class="bg-emerald-50/70 p-4 rounded-xl border border-emerald-100 flex items-center justify-between">
              <div>
                <span class="text-xs text-emerald-800 font-medium block">Total Pembayaran (Auto-Calculated)</span>
                <span class="text-xl font-bold text-emerald-950">{{ formatRupiah(computedOrderTotal) }}</span>
              </div>
              <i class="fa-solid fa-credit-card text-emerald-500 text-2xl"></i>
            </div>

            <button type="submit" :disabled="submitting || computedOrderTotal === 0"
              class="w-full bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-semibold py-3 px-4 rounded-xl transition duration-200 shadow-md shadow-emerald-100 flex items-center justify-center space-x-2">
              <i :class="['fa-solid', submitting ? 'fa-spinner fa-spin' : 'fa-cart-plus']"></i>
              <span>{{ submitting ? 'Memproses Order...' : 'Proses & Terbitkan Order' }}</span>
            </button>
          </form>
        </div>

        <!-- Right: Daftar Order & Actions (Vue Reactive Cards) -->
        <div class="lg:col-span-7 bg-white p-6 rounded-2xl border border-slate-200 shadow-sm space-y-6">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 border-b border-slate-100 pb-4">
            <div>
              <h2 class="text-base font-bold text-slate-900">Daftar Order Masuk</h2>
              <p class="text-xs text-slate-500">Data live dari backend Go Echo</p>
            </div>
            
            <div class="flex items-center space-x-2">
              <label class="text-xs text-slate-400 font-medium">Filter:</label>
              <select v-model="filterStatus" class="text-xs border border-slate-200 rounded-lg px-2.5 py-1.5 bg-slate-50 font-medium focus:outline-none">
                <option value="ALL">Semua Status</option>
                <option value="PENDING">Menunggu Bayar</option>
                <option value="PAID">Sudah Dibayar</option>
                <option value="CANCELLED">Dibatalkan</option>
              </select>
            </div>
          </div>

          <!-- Loading State -->
          <div v-if="loading && orders.length === 0" class="text-center py-12 text-slate-400">
            <i class="fa-solid fa-spinner fa-spin text-2xl mb-2 text-emerald-500"></i>
            <p class="text-sm">Memuat data order...</p>
          </div>

          <!-- Empty State -->
          <div v-else-if="filteredOrders.length === 0" class="text-center py-12 border-2 border-dashed border-slate-200 rounded-2xl">
            <i class="fa-solid fa-inbox text-slate-300 text-3xl mb-2"></i>
            <p class="text-sm font-medium text-slate-500">Belum ada order dengan status ini</p>
            <p class="text-xs text-slate-400 mt-1">Buat orderan baru menggunakan formulir di sebelah kiri</p>
          </div>

          <!-- List Orders (v-for) -->
          <div v-else class="space-y-3">
            <div v-for="order in filteredOrders" :key="order.id"
              class="p-4 rounded-xl border border-slate-200 bg-white hover:border-emerald-300 hover:shadow-sm transition space-y-3">
              
              <div class="flex items-start justify-between gap-3">
                <div>
                  <div class="flex items-center space-x-2">
                    <span class="font-mono text-xs font-bold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded border border-emerald-100">
                      {{ order.id }}
                    </span>
                    <span class="text-xs text-slate-400">&bull; {{ formatDate(order.created_at) }}</span>
                  </div>
                  <h4 class="font-bold text-slate-800 text-sm mt-1">{{ order.customer_name }}</h4>
                  <p class="text-xs text-slate-500">{{ order.customer_email }}</p>
                </div>
                
                <div class="text-right">
                  <span :class="[
                    'inline-flex items-center text-xs font-semibold px-2.5 py-1 rounded-full border',
                    getStatusBadgeClass(order.status)
                  ]">
                    <i :class="['fa-solid mr-1.5', getStatusIcon(order.status)]"></i>
                    {{ order.status }}
                  </span>
                  <p class="font-bold text-slate-900 text-base mt-1">{{ formatRupiah(order.total_amount) }}</p>
                </div>
              </div>

              <!-- Items Break-down -->
              <div class="text-xs text-slate-600 bg-slate-50 p-2.5 rounded-lg space-y-1">
                <div v-for="item in order.items" :key="item.product_id" class="flex justify-between">
                  <span>{{ item.product_name }} <span class="text-slate-400">&times; {{ item.quantity }}</span></span>
                  <span class="font-medium text-slate-700">{{ formatRupiah(item.subtotal) }}</span>
                </div>
              </div>

              <!-- Actions (v-if PENDING) -->
              <div class="flex items-center justify-between pt-1 text-xs">
                <button @click="openModal(order)" class="text-emerald-700 hover:text-emerald-900 font-semibold flex items-center gap-1">
                  <i class="fa-solid fa-circle-info"></i> Detail & Pembayaran
                </button>
                
                <div class="flex items-center space-x-2">
                  <template v-if="order.status === 'PENDING'">
                    <button @click="simulatePayment(order.id)"
                      class="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-semibold transition flex items-center gap-1 shadow-sm">
                      <i class="fa-solid fa-money-bill-wave"></i> Simulasi Bayar
                    </button>
                    <button @click="cancelOrder(order.id)"
                      class="px-3 py-1.5 bg-rose-50 hover:bg-rose-100 text-rose-600 rounded-lg font-semibold transition border border-rose-200">
                      Batalkan
                    </button>
                  </template>
                </div>
              </div>

            </div>
          </div>

        </div>

      </div>

    </main>

    <!-- Modal Detail & Info Payment -->
    <div v-if="selectedOrder" class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-md w-full p-6 space-y-5 shadow-2xl">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 class="font-bold text-slate-900 text-base flex items-center gap-2">
            <i class="fa-solid fa-receipt text-emerald-600"></i> Detail Order & Gateway
          </h3>
          <button @click="selectedOrder = null" class="text-slate-400 hover:text-slate-600 text-lg">
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>

        <div class="space-y-4 text-sm">
          <div class="bg-slate-50 p-3 rounded-xl border border-slate-200 space-y-1.5 text-xs">
            <div class="flex justify-between"><span class="text-slate-500">ID Order:</span><span class="font-mono font-bold text-slate-800">{{ selectedOrder.id }}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Customer:</span><span class="font-bold text-slate-800">{{ selectedOrder.customer_name }}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Email:</span><span class="text-slate-700">{{ selectedOrder.customer_email }}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Status:</span><span class="font-bold text-emerald-700">{{ selectedOrder.status }}</span></div>
            <div class="flex justify-between"><span class="text-slate-500">Waktu Order:</span><span class="text-slate-700">{{ formatDate(selectedOrder.created_at) }}</span></div>
          </div>

          <div class="border border-emerald-100 bg-emerald-50/50 p-3.5 rounded-xl space-y-2">
            <span class="text-xs font-bold text-emerald-950 uppercase tracking-wide flex items-center gap-1.5">
              <i class="fa-solid fa-shield-halved text-emerald-600"></i> Payment Gateway Info
            </span>
            <div class="text-xs space-y-1">
              <p><span class="text-slate-500">Payment ID:</span> <span class="font-mono text-slate-800 font-medium">{{ selectedOrder.payment_id || '-' }}</span></p>
              <p><span class="text-slate-500">Metode Bayar:</span> <span class="font-medium uppercase">{{ selectedOrder.payment_method || 'QRIS / Virtual Account' }}</span></p>
              <p v-if="selectedOrder.payment_url" class="pt-1">
                <a :href="selectedOrder.payment_url" target="_blank" class="inline-flex items-center text-emerald-700 hover:text-emerald-900 font-semibold underline text-xs">
                  Buka Sandbox Payment URL <i class="fa-solid fa-arrow-up-right-from-square ml-1 text-[10px]"></i>
                </a>
              </p>
            </div>
          </div>

          <div>
            <span class="text-xs font-bold text-slate-700 uppercase block mb-1.5">Item Belanja:</span>
            <div class="space-y-1.5 text-xs">
              <div v-for="item in selectedOrder.items" :key="item.product_id" class="flex justify-between py-1 border-b border-slate-100">
                <span>{{ item.product_name }} <span class="text-slate-400">({{ item.quantity }} &times; {{ formatRupiah(item.price) }})</span></span>
                <span class="font-semibold text-slate-800">{{ formatRupiah(item.subtotal) }}</span>
              </div>
              <div class="flex justify-between pt-2 text-sm font-bold text-slate-900">
                <span>Total Tagihan:</span>
                <span class="text-emerald-700">{{ formatRupiah(selectedOrder.total_amount) }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="border-t border-slate-100 pt-3 flex justify-end">
          <button @click="selectedOrder = null" class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-sm font-semibold transition">
            Tutup
          </button>
        </div>
      </div>
    </div>

    <!-- Toast Notification (Vue Transition) -->
    <div v-if="toast.show"
      :class="[
        'fixed bottom-5 right-5 z-50 flex items-center space-x-3 px-4 py-3 rounded-xl shadow-lg text-sm font-medium text-white transition-all duration-300',
        toast.type === 'error' ? 'bg-rose-600' : 'bg-emerald-600'
      ]">
      <i :class="['fa-solid', toast.type === 'error' ? 'fa-circle-exclamation' : 'fa-circle-check']"></i>
      <span>{{ toast.message }}</span>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'

const API_BASE = window.location.origin.includes('localhost') || window.location.origin.includes('127.0.0.1')
  ? `${window.location.origin}/api/v1`
  : 'http://localhost:1323/api/v1';

// Reactive State
const apiConnected = ref(false);
const loading = ref(false);
const submitting = ref(false);
const filterStatus = ref('ALL');
const selectedOrder = ref(null);
const orders = ref([]);

// Toast State
const toast = reactive({
  show: false,
  message: '',
  type: 'success'
});

// Katalog Produk Bawaan
const products = ref([
  { id: 'PROD-001', name: 'Kopi Arabika 250g', price: 75000, stock: 50 },
  { id: 'PROD-002', name: 'Syrup Vanilla 500ml', price: 45000, stock: 20 },
  { id: 'PROD-003', name: 'Paper Cup 12oz (Pack of 50)', price: 35000, stock: 100 }
]);

// Form State
const form = reactive({
  customer_name: '',
  customer_email: '',
  items: {
    'PROD-001': 1,
    'PROD-002': 0,
    'PROD-003': 0
  }
});

// Computed Total Pembayaran
const computedOrderTotal = computed(() => {
  let sum = 0;
  for (const prod of products.value) {
    const qty = form.items[prod.id] || 0;
    sum += qty * prod.price;
  }
  return sum;
});

// Computed Stats
const stats = computed(() => {
  const total = orders.value.length;
  const pending = orders.value.filter(o => o.status === 'PENDING').length;
  const paid = orders.value.filter(o => o.status === 'PAID').length;
  const revenue = orders.value
    .filter(o => o.status === 'PAID')
    .reduce((acc, o) => acc + (o.total_amount || 0), 0);
  return { total, pending, paid, revenue };
});

// Computed Filtered Orders
const filteredOrders = computed(() => {
  if (filterStatus.value === 'ALL') return orders.value;
  return orders.value.filter(o => o.status === filterStatus.value);
});

// Helper Toast
const triggerToast = (msg, type = 'success') => {
  toast.message = msg;
  toast.type = type;
  toast.show = true;
  setTimeout(() => {
    toast.show = false;
  }, 3500);
};

// Format Helpers
const formatRupiah = (num) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num || 0);
};

const formatDate = (dateStr) => {
  if (!dateStr) return '-';
  const d = new Date(dateStr);
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
};

const getStatusBadgeClass = (status) => {
  if (status === 'PAID') return 'bg-emerald-50 text-emerald-700 border-emerald-200';
  if (status === 'CANCELLED') return 'bg-rose-50 text-rose-700 border-rose-200';
  return 'bg-amber-50 text-amber-700 border-amber-200';
};

const getStatusIcon = (status) => {
  if (status === 'PAID') return 'fa-circle-check';
  if (status === 'CANCELLED') return 'fa-circle-xmark';
  return 'fa-clock';
};

// Fetch Orders from Go Echo Backend
const fetchOrders = async () => {
  loading.value = true;
  try {
    const res = await fetch(`${API_BASE}/orders`);
    if (!res.ok) throw new Error('Status: ' + res.status);
    const json = await res.json();
    orders.value = json.data || [];
    apiConnected.value = true;
  } catch (err) {
    console.warn('API error:', err);
    apiConnected.value = false;
    // Fallback demo items
    if (orders.value.length === 0) {
      orders.value = [
        {
          id: "ORD-20261119-DEMO01",
          customer_name: "Ahmad Dahlan",
          customer_email: "ahmad@example.com",
          total_amount: 195000,
          status: "PENDING",
          payment_id: "PAY-DEMO-991",
          payment_url: "/sandbox/pay/PAY-DEMO-991?order_id=ORD-20261119-DEMO01&amount=195000&customer=Ahmad%20Dahlan",
          created_at: new Date().toISOString(),
          items: [
            { product_id: "PROD-001", product_name: "Kopi Arabika 250g", price: 75000, quantity: 2, subtotal: 150000 },
            { product_id: "PROD-002", product_name: "Syrup Vanilla 500ml", price: 45000, quantity: 1, subtotal: 45000 }
          ]
        }
      ];
    }
  } finally {
    loading.value = false;
  }
};

// Submit New Order
const handleCreateOrder = async () => {
  if (computedOrderTotal.value === 0) {
    triggerToast('Pilih minimal satu kuantitas produk!', 'error');
    return;
  }

  const itemsPayload = [];
  for (const [prodId, qty] of Object.entries(form.items)) {
    if (qty > 0) {
      itemsPayload.push({
        product_id: prodId,
        quantity: qty
      });
    }
  }

  submitting.value = true;
  try {
    const res = await fetch(`${API_BASE}/orders`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        customer_name: form.customer_name,
        customer_email: form.customer_email,
        items: itemsPayload
      })
    });

    const json = await res.json();
    if (!res.ok || !json.success) throw new Error(json.message || 'Gagal membuat order');

    triggerToast('Order berhasil dibuat! Menunggu pembayaran.');
    form.customer_name = '';
    form.customer_email = '';
    form.items['PROD-001'] = 1;
    form.items['PROD-002'] = 0;
    form.items['PROD-003'] = 0;
    await fetchOrders();
  } catch (err) {
    triggerToast(err.message, 'error');
  } finally {
    submitting.value = false;
  }
};

// Simulate Payment
const simulatePayment = async (orderId) => {
  try {
    const res = await fetch(`${API_BASE}/orders/${orderId}/pay`, { method: 'POST' });
    const json = await res.json();
    if (!res.ok || !json.success) throw new Error(json.message || 'Gagal simulasi bayar');
    triggerToast(`Order ${orderId} berhasil dibayar!`);
    await fetchOrders();
  } catch (err) {
    const o = orders.value.find(item => item.id === orderId);
    if (o) {
      o.status = 'PAID';
      triggerToast(`Order ${orderId} disimulasikan LUNAS`);
    } else {
      triggerToast(err.message, 'error');
    }
  }
};

// Cancel Order
const cancelOrder = async (orderId) => {
  if (!confirm(`Batalkan pesanan ${orderId}?`)) return;
  try {
    const res = await fetch(`${API_BASE}/orders/${orderId}/cancel`, { method: 'PUT' });
    const json = await res.json();
    if (!res.ok || !json.success) throw new Error(json.message || 'Gagal membatalkan order');
    triggerToast(`Order ${orderId} telah dibatalkan`);
    await fetchOrders();
  } catch (err) {
    const o = orders.value.find(item => item.id === orderId);
    if (o) {
      o.status = 'CANCELLED';
      triggerToast(`Order ${orderId} dibatalkan`);
    } else {
      triggerToast(err.message, 'error');
    }
  }
};

const openModal = (order) => {
  selectedOrder.value = order;
};

onMounted(() => {
  fetchOrders();
});
</script>

<style scoped>
/* App.vue Component Styles */
</style>
