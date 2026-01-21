package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupMarketplaceRoutes(r *gin.Engine) {
	r.GET("/marketplace", 
		middlewares.AuthMiddleware(), 
		controllers.GetMarketplace,
	)
}
