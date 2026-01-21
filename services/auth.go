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

type AuthService struct{}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"role_id" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *AuthService) RegisterUser(input RegisterInput) error {
	roleUUID, err := uuid.Parse(input.RoleID)
	if err != nil {
		return errors.New("format role_id salah (harus UUID)")
	}

	var role models.Role
	if err := database.DB.First(&role, "id = ?", roleUUID).Error; err != nil {
		return errors.New("role tidak ditemukan")
	}

	if role.Name == "Admin" {
		return errors.New("registrasi admin dilarang melalui jalur publik")
	}


	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   roleUUID, 
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return errors.New("gagal register")
	}

	return nil
}

func (s *AuthService) LoginUser(input LoginInput) (LoginResponse, error) {
	var user models.User
	if err := database.DB.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
		return LoginResponse{}, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return LoginResponse{}, errors.New("email atau password salah")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role.Name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

func (s *AuthService) GetAvailableRoles() ([]RoleResponse, error) {
	var roles []models.Role
	
	// Get only Seller and Pelanggan roles (exclude Admin for public registration)
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