package controllers

import (
	"technical-test-backend/database"
	"technical-test-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetProfile godoc
// @Summary Get User Profile
// @Tags Profile
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /profile [get]
func GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	
	var user models.User
	if err := database.DB.Preload("Role").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(200, gin.H{
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role.Name,
		},
	})
}

type UpdateProfileInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateProfile godoc
// @Summary Update User Profile
// @Tags Profile
// @Security BearerAuth
// @Param input body UpdateProfileInput true "Profile Data"
// @Success 200 {object} map[string]interface{}
// @Router /profile [put]
func UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	
	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	
	// Update hanya field yang diisi
	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	
	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to update profile"})
		return
	}
	
	c.JSON(200, gin.H{
		"message": "Profile updated successfully",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword godoc
// @Summary Change User Password
// @Tags Profile
// @Security BearerAuth
// @Param input body ChangePasswordInput true "Password Data"
// @Success 200 {object} map[string]interface{}
// @Router /profile/password [put]
func ChangePassword(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	
	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	
	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(401, gin.H{"error": "Old password is incorrect"})
		return
	}
	
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}
	
	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update password"})
		return
	}
	
	c.JSON(200, gin.H{"message": "Password changed successfully"})
}
