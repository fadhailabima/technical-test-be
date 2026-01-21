package controllers

import (
	"net/http" 
	"technical-test-backend/services"

	"github.com/gin-gonic/gin"
)

var trxService = services.TransactionService{}

// CreateOrder godoc
// @Summary (Pembeli) Buat Pesanan
// @Description Pembeli membuat order ke lapak seller (Status: PENDING)
// @Tags Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body services.CreateOrderInput true "Data Order"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /transactions [post]
func CreateOrder(c *gin.Context) {
	var input services.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	input.UserID = c.GetString("userID")
	
	trx, err := trxService.CreateOrder(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": trx})
}

// ConfirmOrder godoc
// @Summary (Seller) Konfirmasi Pesanan
// @Description Seller memproses order pending. Stok gudang admin akan berkurang di sini.
// @Tags Transaction
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /transactions/{id}/confirm [post]
func ConfirmOrder(c *gin.Context) {
	id := c.Param("id")
	sellerID := c.GetString("userID")

	if err := trxService.ConfirmOrder(id, sellerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Confirmed"})
}

// GetSellerTransactions godoc
// @Summary (Seller) List Semua Transaksi
// @Description Seller melihat semua transaksi dari produk mereka dengan detail buyer dan profit
// @Tags Transaction
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /seller/transactions [get]
func GetSellerTransactions(c *gin.Context) {
	sellerID := c.GetString("userID")

	transactions, err := trxService.GetSellerTransactions(sellerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
// GetCustomerTransactions godoc
// @Summary (Customer) List Transaksi Pembelian
// @Description Customer melihat semua transaksi pembelian mereka dengan detail seller dan produk
// @Tags Transaction
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /customer/transactions [get]
func GetCustomerTransactions(c *gin.Context) {
	customerID := c.GetString("userID")

	transactions, err := trxService.GetCustomerTransactions(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
