package main

import (
	"log"

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
		log.Println("⚠️ No .env file found")
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			"X-ROLE",
		},
	}))


	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	config.ConnectDatabase()
	routes.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
