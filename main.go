package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"bimbel-absensi/config"
	"bimbel-absensi/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	config.ConnectDatabase()

	// Setup router
	r := routes.SetupRouter()

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server
	r.Run(":" + port)
}