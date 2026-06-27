# Midtrans Payment Tester (Go + Gin)

Aplikasi backend testing tool untuk melakukan integrasi dan pengujian alur pembayaran Midtrans — mendukung **Snap Payment (Popup UI)**, **Core API (Direct Integration)**, **Transaction Tools**, serta pengujian **Webhook Notification (termasuk Recurring & GoPay Account Linking)**.

---

## Prasyarat (Prerequisites)

- **Go** versi 1.21 atau lebih baru.
- **ngrok** (opsional, diperlukan untuk menguji notifikasi webhook dari server Midtrans ke localhost Anda).

---

## Panduan Penggunaan & Pengujian Step-by-Step

### Langkah 1: Konfigurasi Environment (`.env`)

1. Buat atau buka file [`.env`](file:///Users/bryanathallah/Bryan/suddenly/go-midtrans/.env) di root direktori project.
2. Isi file tersebut dengan Key dari dashboard Midtrans Anda (bisa didapatkan di **Settings > Access Keys**):

```env
# Server Key & Client Key dari Dashboard Midtrans Anda
MIDTRANS_SERVER_KEY=YOUR_KEY
MIDTRANS_CLIENT_KEY=YOUR_KEY

# Set ke false untuk Sandbox (Uji coba) atau true untuk Production (Live)
MIDTRANS_IS_PRODUCTION=false
```

> [!IMPORTANT]
> Pastikan nilai **`MIDTRANS_SERVER_KEY`** dan **`MIDTRANS_CLIENT_KEY`** diisi tepat sama seperti yang tertera di dashboard Midtrans Anda (jangan ditambahkan awalan `SB-` secara manual jika di dashboard Anda tidak ada awalan tersebut).

---

### Langkah 2: Jalankan Server Go

1. Buka terminal di direktori project `go-midtrans`.
2. Unduh dependencies dan jalankan aplikasi:
    ```bash
    go mod tidy
    go run main.go
    ```
3. Server akan berjalan di port `8080`. Buka browser Anda dan akses:
   **[http://localhost:8080](http://localhost:8080)**

---

### Langkah 3: Setup Webhook Notification (Menggunakan ngrok)

Agar server Midtrans di internet dapat mengirimkan notifikasi status pembayaran (webhook) ke server lokal (`localhost:8080`) Anda, Anda perlu menggunakan **ngrok** sebagai terowongan (tunnel):

1. Buka terminal baru, lalu jalankan ngrok:
    ```bash
    ngrok http 8080
    ```
2. Salin URL Forwarding HTTPS yang diberikan oleh ngrok (contoh: `https://abcd-123-45.ngrok-free.app`).
3. Masuk ke **Dashboard Midtrans** dan tempelkan URL tersebut ke konfigurasi notifikasi:
    - **Payment Notification URL** (di _Settings > Payment > Notification URL_):
      Isi dengan: `https://abcd-123-45.ngrok-free.app/notification`
    - **Recurring Payment Notification URL** (di _Settings > Recurring_:
      Isi dengan: `https://abcd-123-45.ngrok-free.app/notification`
    - **Account Linking Notification URL** (di _Settings > GoPay Account Linking_):
      Isi dengan: `https://abcd-123-45.ngrok-free.app/notification`

---

### Langkah 4: Uji Coba Pembayaran via SNAP (Popup UI) — _Direkomendasikan_

1. Buka Halaman Tester di browser Anda: **`http://localhost:8080`**.
2. **Isi Informasi Transaksi** di bagian atas (kosongkan kolom **Order ID** agar server Go membuat ID unik secara otomatis).
3. Scroll ke bagian **Snap Payment** dan klik tombol **"Pay via Snap (Popup)"**.
4. Layar popup pilihan pembayaran Midtrans akan muncul.
5. Pilih metode pembayaran di dalam popup tersebut, misalnya **Virtual Account > BCA**.
6. Catat nomor Virtual Account (VA) yang ditampilkan di layar popup.

---

### Langkah 5: Simulasi Pembayaran Sukses (Sandbox Simulator)

Karena menggunakan mode Sandbox, Anda bisa membayar secara gratis menggunakan simulator:

- **Untuk Virtual Account (BCA, BNI, BRI, Permata)**:
    1. Buka website **[Midtrans Sandbox Simulator - Virtual Account](https://simulator.sandbox.midtrans.com/bca/va/index)**.
    2. Masukkan nomor Virtual Account yang Anda salin pada **Langkah 4**.
    3. Klik **Inquire**, kemudian klik **Pay**.
- **Untuk QRIS / Gopay / ShopeePay**:
    1. Buka website **[Midtrans Sandbox Simulator - QRIS](https://simulator.sandbox.midtrans.com/qris/index)**.
    2. Unggah/upload screenshot QR Code yang Anda peroleh, lalu selesaikan simulasi pembayarannya.

---

### Langkah 6: Verifikasi Status Transaksi & Notifikasi Webhook

Setelah simulasi pembayaran berhasil dilakukan:

1. **Cek Notifikasi Webhook di Terminal**:
   Periksa log terminal tempat Anda menjalankan server Go. Webhook notifikasi sukses dari Midtrans akan otomatis tercetak di sana:
    ```
    [Notification] ✅ PAID (capture+accept): ORDER-xxxx
    ```
2. **Cek Status Manual via UI Tester**:
   Kembali ke halaman tester `http://localhost:8080`, gulir ke bagian paling bawah (**Transaction Tools**), pastikan Order ID yang sesuai sudah terisi, lalu klik **"Check Status"**. Status transaksi JSON akan berubah menjadi `"transaction_status": "settlement"` (Sukses Terbayar).

---

## Catatan Penting Mengenai Uji Coba Direct Core API

Jika Anda mencoba mengklik tombol pembayaran langsung (seperti **BCA VA**, **Mandiri Bill**, atau **QRIS** di bawah kategori bank transfer/ewallet) dan mendapatkan error **402 Payment channel is not activated**:

- Hal ini dikarenakan secara default, akses **Direct Core API (Direct Integration)** dinonaktifkan oleh Midtrans untuk akun baru demi kepatuhan keamanan kartu/finansial.
- Akses Direct API tersebut hanya bisa dibuka dengan menghubungi pihak Support Midtrans.
- Untuk keperluan pengujian standar, **Snap Payment (Popup)** sudah mencakup semua metode pembayaran tersebut dan bisa digunakan secara langsung tanpa kendala.

---

## Daftar Kartu Kredit Uji Coba (Sandbox Test Cards)

Gunakan detail kartu berikut untuk menguji pembayaran kartu kredit di mode Sandbox:

| Nomor Kartu           | CVV   | Tanggal Kadaluarsa                | Hasil Simulasi              |
| :-------------------- | :---- | :-------------------------------- | :-------------------------- |
| `4811 1111 1111 1114` | `123` | Bulan & Tahun Bebas di Masa Depan | ✅ Success                  |
| `4911 1111 1111 1113` | `123` | Bulan & Tahun Bebas di Masa Depan | ❌ Failure                  |
| `4411 1111 1111 1118` | `123` | Bulan & Tahun Bebas di Masa Depan | ⚠️ Challenge (Butuh Review) |
