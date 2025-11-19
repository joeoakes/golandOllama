package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaRequestIntel struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponseIntel struct {
	Response string `json:"response"`
}

func main() {
	reqBody := OllamaRequestIntel{
		Model: "llama3", // or whatever model you pulled
		Prompt: "What is the current stock price of Intel (ticker INTC) in US dollars? " +
			"Only output the number (no currency symbol, no words).",
		Stream: false, // single JSON response
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result OllamaResponseIntel
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Ollama says Intel price is:", result.Response)
}
