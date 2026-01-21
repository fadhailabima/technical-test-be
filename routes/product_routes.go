package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.Engine) {
	// Product Master (Gudang Pusat)
	r.GET("/products", 
		middlewares.AuthMiddleware(), 
		controllers.FindAllProducts,
	)
	
	r.POST("/products", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Admin"), 
		controllers.CreateProduct,
	)
	
	r.DELETE("/products/:id", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Admin"), 
		controllers.DeleteProduct,
	)
	
	r.PUT("/products/:id",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Admin"),
		controllers.UpdateProduct,
	)
}
