package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	// User Management endpoints (Admin only)
	r.GET("/users", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Admin"), 
		controllers.FindUsers,
	)
	
	r.POST("/users/admin", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Admin"), 
		controllers.CreateAdmin,
	)
	
	r.DELETE("/users/:id", 
		middlewares.AuthMiddleware(), 
		middlewares.RoleMiddleware("Admin"), 
		controllers.DeleteUser,
	)
	
	r.GET("/users/:id",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Admin"),
		controllers.GetUserDetail,
	)
	
	r.PUT("/users/:id",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Admin"),
		controllers.UpdateUser,
	)
}
