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

type UpdateTypeInput struct {
	Name string `json:"name" binding:"required"`
}

// UpdateType godoc
// @Summary Update Kategori (Admin)
// @Tags Product Type
// @Security BearerAuth
// @Param id path string true "Product Type ID"
// @Param input body UpdateTypeInput true "Nama Kategori Baru"
// @Success 200 {object} map[string]interface{}
// @Router /product-types/{id} [put]
func UpdateType(c *gin.Context) {
	id := c.Param("id")
	var input UpdateTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	res, err := typeService.Update(id, input.Name)
	if err != nil {
		c.JSON(404, gin.H{"error": "Product type not found"}); return
	}
	c.JSON(200, gin.H{"data": res})
}

// DeleteType godoc
// @Summary Hapus Kategori (Admin)
// @Tags Product Type
// @Security BearerAuth
// @Param id path string true "Product Type ID"
// @Success 200 {object} map[string]interface{}
// @Router /product-types/{id} [delete]
func DeleteType(c *gin.Context) {
	id := c.Param("id")
	if err := typeService.Delete(id); err != nil {
		c.JSON(404, gin.H{"error": "Product type not found"}); return
	}
	c.JSON(200, gin.H{"message": "Product type deleted"})
}