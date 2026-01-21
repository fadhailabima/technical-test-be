package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter adalah main router yang mengimport semua route modules
func SetupRouter(r *gin.Engine) {
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Setup all routes by entity
	SetupAuthRoutes(r)
	SetupProductRoutes(r)
	SetupProductTypeRoutes(r)
	SetupMarketplaceRoutes(r)
	SetupSellerRoutes(r)
	SetupCustomerRoutes(r)
	SetupTransactionRoutes(r)
	SetupDashboardRoutes(r)
	SetupUserRoutes(r)
}