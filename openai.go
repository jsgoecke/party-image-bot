package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// PromptsImages is a struct for sending prompts and images to the client
type PromptsImages struct {
	MessageType string `json:"message_type,omitempty"`
	HumanPrompt string `json:"human_prompt,omitempty"`
	AiPrompt    string `json:"ai_prompt,omitempty"`
	AiImage     string `json:"ai_image,omitempty"`
	From        string `json:"from,omitempty"`
	Status      string `json:"status,omitempty"`
	mx          sync.Mutex
}

func (pi *PromptsImages) SaveToDB() {
	pi.mx.Lock()
	defer pi.mx.Unlock()

	jsonData, err := json.Marshal(pi)
	if err != nil {
		log.Fatalf("Error marshaling data to JSON: %s", err)
	}
	err = rdb.Set(uuid.New().String(), jsonData, 0).Err()
	if err != nil {
		log.Fatalf("Error setting key in Redis: %s", err)
	}
}

// SendMessage sends a message to the client
func sendMessage() {
	for event := range promptsImagesChan {
		log.Print("Received promptsImagesChan...")
		log.Print(event)
		body, err := json.Marshal(event)
		if err != nil {
			log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
		}
		if ws != nil {
			ws.WriteMessage(1, []byte(body))
			log.Println("Sent message to client: " + string(body))
		} else {
			log.Println("Dropped message to client, websocket closed: " + string(body))
		}
		if event.Status == "IMAGES-GENERATED" {
			time.Sleep(MessageDelay * time.Second)
		}
	}
}

// POST /api/v1/sms
func processSMS(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	promptsImages := &PromptsImages{}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	promptsImages.HumanPrompt = r.Form.Get("Body")
	promptsImages.From = r.Form.Get("From")
	log.Println("Generating ai prompt...")
	promptsImages.AiPrompt = embellishPrompt(promptsImages.HumanPrompt)
	log.Println("Generating ai prompted image...")
	promptsImages.AiImage = createImage(promptsImages.AiPrompt)
	promptsImages.Status = "IMAGES-GENERATED"
	promptsImages.SaveToDB()
	if promptsImages.AiImage != "" {
		promptsImagesChan <- *promptsImages
	} else {
		log.Println("Error generating image with embellished prompt: " + promptsImages.AiPrompt)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func isValidURL(url string) bool {
	re := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return re.MatchString(url)
}

func downloadImage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	guid := uuid.New()
	filepath := "web/images/" + guid.String() + ".jpg"

	file, err := os.Create(filepath)
	if err != nil {
		return err.Error()
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println("Error downloading image: " + err.Error())
		return err.Error()
	}
	return filepath
}

// Create an image from a prompt
func createImage(prompt string) string {
	c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	reqUrl := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	respUrl, err := c.CreateImage(ctx, reqUrl)
	if err != nil {
		log.Println("Image creation error: %v\n", err)
		return err.Error()
	}

	log.Println(respUrl.Data[0])
	if !isValidURL(respUrl.Data[0].URL) {
		return ""
	} else {
		image := downloadImage(respUrl.Data[0].URL)
		log.Println("Image URL: " + image)
		return image
	}
}

// Embellish a prompt
func embellishPrompt(prompt string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	messages := make([]openai.ChatCompletionMessage, 0)

	prompt = SystemPrompt + prompt
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4TurboPreview,
			Messages: messages,
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return err.Error()
	}

	content := resp.Choices[0].Message.Content
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})
	if len(content) < 1000 {
		return content
	}
	return `
		You are a master at creating prompts for DALL-E images in the genre of the
		Hitchhiker's Guide to the Galaxy by Douglas Adams. Generate an image that 
		would be at home in the Hitchhiker's guide to the galaxy. Use cats as the
		content of the image.
	`
}
