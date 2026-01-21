package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(r *gin.Engine) {
	r.GET("/customer/transactions",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Pelanggan"),
		controllers.GetCustomerTransactions,
	)
}
