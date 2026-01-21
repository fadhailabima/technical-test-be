package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupDashboardRoutes(r *gin.Engine) {
	// Dashboard endpoint - accessible by all authenticated users
	// Response berbeda tergantung role (Admin/Seller/Pelanggan)
	r.GET("/dashboard", 
		middlewares.AuthMiddleware(), 
		controllers.GetDashboard,
	)
}
