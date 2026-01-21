package services

import (
	"errors"
	"technical-test-backend/database"
	"technical-test-backend/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

// Input khusus buat Admin nambah Admin baru
type CreateAdminInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CREATE ADMIN (Hanya bisa dilakukan oleh Admin lain)
func (s *UserService) CreateAdmin(input CreateAdminInput) error {
	// Cari Role ID untuk "Admin"
	var role models.Role
	if err := database.DB.Where("name = ?", "Admin").First(&role).Error; err != nil {
		return errors.New("role admin tidak ditemukan di database")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	admin := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   role.ID, // Pakai UUID Role Admin
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		return errors.New("gagal membuat admin (email mungkin duplikat)")
	}
	return nil
}

// DELETE USER (Fitur Admin)
func (s *UserService) DeleteUser(userID string) error {
	if err := database.DB.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

// LIST USERS
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	// Tampilkan semua user kecuali password
	if err := database.DB.Preload("Role").Omit("password").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}