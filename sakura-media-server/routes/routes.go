package routes

import (
	"net/http"
	"sakura/controllers"
)

const APIBase = "/api"


func RegisterRoutes() {
	http.HandleFunc( APIBase + "/mediaserver", controllers.MediaServerHandler)
}