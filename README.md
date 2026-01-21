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
✅ Demo Users Berhasil Dibuat! (Password semua user: password123)
   - 2 Admin: admin@example.com, admin.staff@example.com
   - 3 Seller: seller1-3@example.com
   - 3 Pelanggan: customer1-3@example.com
✅ Sample Products Berhasil Dibuat! (24 produk)
```

Server berjalan di: `http://localhost:8080`

## ✨ Fitur yang Sudah Dikerjakan

### 1. **Authentication & Authorization**

- ✅ Register user dengan role (Seller, Pelanggan)
- ✅ Login dengan JWT token (expired 24 jam)
- ✅ Middleware autentikasi untuk validasi token
- ✅ Middleware role-based authorization (Admin, Seller, Pelanggan)
- ✅ Password hashing dengan bcrypt
- ✅ Proteksi registrasi Admin (tidak bisa register publik)

### 2. **User Management (Admin Only)**

- ✅ CRUD semua user
- ✅ Get all users dengan pagination
- ✅ Update user (name, email, role)
- ✅ Delete user (hard delete)
- ✅ Update profile sendiri (all roles)
- ✅ Change password

### 3. **Product Management - Gudang Pusat (Admin)**

- ✅ CRUD Product master
- ✅ Stock management (add/reduce stock)
- ✅ Low stock alerts (threshold: 10)
- ✅ Product dengan UUID sebagai ID
- ✅ Relasi dengan Product Type (kategorisasi)
- ✅ Validasi stok tidak boleh negatif

### 4. **Product Types (Admin)**

- ✅ CRUD Product Types
- ✅ Kategorisasi produk (Elektronik, Pakaian, Makanan, dll)
- ✅ Relasi one-to-many dengan Product

### 5. **Seller Catalog (Marketplace)**

- ✅ Seller add produk dari gudang pusat ke etalase
- ✅ Sistem markup harga (Selling Price = Base Price + Markup)
- ✅ Validasi harga jual >= harga modal
- ✅ CRUD seller products (get, update price, activate/deactivate, delete)
- ✅ Toggle active/inactive produk di marketplace

### 6. **Marketplace (Public with Search & Filter)**

- ✅ Browse all active products dari semua seller
- ✅ Search by product name (case insensitive)
- ✅ Filter by category (product type)
- ✅ Filter by price range (min-max)
- ✅ Menampilkan: product name, seller name, price, available stock

### 7. **Transaction Management**

- ✅ Customer create order (beli produk dari marketplace)
- ✅ Validasi stok tersedia saat order
- ✅ Validasi produk aktif saat order
- ✅ Kalkulasi otomatis: Total Price, Admin Fee, Seller Profit
- ✅ Seller confirm order (status: PENDING → COMPLETED)
- ✅ Auto stock reduction dari gudang pusat saat confirm
- ✅ Database locking untuk prevent race condition
- ✅ Customer cancel order (hanya status PENDING)
- ✅ Status tracking: PENDING, COMPLETED, CANCELLED
- ✅ Get customer transactions history
- ✅ Get seller transactions history
- ✅ Get transaction detail (full info buyer, seller, product)

### 8. **Dashboard (Multi-Role)**

- ✅ **Customer Dashboard:**
  - Total orders, pending orders, completed orders
  - Total spent (hanya transaksi COMPLETED)
  - Recent 3 orders
- ✅ **Seller Dashboard:**
  - Products in marketplace count
  - Total sales revenue (hanya COMPLETED)
  - Total transactions, pending orders, completed orders
  - Total profit, profit margin
  - Top 3 products (best sellers)
- ✅ **Admin Dashboard:**
  - Total products, total product types
  - Total sellers, total customers
  - Transactions today
  - Platform income (total admin fees dari COMPLETED transactions)

### 9. **Reports (Admin Only)**

- ✅ Sales report by date range (daily breakdown)
- ✅ Top products report (by quantity sold)
- ✅ Top sellers report (by total sales)
- ✅ Semua report hanya count transaksi COMPLETED
- ✅ Configurable limit (default 10)

### 10. **Database**

- ✅ PostgreSQL integration
- ✅ GORM ORM dengan relasi lengkap
- ✅ Auto-migration semua models
- ✅ UUID sebagai Primary Key (semua table)
- ✅ Hard Delete implementation (no soft delete)
- ✅ Seeding data awal:
  - 3 Roles: Admin, Seller, Pelanggan
  - 5 Product Types
  - 8 Demo Users (2 Admin, 3 Seller, 3 Pelanggan)
  - 24 Sample Products (berbagai kategori)
