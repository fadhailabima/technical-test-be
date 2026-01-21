package database

import (
	"fmt"
	"log"
	"os"
	"technical-test-backend/models"

	"golang.org/x/crypto/bcrypt" 
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase - Fungsi utama untuk koneksi database dan inisialisasi
// Alur: Load env -> Koneksi PostgreSQL -> Auto migrate -> Seeding data
func ConnectDatabase() {
	// 1. Load konfigurasi database dari environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	// 2. Build connection string (DSN)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	// 3. Open connection ke PostgreSQL dengan GORM
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// 4. Auto migrate semua model (create tables jika belum ada)
	// Urutan penting: Role -> ProductType -> User -> Product -> SellerProduct -> Transaction
	err = database.AutoMigrate(
		&models.Role{}, 
		&models.ProductType{}, 
		&models.User{}, 
		&models.Product{},
		&models.SellerProduct{},
    	&models.Transaction{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}
	fmt.Println("✅ Migrasi Database Berhasil!")

	// 5. Seeding data awal untuk development/testing
	seedDatabase(database)

	// 6. Assign database connection ke global variable
	DB = database
}

// seedDatabase - Fungsi untuk mengisi data awal (roles, users, products)
// Hanya jalan jika table masih kosong (idempotent)
func seedDatabase(db *gorm.DB) {
	// --- SEEDING ROLES ---
	// Buat 3 role utama: Admin, Seller, Pelanggan
	var countRoles int64
	db.Model(&models.Role{}).Count(&countRoles)
	if countRoles == 0 {
		roles := []models.Role{
			{Name: "Admin"}, 
			{Name: "Seller"}, 
			{Name: "Pelanggan"},
		}
		db.Create(&roles)
		fmt.Println("✅ Data Roles Berhasil Dibuat!")
	}

	// --- SEEDING PRODUCT TYPES ---
	// Buat 5 kategori produk untuk marketplace
	var countTypes int64
	db.Model(&models.ProductType{}).Count(&countTypes)
	if countTypes == 0 {
		types := []models.ProductType{
			{Name: "Elektronik"},
			{Name: "Pakaian"},
			{Name: "Makanan"},
			{Name: "Furniture"},
			{Name: "Olahraga"},
		}
		db.Create(&types)
		fmt.Println("✅ Data Product Types Berhasil Dibuat!")
	}

	// --- GET ROLES ---
	// Ambil role ID untuk digunakan saat seeding users
	var adminRole, sellerRole, pelangganRole models.Role
	db.Where("name = ?", "Admin").First(&adminRole)
	db.Where("name = ?", "Seller").First(&sellerRole)
	db.Where("name = ?", "Pelanggan").First(&pelangganRole)

	// --- SEEDING USERS FOR EACH ROLE ---
	// Buat demo users: 2 Admin, 3 Seller, 3 Pelanggan
	// Password semua user: password123
	var countUsers int64
	db.Model(&models.User{}).Count(&countUsers)
	
	if countUsers == 0 {
		fmt.Println("⚠️ Membuat demo users untuk setiap role...")
		
		// Hash password yang sama untuk semua demo user (bcrypt)
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		
		users := []models.User{
			// Admin Users
			{
				Name:     "Super Admin",
				Email:    "admin@example.com",
				Password: string(passwordHash),
				RoleID:   adminRole.ID,
			},
			{
				Name:     "Admin Staff",
				Email:    "admin.staff@example.com",
				Password: string(passwordHash),
				RoleID:   adminRole.ID,
			},
			// Seller Users
			{
				Name:     "Toko Elektronik Jaya",
				Email:    "seller1@example.com",
				Password: string(passwordHash),
				RoleID:   sellerRole.ID,
			},
			{
				Name:     "Fashion Store",
				Email:    "seller2@example.com",
				Password: string(passwordHash),
				RoleID:   sellerRole.ID,
			},
			{
				Name:     "Food Corner",
				Email:    "seller3@example.com",
				Password: string(passwordHash),
				RoleID:   sellerRole.ID,
			},
			// Pelanggan Users
			{
				Name:     "Budi Santoso",
				Email:    "customer1@example.com",
				Password: string(passwordHash),
				RoleID:   pelangganRole.ID,
			},
			{
				Name:     "Siti Rahayu",
				Email:    "customer2@example.com",
				Password: string(passwordHash),
				RoleID:   pelangganRole.ID,
			},
			{
				Name:     "Ahmad Hidayat",
				Email:    "customer3@example.com",
				Password: string(passwordHash),
				RoleID:   pelangganRole.ID,
			},
		}

		if err := db.Create(&users).Error; err != nil {
			log.Fatal("Gagal seeding users:", err)
		}
		fmt.Println("✅ Demo Users Berhasil Dibuat! (Password semua user: password123)")
		fmt.Println("   - 2 Admin: admin@example.com, admin.staff@example.com")
		fmt.Println("   - 3 Seller: seller1-3@example.com")
		fmt.Println("   - 3 Pelanggan: customer1-3@example.com")
	}

	// --- SEEDING PRODUCTS ---
	// Buat 24 sample products untuk testing
	// 5 Elektronik, 5 Pakaian, 5 Makanan, 4 Furniture, 5 Olahraga
	var countProducts int64
	db.Model(&models.Product{}).Count(&countProducts)
	
	if countProducts == 0 {
		fmt.Println("⚠️ Membuat sample products...")
		
		// Ambil ID product types yang sudah di-seed sebelumnya
		var elektronikType, pakaianType, makananType, furnitureType, olahragaType models.ProductType
		db.Where("name = ?", "Elektronik").First(&elektronikType)
		db.Where("name = ?", "Pakaian").First(&pakaianType)
		db.Where("name = ?", "Makanan").First(&makananType)
		db.Where("name = ?", "Furniture").First(&furnitureType)
		db.Where("name = ?", "Olahraga").First(&olahragaType)
		
		products := []models.Product{
			// Elektronik
			{Name: "Laptop ASUS ROG", ProductTypeID: elektronikType.ID, Price: 15000000, Stock: 10},
			{Name: "iPhone 15 Pro", ProductTypeID: elektronikType.ID, Price: 18000000, Stock: 15},
			{Name: "Samsung Galaxy S24", ProductTypeID: elektronikType.ID, Price: 12000000, Stock: 20},
			{Name: "Headphone Sony WH-1000XM5", ProductTypeID: elektronikType.ID, Price: 4500000, Stock: 30},
			{Name: "Mouse Logitech MX Master 3", ProductTypeID: elektronikType.ID, Price: 1200000, Stock: 50},
			
			// Pakaian
			{Name: "Kemeja Batik Premium", ProductTypeID: pakaianType.ID, Price: 350000, Stock: 40},
			{Name: "Celana Jeans Levi's", ProductTypeID: pakaianType.ID, Price: 800000, Stock: 35},
			{Name: "Jaket Kulit", ProductTypeID: pakaianType.ID, Price: 1500000, Stock: 15},
			{Name: "Sepatu Nike Air Max", ProductTypeID: pakaianType.ID, Price: 2000000, Stock: 25},
			{Name: "Tas Ransel Premium", ProductTypeID: pakaianType.ID, Price: 450000, Stock: 30},
			
			// Makanan
			{Name: "Kopi Arabica 1kg", ProductTypeID: makananType.ID, Price: 150000, Stock: 100},
			{Name: "Coklat Belgia Premium", ProductTypeID: makananType.ID, Price: 250000, Stock: 60},
			{Name: "Madu Murni 500ml", ProductTypeID: makananType.ID, Price: 120000, Stock: 80},
			{Name: "Teh Hijau Organik", ProductTypeID: makananType.ID, Price: 80000, Stock: 90},
			{Name: "Snack Mix Premium", ProductTypeID: makananType.ID, Price: 50000, Stock: 150},
			
			// Furniture
			{Name: "Kursi Gaming", ProductTypeID: furnitureType.ID, Price: 3500000, Stock: 12},
			{Name: "Meja Kerja Minimalis", ProductTypeID: furnitureType.ID, Price: 2500000, Stock: 8},
			{Name: "Lemari Pakaian", ProductTypeID: furnitureType.ID, Price: 4000000, Stock: 6},
			{Name: "Sofa 3 Seater", ProductTypeID: furnitureType.ID, Price: 6000000, Stock: 5},
			
			// Olahraga
			{Name: "Sepeda Gunung MTB", ProductTypeID: olahragaType.ID, Price: 5000000, Stock: 10},
			{Name: "Raket Badminton Yonex", ProductTypeID: olahragaType.ID, Price: 800000, Stock: 25},
			{Name: "Bola Sepak Adidas", ProductTypeID: olahragaType.ID, Price: 300000, Stock: 40},
			{Name: "Matras Yoga", ProductTypeID: olahragaType.ID, Price: 250000, Stock: 50},
			{Name: "Dumbbell Set 20kg", ProductTypeID: olahragaType.ID, Price: 1200000, Stock: 20},
		}

		if err := db.Create(&products).Error; err != nil {
			log.Fatal("Gagal seeding products:", err)
		}
		fmt.Println("✅ Sample Products Berhasil Dibuat! (24 produk)")
	}
}