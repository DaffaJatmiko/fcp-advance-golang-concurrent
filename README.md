# Advanced Student Portal: Concurrent Request Handling

## Description

Proyek ini merupakan pengembangan dari Final Project pada Dasar Pemrograman Backend, yaitu Student Portal 3, dengan beberapa tambahan fitur sebagai berikut:

1. **Import Student**: Registrasi banyak data mahasiswa sekaligus menggunakan teknik concurrency.
2. **Submit Assignment**: Mahasiswa bisa submit tugas menggunakan teknik job queue.
3. **Update Welcome Message on Login**: Menambahkan program studi di welcome message pada function Login.
4. **Maximum Login Attempts**: Melakukan pengecekan percobaan login dengan ID yang sama, jika sudah melebihi batas maksimum, ID tersebut akan diblokir dari proses login.

## Features

### 1. Import Student

Registrasi banyak data mahasiswa sekaligus dengan teknik concurrency. Data mahasiswa diambil dari file CSV dan didaftarkan ke dalam sistem dengan cepat dan tepat.

### 2. Submit Assignment

Pengiriman tugas oleh mahasiswa menggunakan teknik job queue untuk meng-handle pengiriman tugas yang besar. Menggunakan 3 goroutine untuk meng-handle pengiriman tugas dari siswa secara efisien.

### 3. Update Welcome Message on Login

Menambahkan informasi program studi pada welcome message setelah login berhasil.

### 4. Maximum Login Attempts

Membatasi percobaan login maksimum 3 kali. Jika melebihi batas tersebut, ID mahasiswa akan diblokir dari proses login.

## Project Structure

```bash
.
├── README.md
├── go.mod
├── go.sum
├── helper
│ └── helper.go
├── main.go
├── model
│ └── model.go
├── repository
│ ├── class.go
│ ├── session.go
│ ├── student.go
│ └── user.go
├── service
│ ├── class.go
│ ├── session.go
│ ├── student.go
│ └── user.go
├── students1.csv
├── students2.csv
└── students3.csv
```

## How to Run

1. Clone repository ini.
2. Pastikan Go (Golang) sudah terinstall di sistem Anda.
3. Jalankan perintah berikut untuk menjalankan aplikasi:

```bash
go run main.go
```

## Usage

Berikut adalah beberapa fitur utama yang dapat digunakan:

1. Login
2. Register
3. Get Study Program
4. Modify Student
5. Bulk Import Student
6. Submit Assignment
7. Exit
