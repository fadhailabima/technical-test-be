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

// GET USER BY ID
func (s *UserService) GetUserByID(userID string) (models.User, error) {
	var user models.User
	if err := database.DB.Preload("Role").Omit("password").Where("id = ?", userID).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// UPDATE USER (Admin can update user role/status)
type UpdateUserInput struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	RoleID *string `json:"role_id"`
}

func (s *UserService) UpdateUser(userID string, input UpdateUserInput) (models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return user, err
	}

	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.RoleID != nil {
		updates["role_id"] = *input.RoleID
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		return user, err
	}

	// Reload with Role
	database.DB.Preload("Role").Omit("password").First(&user, "id = ?", userID)
	return user, nil
}
