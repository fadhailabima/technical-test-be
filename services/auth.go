package services

import (
	"errors"
	"os"
	"technical-test-backend/database"
	"technical-test-backend/models"
	"time"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService menangani semua logika bisnis terkait autentikasi dan otorisasi
type AuthService struct{}

// RegisterInput adalah struktur data untuk input registrasi user baru
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"role_id" binding:"required"` // UUID dari role (Seller/Pelanggan)
}

// LoginInput adalah struktur data untuk input login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse adalah struktur response setelah login berhasil
type LoginResponse struct {
	Token string      `json:"token"` // JWT token untuk autentikasi
	User  models.User `json:"user"`  // Data user yang login
}

// RoleResponse adalah struktur untuk menampilkan role yang tersedia untuk registrasi
type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RegisterUser - Proses registrasi user baru
// Alur: Validasi role -> Hash password -> Simpan ke database
func (s *AuthService) RegisterUser(input RegisterInput) error {
	// 1. Parse dan validasi UUID dari role_id
	roleUUID, err := uuid.Parse(input.RoleID)
	if err != nil {
		return errors.New("format role_id salah (harus UUID)")
	}

	// 2. Cek apakah role tersebut exist di database
	var role models.Role
	if err := database.DB.First(&role, "id = ?", roleUUID).Error; err != nil {
		return errors.New("role tidak ditemukan")
	}

	// 3. Proteksi: Tidak boleh register sebagai Admin melalui jalur publik
	if role.Name == "Admin" {
		return errors.New("registrasi admin dilarang melalui jalur publik")
	}

	// 4. Hash password menggunakan bcrypt untuk keamanan
	// 4. Hash password menggunakan bcrypt untuk keamanan
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	
	// 5. Buat instance user baru
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   roleUUID, 
	}

	// 6. Simpan ke database
	if err := database.DB.Create(&user).Error; err != nil {
		return errors.New("gagal register")
	}

	return nil
}

// LoginUser - Proses autentikasi user dan generate JWT token
// Alur: Cek email -> Validasi password -> Generate JWT token
func (s *AuthService) LoginUser(input LoginInput) (LoginResponse, error) {
	// 1. Cari user berdasarkan email dan load data role-nya
	var user models.User
	if err := database.DB.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
		return LoginResponse{}, errors.New("email atau password salah")
	}

	// 2. Verifikasi password dengan bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return LoginResponse{}, errors.New("email atau password salah")
	}

	// 3. Generate JWT token dengan claims: user_id, role, expiry (24 jam)
	// 3. Generate JWT token dengan claims: user_id, role, expiry (24 jam)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,              // Subject: User ID
		"role": user.Role.Name,       // Role untuk authorization
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expired dalam 24 jam
	})

	// 4. Sign token dengan JWT_SECRET dari environment variable
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return LoginResponse{}, err
	}

	// 5. Return token dan data user
	return LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

// GetAvailableRoles - Mendapatkan list role yang tersedia untuk registrasi publik
// Hanya mengembalikan role Seller dan Pelanggan (Admin tidak bisa register publik)
func (s *AuthService) GetAvailableRoles() ([]RoleResponse, error) {
	var roles []models.Role
	
	// Ambil hanya role Seller dan Pelanggan (exclude Admin)
	if err := database.DB.Where("name IN ?", []string{"Seller", "Pelanggan"}).Find(&roles).Error; err != nil {
		return nil, err
	}

	var roleResponses []RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, RoleResponse{
			ID:   role.ID.String(),
			Name: role.Name,
		})
	}

	return roleResponses, nil
}