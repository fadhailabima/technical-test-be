package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupProductTypeRoutes(r *gin.Engine) {
	r.GET("/product-types", 
		middlewares.AuthMiddleware(), 
		controllers.GetTypes,
	)
	
	r.POST("/product-types", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Admin"), 
		controllers.CreateType,
	)
}
