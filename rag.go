package main

import (
	"context"
	"fmt"
	"os"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)


func ProcessFileWithRAG(filePath, question string) (string, error) {
  ctx := context.Background()
	
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(question))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprint(resp.Candidates[0].Content.Parts[0]), nil
	}

	return "", fmt.Errorf("no content generated")
}