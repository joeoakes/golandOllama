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

// Payload sent to Ollama
type OllamaCalcRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Non-streaming response from Ollama
type OllamaCalcResponse struct {
	Response string `json:"response"`
}

// JSON we EXPECT Ollama to output
// Example: {"a": 2.5, "b": 7, "op": "add"}
type ParsedExpression struct {
	A  float64 `json:"a"`
	B  float64 `json:"b"`
	Op string  `json:"op"` // "add", "sub", "mul", "div"
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Ollama Calculator ===")
	fmt.Println("Type an expression like:")
	fmt.Println("  - add 2 and 3")
	fmt.Println("  - 10 / 4")
	fmt.Println("  - subtract 5 from 12")
	fmt.Println("Type 'exit' to quit.")
	fmt.Println()

	for {
		fmt.Print("Enter expression: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		line = strings.TrimSpace(line)
		if strings.EqualFold(line, "exit") || line == "" {
			fmt.Println("Goodbye.")
			return
		}

		parsed, err := askOllamaToParse(line)
		if err != nil {
			fmt.Println("Error parsing expression with Ollama:", err)
			continue
		}

		result, err := compute(parsed)
		if err != nil {
			fmt.Println("Error computing result:", err)
			continue
		}

		fmt.Printf("Result: %.6f\n\n", result)
	}
}

// askOllamaToParse sends the user's expression to Ollama and expects JSON back.
//
// Example prompt to Ollama:
//
//	User input: "add 2 and 3"
//	Ollama should respond: {"a":2,"b":3,"op":"add"}
func askOllamaToParse(userInput string) (*ParsedExpression, error) {
	// Prompt engineered so the model only outputs JSON we can parse.
	prompt := fmt.Sprintf(`
You are a math expression parser.

The user will provide an arithmetic expression in plain English or with symbols.
You must convert it into a JSON object with this exact structure:

{
  "a": <number>,
  "b": <number>,
  "op": "<operation>"
}

Where:
- "op" is one of: "add", "sub", "mul", "div".
- "a" and "b" are numbers (can be integers or decimals).

Rules:
- Do NOT explain anything.
- Do NOT add any text before or after the JSON.
- Respond with ONLY valid JSON.

User input: %q
`, userInput)

	reqBody := OllamaCalcRequest{
		Model:  "llama3", // change if you're using a different model
		Prompt: prompt,
		Stream: false, // single JSON response
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama status: %s", resp.Status)
	}

	var oResp OllamaCalcResponse
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return nil, fmt.Errorf("decode Ollama response: %w", err)
	}

	// The "response" field itself should be JSON (as a string), so we parse it
	var parsed ParsedExpression
	trimmed := strings.TrimSpace(oResp.Response)
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal parsed expression JSON from '%s': %w", trimmed, err)
	}

	return &parsed, nil
}

// compute performs the actual arithmetic in Go.
func compute(p *ParsedExpression) (float64, error) {
	op := strings.ToLower(strings.TrimSpace(p.Op))

	switch op {
	case "add":
		return p.A + p.B, nil
	case "sub", "subtract":
		return p.A - p.B, nil
	case "mul", "multiply":
		return p.A * p.B, nil
	case "div", "divide":
		if p.B == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return p.A / p.B, nil
	default:
		return 0, fmt.Errorf("unknown op: %q", p.Op)
	}
}
