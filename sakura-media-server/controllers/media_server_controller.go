package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sakura/models"
	"sakura/roommanager"
	"sakura/sfu"

	"github.com/pion/webrtc/v4"
)

func MediaServerHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
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
	var requestData models.SDPRequest
	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Printf("Error parsing JSON: %s", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Managing room
	roomID := requestData.RoomID
	room := roommanager.GetOrCreateRoom(roomID)

	// Create a new participant
	participant := &models.Participant{
		ClientID: requestData.ClientID,
		Tracks:   make(map[string]*webrtc.TrackLocalStaticRTP),
	}

	// Add participant to the room
	room.Mutex.Lock()
	room.Participants[participant.ClientID] = participant
	room.Mutex.Unlock()

	// Run RunReflectServer in a goroutine
	go sfu.RunReflectServer(requestData, room, participant)

	// Create a response indicating successful processing
	response := models.Response{Message: "SDP received and processed successfully"}

	// Encode the response object as JSON and send it
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %s", err)
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}
