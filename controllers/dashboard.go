package controllers

import (
	"net/http"
	"technical-test-backend/services"
	"github.com/gin-gonic/gin"
)

var dashService = services.DashboardService{}

// GetDashboard godoc
// @Summary Dashboard Statistik (Multi-Role)
// @Description Menampilkan statistik keuangan dan transaksi berdasarkan Role User yang login
// @Tags Dashboard
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /dashboard [get]
func GetDashboard(c *gin.Context) {
	role := c.GetString("role")
	uid := c.GetString("userID")

	switch role {
	case "Pelanggan":
		c.JSON(200, gin.H{"role": role, "data": dashService.GetBuyerStats(uid)})
	case "Seller":
		c.JSON(200, gin.H{"role": role, "data": dashService.GetSellerStats(uid)})
	case "Admin":
		c.JSON(200, gin.H{"role": role, "data": dashService.GetAdminStats()})
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Unknown role"})
	}
}