package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sakura/models"
	"sakura/roommanager"
	"sakura/sfu"
	"sakura/utils"

	"github.com/pion/webrtc/v4"
)

var mediaManagerURL string = sfu.GetMediaManagerURL()
var audioURL = mediaManagerURL + "media/_Fo6n3nl_Sk"

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
	var requestData models.SDPRenegot
	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Printf("Error parsing JSON: %s", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("inside renegot hanlder %s", requestData.SDP)

	curParticipant, curRoom := roommanager.FindParticipantByClientID(requestData.ClientID)
	if curParticipant == nil || curRoom == nil {
		panic(fmt.Sprintf("Participant or Room not found for ClientID: %s", requestData.ClientID))
	}

	fmt.Printf("Participant details: %+v\n", curParticipant)

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	utils.Decode(requestData.SDP, &offer)

	// Set the remote description
	if err := curParticipant.PeerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// Create a music track (audio) to add to the peer connection
	musicTrack := sfu.CreateAudioTrack()

	// Add the track to the peer connection
	rtpSender, err := curParticipant.PeerConnection.AddTrack(musicTrack)
	if err != nil {
		panic(err)
	}

	go sfu.WriteAudioToTrack(audioURL, musicTrack)

	// Wait for ICE gathering to complete
	gatherComplete := webrtc.GatheringCompletePromise(curParticipant.PeerConnection)
	<-gatherComplete

	// Create an SDP answer and set it as the local description
	answer, err := curParticipant.PeerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	if err := curParticipant.PeerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Block until ICE gathering completes
	<-gatherComplete

	// Send the local SDP answer back to the client
	// Output the answer in base64 so we can paste it in browser
	encodedAnswer := utils.Encode(curParticipant.PeerConnection.LocalDescription())

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

	// Handle RTCP packets (optional, depending on your application)
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()
	// You can also add logic to stream video/audio from disk here, similar to the play-from-disk renegotiation
	// Use a function similar to writeVideoToTrack, where you can stream video or audio data to the track.
}

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
	var requestData models.BotRequest
	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Printf("Error parsing JSON: %s", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Managing room
	songId := requestData.SongId
	fmt.Printf("Bot request: %s\n", songId)

	go sfu.StartMusicEverywhere()
}

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
		Tracks:   make(map[string]webrtc.TrackLocal),
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