- ✅ Timestamps (CreatedAt, UpdatedAt) otomatis

### 11. **API Documentation**

- ✅ Swagger/OpenAPI 2.0 integration
- ✅ Interactive API docs di `/swagger/index.html`
- ✅ Dokumentasi lengkap 37 endpoints
- ✅ Request/Response examples
- ✅ Authentication bearer token support

### 12. **Security & Best Practices**

- ✅ JWT token dengan expiry
- ✅ Bcrypt password hashing (cost 10)
- ✅ SQL injection protection (GORM parameterized queries)
- ✅ CORS configuration (allow all origins)
- ✅ Error handling konsisten
- ✅ Input validation dengan Gin binding
- ✅ Database transaction untuk operasi kritis
- ✅ Komentar lengkap di kode untuk dokumentasi

### 13. **Code Quality**

- ✅ Struktur project terorganisir (MVC pattern)
- ✅ Service layer untuk business logic
- ✅ Controller layer untuk request handling
- ✅ Middleware untuk cross-cutting concerns
- ✅ Models terpisah untuk setiap entity
- ✅ Komentar detail menjelaskan alur logika
- ✅ Type-safe UUID handling
- ✅ Environment variables untuk konfigurasi

---

## 📊 Statistik Project

- **Total Endpoints**: 37
- **Total Models**: 7 (Role, User, ProductType, Product, SellerProduct, Transaction, Base)
- **Total Services**: 7 (Auth, User, Product, ProductType, Catalog, Transaction, Dashboard, Report)
- **Total Controllers**: 8
- **Total Middlewares**: 2 (Auth, Role)
- **Lines of Code**: ~3000+ (tanpa generated files)

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

Total **37 Endpoints** tersedia:

- **3** Authentication endpoints (Public)
- **3** User Profile endpoints
- **6** Product Management endpoints (Admin)
- **4** Product Types endpoints (Admin)
- **1** Marketplace endpoint (with search & filter)
- **5** Seller Catalog endpoints
- **5** Transaction endpoints
- **3** Reports endpoints (Admin only)
- **1** Dashboard endpoint (Multi-role)
- **6** User Management endpoints (Admin only)

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

### � User Profile (All Roles)

#### 1. Get Profile

```
GET /profile
Authorization: Bearer <token>

Response 200:
{
  "data": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "role": "Admin|Seller|Pelanggan"
  }
}
```

#### 2. Update Profile

```
PUT /profile
Authorization: Bearer <token>
Content-Type: application/json

Body (all fields optional):
{
  "name": "string",
  "email": "string"
}

Response 200:
{
  "message": "Profile updated successfully",
  "data": { updated user object }
}
```

#### 3. Change Password

```
PUT /profile/password
Authorization: Bearer <token>
Content-Type: application/json

Body:
{
  "old_password": "string",
  "new_password": "string (min 6 chars)"
}

Response 200:
{
  "message": "Password changed successfully"
}
```

---

### �📦 Products (Master/Gudang Pusat)

#### 1. Get All Products

```
GET /products?search=keyword&product_type_id=uuid
Authorization: Bearer <token>

Query Parameters (optional):
- search: Search by product name
- product_type_id: Filter by product type

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

#### 5. Get Low Stock Products (Admin Only)

```
GET /products/low-stock?threshold=10
Authorization: Bearer <admin_token>

Query Parameters:
- threshold: Stock threshold (default: 10)

