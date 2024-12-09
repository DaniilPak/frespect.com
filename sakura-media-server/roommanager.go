package sakura

import (
	"fmt"
	"sync"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type Room struct {
	ID           string
	Participants map[string]*Participant
	Mutex        sync.RWMutex
}

type Participant struct {
	ClientID       string
	PeerConnection *webrtc.PeerConnection
	Tracks         map[string]webrtc.TrackLocal
	Mutex          sync.RWMutex
}

type RoomManager struct {
	rooms            map[string]*Room
	clientIDToRoomID map[string]string
	roomsMutex       sync.RWMutex
	renegotiationSvc *RenegotiationService
}

// Singleton instance and sync.Once
var (
	roomManagerInstance *RoomManager
	once                sync.Once
)

// GetRoomManager returns the singleton instance of RoomManager.
func GetRoomManager() *RoomManager {
	once.Do(func() {
		renegotiationService := &RenegotiationService{}

		roomManagerInstance = &RoomManager{
			rooms:            make(map[string]*Room),
			renegotiationSvc: renegotiationService,
		}
	})
	return roomManagerInstance
}

// GetRoomByParticipantID finds the room where the participant with the given clientID is located.
func (rm *RoomManager) GetRoomByParticipantID(clientID string) *Room {
	rm.roomsMutex.RLock()
	defer rm.roomsMutex.RUnlock()

	// Check if the clientID exists in the mapping
	roomID, exists := rm.clientIDToRoomID[clientID]
	if !exists {
		return nil
	}

	// Retrieve the room using roomID
	room, roomExists := rm.rooms[roomID]
	if !roomExists {
		return nil
	}
	return room
}

// AddParticipantToRoom adds a participant to the room and updates the mapping.
func (rm *RoomManager) AddParticipantToRoom(roomID, clientID string, participant *Participant) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	// Get or create the room
	room := rm.GetOrCreateRoom(roomID)

	// Add the participant to the room
	room.Mutex.Lock()
	room.Participants[clientID] = participant
	room.Mutex.Unlock()

	// Update the mapping from clientID to roomID
	rm.clientIDToRoomID[clientID] = roomID
}

// FindParticipantByClientID locates a participant and its room by client ID.
func (rm *RoomManager) FindParticipantByClientID(clientID string) (*Participant, *Room) {
	rm.roomsMutex.RLock()
	defer rm.roomsMutex.RUnlock()

	for _, room := range rm.rooms {
		room.Mutex.RLock()
		participant, exists := room.Participants[clientID]
		room.Mutex.RUnlock()

		if exists {
			return participant, room
		}
	}
	return nil, nil
}

// AddTrackToAllParticipants adds a track to all participants in all rooms.
func (rm *RoomManager) AddTrackToAllParticipants(track webrtc.TrackLocal) {
	rm.roomsMutex.RLock()
	defer rm.roomsMutex.RUnlock()

	for _, room := range rm.rooms {
		room.Mutex.Lock()
		for _, participant := range room.Participants {
			participant.Mutex.Lock()
			participant.PeerConnection.AddTrack(track)
			participant.Tracks[track.ID()] = track
			participant.Mutex.Unlock()
		}
		room.Mutex.Unlock()
	}
}

// Wrtp writes an RTP packet to all tracks of all participants.
func (rm *RoomManager) Wrtp(rtp *rtp.Packet, room *Room) {
	rm.roomsMutex.RLock()
	defer rm.roomsMutex.RUnlock()

	room.Mutex.RLock()
	for _, participant := range room.Participants {
		participant.Mutex.RLock()
		for _, track := range participant.Tracks {
			if staticTrack, ok := track.(*webrtc.TrackLocalStaticRTP); ok {
				err := staticTrack.WriteRTP(rtp)
				if err != nil {
					fmt.Printf("Failed to write RTP to track %s: %v\n", staticTrack.ID(), err)
				}
			}
		}
		participant.Mutex.RUnlock()
	}
	room.Mutex.RUnlock()
}

func (rm *RoomManager) RenegotClientsAround() {
	rm.roomsMutex.RLock()
	defer rm.roomsMutex.RUnlock()

	for _, room := range rm.rooms {
		rm.renegotiationSvc.RenegotiateParticipants(room)
	}
}

// GetOrCreateRoom gets or creates a room by ID.
func (rm *RoomManager) GetOrCreateRoom(roomID string) *Room {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		room = &Room{
			ID:           roomID,
			Participants: make(map[string]*Participant),
		}
		rm.rooms[roomID] = room
	}
	return room
}

// DeleteRoom removes a room by ID.
func (rm *RoomManager) DeleteRoom(roomID string) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	delete(rm.rooms, roomID)
}

// RemoveParticipantFromRoom removes a participant and updates the clientIDToRoomID map.
func (rm *RoomManager) RemoveParticipantFromRoom(roomID, clientID string) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return
	}

	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	// Remove the participant from the room
	delete(room.Participants, clientID)

	// Remove the mapping from clientID to roomID
	delete(rm.clientIDToRoomID, clientID)

	// If the room is empty, delete it
	if len(room.Participants) == 0 {
		delete(rm.rooms, roomID)
	}
}
