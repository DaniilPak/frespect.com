package sakura

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/webrtc/v4"

	"github.com/pion/webrtc/v4/pkg/media"
	"github.com/pion/webrtc/v4/pkg/media/oggreader"
)

var serverURL string = GetServerURL("answer")
var serverURLrenegot string = GetServerURL("renegot")

const (
	oggPageDuration = time.Millisecond * 20
)

func StartSession(sdpRequest SDPRequest, room *Room, participant *Participant) {
	m := &webrtc.MediaEngine{}

	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 2, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}

	i := &interceptor.Registry{}

	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		panic(err)
	}

	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		panic(err)
	}
	i.Add(intervalPliFactory)

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:global.stun.twilio.com:3478"},
			},
		},
	}
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	outputTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		panic(err)
	}

	rtpSender, err := peerConnection.AddTrack(outputTrack)
	if err != nil {
		panic(err)
	}

	participant.Tracks[outputTrack.ID()] = outputTrack

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	offer := webrtc.SessionDescription{}
	Decode(sdpRequest.SDP, &offer)

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	participant.Mutex.Lock()
	participant.PeerConnection = peerConnection
	participant.Mutex.Unlock()

	roomManager := GetRoomManager()

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) { //nolint: revive
		fmt.Printf("Track has started, of type %d: %s \n", track.PayloadType(), track.Codec().MimeType)
		for {
			rtp, _, readErr := track.ReadRTP()
			if readErr != nil {
				panic(readErr)
			}

			roomManager.Wrtp(rtp)
		}
	})

	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateClosed || s == webrtc.PeerConnectionStateFailed {
			roomManager.RemoveParticipantFromRoom(sdpRequest.RoomID, participant.ClientID)
		}
	})

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	encodedAnswer := Encode(peerConnection.LocalDescription())

	// Send remote SDP data to signalling server
	// Create JSON payload to send to localhost:4000
	payload := map[string]string{
		"clientId":  sdpRequest.ClientID,
		"sdpAnswer": encodedAnswer,
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

	// Block forever
	select {}
}

func BotStart() *webrtc.TrackLocalStaticSample {
	// Fetch the audio file from HTTP API
	resp, err := http.Get(audioURL)
	if err != nil {
		fmt.Println("Failed to fetch audio file.")
		return nil
	}

	// Create an audio track for music
	musicTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio-music", "pion-music")
	if err != nil {
		panic(err)
	}

	go func() {
		ogg, _, err := oggreader.NewWith(resp.Body)
		if err != nil {
			panic(err)
		}

		ticker := time.NewTicker(oggPageDuration)
		defer ticker.Stop()

		var lastGranule uint64
		for range ticker.C {
			pageData, pageHeader, err := ogg.ParseNextPage()
			if errors.Is(err, io.EOF) {
				fmt.Println("All audio pages parsed.")
				break
			}
			if err != nil {
				panic(err)
			}

			sampleCount := float64(pageHeader.GranulePosition - lastGranule)
			lastGranule = pageHeader.GranulePosition
			sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

			if writeErr := musicTrack.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); writeErr != nil {
				panic(writeErr)
			}
		}

		// Close the response body once the audio processing is done
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	return musicTrack
}

func DoRenegotiationAll() {
	roomManager := GetRoomManager()
	roomManager.RenegotAll(serverURLrenegot)
}

func StartMusicEverywhere() {
	DoRenegotiationAll()
}

func CreateAudioTrack() *webrtc.TrackLocalStaticSample {
	// Create an audio track for music
	musicTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio-music", "pion-music")
	if err != nil {
		panic(err)
	}
	return musicTrack
}

func WriteAudioToTrack(audioURL string, musicTrack *webrtc.TrackLocalStaticSample) {
	// Fetch the audio file from HTTP API
	resp, err := http.Get(audioURL)
	if err != nil {
		fmt.Println("Failed to fetch audio file.")
		return
	}
	defer resp.Body.Close()

	ogg, _, err := oggreader.NewWith(resp.Body)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(oggPageDuration)
	defer ticker.Stop()

	var lastGranule uint64
	for range ticker.C {
		pageData, pageHeader, err := ogg.ParseNextPage()
		if errors.Is(err, io.EOF) {
			fmt.Println("All audio pages parsed.")
			break
		}
		if err != nil {
			panic(err)
		}

		sampleCount := float64(pageHeader.GranulePosition - lastGranule)
		lastGranule = pageHeader.GranulePosition
		sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

		if writeErr := musicTrack.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); writeErr != nil {
			panic(writeErr)
		}
	}
}
