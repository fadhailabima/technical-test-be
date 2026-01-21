package controllers

import (
	"net/http"
	"strconv"
	"technical-test-backend/services"

	"github.com/gin-gonic/gin"
)

var catService = services.CatalogService{}

// AddToEtalase godoc
// @Summary (Seller) Pajang Barang & Markup Harga
// @Description Seller memilih barang dari gudang admin dan menentukan harga jual sendiri
// @Tags Seller Catalog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body services.AddToEtalaseInput true "Data Markup"
// @Success 201 {object} map[string]interface{}
// @Router /seller/products [post]
func AddToEtalase(c *gin.Context) {
	var input services.AddToEtalaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// Gunakan http.StatusBadRequest alih-alih 400
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil UserID dari Middleware
	userID := c.GetString("userID")
	
	res, err := catService.AddToEtalase(userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Gunakan http.StatusCreated alih-alih 201
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

// GetMarketplace godoc
// @Summary (Pembeli) Lihat Marketplace
// @Description Melihat daftar barang yang dijual oleh Seller
// @Tags Marketplace
// @Security BearerAuth
// @Param search query string false "Search product name"
// @Param category query string false "Category/Product Type ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {object} map[string]interface{}
// @Router /marketplace [get]
func GetMarketplace(c *gin.Context) {
	// Get query parameters
	search := c.Query("search")
	categoryID := c.Query("category")
	minPrice := c.DefaultQuery("min_price", "0")
	maxPrice := c.DefaultQuery("max_price", "0")
	
	// Convert price strings to float64
	var minPriceFloat, maxPriceFloat float64
	if minPrice != "0" {
		if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
			minPriceFloat = val
		}
	}
	if maxPrice != "0" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			maxPriceFloat = val
		}
	}
	
	items, err := catService.GetMarketplaceItems(search, categoryID, minPriceFloat, maxPriceFloat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetSellerProducts godoc
// @Summary (Seller) Lihat Daftar Produk Sendiri
// @Description Melihat daftar produk yang dijual oleh seller yang sedang login
// @Tags Seller Catalog
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /seller/products [get]
func GetSellerProducts(c *gin.Context) {
	sellerID := c.GetString("userID")
	
	items, err := catService.GetSellerProducts(sellerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// UpdateSellerProduct godoc
// @Summary (Seller) Update Harga Produk di Marketplace
// @Description Seller dapat memperbarui harga jual atau status aktif produk mereka
// @Tags Seller Catalog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Seller Product ID (UUID)"
// @Param input body services.UpdateSellerProductInput true "Data Update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /seller/products/{id} [put]
func UpdateSellerProduct(c *gin.Context) {
	var input services.UpdateSellerProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	sellerID := c.GetString("userID")
	sellerProductID := c.Param("id")
	
	product, err := catService.UpdateSellerProduct(sellerProductID, sellerID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": product})
}

// DeleteSellerProduct godoc
// @Summary (Seller) Hapus Produk dari Marketplace
// @Description Seller menonaktifkan produk mereka dari marketplace (soft delete)
// @Tags Seller Catalog
// @Security BearerAuth
// @Param id path string true "Seller Product ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /seller/products/{id} [delete]
func DeleteSellerProduct(c *gin.Context) {
	sellerID := c.GetString("userID")
	sellerProductID := c.Param("id")
	
	if err := catService.DeleteSellerProduct(sellerProductID, sellerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Product removed from marketplace"})
}