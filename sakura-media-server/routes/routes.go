package routes

import (
	"net/http"
	"sakura/controllers"
)

const APIBase = "/api"

func RegisterRoutes() {
	http.HandleFunc(APIBase+"/mediaserver", controllers.MediaServerHandler)
	http.HandleFunc(APIBase+"/bot", controllers.BotHandler)
	http.HandleFunc(APIBase+"/renegot", controllers.RenegotHandler)
}
