package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// MessageDelay is the delay between messages in seconds
const MessageDelay = 120
const AiImageSize = "1024px"
const IconImageSize = "180px"
const QRCodeImageSize = "240px"
const SystemPrompt = `
You are a master at creating prompts for DALL-E images in the genre of the
Hitchhiker's Guide to the Galaxy by Douglas Adams.  

You are also a DALL-E-3 prompt engineer.  You are tasked with creating a prompt
that is no more than 1000 characters long and does not violate the safety guidelines of 
DALL-E-3.

Now please create an image prompt that would be at home in the Hitchhiker's 
guide based on the following prompt details: 
`

var promptsImagesChan chan PromptsImages
var ws *websocket.Conn
var rdb *redis.Client

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
	// Initialize our message queue
	promptsImagesChan = make(chan PromptsImages, 1000)

	// Initialize Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DATABASE"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
	})

	go sendMessage()

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
