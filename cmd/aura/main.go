package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Mercury1565/Aura/internal/ai"
)

func main() {
	ctx := context.Background()
	modelName := os.Getenv("MODEL_NAME")

	if modelName == "" {
		modelName = "llama-3.3-70b-versatile"
	}

	llm, err := ai.NewGroqClient(modelName)
	if err != nil {
		panic(err)
	}

	stream, err := llm.Chat(ctx, ai.ChatRequest{
		Messages: []ai.Message{
			{Role: "user", Content: "Explain goroutines"},
		},
		Temperature: 0.5,
		Stream:      true,
	})
	if err != nil {
		panic(err)
	}

	for chunk := range stream {
		if chunk.Err != nil {
			panic(chunk.Err)
		}
		if chunk.Done {
			break
		}
		fmt.Print(chunk.Content)
	}

}
