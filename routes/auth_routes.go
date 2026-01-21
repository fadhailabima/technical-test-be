package routes

import (
	"technical-test-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	r.POST("/auth/register", controllers.Register)
	r.POST("/auth/login", controllers.Login)
	r.GET("/auth/roles", controllers.GetRoles)
}