Response 200:
{
  "threshold": 10,
  "count": 5,
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

#### 3. Update Product Type (Admin Only)

```
PUT /product-types/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

Body:
{
  "name": "string"
}

Response 200:
{
  "data": { updated product_type object }
}
```

#### 4. Delete Product Type (Admin Only)

```
DELETE /product-types/:id
Authorization: Bearer <admin_token>

Response 200:
{
  "message": "Product type deleted"
}
```

---

### 🛍️ Marketplace

#### 1. Get Marketplace Items

```
GET /marketplace?search=keyword&category=uuid&min_price=1000&max_price=50000
Authorization: Bearer <token>

Query Parameters (optional):
- search: Search by product name
- category: Filter by product type ID
- min_price: Minimum price filter
- max_price: Maximum price filter

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

#### 4. Get Transaction Detail

```
GET /transactions/:id
Authorization: Bearer <token>

Response 200:
{
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "buyer_name": "string",
    "buyer_email": "string",
    "seller_product_id": "uuid",
    "product_name": "string",
    "seller_name": "string",
    "seller_email": "string",
    "quantity": 0,
    "total_price": 0,
    "admin_fee": 0,
    "seller_profit": 0,
    "status": "PENDING|CONFIRMED|CANCELLED",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

#### 5. Cancel Transaction (Customer Only)

```
POST /transactions/:id/cancel
Authorization: Bearer <pelanggan_token>

Response 200:
{
  "message": "Transaction cancelled successfully"
}

Note: Only PENDING transactions can be cancelled
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

### � Reports (Admin Only)

#### 1. Sales Report

```
GET /reports/sales?start_date=2026-01-01&end_date=2026-01-31
Authorization: Bearer <admin_token>

Query Parameters:
- start_date: Start date (YYYY-MM-DD) - optional (default: last 30 days)
- end_date: End date (YYYY-MM-DD) - optional (default: today)

Response 200:
{
  "data": [
    {
      "date": "2026-01-31",
      "total_orders": 5,
      "total_revenue": 1500000,
      "admin_income": 75000,
      "seller_income": 1425000
    },
    {
      "date": "2026-01-30",
      "total_orders": 8,
      "total_revenue": 2400000,
      "admin_income": 120000,
      "seller_income": 2280000
    }
  ]
}

Note: Only COMPLETED transactions are included. Results ordered by date DESC.
```

#### 2. Top Products Report

```
GET /reports/top-products?limit=10
Authorization: Bearer <admin_token>

Query Parameters:
- limit: Number of products (default: 10)

Response 200:
{
  "limit": 10,
  "data": [
    {
      "product_name": "string",
      "category": "string",
      "total_sold": 50,
      "total_revenue": 15000000,
      "total_transactions": 25
    }
  ]
}

Note: Only COMPLETED transactions are counted. Sorted by total_sold DESC.
```

#### 3. Top Sellers Report

```
GET /reports/top-sellers?limit=10
Authorization: Bearer <admin_token>

Query Parameters:
- limit: Number of sellers (default: 10)

Response 200:
{
  "limit": 10,
  "data": [
    {
      "seller_name": "string",
      "seller_email": "string",
      "total_products": 8,
      "total_sales": 10500000,
      "total_profit": 2100000,
      "total_transactions": 35
    }
  ]
}

Note: Only COMPLETED transactions are counted. Sorted by total_sales DESC.
```

---

### �👥 User Management (Admin Only)

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

#### 3. Get User Detail

```
GET /users/:id
Authorization: Bearer <admin_token>

Response 200:
{
  "data": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "role": {
      "id": "uuid",
      "name": "string"
    },
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

#### 4. Update User

```
PUT /users/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

Body (all fields optional):
{
  "name": "string",
  "email": "string",
  "role_id": "uuid"
}

Response 200:
{
  "message": "User updated successfully",
  "data": { updated user object }
}
```

#### 5. Delete User

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
// Login (gunakan salah satu demo user)
const response = await fetch("http://localhost:8080/auth/login", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    email: "admin@example.com", // atau seller1@example.com, customer1@example.com, dll
    password: "password123",
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
| GET /profile                   | ✅    | ✅     | ✅        |
| PUT /profile                   | ✅    | ✅     | ✅        |
| PUT /profile/password          | ✅    | ✅     | ✅        |
| GET /products                  | ✅    | ✅     | ✅        |
| POST /products                 | ✅    | ❌     | ❌        |
| PUT /products/:id              | ✅    | ❌     | ❌        |
| DELETE /products/:id           | ✅    | ❌     | ❌        |
| GET /products/low-stock        | ✅    | ❌     | ❌        |
| GET /product-types             | ✅    | ✅     | ✅        |
| POST /product-types            | ✅    | ❌     | ❌        |
| PUT /product-types/:id         | ✅    | ❌     | ❌        |
| DELETE /product-types/:id      | ✅    | ❌     | ❌        |
| GET /marketplace               | ✅    | ✅     | ✅        |
| POST /seller/products          | ❌    | ✅     | ❌        |
| GET /seller/products           | ❌    | ✅     | ❌        |
| PUT /seller/products/:id       | ❌    | ✅     | ❌        |
| DELETE /seller/products/:id    | ❌    | ✅     | ❌        |
| GET /seller/transactions       | ❌    | ✅     | ❌        |
| POST /transactions             | ❌    | ❌     | ✅        |
| GET /transactions/:id          | ✅    | ✅     | ✅        |
| POST /transactions/:id/confirm | ❌    | ✅     | ❌        |
| POST /transactions/:id/cancel  | ❌    | ❌     | ✅        |
| GET /customer/transactions     | ❌    | ❌     | ✅        |
| GET /dashboard                 | ✅    | ✅     | ✅        |
| GET /reports/sales             | ✅    | ❌     | ❌        |
| GET /reports/top-products      | ✅    | ❌     | ❌        |
| GET /reports/top-sellers       | ✅    | ❌     | ❌        |
| GET /users                     | ✅    | ❌     | ❌        |
| GET /users/:id                 | ✅    | ❌     | ❌        |
| POST /users/admin              | ✅    | ❌     | ❌        |
| PUT /users/:id                 | ✅    | ❌     | ❌        |
| DELETE /users/:id              | ✅    | ❌     | ❌        |

---

## 🔒 Protected Endpoints (Require Authentication)

Semua endpoint di atas kecuali `/auth/*` dan `/swagger/*` memerlukan header:

```
Authorization: Bearer <jwt_token>
```

---

## 👤 Default Users

Setelah aplikasi pertama kali dijalankan, akan ada 8 demo users yang otomatis dibuat:

### Admin Users (2 users)

1. **Super Admin**
   - Email: `admin@example.com`
   - Password: `password123`
   - Role: Admin

2. **Admin Staff**
   - Email: `admin.staff@example.com`
   - Password: `password123`
   - Role: Admin

### Seller Users (3 users)

1. **Toko Elektronik Jaya**
   - Email: `seller1@example.com`
   - Password: `password123`
   - Role: Seller

2. **Fashion Store**
   - Email: `seller2@example.com`
   - Password: `password123`
   - Role: Seller

3. **Food Corner**
   - Email: `seller3@example.com`
   - Password: `password123`
   - Role: Seller

### Pelanggan Users (3 users)

1. **Budi Santoso**
   - Email: `customer1@example.com`
   - Password: `password123`
   - Role: Pelanggan

2. **Siti Rahayu**
   - Email: `customer2@example.com`
   - Password: `password123`
   - Role: Pelanggan

3. **Ahmad Hidayat**
   - Email: `customer3@example.com`
   - Password: `password123`
   - Role: Pelanggan

**⚠️ PENTING:** Semua demo users menggunakan password yang sama: `password123`. Segera ganti password di production!

## 🗃 Database Schema

### Tables

- **roles** - Role management (Admin, Seller, Pelanggan)
- **product_types** - Kategori produk (Elektronik, Pakaian, Makanan, Furniture, Olahraga)
- **users** - Data user dengan role (8 demo users di-seed otomatis)
- **products** - Master produk (gudang pusat, 24 produk sample di-seed otomatis)
- **seller_products** - Katalog marketplace seller dengan markup
- **transactions** - Transaksi pembelian

### Seeded Data

Aplikasi akan otomatis men-seed data berikut saat pertama kali dijalankan:

**Roles (3):**

- Admin
- Seller
- Pelanggan

**Product Types (5):**

- Elektronik
- Pakaian
- Makanan
- Furniture
- Olahraga

**Users (8):**

- 2 Admin users
- 3 Seller users
- 3 Pelanggan users
- Semua dengan password: `password123`

**Products (24):**

- 5 produk Elektronik (Laptop ASUS ROG, iPhone 15 Pro, Samsung Galaxy S24, Headphone Sony, Mouse Logitech)
- 5 produk Pakaian (Kemeja Batik, Celana Jeans, Jaket Kulit, Sepatu Nike, Tas Ransel)
- 5 produk Makanan (Kopi Arabica, Coklat Belgia, Madu Murni, Teh Hijau, Snack Mix)
- 4 produk Furniture (Kursi Gaming, Meja Kerja, Lemari Pakaian, Sofa)
- 5 produk Olahraga (Sepeda MTB, Raket Badminton, Bola Sepak, Matras Yoga, Dumbbell Set)

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
# Login sebagai Admin
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'

# Login sebagai Seller
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "seller1@example.com",
    "password": "password123"
  }'

# Login sebagai Pelanggan
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer1@example.com",
    "password": "password123"
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

## 🐛 Troubleshooting

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
