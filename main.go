package main

import (
	"log"
	"fmt"

	"spacebook/config"
	"spacebook/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDatabase()
	fmt.Println("Connection is successful")

	// Initialize Echo app
	e := echo.New()

	// Adding CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Accept", "Origin"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(e)

	// Start the server
	err = e.Start(":8000")
	if err != nil {
		panic("could not start server")
	}
}
