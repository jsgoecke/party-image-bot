package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/sashabaranov/go-openai"
)

// PromptsImages is a struct for sending prompts and images to the client
type PromptsImages struct {
	MessageType string `json:"message_type,omitempty"`
	HumanPrompt string `json:"human_prompt,omitempty"`
	HumanImage  string `json:"human_image,omitempty"`
	AiPrompt    string `json:"ai_prompt,omitempty"`
	AiImage     string `json:"ai_image,omitempty"`
	From        string `json:"from,omitempty"`
	Status      string `json:"status,omitempty"`
	mx          sync.Mutex
}

// SendMessage sends a message to the client
func sendMessage(promptImages PromptsImages) {
	body, err := json.Marshal(promptImages)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	ws.WriteMessage(1, []byte(body))
	log.Println("Sent message to client: " + string(body))
}

// POST /api/v1/sms
func processSMS(w http.ResponseWriter, r *http.Request) {
	promptsImages := &PromptsImages{}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	promptsImages.HumanPrompt = r.Form.Get("Body")
	promptsImages.From = r.Form.Get("From")
	promptsImages.Status = "SMS-RECEIVED"
	log.Println("Received SMS from -> " + promptsImages.From + " with prompt -> " + promptsImages.HumanPrompt)
	sendMessage(*promptsImages)

	log.Println("Generating human prompted image...")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		promptsImages.HumanImage = createImage(promptsImages.HumanPrompt)
	}()

	log.Println("Generating ai prompt...")
	promptsImages.AiPrompt = embellishPrompt(promptsImages.HumanPrompt)
	log.Println("Generating ai prompted image...")
	promptsImages.AiImage = createImage(promptsImages.AiPrompt)
	promptsImages.Status = "IMAGES-GENERATED"
	wg.Wait()

	log.Println("Sending images to client...")
	sendMessage(*promptsImages)

	w.WriteHeader(http.StatusOK)
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
		log.Printf("Image creation error: %v\n", err)
		return err.Error()
	}

	return respUrl.Data[0].URL
}

// Embellish a prompt
func embellishPrompt(prompt string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	messages := make([]openai.ChatCompletionMessage, 0)

	prompt = "Create a prompt for DALL-E with: " + prompt + ", in the style of Hitchhiker's Guide to the Galaxy by Douglas Adams. The prompt should be equal to or less than 1000 characters."
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
	return content
}
