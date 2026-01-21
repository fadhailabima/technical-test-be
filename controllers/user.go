package controllers

import (
	"technical-test-backend/services"
	"github.com/gin-gonic/gin"
)

var userService = services.UserService{}

// CreateAdmin godoc
// @Summary Tambah Admin Baru (Super Admin Only)
// @Description Menambahkan user dengan role Admin
// @Tags User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body services.CreateAdminInput true "Data Admin"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users/admin [post]
func CreateAdmin(c *gin.Context) {
	var input services.CreateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	if err := userService.CreateAdmin(input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	c.JSON(201, gin.H{"message": "Admin created"})
}

// DeleteUser godoc
// @Summary Hapus User
// @Description Menghapus user berdasarkan ID (Hard Delete)
// @Tags User Management
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	if err := userService.DeleteUser(c.Param("id")); err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete"}); return
	}
	c.JSON(200, gin.H{"message": "User deleted"})
}

// FindUsers godoc
// @Summary Lihat Daftar Semua User
// @Description Mengambil list semua user beserta rolenya
// @Tags User Management
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /users [get]
func FindUsers(c *gin.Context) {
	users, _ := userService.GetAllUsers()
	c.JSON(200, gin.H{"data": users})
}