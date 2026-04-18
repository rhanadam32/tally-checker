# 📊 Tally Checker - Excel Database System

Aplikasi berbasis web sederhana menggunakan bahasa pemrograman Go (Golang) untuk melakukan pengecekan data, input tally-in, dan perhitungan selisih otomatis menggunakan file Excel sebagai database.

## 🚀 Fitur
- **Pencarian Data**: Cari data berdasarkan ID, Nama, atau Keterangan.
- **Munculkan Semua Data**: Menampilkan seluruh isi database excel saat pertama kali dibuka.
- **Input Tally-In**: Menginput data secara langsung dari web ke dalam file Excel.
- **Kalkulasi Otomatis**: Menghitung selisih antara jumlah target dan jumlah tally-in secara otomatis.
- **Desain Modern**: Menggunakan Tailwind CSS dengan estetika minimalis ala Shadcn UI.

## 🛠️ Prasyarat
Sebelum menjalankan program, pastikan Anda sudah menginstal:
- [Go (Golang)](https://go.dev/dl/) (Versi terbaru direkomendasikan)
- File Excel bernama `database.xlsx` (lihat bagian struktur database)

## 📂 Struktur Database Excel
Pastikan file `database.xlsx` memiliki satu sheet bernama **`Sheet1`** dengan struktur kolom sebagai berikut:

| Kolom A (ID) | Kolom B (Target Qty) | Kolom C (Tally-In) | Kolom D (Selisih) |
| :--- | :--- | :--- | :--- |
| ID-001 | 100 | (Kosongkan) | (Otomatis) |
| ID-002 | 50 | (Kosongkan) | (Otomatis) |

*Catatan: Baris pertama harus berupa header.*

## 🏁 Langkah-Langkah Menjalankan

### 1. Persiapkan Folder Proyek
Ekstrak atau buat struktur folder seperti berikut:
```text
.
├── main.go
├── database.xlsx
└── templates/
    └── index.html
```

### 2. Inisialisasi Modul Go
Buka terminal/command prompt di folder proyek, lalu jalankan perintah berikut:
```bash
go mod init tally-checker
```

### 3. Install Library Excelize
Program ini menggunakan library `excelize` untuk mengelola file Excel. Install dengan perintah:
```bash
go get github.com/xuri/excelize/v2
```

### 4. Jalankan Aplikasi
Jalankan program dengan perintah:
```bash
go run main.go
```

### 5. Akses Aplikasi
Buka browser favorit Anda dan akses alamat berikut:
👉 **`http://localhost:8080`**

## 📖 Cara Penggunaan
1. **Melihat Semua Data**: Cukup buka halaman utama, semua data dari Excel akan muncul.
2. **Mencari Data**: Ketik kata kunci pada kolom pencarian dan klik tombol **"Cari Data"**.
3. **Input Tally**: 
   - Masukkan angka pada kolom **Tally-In (Input)** di baris data yang diinginkan.
   - Klik tombol **"Simpan"**.
   - Program akan otomatis menyimpan angka tersebut ke Excel dan menghitung selisihnya.
4. **Reset Tampilan**: Klik tombol **"Munculkan Data"** untuk kembali menampilkan seluruh data tanpa filter pencarian.
