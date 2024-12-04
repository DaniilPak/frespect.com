package roommanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sakura/models"
	"sync"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

var (
	rooms      = make(map[string]*models.Room)
	roomsMutex sync.RWMutex
)

func FindParticipantByClientID(clientId string) (*models.Participant, *models.Room) {
	for _, room := range rooms {
		room.Mutex.RLock()
		participant, exists := room.Participants[clientId]
		room.Mutex.RUnlock()

		if exists {
			return participant, room // Return the participant and the room it's in
		}
	}
	return nil, nil // Participant not found in any room
}

func AddTrackToAllParticipants(track webrtc.TrackLocal) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	for _, room := range rooms {
		room.Mutex.Lock()
		for _, participant := range room.Participants {
			participant.Mutex.Lock()
			participant.PeerConnection.AddTrack(track)
			participant.Tracks[track.ID()] = track
			participant.Mutex.Unlock()
		}
		room.Mutex.Unlock()
	}

	PrintRooms()
}

func Wrtp(rtp *rtp.Packet) {
	for _, room := range rooms {
		for _, participant := range room.Participants {
			for _, track := range participant.Tracks {
				// Type assertion to ensure track is TrackLocalStaticRTP
				if staticTrack, ok := track.(*webrtc.TrackLocalStaticRTP); ok {
					err := staticTrack.WriteRTP(rtp)
					if err != nil {
						fmt.Printf("Failed to write RTP to track %s: %v\n", staticTrack.ID(), err)
					}
				} else {
					fmt.Printf("Track %s is not a TrackLocalStaticRTP\n", track.ID())
				}
			}
		}
	}
}

func RenegotAll(serverURL string) {
	for _, room := range rooms {
		for _, participant := range room.Participants {
			payload := map[string]string{
				"clientId": participant.ClientID,
			}

			// Convert payload to JSON format
			jsonData, err := json.Marshal(payload)
			if err != nil {
				panic(err)
			}

			// Create a new HTTP POST request
			req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("Failed to create POST request: %v\n", err)
				return
			}

			// Set the Content-Type and x-api-key headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("x-api-key", "your-secret-key")

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Failed to send POST request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Check response status
			if resp.StatusCode == http.StatusOK {
				fmt.Println("Successfully sent answer to server")
			} else {
				fmt.Printf("Failed to send answer with status code: %d\n", resp.StatusCode)
			}
		}
	}
}

func PrintRooms() {
	roomsMutex.RLock()
	defer roomsMutex.RUnlock()

	for roomID, room := range rooms {
		fmt.Printf("Room ID: %s\n", roomID)

		room.Mutex.RLock()
		for clientID, participant := range room.Participants {
			fmt.Printf("  Participant ID: %s\n", clientID)
			fmt.Printf("    PeerConnection: %v\n", participant.PeerConnection)

			participant.Mutex.RLock()
			for trackID, track := range participant.Tracks {
				fmt.Printf("      Track ID: %s\n", trackID)
				fmt.Printf("      Track Info: %v\n", track)
			}
			participant.Mutex.RUnlock()
		}
		room.Mutex.RUnlock()
	}
}

func GetOrCreateRoom(roomID string) *models.Room {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		room = &models.Room{
			ID:           roomID,
			Participants: make(map[string]*models.Participant),
		}
		rooms[roomID] = room
	}
	return room
}

func DeleteRoom(roomID string) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	delete(rooms, roomID)
}

// RemoveParticipantFromRoom removes a participant from a room
func RemoveParticipantFromRoom(roomID string, clientID string) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		return // Room doesn't exist
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	delete(room.Participants, clientID)

	// If the room is empty after removing the participant, delete the room
	if len(room.Participants) == 0 {
		delete(rooms, roomID)
	}
}
