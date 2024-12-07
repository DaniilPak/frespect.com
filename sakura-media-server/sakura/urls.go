package sakura

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetServerURL(prefix string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	signalHost := os.Getenv("SIGNAL_SERVER_HOST")
	signalPort := os.Getenv("SIGNAL_SERVER_PORT")

	if signalHost == "" && signalPort == "" {
		log.Fatalf("Environment variables are not set")
	}

	return fmt.Sprintf("http://%s:%s/media-server/%s", signalHost, signalPort, prefix)
}

func GetMediaManagerURL() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	signalHost := os.Getenv("NODE_API_HOST")
	signalPort := os.Getenv("NODE_API_PORT")

	if signalHost == "" && signalPort == "" {
		log.Fatalf("Environment variables are not set")
	}

	fmt.Printf("Media manager server %s\n", fmt.Sprintf("http://%s:%s/", signalHost, signalPort))

	return fmt.Sprintf("http://%s:%s/", signalHost, signalPort)
}
