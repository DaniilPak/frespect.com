package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sakura"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	appPort := os.Getenv("APP_PORT")

	if appPort == "" {
		log.Fatalf("Environment variables are not set")
	}

	sakura.RegisterRoutes()

	port := ":4000"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
