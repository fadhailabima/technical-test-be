package routes

import (
	"technical-test-backend/controllers"
	"technical-test-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupReportRoutes(r *gin.Engine) {
	r.GET("/reports/sales",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Admin"),
		controllers.GetSalesReport,
	)

	r.GET("/reports/top-products",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Admin"),
		controllers.GetTopProducts,
	)

	r.GET("/reports/top-sellers",
		middlewares.AuthMiddleware(),
		middlewares.RoleMiddleware("Admin"),
		controllers.GetTopSellers,
	)
}
