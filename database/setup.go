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

func ConnectDatabase() {
	// Environment Variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	// Buka Koneksi
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// AutoMigrate
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

	seedDatabase(database)

	DB = database
}

func seedDatabase(db *gorm.DB) {
	// --- SEEDING ROLES ---
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
	var countTypes int64
	db.Model(&models.ProductType{}).Count(&countTypes)
	if countTypes == 0 {
		types := []models.ProductType{
			{Name: "Elektronik"},
			{Name: "Pakaian"},
			{Name: "Makanan"},
		}
		db.Create(&types)
		fmt.Println("✅ Data Product Types Berhasil Dibuat!")
	}

	// --- SEEDING DEFAULT ADMIN USER ---
	var adminRole models.Role
	err := db.Where("name = ?", "Admin").First(&adminRole).Error
	
	if err == nil {
		var countUser int64
		db.Model(&models.User{}).Where("role_id = ?", adminRole.ID).Count(&countUser)

		if countUser == 0 {
			fmt.Println("⚠️ Admin belum ada. Membuat Super Admin default...")
			
			passwordHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

			admin := models.User{
				Name:     "Super Admin",
				Email:    "admin@example.com",
				Password: string(passwordHash),
				RoleID:   adminRole.ID,
			}

			if err := db.Create(&admin).Error; err != nil {
				log.Fatal("Gagal seeding admin:", err)
			}
			fmt.Println("✅ Super Admin Berhasil Dibuat! (Email: admin@example.com / Pass: admin123)")
		}
	}
}