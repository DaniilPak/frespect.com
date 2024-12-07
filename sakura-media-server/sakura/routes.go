package sakura

import (
	"net/http"
)

const APIBase = "/api"

func RegisterRoutes() {
	http.HandleFunc(APIBase+"/mediaserver", SessionHandler)
	http.HandleFunc(APIBase+"/bot", BotHandler)
	http.HandleFunc(APIBase+"/renegot", RenegotHandler)
}
