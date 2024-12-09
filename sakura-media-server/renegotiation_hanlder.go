package sakura

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pion/webrtc/v4"
)

func RenegotHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Renegot hanlder invoked")
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
	var requestData SDPRenegotiateRequest
	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Printf("Error parsing JSON: %s", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("inside renegot hanlder %s", requestData.SDP)

	roomManager := GetRoomManager()
	curParticipant, curRoom := roomManager.FindParticipantByClientID(requestData.ClientID)
	if curParticipant == nil || curRoom == nil {
		panic(fmt.Sprintf("Participant or Room not found for ClientID: %s", requestData.ClientID))
	}

	fmt.Printf("Participant details: %+v\n", curParticipant)

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	Decode(requestData.SDP, &offer)

	// Set the remote description
	if err := curParticipant.peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// bot, err := NewBot()
	// if err != nil {
	// 	panic("Bot failed to create")
	// }

	// // Create a music track (audio) to add to the peer connection
	// musicTrack := bot.CreateAudioTrack()

	// // Add the track to the peer connection
	// rtpSender, err := curParticipant.PeerConnection.AddTrack(musicTrack)
	// if err != nil {
	// 	panic(err)
	// }

	// go bot.WriteAudioToTrack(audioURL, musicTrack)

	// Wait for ICE gathering to complete
	gatherComplete := webrtc.GatheringCompletePromise(curParticipant.peerConnection)
	<-gatherComplete

	// Create an SDP answer and set it as the local description
	answer, err := curParticipant.peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	if err := curParticipant.peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Block until ICE gathering completes
	<-gatherComplete

	// Send the local SDP answer back to the client
	// Output the answer in base64 so we can paste it in browser
	encodedAnswer := Encode(curParticipant.peerConnection.LocalDescription())

	payload := map[string]string{
		"sdp": encodedAnswer,
	}

	response, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Write the response
	if _, err := w.Write(response); err != nil {
		panic(err)
	}

	// // Handle RTCP packets (optional, depending on your application)
	// go func() {
	// 	rtcpBuf := make([]byte, 1500)
	// 	for {
	// 		if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
	// 			return
	// 		}
	// 	}
	// }()
}
