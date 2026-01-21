package controllers

import (
	"strconv"
	"technical-test-backend/models"
	"technical-test-backend/services"
	"github.com/gin-gonic/gin"
)

var prodService = services.ProductService{}

// CreateProduct godoc
// @Summary Input Barang ke Gudang (Admin)
// @Description Admin memasukkan master data produk
// @Tags Product Master (Gudang)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body services.CreateProductInput true "Data Produk"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var input services.CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	res, err := prodService.Create(input)
	if err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	c.JSON(201, gin.H{"data": res})
}

// FindAllProducts godoc
// @Summary Lihat Daftar Barang Gudang
// @Description Melihat semua master produk (Admin & Seller bisa lihat)
// @Tags Product Master (Gudang)
// @Security BearerAuth
// @Param search query string false "Search product name"
// @Param product_type_id query string false "Product Type ID filter"
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func FindAllProducts(c *gin.Context) {
	search := c.Query("search")
	productTypeID := c.Query("product_type_id")
	
	var products []models.Product
	var err error
	
	if search != "" || productTypeID != "" {
		products, err = prodService.FindAllWithFilters(search, productTypeID)
	} else {
		products, err = prodService.FindAll()
	}
	
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"data": products})
}

// DeleteProduct godoc
// @Summary Hapus Barang Gudang (Admin)
// @Tags Product Master (Gudang)
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	if err := prodService.Delete(c.Param("id")); err != nil {
		c.JSON(400, gin.H{"error": "Failed delete"}); return
	}
	c.JSON(200, gin.H{"message": "Deleted"})
}
// UpdateProduct godoc
// @Summary Update Barang Gudang (Admin)
// @Description Admin dapat memperbarui data produk master (nama, stock, harga, kategori)
// @Tags Product Master (Gudang)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param input body services.UpdateProductInput true "Data Update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	var input services.UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	product, err := prodService.Update(c.Param("id"), input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"data": product})
}

// GetLowStock godoc
// @Summary Get Low Stock Products (Admin)
// @Description Get products with stock below threshold
// @Tags Product Master (Gudang)
// @Security BearerAuth
// @Param threshold query int false "Stock threshold (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /products/low-stock [get]
func GetLowStock(c *gin.Context) {
	threshold := 10 // default
	if t := c.Query("threshold"); t != "" {
		if val, err := strconv.Atoi(t); err == nil {
			threshold = val
		}
	}
	
	products, err := prodService.GetLowStock(threshold)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"threshold": threshold,
		"count":     len(products),
		"data":      products,
	})
}

