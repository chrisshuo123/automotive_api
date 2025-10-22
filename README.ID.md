# Automotive APIs
Projek ini adalah untuk membuktikan kemampuan saya dalam membuat Pointer API menggunakan Golang MVC dan ditampilkan lewat HTML Template Frontend serta Postman API.
Saya akan update mengenai caranya untuk mengsetup Postman API serta MySQL DB menggunakan XAMPP.

## Rilis Mendatang:
- Aset 3D per mobil lowpoly 3D (diawali dari models, lalu dapat dimainkan lewat UNITY 3D. Catatan: Hanya Aset Lowpoly ya) ⏹🚗🚘

## 🚧 Learning Project: [Automotive API]

**Status:** Dalam perkembangan - Fitur inti telah tuntas, optimasi sedang berlangsung

**Apa yang dapat bekerja:**
- ✅ Migrasi & Seeding Database
- ✅ HTML FrontEnd API
- ✅ RESTful API end point

**Sedang Belajar:**
- 🔄 Menambahkan fitur Update di Postman API
- 🔄 Dokumentasi API
- 🔄 Golang Approval Workflow

**Perbaikan di Masa Depan:**
- 💡 Menambahkan Image pada setiap API Kendaraan
- 💡 Menambahkan 'Approval Workflow' _(saat user menambahkan atau mengupdate data baru, database perlu divalidasi oleh admin sebelum dikirim ke DB MySQL)_
- 💡 Menambahkan Model Mobil 3D Lowpoly _(Mengapa Lowpoly? Karena kami merupakan tim kecil)_

**Tech Stack:** Golang, MySQL, HTML, CSS, JS

# Instruksi

Pada halaman IDE VSCode dibagian Terminal, input sebagai berikut:
"git pull https://github.com/chrisshuo123/automotive_api.git main"
_Jika tidak bisa: Download Repository kami_

## PART 1: Impor Database & Download XAMPP

**Sebelum menjalankan server, bisa download XAMPP terlebih dahulu _(sesuaikan versi komputermu: windows, Linux, atau Mac)_ ikuti petunjuk lengkapnya** <br>
<br>
**Setelah instalasi selesai, bisa klik "start" di bagian Apache dan MySQL, lanjut ke browser dan:**
1. input "localhost/phpmyadmin/" pada url
2. Panel kiri, klik "new", berikan nama "automotive_api"
3. Pada scheme automotive_api dalam mode terpilih, di atas panel klik "import"
4. Klik "choose file", lalu pada direktori backend, klik "automotive_api.sql". _Note: jika toggle Partial Import, Other Options yg Foreign Key (FK) Checks nyala, bisa dimatikan togglenya._ Lalu, klik "import"<br>
<br>
**Kembali ke Vscode, pada vscode terminal jalankan server _pakai Golang_, input ini di terminal:**
1. "cd backend"
2. "go run main.go"
Jika benar, maka akan menampilkan "ECHO" yg tulisannya besar _(DISC: jika tulisan error di DB, abaikan saja)_

**NOTE:** Golang jika terdapat error, maka ini langkahnya, ke terminal di vscode:
1. cd backend
2. go mod tidy
_Tujuan dari "go mod tidy" adalah mengimport package .go yg esensial dalam menjalankan server dari Github Community agar Golang dapat berjalan dengan lancar_

## PART 2: Postman API Link
Jika belum memiliki PostmanAPI, dapat didownload dulu via link:
https://www.postman.com/downloads/ 

#### Usai membuka Postman, silahkan membuat sebuah file, pilih jenis pointer (GET atau PUT), lalu copy & paste link tersedia kepada kolom URL:
👉 Pointer GET Cars: 
http://localhost:8000/api/cars
👉 Pointer GET brands:
http://localhost:8000/api/brands
👉 Pointer GET types:
http://localhost:8000/api/types
👉 Pointer PUT Update Cars:
http://localhost:8000/api/cars/_id_

**Disclaimer:** Fitur Update Cars masih dalam perkembangan 

## PART 3: Website Template: API & CRUD Editor panel
**Selama servernya running, dapat menuju ke direktori _frontend > views_, lalu:**
1. Klik kanan index.html pada left side bar VSCode, Open with live server, untuk tampilan list API Database yang dapat ditampilkan pada Website
2. Klik kanan addCar.html pada left side bar VSCode, Open with live server, disini merupakan CRUD Panel yang dapat melakukan Create, Update, Delete DB API via Frontend





