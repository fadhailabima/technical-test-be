package controllers

import (
	"net/http"
	"technical-test-backend/services"

	"github.com/gin-gonic/gin"
)

var authService = services.AuthService{}

// @Summary Register User Baru
// @Description Mendaftarkan pengguna baru (Role ID dikirim sebagai String UUID)
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body services.RegisterInput true "Input Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input services.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := authService.RegisterUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil!"})
}

// @Summary Login User
// @Description Masuk menggunakan Email & Password untuk dapat Token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body services.LoginInput true "Input Data"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input services.LoginInput

	// Validasi Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Panggil Service
	result, err := authService.LoginUser(input)
	if err != nil {
		// Jika error (misal password salah), return 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Response Sukses (Token + Data User)
	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user": gin.H{
			"id":   result.User.ID,        // UUID otomatis terkonversi jadi string di JSON
			"name": result.User.Name,
			"role": result.User.Role.Name, // Mengambil nama role dari relasi
		},
	})
}

// @Summary Get Available Roles
// @Description Mendapatkan list role yang tersedia untuk registrasi (Seller dan Pelanggan saja, Admin tidak termasuk)
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/roles [get]
func GetRoles(c *gin.Context) {
	roles, err := authService.GetAvailableRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}