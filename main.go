package main

import (
	"log"
	"os"
	"technical-test-backend/database"
	"technical-test-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	_ "technical-test-backend/docs"
)

// @title           Inventory Management API
// @version         1.0
// @description     API Technical Test Junior Fullstack (Go + Postgres + UUID)
// @termsOfService  http://swagger.io/terms/

// @contact.name    Fadhail Athaillah Bima Dharmawan
// @contact.email   bimadharmawan6@gmail.com

// @host            localhost:8080
// @BasePath        /
func main() {
	// Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect & Migrate Database
	database.ConnectDatabase()

	// Setup Gin Engine
	r := gin.Default()
	
	// Disable automatic redirect for trailing slash
	r.RedirectTrailingSlash = false

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Panggil Routing dari folder routes/api.go
	routes.SetupRouter(r)

	// alankan Server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}