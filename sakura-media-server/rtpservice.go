package sakura

import (
	"fmt"
	"sync"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type RTPService struct{}

// Singleton instance and sync.Once for RTPService
var (
	rtpServiceInstance *RTPService
	rtpServiceOnce     sync.Once
)

// GetRTPService returns the singleton instance of RTPService.
func GetRTPService() *RTPService {
	rtpServiceOnce.Do(func() {
		rtpServiceInstance = &RTPService{}
	})
	return rtpServiceInstance
}

func (r *RTPService) WriteRTPToParticipants(rtp *rtp.Packet, room *Room) {
	room.Mutex.RLock()
	defer room.Mutex.RUnlock()

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
}
