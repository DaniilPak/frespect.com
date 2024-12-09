package sakura

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RenegotiationService struct{}

var serverURLrenegot string = GetServerURL("renegot")

func (r *RenegotiationService) RenegotiateParticipants(room *Room) {
	room.mutex.RLock()
	for _, participant := range room.participants {
		payload := map[string]string{
			"clientId": participant.clientID,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", serverURLrenegot, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Failed to create POST request: %v\n", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "your-secret-key")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send POST request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Successfully sent answer to server")
		} else {
			fmt.Printf("Failed to send answer with status code: %d\n", resp.StatusCode)
		}
	}
	room.mutex.RUnlock()
}
