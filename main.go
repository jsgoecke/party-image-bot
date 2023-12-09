package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var ws *websocket.Conn

// Log HTTP requests
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// Websocket Upgrader
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Websocket endpoint
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
	ws, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Websocket Client Connected")
	ws.WriteMessage(1, []byte("{ \"status\": \"CONNECTED\" }"))
}

func main() {
	// PWA
	app.Route("/", &IndexPage{})
	app.RunWhenOnBrowser()

	// Standard HTTP routing (server-side):
	http.Handle("/", &app.Handler{
		Name:        "Party Image Bot",
		Description: "A bot for generating images based on text prompts.",
	})

	// REST API
	http.Handle("/api/v1/sms", logRequest(http.HandlerFunc(processSMS)))

	// WS API
	http.Handle("/api/v1/ws", logRequest(http.HandlerFunc(wsEndpoint)))

	// Serve HTTP
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
