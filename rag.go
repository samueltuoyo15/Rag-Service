package main

import (
	"context"
	"fmt"
	"os"
	"log"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"github.com/joho/godotenv"
)


func ProcessFileWithRAG(filePath, question string) (string, error) {
  err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found.")
		}
		
  ctx := context.Background()
	
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", err
	}
	defer client.Close()
 
 	// to read file content 
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fullPrompt := fmt.Sprintf("Here is a document:\n\n%s\n\nQuestion: %s", string(fileContent), question)
	
	model := client.GenerativeModel("gemini-2.0-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprint(resp.Candidates[0].Content.Parts[0]), nil
	}

	return "", fmt.Errorf("no content generated")
}



