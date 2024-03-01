package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var DEBUG = false // Set to true to enable more verbose logging

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}
type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	// Logprobs interface{} `json:"logprobs"`
}
type Response struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	SystemFingerprint  interface{}   `json:"system_fingerprint"`	
	Usage   struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Define additional structs to handle streaming response format
type Delta struct {
	Role        string `json:"role,omitempty"`
	Content     string `json:"content,omitempty"`
}
type StreamedChoice struct {
	Index         int    `json:"index"`
	Delta         Delta  `json:"delta"`
	// Logprobs      interface{} `json:"logprobs"`
	FinishReason  string `json:"finish_reason,omitempty"`
}
type StreamedResponse struct {
	ID                string           `json:"id"`
	Object            string           `json:"object"`
	Created           int64            `json:"created"`
	Model             string           `json:"model"`
	SystemFingerprint string           `json:"system_fingerprint,omitempty"`
	Choices           []StreamedChoice `json:"choices"`
}


// GenerateRandomID generates a random ID similar to the given format
func GenerateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 30)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "chatcmpl-" + string(b)
}

// CalculatePromptTokens calculates the total length of all "content" fields in the messages divided by 3.5
func CalculatePromptTokens(messages []map[string]interface{}) int {
	totalLength := 0
	for _, message := range messages {
		if content, ok := message["content"].(string); ok {
			totalLength += len(content)
		}
	}
	// Divide by 3.5 and round to nearest integer
	return int(float64(totalLength) / 3.5)
}

// Handler for the "/v1/chat/completions" endpoint
func chatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	if DEBUG {
		log.Printf("Received %s request for %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	}

	var requestBody map[string]interface{}

	// Check if the request body is not empty
	if r.Body != nil && r.ContentLength > 0 {
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
					log.Printf("Error decoding request body: %v", err)
					// You might choose not to return an error here if the body is optional
					// and just proceed with default values instead.
			}
	}

	if DEBUG {
		log.Println("Request Body:", requestBody)
	}

	if stream, ok := requestBody["stream"].(bool); ok && stream {
			// Handle streaming response
			streamResponse(w, requestBody)
	} else {
			// Handle regular response
			regularResponse(w, requestBody)
	}
}


func streamResponse(w http.ResponseWriter, requestBody map[string]interface{}) {
	// Simulate a streaming response
	id := GenerateRandomID()
	streamedResponse := StreamedResponse{
			ID:                id,
			Object:            "chat.completion.chunk",
			Created:           time.Now().Unix(),
			Model:             "gpt-3.5-turbo", // Default or extracted from requestBody
			SystemFingerprint: "fp_44709d6fcb", // Example value
			Choices: []StreamedChoice{
					{
							Index: 0,
							Delta: Delta{
									Content: "Blank response from OpenAI API emulator.",
							},
							FinishReason: "STOP",
					},
			},
	}

	// Set Content-Type for streaming
	w.Header().Set("Content-Type", "application/json")
	// Encode and send the streamed response
	if err := json.NewEncoder(w).Encode(streamedResponse); err != nil {
			log.Printf("Error streaming response: %v", err)
			http.Error(w, "Error streaming response", http.StatusInternalServerError)
	}
}

func regularResponse(w http.ResponseWriter, requestBody map[string]interface{}) {
	// Prepare the response structure
	response := Response{
			ID:      GenerateRandomID(),
			Choices: []Choice{
					{
							FinishReason: "STOP",
							Index:        0,
							Message: Message{
									Content: "Blank response from OpenAI API emulator.",
									Role:    "assistant",
							},
					},
			},
			Created: time.Now().Unix(),
			Model:   "gpt-3.5-turbo", // Default model value, could be overridden below
			Object:  "chat.completion",
			Usage: struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
			}{
					PromptTokens:     57, // Default prompt tokens, could be recalculated below
					CompletionTokens: 8,  // Fixed completion tokens
					TotalTokens:      65, // Default total tokens, could be recalculated below
			},
	}

	// Extract and set model from request if present
	if model, ok := requestBody["model"].(string); ok {
			response.Model = model
	}

	// Calculate prompt tokens if messages are present
	if messages, ok := requestBody["messages"].([]interface{}); ok {
			mappedMessages := make([]map[string]interface{}, len(messages))
			for i, msg := range messages {
					if msgMap, ok := msg.(map[string]interface{}); ok {
							mappedMessages[i] = msgMap
					}
			}
			promptTokens := CalculatePromptTokens(mappedMessages)
			response.Usage.PromptTokens = promptTokens
			response.Usage.TotalTokens = promptTokens + response.Usage.CompletionTokens
	}

	// Set Content-Type and return the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	} else {
		if DEBUG {
			log.Printf("Response sent: %v", response)
		}
	}
}

func main() {
	http.HandleFunc("/v1/chat/completions", chatCompletionsHandler)

	log.Println("Starting server on :8383")
	if err := http.ListenAndServe(":8383", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}