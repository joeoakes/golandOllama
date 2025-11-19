package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"` // <- NEW
}

type GenerateResponse struct {
	Response string `json:"response"`
	// You can add other fields if you want, e.g. Done bool `json:"done"`
}

func main() {
	url := "http://localhost:11434/api/generate"

	reqBody := GenerateRequest{
		Model:  "llama3",
		Prompt: "Explain photosynthesis in 2 sentences.",
		Stream: false, // <- IMPORTANT
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("bad status: " + resp.Status)
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("AI Response:", result.Response)
}
