package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Request sent to Ollama
type TravelRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Non-streaming response from Ollama
type TravelResponse struct {
	Response string `json:"response"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== AI Travel Guide ===")
	fmt.Print("Enter a destination (city, country, or region): ")

	dest, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	dest = strings.TrimSpace(dest)
	if dest == "" {
		fmt.Println("No destination entered, exiting.")
		return
	}

	// Prompt for Ollama
	prompt := fmt.Sprintf(
		"You are a helpful travel guide. "+
			"List 5–7 interesting places to see in %s. "+
			"Use short bullet points, each on its own line.",
		dest,
	)

	reqBody := TravelRequest{
		Model:  "llama3", // change if you use a different model
		Prompt: prompt,
		Stream: false, // single JSON object back
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Error calling Ollama:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ollama returned non-OK status:", resp.Status)
		return
	}

	var result TravelResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println()
	fmt.Printf("Top sights in %s:\n", dest)
	fmt.Println("--------------------------------")
	fmt.Println(result.Response)
}
