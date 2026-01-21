package controllers

import (
	"strconv"
	"technical-test-backend/services"

	"github.com/gin-gonic/gin"
)

var reportService = services.ReportService{}

// GetSalesReport godoc
// @Summary Sales Report (Admin)
// @Tags Reports
// @Security BearerAuth
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /reports/sales [get]
func GetSalesReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	report, err := reportService.GetSalesReport(startDate, endDate)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": report})
}

// GetTopProducts godoc
// @Summary Top Products Report (Admin)
// @Tags Reports
// @Security BearerAuth
// @Param limit query int false "Limit (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /reports/top-products [get]
func GetTopProducts(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	report, err := reportService.GetTopProducts(limit)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"limit": limit,
		"data":  report,
	})
}

// GetTopSellers godoc
// @Summary Top Sellers Report (Admin)
// @Tags Reports
// @Security BearerAuth
// @Param limit query int false "Limit (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /reports/top-sellers [get]
func GetTopSellers(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	report, err := reportService.GetTopSellers(limit)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"limit": limit,
		"data":  report,
	})
}
