package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func main() {
	url := "http://localhost:11434/api/generate"

	// ***** Change these numbers *****
	a := 7
	b := 13

	prompt := fmt.Sprintf("Add %d and %d. Only output the number.", a, b)

	reqBody := OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false, // return single JSON response
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result OllamaResponse
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Printf("AI says: %s\n", result.Response)
}
