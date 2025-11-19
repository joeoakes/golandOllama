package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Request body sent to Ollama
type TranslateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Response body from Ollama (non-streaming)
type TranslateResponse struct {
	Response string `json:"response"`
}

func main() {
	// Prompt asks for ONLY the German word
	prompt := "Translate the English word 'Hello' into German. " +
		"Respond with only the German word and nothing else."

	reqBody := TranslateRequest{
		Model:  "llama3", // or whichever model you have pulled
		Prompt: prompt,
		Stream: false, // single JSON response (no streaming)
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Ollama returned status: " + resp.Status)
	}

	var result TranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("German for 'Hello' is:", result.Response)
}
