package reviewer

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mercury1565/Aura/internal/ai"
)

type LLMReviewer struct {
	llm ai.LLMClient
}

func NewLLMReviewer(llm ai.LLMClient) *LLMReviewer {
	return &LLMReviewer{llm: llm}
}

// takes a git diff and returns the LLM's feedback
func (r *LLMReviewer) ReviewDiff(ctx context.Context, diff string) (string, error) {
	prompt := fmt.Sprintf(
		ai.TestPrompt,
		diff,
	)

	req := ai.ChatRequest{
		Messages: []ai.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.2, // low temperature for deterministic analysis
		Stream:      false,
	}

	stream, err := r.llm.Chat(ctx, req)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for chunk := range stream {
		if chunk.Err != nil {
			return "", chunk.Err
		}
		builder.WriteString(chunk.Content)
		if chunk.Done {
			break
		}
	}

	return builder.String(), nil
}
