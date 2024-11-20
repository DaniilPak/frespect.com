package roommanager

import (
	"fmt"
	"sakura/models"
	"sync"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

var (
	rooms      = make(map[string]*models.Room)
	roomsMutex sync.RWMutex
)

func AddTrackToAllParticipants(track webrtc.TrackLocal) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	for _, room := range rooms {
		room.Mutex.Lock()
		for _, participant := range room.Participants {
			participant.Mutex.Lock()
			participant.PeerConnection.AddTrack(track)
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
