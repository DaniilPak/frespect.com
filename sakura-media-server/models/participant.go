package models

import (
	"sync"

	"github.com/pion/webrtc/v4"
)

type Participant struct {
	ClientID       string
	PeerConnection *webrtc.PeerConnection
	Tracks         map[string]*webrtc.TrackLocalStaticRTP
	Mutex          sync.RWMutex
}
