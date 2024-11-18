package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sakura/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Access environment variables
	appPort := os.Getenv("APP_PORT")

	if appPort == "" {
		log.Fatalf("Environment variables are not set")
	}

	// Register the routes using the routes package
	routes.RegisterRoutes()

	// Define the port and start the server on localhost:3000
	port := ":4000"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
