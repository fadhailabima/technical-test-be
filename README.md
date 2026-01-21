# Inventory Management API

API Technical Test Junior Fullstack - Sistem Manajemen Inventori dengan Go, PostgreSQL, dan UUID

## 📋 Daftar Isi

- [Deskripsi](#deskripsi)
- [Teknologi](#teknologi)
- [Instalasi](#instalasi)
- [Setup Database](#setup-database)
- [Fitur yang Sudah Dikerjakan](#fitur-yang-sudah-dikerjakan)
- [Endpoint API](#endpoint-api)
- [Default User](#default-user)

## 📝 Deskripsi

Aplikasi backend untuk manajemen inventori yang mendukung multi-role (Admin, Seller, Pelanggan) dengan fitur marketplace, manajemen produk, dan sistem transaksi.

## 🛠 Teknologi

- **Backend**: Go 1.25.5
- **Framework**: Gin (Web Framework)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Token)
- **Password Hashing**: bcrypt
- **Documentation**: Swagger/OpenAPI
- **UUID**: google/uuid

## 🚀 Instalasi

### Prerequisites

Pastikan sistem Anda sudah terinstall:

- Go 1.25.5 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi
- Git

### Langkah-langkah Instalasi

1. **Clone Repository**

   ```bash
   git clone <repository-url>
   cd technical-test
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   go mod tidy
   ```

3. **Setup Environment Variables**

   Buat file `.env` di root project dengan konfigurasi berikut:

   ```env
   DB_HOST=localhost
   DB_USER=your_postgres_user
   DB_PASSWORD=your_postgres_password
   DB_NAME=technical_test_db
   DB_PORT=5432
   DB_SSLMODE=disable
   DB_TIMEZONE=Asia/Jakarta

   SERVER_PORT=8080
   JWT_SECRET=your_secret_key_here_make_it_long_and_secure
   ```

## 🗄 Setup Database

### 1. Install PostgreSQL

**macOS (menggunakan Homebrew):**

```bash
brew install postgresql@15
brew services start postgresql@15
```

**Linux (Ubuntu/Debian):**

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### 2. Buat Database

Login ke PostgreSQL:

```bash
psql -U postgres
```

Buat database baru:

```sql
CREATE DATABASE technical_test_db;
```

Buat user baru (opsional):

```sql
CREATE USER your_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE technical_test_db TO your_user;
```

Keluar dari psql:

```sql
\q
```

### 3. Jalankan Aplikasi

```bash
go run main.go
```

Aplikasi akan:

- Terhubung ke database PostgreSQL
- Melakukan auto-migration untuk semua tabel
- Seeding data awal (Roles, Product Types, Admin User)
- Menjalankan server di port 8080 (atau sesuai `.env`)

### 4. Verifikasi Database

Server akan menampilkan pesan:

```
✅ Migrasi Database Berhasil!
✅ Data Roles Berhasil Dibuat!
✅ Data Product Types Berhasil Dibuat!
✅ Super Admin Berhasil Dibuat! (Email: admin@example.com / Pass: admin123)
```

Server berjalan di: `http://localhost:8080`

## ✨ Fitur yang Sudah Dikerjakan

### 1. **Authentication & Authorization**

- ✅ Register user dengan role (Admin, Seller, Pelanggan)
- ✅ Login dengan JWT token
- ✅ Middleware autentikasi
- ✅ Middleware role-based authorization

### 2. **User Management**

- ✅ Model User dengan UUID
- ✅ Password hashing dengan bcrypt
- ✅ Relasi User dengan Role
- ✅ Auto-seeding Super Admin default

### 3. **Product Management (Master/Gudang Pusat)**

- ✅ Model Product dengan UUID
- ✅ Product Type categorization
- ✅ Stock management
- ✅ CRUD Product (untuk Admin)

### 4. **Seller Catalog (Marketplace)**

- ✅ Model SellerProduct untuk etalase seller
- ✅ Seller dapat add produk ke marketplace dengan markup harga
- ✅ Sistem harga: Base Price + Markup = Selling Price

### 5. **Transaction Management**

- ✅ Model Transaction untuk order
- ✅ Pelanggan dapat membuat pesanan
- ✅ Seller dapat mengonfirmasi pesanan
- ✅ Status tracking (Pending, Confirmed)
- ✅ Stock reduction otomatis saat konfirmasi

### 6. **Database**

- ✅ PostgreSQL integration
- ✅ GORM ORM
- ✅ Auto-migration
- ✅ Seeding data awal
- ✅ UUID sebagai Primary Key
- ✅ Soft Delete implementation

### 7. **API Documentation**

- ✅ Swagger/OpenAPI integration
- ✅ Dokumentasi endpoint otomatis
- ✅ Akses via `/swagger/index.html`

### 8. **CORS Configuration**

- ✅ CORS middleware untuk frontend integration

## 📡 Endpoint API

### 🌐 Base URL

```
http://localhost:8080
```

### 🔒 Authentication

Semua endpoint kecuali `/auth/*` dan `/swagger/*` memerlukan JWT token di header:

```
Authorization: Bearer <your_jwt_token>
```

### ⚠️ CORS

API sudah dikonfigurasi untuk menerima request dari origin manapun (`*`). Frontend bisa langsung consume API tanpa masalah CORS.

---

## 📋 List Endpoint

### 📋 Summary

Total **22 Endpoints** tersedia:

- **3** Authentication endpoints (Public)
- **4** Product Management endpoints (Admin)
- **2** Product Types endpoints
- **1** Marketplace endpoint
- **5** Seller Catalog endpoints
- **3** Transaction endpoints
- **1** Dashboard endpoint (Multi-role)
- **3** User Management endpoints (Admin only)

---

### 🔓 Authentication (Public - No Auth Required)

#### 1. Get Available Roles

```
GET /auth/roles

Response 200:
{
  "roles": [
    {
      "id": "uuid",
      "name": "Seller"
    },
    {
      "id": "uuid",
      "name": "Pelanggan"
    }
  ]
}
```

#### 2. Register User

```
POST /auth/register
Content-Type: application/json

Body:
{
  "name": "string",
  "email": "string",
  "password": "string (min 6 chars)",
  "role_id": "uuid (Admin/Seller/Pelanggan)"
}

Response 201:
{
  "message": "Registrasi berhasil!"
}
```

#### 3. Login

```
POST /auth/login
Content-Type: application/json

Body:
{
  "email": "string",
  "password": "string"
}

Response 200:
{
  "token": "jwt_token_string",
  "user": {
    "id": "uuid",
    "name": "string",
    "role": "Admin|Seller|Pelanggan"
  }
}
```

---

### 📦 Products (Master/Gudang Pusat)

#### 1. Get All Products

```
GET /products
Authorization: Bearer <token>

Response 200:
{
  "data": [
    {
      "id": "uuid",
      "name": "string",
      "product_type_id": "uuid",
      "price": 0,
      "stock": 0,
      "created_at": "timestamp"
    }
  ]
}
```

#### 2. Create Product (Admin Only)

```
POST /products
Authorization: Bearer <admin_token>
Content-Type: application/json

Body:
{
  "name": "string",
  "product_type_id": "uuid",
  "price": 0,
  "stock": 0
}

Response 201:
{
  "data": { product object }
}
```

#### 3. Delete Product (Admin Only)

```
DELETE /products/:id
Authorization: Bearer <admin_token>

Response 200:
{
  "message": "Deleted"
}
```

#### 4. Update Product (Admin Only)

```
PUT /products/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

Body (all fields optional):
{
  "name": "string",
  "product_type_id": "uuid",
  "price": 0,
  "stock": 0
}

Response 200:
{
  "data": { updated product object }
}
```

---

### 📂 Product Types

#### 1. Get All Product Types

```
GET /product-types
Authorization: Bearer <token>

Response 200:
{
  "data": [
    {
      "id": "uuid",
      "name": "string"
    }
  ]
}
```

#### 2. Create Product Type (Admin Only)

```
POST /product-types
Authorization: Bearer <admin_token>
Content-Type: application/json

Body:
{
  "name": "string"
}

Response 201:
{
  "data": { product_type object }
}
```

---

### 🛍️ Marketplace

#### 1. Get Marketplace Items

```
GET /marketplace
Authorization: Bearer <token>

Response 200:
{
  "data": [
    {
      "seller_product_id": "uuid",
      "product_name": "string",
      "category": "string",
      "seller_name": "string",
      "price": 0,
      "stock_available": 0
    }
  ]
}
```

---

### 🛒 Seller Catalog

#### 1. Add Product to Marketplace (Seller Only)

```
POST /seller/products
Authorization: Bearer <seller_token>
Content-Type: application/json

Body:
{
  "product_id": "uuid",
  "selling_price": 0
}

Response 201:
{
  "data": { seller_product object }
}
```

#### 2. Get Seller Products (Seller Only)

```
GET /seller/products
Authorization: Bearer <seller_token>

Response 200:
{
  "data": [
    {
      "id": "uuid",
      "product_name": "string",
      "category": "string",
      "base_price": 0,
      "selling_price": 0,
      "profit_margin": 0,
      "stock": 0,
      "is_active": true
    }
  ]
}
```

#### 3. Update Seller Product (Seller Only)

```
PUT /seller/products/:id
Authorization: Bearer <seller_token>
Content-Type: application/json

Body (all fields optional):
{
  "selling_price": 0,
  "is_active": true
}

Response 200:
{
  "data": { updated seller_product object }
}
```

#### 4. Delete Seller Product (Seller Only)

```
DELETE /seller/products/:id
Authorization: Bearer <seller_token>

Response 200:
{
  "message": "Product removed from marketplace"
}
```

#### 5. Get Seller Transactions (Seller Only)

```
GET /seller/transactions
Authorization: Bearer <seller_token>

Response 200:
{
  "data": [
    {
      "id": "uuid",
      "product_name": "string",
      "buyer_name": "string",
      "buyer_email": "string",
      "quantity": 0,
      "total_price": 0,
      "seller_profit": 0,
      "status": "PENDING|CONFIRMED|CANCELLED",
      "created_at": "timestamp"
    }
  ]
}
```

---

### 💰 Transactions

#### 1. Create Order (Pelanggan Only)

```
POST /transactions
Authorization: Bearer <pelanggan_token>
Content-Type: application/json

Body:
{
  "seller_product_id": "uuid",
  "quantity": 1
}

Response 201:
{
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "seller_product_id": "uuid",
    "quantity": 1,
    "status": "PENDING",
    "total_price": 0,
    "admin_fee": 0,
    "seller_profit": 0
  }
}
```

#### 2. Confirm Order (Seller Only)

```
POST /transactions/:id/confirm
Authorization: Bearer <seller_token>

Response 200:
{
  "message": "Confirmed"
}
```

#### 3. Get Customer Transactions (Pelanggan Only)

```
GET /customer/transactions
Authorization: Bearer <pelanggan_token>

Response 200:
{
  "data": [
    {
      "id": "uuid",
      "product_name": "string",
      "seller_name": "string",
      "seller_email": "string",
      "quantity": 0,
      "total_price": 0,
      "admin_fee": 0,
      "status": "PENDING|CONFIRMED|CANCELLED",
      "created_at": "timestamp"
    }
  ]
}
```

---

### � Dashboard

#### 1. Get Dashboard Stats (All Roles)

```
GET /dashboard
Authorization: Bearer <token>

Response 200 (Admin):
{
  "role": "Admin",
  "data": {
    "total_products": 0,
    "total_product_types": 0,
    "total_sellers": 0,
    "total_customers": 0,
    "transactions_today": 0,
    "platform_income": 0
  }
}

Response 200 (Seller):
{
  "role": "Seller",
  "data": {
    "products_in_marketplace": 0,
    "total_sales_revenue": 0,
    "total_transactions": 0,
    "pending_orders": 0,
    "confirmed_orders": 0,
    "total_profit": 0,
    "profit_margin_percentage": 0,
    "top_products": [
      {
        "product_name": "string",
        "transaction_count": 0,
        "total_quantity": 0,
        "total_revenue": 0
      }
    ]
  }
}

Response 200 (Pelanggan):
{
  "role": "Pelanggan",
  "data": {
    "total_orders": 0,
    "pending_orders": 0,
    "confirmed_orders": 0,
    "total_spent": 0,
    "recent_orders": [
      {
        "id": "uuid",
        "product_name": "string",
        "seller_name": "string",
        "quantity": 0,
        "total_price": 0,
        "status": "PENDING|CONFIRMED",
        "created_at": "timestamp"
      }
    ]
  }
}
```

---

### 👥 User Management (Admin Only)

#### 1. Get All Users

```
GET /users
Authorization: Bearer <admin_token>

Response 200:
{
  "data": [
    {
      "id": "uuid",
      "name": "string",
      "email": "string",
      "role": {
        "id": "uuid",
        "name": "string"
      }
    }
  ]
}
```

#### 2. Create Admin User

```
POST /users/admin
Authorization: Bearer <admin_token>
Content-Type: application/json

Body:
{
  "name": "string",
  "email": "string",
  "password": "string (min 6 chars)"
}

Response 201:
{
  "message": "Admin created"
}
```

#### 3. Delete User

```
DELETE /users/:id
Authorization: Bearer <admin_token>

Response 200:
{
  "message": "User deleted"
}
```

---

### �📚 Documentation

#### Swagger UI

```
GET /swagger/index.html
Browser: http://localhost:8080/swagger/index.html
```

---

## 🔑 Cara Menggunakan API (Frontend Integration)

### 1. Get Roles untuk Registrasi

```javascript
// Fetch available roles
const rolesResponse = await fetch("http://localhost:8080/auth/roles");
const rolesData = await rolesResponse.json();

console.log(rolesData.roles);
// Output: [
//   { id: "uuid-seller", name: "Seller" },
//   { id: "uuid-pelanggan", name: "Pelanggan" }
// ]

// Gunakan role.id saat register
const selectedRoleId = rolesData.roles[0].id; // Seller
```

### 2. Register User

```javascript
const registerResponse = await fetch("http://localhost:8080/auth/register", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    name: "John Seller",
    email: "john@seller.com",
    password: "password123",
    role_id: selectedRoleId, // UUID dari endpoint /auth/roles
  }),
});

const registerData = await registerResponse.json();
console.log(registerData.message); // "Registrasi berhasil!"
```

### 3. Login & Simpan Token

```javascript
// Login
const response = await fetch("http://localhost:8080/auth/login", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    email: "admin@example.com",
    password: "admin123",
  }),
});

const data = await response.json();
const token = data.token;

// Simpan token (localStorage/sessionStorage/cookie)
localStorage.setItem("token", token);
```

### 4. Request dengan Token

```javascript
const token = localStorage.getItem("token");

const response = await fetch("http://localhost:8080/products", {
  method: "GET",
  headers: {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  },
});

const products = await response.json();
```

---

## 🎯 Role-Based Access

| Endpoint                       | Admin | Seller | Pelanggan |
| ------------------------------ | ----- | ------ | --------- |
| GET /products                  | ✅    | ✅     | ✅        |
| POST /products                 | ✅    | ❌     | ❌        |
| PUT /products/:id              | ✅    | ❌     | ❌        |
| DELETE /products/:id           | ✅    | ❌     | ❌        |
| GET /product-types             | ✅    | ✅     | ✅        |
| POST /product-types            | ✅    | ❌     | ❌        |
| GET /marketplace               | ✅    | ✅     | ✅        |
| POST /seller/products          | ❌    | ✅     | ❌        |
| GET /seller/products           | ❌    | ✅     | ❌        |
| PUT /seller/products/:id       | ❌    | ✅     | ❌        |
| DELETE /seller/products/:id    | ❌    | ✅     | ❌        |
| GET /seller/transactions       | ❌    | ✅     | ❌        |
| POST /transactions             | ❌    | ❌     | ✅        |
| POST /transactions/:id/confirm | ❌    | ✅     | ❌        |
| GET /customer/transactions     | ❌    | ❌     | ✅        |
| GET /dashboard                 | ✅    | ✅     | ✅        |
| GET /users                     | ✅    | ❌     | ❌        |
| POST /users/admin              | ✅    | ❌     | ❌        |
| DELETE /users/:id              | ✅    | ❌     | ❌        |

---

## 🔒 Protected Endpoints (Require Authentication)

Semua endpoint di atas kecuali `/auth/*` dan `/swagger/*` memerlukan header:

```
Authorization: Bearer <jwt_token>
```

---

## 👤 Default User

Setelah aplikasi pertama kali dijalankan, akan ada user admin default:

**Super Admin**

- Email: `admin@example.com`
- Password: `admin123`
- Role: Admin

**Catatan:** Segera ganti password setelah login pertama kali di production!

## 🗃 Database Schema

### Tables

- **roles** - Role management (Admin, Seller, Pelanggan)
- **product_types** - Kategori produk (Elektronik, Pakaian, Makanan)
- **users** - Data user dengan role
- **products** - Master produk (gudang pusat)
- **seller_products** - Katalog marketplace seller dengan markup
- **transactions** - Transaksi pembelian

## 🔧 Swagger Documentation

Akses dokumentasi API lengkap dengan Swagger UI:

```
http://localhost:8080/swagger/index.html
```

Untuk regenerate Swagger docs setelah update komentar API:

```bash
swag init
```

## 📝 Contoh Request

### 1. Get Available Roles

```bash
curl -X GET http://localhost:8080/auth/roles
```

### 2. Register User

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Seller",
    "email": "john@seller.com",
    "password": "password123",
    "role_id": "uuid-role-seller"
  }'
```

### 3. Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### 4. Get All Products (Gudang)

```bash
curl -X GET http://localhost:8080/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Create Product (Admin)

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Laptop ASUS ROG",
    "product_type_id": "uuid-product-type",
    "price": 15000000,
    "stock": 10
  }'
```

### 6. Get Product Types

```bash
curl -X GET http://localhost:8080/product-types \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 7. Get Marketplace

```bash
curl -X GET http://localhost:8080/marketplace \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 8. Add Product to Marketplace (Seller)

```bash
curl -X POST http://localhost:8080/seller/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "product_id": "uuid-product-id",
    "selling_price": 16000000
  }'
```

### 9. Get Seller Transactions (Seller)

```bash
curl -X GET http://localhost:8080/seller/transactions \
  -H "Authorization: Bearer YOUR_SELLER_TOKEN"
```

### 10. Create Order (Pelanggan)

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "seller_product_id": "uuid-seller-product-id",
    "quantity": 2
  }'
```

### 11. Confirm Order (Seller)

```bash
curl -X POST http://localhost:8080/transactions/uuid-transaction-id/confirm \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## � Deploy ke Railway

Railway adalah platform hosting modern yang mudah digunakan untuk deploy aplikasi Go. Berikut langkah-langkahnya:

### 1. Persiapan

Pastikan project sudah di-push ke GitHub repository.

### 2. Setup Railway

1. Buat akun di [Railway.app](https://railway.app)
2. Klik **"New Project"**
3. Pilih **"Deploy from GitHub repo"**
4. Pilih repository project ini
5. Railway akan otomatis detect Go project dan mulai build

### 3. Tambah PostgreSQL Database

1. Di Railway dashboard project, klik **"New"** → **"Database"** → **"Add PostgreSQL"**
2. Railway akan otomatis provision database PostgreSQL
3. Database credentials akan tersedia sebagai environment variables

### 4. Setup Environment Variables

Di Railway dashboard, buka **"Variables"** tab dan tambahkan:

```env
DB_HOST=${{Postgres.PGHOST}}
DB_USER=${{Postgres.PGUSER}}
DB_PASSWORD=${{Postgres.PGPASSWORD}}
DB_NAME=${{Postgres.PGDATABASE}}
DB_PORT=${{Postgres.PGPORT}}
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Jakarta

SERVER_PORT=${{PORT}}
JWT_SECRET=your_production_secret_key_here_make_it_very_long_and_secure
```

**Note:** Railway otomatis inject `PORT` variable, jadi gunakan `${{PORT}}` untuk `SERVER_PORT`.

### 5. Deploy

1. Railway akan otomatis deploy setiap kali ada push ke GitHub
2. Build process akan:
   - Download Go dependencies (`go mod download`)
   - Build binary (`go build -o main .`)
   - Run migration & seeding otomatis saat startup
   - Start server

### 6. Akses Aplikasi

Setelah deploy berhasil:

1. Railway akan memberikan public URL (contoh: `https://your-app.up.railway.app`)
2. Akses API di: `https://your-app.up.railway.app`
3. Swagger docs: `https://your-app.up.railway.app/swagger/index.html`

### 7. Monitoring & Logs

- **Logs**: Buka tab "Deployments" di Railway dashboard untuk melihat real-time logs
- **Metrics**: Railway menyediakan CPU, Memory, dan Network metrics
- **Redeploy**: Push ke GitHub atau klik "Redeploy" di Railway dashboard

### Troubleshooting Railway Deployment

#### Build Failed
- Pastikan `go.mod` dan `go.sum` sudah ter-commit
- Check build logs di Railway dashboard

#### Database Connection Error
- Pastikan PostgreSQL service sudah running
- Verifikasi environment variables menggunakan `${{Postgres.*}}` syntax

#### Migration Failed
- Check logs untuk error detail
- Pastikan database credentials benar
- Database akan otomatis di-migrate saat aplikasi pertama kali start

#### Port Binding Error
- Pastikan menggunakan `SERVER_PORT=${{PORT}}` bukan hardcode port
- Railway akan inject PORT variable secara otomatis

---

## 🐛 Troubleshooting (Local Development)

### Database Connection Error

- Pastikan PostgreSQL sudah berjalan
- Cek kredensial database di file `.env`
- Pastikan database sudah dibuat

### Port Already in Use

Ubah `SERVER_PORT` di file `.env` ke port lain:

```env
SERVER_PORT=8081
```

### JWT Token Invalid

- Pastikan `JWT_SECRET` di `.env` tidak berubah
- Token mungkin sudah expired, login ulang untuk mendapat token baru

## 📄 License

Copyright © 2026 Fadhail Athaillah Bima Dharmawan

---

**Contact:**

- Name: Fadhail Athaillah Bima Dharmawan
- Email: bimadharmawan6@gmail.com
