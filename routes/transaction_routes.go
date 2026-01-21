package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(r *gin.Engine) {
	r.POST("/transactions", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Pelanggan"), 
		controllers.CreateOrder,
	)
	
	r.POST("/transactions/:id/confirm", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Seller"), 
		controllers.ConfirmOrder,
	)
	
	r.GET("/transactions/:id",
		middlewares.AuthMiddleware(),
		controllers.GetTransactionDetail,
	)
	
	r.POST("/transactions/:id/cancel",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Pelanggan"),
		controllers.CancelTransaction,
	)
}
