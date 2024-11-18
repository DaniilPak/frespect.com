package models

type SDPRequest struct {
	SDP      string `json:"sdp"`
	ClientID string `json:"clientId"`
	RoomID   string `json:"roomId"`
}
