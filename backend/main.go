package main

import (
	"hero-hunter/config"
	"hero-hunter/handlers"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize configuration (Redis, MongoDB)
	config.InitRedis()
	config.InitMongoDB()

	// Initialize Gin Router
	router := gin.Default()

	// Enable CORS with specific origin (for security, don't allow all origins in production)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://hero-huntr.vercel.app/", "https://hero-huntr-9jlg0lpab-demis-projects-34d8549b.vercel.app/"}, // Allow the React app's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Set up routes
	router.GET("/api/search", handlers.SearchHandler)

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if the environment variable isn't set
	}

	// Run the server
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
