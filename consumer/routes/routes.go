package routes

import (
	"consumer/handler"
	"net/http"
)

func RegisterRoutes() {
	// accepts a message and pushes it to kafka topic along with other details
	http.HandleFunc("/getValueForId", handler.GetValueForId)
}
