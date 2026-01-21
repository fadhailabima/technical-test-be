package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupProfileRoutes(r *gin.Engine) {
	r.GET("/profile",
		middlewares.AuthMiddleware(),
		controllers.GetProfile,
	)
	
	r.PUT("/profile",
		middlewares.AuthMiddleware(),
		controllers.UpdateProfile,
	)
	
	r.PUT("/profile/password",
		middlewares.AuthMiddleware(),
		controllers.ChangePassword,
	)
}
