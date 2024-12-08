package sakura

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func BotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %s", err)
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}

	// Parse the JSON request body
	var requestData BotRequest
	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Printf("Error parsing JSON: %s", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Managing room
	songId := requestData.SongId
	fmt.Printf("Bot request: %s\n", songId)

	go StartMusicEverywhere()
}
