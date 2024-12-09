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
	clientID := requestData.ClientID
	fmt.Printf("Bot request: %s %s\n", songId, clientID)

	// Create a music track (audio) to add to the peer connection
	roomManager := GetRoomManager()
	_, room := roomManager.FindParticipantByClientID(clientID)

	bot, err := NewBot()
	if err != nil {
		panic("Bot failed to create")
	}
	bot.SetRoom(room)
	musicTrack := bot.CreateAudioTrack()

	err2 := roomManager.AddTrackToParticipant(clientID, bot.botID, musicTrack)
	if err2 != nil {
		panic(fmt.Sprintf("AddTrack: Failed adding tracks: %v", err))
	}

	go bot.WriteAudioToTrack(audioURL, musicTrack)
	roomManager.RenegotiateAllClients()
}
