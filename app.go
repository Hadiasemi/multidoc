package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	claude "github.com/liushuangls/go-anthropic/v2"
	openai "github.com/sashabaranov/go-openai"
	gemini "google.golang.org/genai"
)

// API keys loaded from environment variables
var (
	openAIKey string
	geminiKey string
	claudeKey string
	configDir string
)

func init() {
	// Set up config directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}
	
	configDir = filepath.Join(homeDir, ".config", "multidoc")
	
	// Check if config directory exists, if not create it
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			log.Fatalf("Error creating config directory: %v", err)
		}
		
		// Create default .env file in the config directory
		envPath := filepath.Join(configDir, ".env")
		envContent := `OPENAI_API_KEY=your_openai_api_key
GEMINI_API_KEY=your_gemini_api_key
CLAUDE_API_KEY=your_claude_api_key`
		
		err = os.WriteFile(envPath, []byte(envContent), 0644)
		if err != nil {
			log.Fatalf("Error creating default .env file: %v", err)
		}
		
		fmt.Printf("Created default config at: %s\nPlease add your API keys to this file.\n", envPath)
		os.Exit(0) // Exit after creating the config file and showing the message
	}
	
	// Load .env file from the config directory
	envPath := filepath.Join(configDir, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: Error loading .env file from %s: %v", envPath, err)
		
		// Try loading from current directory as fallback
		err = godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	// Load API keys from environment variables
	openAIKey = os.Getenv("OPENAI_API_KEY")
	geminiKey = os.Getenv("GEMINI_API_KEY")
	claudeKey = os.Getenv("CLAUDE_API_KEY")

	// Validate that required API keys are set
	if openAIKey == "" || geminiKey == "" || claudeKey == "" {
		log.Fatalf("Missing required API keys in .env file. Please set OPENAI_API_KEY, GEMINI_API_KEY, and CLAUDE_API_KEY in %s", envPath)
	}
}

func main() {
	startTime := time.Now()
	
	// Step 1: Read input from stdin
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}
	input := string(inputBytes)
	
	// Validate input
	if strings.TrimSpace(input) == "" {
		fmt.Fprintf(os.Stderr, "Error: Input cannot be empty\n")
		os.Exit(1)
	}

	// Step 2: Define the system prompt for all models
	systemPrompt := "Process the following input:"

	// Step 3: List of models to query
	modelNames := []string{"gpt-4o-2024-08-06", "o3-mini", "o1-mini", "gemini-2.0-flash", "claude-3-7-sonnet-20250219"}

	fmt.Printf("Processing input with %d different AI models...\n", len(modelNames))

	// Step 4: Set up concurrency with WaitGroup and a slice for outputs
	var wg sync.WaitGroup
	wg.Add(len(modelNames))
	allOutputs := make([]string, len(modelNames))
	modelTimes := make([]time.Duration, len(modelNames))
	var mu sync.Mutex

	// Step 5: Launch goroutines to call each model's API concurrently
	for i, modelName := range modelNames {
		go func(i int, modelName string) {
			defer wg.Done()
			
			modelStart := time.Now()
			
			// Combine system prompt with user input
			fullPrompt := systemPrompt + " " + input
			key := getKey(modelName)
			
			// Call the appropriate API based on the model
			output, err := callAPI(modelName, fullPrompt, key)
			
			modelDuration := time.Since(modelStart)
			
			mu.Lock()
			modelTimes[i] = modelDuration
			if err != nil {
				output = fmt.Sprintf("Error from %s: %v", modelName, err)
				fmt.Fprintf(os.Stderr, "Error from %s: %v\n", modelName, err)
			} else {
				fmt.Printf("Received response from %s (%.2fs)\n", modelName, modelDuration.Seconds())
			}
			
			// Safely store the output in the slice
			allOutputs[i] = output
			mu.Unlock()
		}(i, modelName)
	}

	// Step 6: Wait for all goroutines to complete
	wg.Wait()

	fmt.Printf("All models responded. Generating summary...\n")

	// Step 7: Combine all outputs into a new prompt for gpt-o1
	combinedPrompt := "Summarize the following outputs from different AI models:\n"
	for i, output := range allOutputs {
		combinedPrompt += fmt.Sprintf("Output from %s (%.2fs): %s\n\n", 
			modelNames[i], modelTimes[i].Seconds(), output)
	}

	// Step 8: Call gpt-o1 with the combined prompt
	summaryStart := time.Now()
	finalOutput, err := callOpenAIAPI("o1-mini", combinedPrompt, openAIKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting final output: %v\n", err)
		os.Exit(1)
	}
	summaryDuration := time.Since(summaryStart)

	// Step 9: Display the final output
	fmt.Printf("\n--- Summary (%.2fs) ---\n", summaryDuration.Seconds())
	fmt.Println(finalOutput)
	
	// Display total execution time
	totalDuration := time.Since(startTime)
	fmt.Printf("\nTotal execution time: %.2f seconds\n", totalDuration.Seconds())
}

