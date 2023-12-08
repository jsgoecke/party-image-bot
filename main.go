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
	app.Route("/", &IndexPage{
		HumanPrompt: "Create a flying dog",
		HumanImage:  "https://oaidalleapiprodscus.blob.core.windows.net/private/org-YzgP90GpZ9ZO8PPVVTNBx5dG/user-y3BQa2xR6DtHzd5qCYyZUoRa/img-ClT8nAJ8L0Y7tKvhAOParNAl.png?st=2023-12-07T13%3A59%3A10Z\u0026se=2023-12-07T15%3A59%3A10Z\u0026sp=r\u0026sv=2021-08-06\u0026sr=b\u0026rscd=inline\u0026rsct=image/png\u0026skoid=6aaadede-4fb3-4698-a8f6-684d7786b067\u0026sktid=a48cca56-e6da-484e-a814-9c849652bcb3\u0026skt=2023-12-07T02%3A03%3A35Z\u0026ske=2023-12-08T02%3A03%3A35Z\u0026sks=b\u0026skv=2021-08-06\u0026sig=6iz8quqNFKNJYxu6PWrVEgvWfBTnYNLnhSftXLIPVxU%3D",
		AiPrompt:    "Imagine a whimsically surreal scene straight from the universe of 'The Hitchhiker's Guide to the Galaxy' by Douglas Adams: A quirky, intelligent, and improbably flying dog, with a bewildered expression and oversized, flapping ears, soaring through the star-speckled cosmos. The dog should have a digital wristwatch clumsily strapped to one paw, and a towel draped over its back - an indispensable item for any galactic hitchhiker. Around it, robots and peculiar alien craft hover, depicted with a mix of absurdity and British dry humor characteristic of Adams' storytelling. Convey the scene with vibrant, yet slightly faded comic-like illustrations that capture the spirit of the original book covers and illustrations associated with the series.",
		AiImage:     "https://oaidalleapiprodscus.blob.core.windows.net/private/org-YzgP90GpZ9ZO8PPVVTNBx5dG/user-y3BQa2xR6DtHzd5qCYyZUoRa/img-IpLbfU9P6Uk9gAqnNecBmOUL.png?st=2023-12-07T14%3A04%3A16Z\u0026se=2023-12-07T16%3A04%3A16Z\u0026sp=r\u0026sv=2021-08-06\u0026sr=b\u0026rscd=inline\u0026rsct=image/png\u0026skoid=6aaadede-4fb3-4698-a8f6-684d7786b067\u0026sktid=a48cca56-e6da-484e-a814-9c849652bcb3\u0026skt=2023-12-07T02%3A07%3A57Z\u0026ske=2023-12-08T02%3A07%3A57Z\u0026sks=b\u0026skv=2021-08-06\u0026sig=KMfe3WWHHL0QRFUG1s12eRHIdl6hx611MD9Fv17UXKQ%3D",
	})
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
