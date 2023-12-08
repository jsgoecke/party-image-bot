package main

import (
	"testing"
)

func TestPromptsImages(t *testing.T) {
	// Create an instance of PromptsImages
	p := PromptsImages{
		HumanPrompt: "TestHumanPrompt",
		HumanImage:  "TestHumanImage.jpg",
		AiPrompt:    "TestAiPrompt",
		AiImage:     "TestAiImage.jpg",
	}

	// Check if the fields have the expected values
	if p.HumanPrompt != "TestHumanPrompt" {
		t.Errorf("Expected HumanPrompt to be 'TestHumanPrompt', got '%s'", p.HumanPrompt)
	}
	if p.HumanImage != "TestHumanImage.jpg" {
		t.Errorf("Expected HumanImage to be 'TestHumanImage.jpg', got '%s'", p.HumanImage)
	}
	if p.AiPrompt != "TestAiPrompt" {
		t.Errorf("Expected AiPrompt to be 'TestAiPrompt', got '%s'", p.AiPrompt)
	}
	if p.AiImage != "TestAiImage.jpg" {
		t.Errorf("Expected AiImage to be 'TestAiImage.jpg', got '%s'", p.AiImage)
	}
}
