package sakura

import (
	"fmt"
	"sync"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type Room struct {
	ID           string
	participants map[string]*Participant
	// bots         map[string]*Bot
	mutex sync.RWMutex
}

type Participant struct {
	clientID       string
	peerConnection *webrtc.PeerConnection
	tracks         map[string]webrtc.TrackLocal
	rtpSenders     map[string]*webrtc.RTPSender
	mutex          sync.RWMutex
}

type RoomManager struct {
	rooms            map[string]*Room
	clientIDToRoomID sync.Map
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

// AddTrackToParticipant adds a track to a specific participant identified by clientID.
func (rm *RoomManager) AddTrackToParticipant(clientID string, botID string, track webrtc.TrackLocal) error {
	// Find the participant and their room
	participant, room := rm.FindParticipantByClientID(clientID)
	if participant == nil || room == nil {
		return fmt.Errorf("participant with clientID %s not found", clientID)
	}

	if participant.tracks == nil {
		fmt.Println("Participant.tracks is nil!")
	}

	for _, participant := range room.participants {
		// Lock the participant and add the track
		participant.mutex.Lock()

		rtpSender, err := participant.peerConnection.AddTrack(track)
		if err != nil {
			participant.mutex.Unlock() // Unlock immediately if there was an error
			return fmt.Errorf("failed to add track to participant: %v", err)
		}

		participant.rtpSenders[botID] = rtpSender
		participant.tracks[track.ID()] = track

		participant.mutex.Unlock() // Unlock after processing the participant
	}

	return nil
}

// RemoveTrackFromRoom removes a track by rtpSender from a specific room using botID.
func (rm *RoomManager) RemoveTrackFromRoom(room *Room, botID string) error {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	// Iterate over participants in the room
	for _, participant := range room.participants {
		participant.mutex.Lock()
		defer participant.mutex.Unlock()

		// Find the rtpSender for the given botID
		rtpSender, exists := participant.rtpSenders[botID]
		if !exists {
			participant.mutex.Unlock()
			return fmt.Errorf("rtpSender with botID %s not found", botID)
		}

		// Remove the track from the participant's track list and sender list
		delete(participant.rtpSenders, botID)
		delete(participant.tracks, rtpSender.Track().ID())

		participant.peerConnection.RemoveTrack(rtpSender)

		// Close the rtpSender if needed
		err := rtpSender.Stop()
		if err != nil {
			participant.mutex.Unlock()
			return fmt.Errorf("failed to stop rtpSender for botID %s: %v", botID, err)
		}
	}

	return nil
}

// FindParticipantByClientID locates a participant and its room by client ID.
func (rm *RoomManager) FindParticipantByClientID(clientID string) (*Participant, *Room) {
	rm.roomsMutex.RLock()
	defer rm.roomsMutex.RUnlock()

	for _, room := range rm.rooms {
		room.mutex.RLock()
		participant, exists := room.participants[clientID]
		room.mutex.RUnlock()

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
		room.mutex.Lock()
		for _, participant := range room.participants {
			participant.mutex.Lock()
			participant.peerConnection.AddTrack(track)
			participant.tracks[track.ID()] = track
			participant.mutex.Unlock()
		}
		room.mutex.Unlock()
	}
}

// Wrtp writes an RTP packet to all tracks of all participants.
func (rm *RoomManager) Wrtp(rtp *rtp.Packet, room *Room) {
	rm.roomsMutex.RLock()
	defer rm.roomsMutex.RUnlock()

	room.mutex.RLock()
	for _, participant := range room.participants {
		participant.mutex.RLock()
		for _, track := range participant.tracks {
			if staticTrack, ok := track.(*webrtc.TrackLocalStaticRTP); ok {
				err := staticTrack.WriteRTP(rtp)
				if err != nil {
					fmt.Printf("Failed to write RTP to track %s: %v\n", staticTrack.ID(), err)
				}
			}
		}
		participant.mutex.RUnlock()
	}
	room.mutex.RUnlock()
}

func (rm *RoomManager) RenegotiateAllClients() {
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
			participants: make(map[string]*Participant),
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

	room.mutex.Lock()
	defer room.mutex.Unlock()

	// Remove the participant from the room
	delete(room.participants, clientID)

	// If the room is empty, delete it
	if len(room.participants) == 0 {
		delete(rm.rooms, roomID)
	}
}
