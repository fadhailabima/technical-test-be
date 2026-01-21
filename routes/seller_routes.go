package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupSellerRoutes(r *gin.Engine) {
	r.POST("/seller/products", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Seller"), 
		controllers.AddToEtalase,
	)
	
	r.GET("/seller/products",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Seller"),
		controllers.GetSellerProducts,
	)
	
	r.PUT("/seller/products/:id",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Seller"),
		controllers.UpdateSellerProduct,
	)
	
	r.DELETE("/seller/products/:id",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Seller"),
		controllers.DeleteSellerProduct,
	)
	
	r.GET("/seller/transactions",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Seller"),
		controllers.GetSellerTransactions,
	)
}
