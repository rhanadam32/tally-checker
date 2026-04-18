# 📊 Tally Checker - Excel Database System

  Aplikasi berbasis web sederhana menggunakan bahasa pemrograman Go (Golang) untuk melakukan
  pengecekan data, input tally-in, dan perhitungan selisih otomatis menggunakan file Excel
  sebagai database.

  ## 🚀 Fitur Utama
  - **Pencarian Data Dinamis**: Cari data berdasarkan ID, Nama, atau Keterangan dengan respon
  cepat.
  - **Dashboard Master Data**: Menampilkan seluruh data master dari Excel secara otomatis saat
  aplikasi dibuka.
  - **Input Tally-In Real-time**: Menginput data tally-in langsung dari web ke file Excel tanpa
  reload halaman (AJAX).
  - **Kalkulasi Selisih Otomatis**: Menghitung selisih antara jumlah Tally-In dan Stok Master
  secara otomatis dengan rumus: `Tally-In - Stok Master`.
  - **Ekstraksi Stok Pintar**: Mampu membaca jumlah target dari format teks seperti `"Produk A =
   100 mc"` secara otomatis.
  - **UI Modern (Shadcn-like)**: Antarmuka minimalis dan profesional menggunakan Tailwind CSS
  dengan notifikasi toast.

  ## 🛠️ Prasyarat
  Sebelum menjalankan program, pastikan Anda sudah menginstal:
  - [Go (Golang)](https://go.dev/dl/) (Versi terbaru direkomendasikan)
  - File Excel bernama `database.xlsx`

  ## 📂 Struktur Database Excel
  Pastikan file `database.xlsx` memiliki satu sheet bernama **`Sheet1`** dengan struktur kolom
  sebagai berikut:

  | Kolom A (ID) | Kolom B (Produk & Target Qty) | Kolom C (Tally-In) | Kolom D (Selisih) |
  | :--- | :--- | :--- | :--- |
  | ID-001 | Produk A = 100 mc | (Kosongkan) | (Otomatis) |
  | ID-002 | Produk B = 50 mc | (Kosongkan) | (Otomatis) |

  *Catatan: Baris pertama harus berupa header. Format di Kolom B harus mengandung tanda `=` agar
   angka stok dapat diekstrak.*

  ## 🏁 Langkah-Langkah Menjalankan

  ### 1. Persiapkan Folder Proyek
  Pastikan struktur folder Anda seperti berikut:
  ```text
  .
  ├── main.go
  ├── database.xlsx
  └── templates/
      └── index.html

  2. Inisialisasi Modul Go

  Buka terminal di folder proyek, lalu jalankan:
  go mod init tally-checker

  3. Install Library Excelize

  Install library untuk manajemen file Excel:
  go get github.com/xuri/excelize/v2

  4. Jalankan Aplikasi

  Jalankan program dengan perintah:
  go run main.go

  5. Akses Aplikasi

  Buka browser dan akses alamat berikut:
  👉 http://localhost:8080

  📖 Cara Penggunaan

  1. Melihat Data: Saat dibuka, aplikasi akan langsung menampilkan seluruh daftar data master.
  2. Pencarian: Gunakan Panel Pencarian untuk memfilter data. Klik "Reset" untuk menampilkan
  kembali semua data.
  3. Input & Update:
    - Masukkan angka pada kolom Tally-In (Input).
    - Klik tombol "Simpan".
    - Sistem akan secara otomatis:
        - Menyimpan angka ke Kolom C di Excel.
      - Menghitung selisih (Tally-In dikurangi Stok Master).
      - Menyimpan hasil selisih ke Kolom D di Excel.
      - Memperbarui angka selisih di layar secara instan tanpa refresh halaman.
  4. Notifikasi: Perhatikan pojok kanan bawah untuk konfirmasi "Data berhasil disimpan!".

  📝 Log Pembaruan (Review)

  - v1.0: Implementasi dasar pembacaan Excel dan tampilan Bootstrap.
  - v1.1: Migrasi desain ke Tailwind CSS (Shadcn UI style) untuk tampilan lebih modern.
  - v1.2: Penambahan fitur "Munculkan Semua Data" dan perbaikan navigasi.
  - v1.3: Penambahan fitur "Stok Master" (ekstraksi angka otomatis dari teks kolom B).
  - v1.4: Implementasi Update AJAX (Simpan tanpa refresh), Notifikasi Toast, dan sinkronisasi
  kalkulasi selisih Tally-In - Stok Master.
  - v1.5: Perbaikan bug undefined: finalDiff dan optimasi redirect halaman.
