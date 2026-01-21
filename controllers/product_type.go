package controllers

import (
	"technical-test-backend/services"
	"github.com/gin-gonic/gin"
)

var typeService = services.ProductTypeService{}

// GetTypes godoc
// @Summary Lihat Daftar Kategori
// @Tags Product Type
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /product-types [get]
func GetTypes(c *gin.Context) {
	types, _ := typeService.GetAll()
	c.JSON(200, gin.H{"data": types})
}

type CreateTypeInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateType godoc
// @Summary Tambah Kategori (Admin)
// @Tags Product Type
// @Security BearerAuth
// @Param input body CreateTypeInput true "Nama Kategori"
// @Success 201 {object} map[string]interface{}
// @Router /product-types [post]
func CreateType(c *gin.Context) {
	var input CreateTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	res, _ := typeService.Create(input.Name)
	c.JSON(201, gin.H{"data": res})
}