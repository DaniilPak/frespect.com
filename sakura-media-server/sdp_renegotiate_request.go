package sakura

type SDPRenegotiateRequest struct {
	SDP      string `json:"sdp"`
	ClientID string `json:"clientId"`
}