// getKey returns the appropriate API key based on the model
func getKey(model string) string {
	switch {
	case model == "gemini-2.0-flash":
		return geminiKey
	case model == "claude-3-7-sonnet-20250219":
		return claudeKey
	case model == "gpt-4o-2024-08-06" || model == "o3-mini" || model == "o1-mini":
		return openAIKey
	default:
		return openAIKey
	}
}

// callAPI routes the request to the correct API based on the model
func callAPI(model, prompt, key string) (string, error) {
	switch {
	case model == "gemini-2.0-flash":
		return callGeminiAPI(model, prompt, key)
	case model == "claude-3-7-sonnet-20250219":
		return callClaudeAPI(prompt, key)
	case model == "gpt-4o-2024-08-06" || model == "o3-mini" || model == "o1-mini":
		return callOpenAIAPI(model, prompt, key)
	default:
		return "", fmt.Errorf("unsupported model: %s", model)
	}
}

// callOpenAIAPI handles API calls to OpenAI models
func callOpenAIAPI(model, prompt, key string) (string, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("error calling OpenAI %s: %v", model, err)
	}
	
	if len(resp.Choices) == 0 {
		return "", errors.New("no response choices returned from OpenAI")
	}
	
	return resp.Choices[0].Message.Content, nil
}

// callGeminiAPI is a placeholder for the Gemini API call
func callGeminiAPI(model, prompt, key string) (string, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Split the prompt into system and user parts (Gemini doesn't support system prompts directly)
	parts := strings.SplitN(prompt, " ", 4) // "Process the following input: actual-input"
	systemPrompt := strings.Join(parts[:3], " ")
	userPrompt := ""
	if len(parts) > 3 {
		userPrompt = parts[3]
	}
	
	// Combine system and user prompts for Gemini
	geminiPrompt := systemPrompt
	if userPrompt != "" {
		geminiPrompt = fmt.Sprintf("%s\n\n%s", systemPrompt, userPrompt)
	}
	
	client, err := gemini.NewClient(ctx, &gemini.ClientConfig{
		APIKey:   key,
		Backend:  gemini.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("error creating Gemini client: %v", err)
	}
	
	// Generate content without config to avoid type issues
	result, err := client.Models.GenerateContent(ctx, model, gemini.Text(geminiPrompt), nil)
	if err != nil {
		return "", fmt.Errorf("error generating content: %v", err)
	}
	
	if result == nil {
		return "", errors.New("received nil result from Gemini API")
	}
	
	resultText := result.Text()
	
	return resultText, nil
}

// callClaudeAPI is a placeholder for the Claude API call
func callClaudeAPI(prompt, key string) (string, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Split the prompt into system and user parts
	parts := strings.SplitN(prompt, " ", 4) // "Process the following input: actual-input"
	systemPrompt := strings.Join(parts[:3], " ")
	userPrompt := ""
	if len(parts) > 3 {
		userPrompt = parts[3]
	}

	client := claude.NewClient(key)
	
	// Create pointer to user message for Claude API
	userPromptPtr := userPrompt
	
	// Create request with system and user content
	req := claude.MessagesRequest{
		Model: "claude-3-7-sonnet-20250219",
		Messages: []claude.Message{
			claude.NewUserTextMessage(userPromptPtr),
		},
		MaxTokens: 1000,
	}
	
	// Add system message if needed
	if systemPrompt != "" {
		req.System = systemPrompt
	}
	
	resp, err := client.CreateMessages(ctx, req)
	
	if err != nil {
		var e *claude.APIError
		if errors.As(err, &e) {
			return "", fmt.Errorf("claude API error, type: %s, message: %s", e.Type, e.Message)
		}
		return "", fmt.Errorf("claude API error: %v", err)
	}
	
	if len(resp.Content) == 0 {
		return "", errors.New("no response from Claude API")
	}
	
	return resp.Content[0].GetText(), nil
}