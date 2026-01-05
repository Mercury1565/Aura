package reviewer

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/Mercury1565/Aura/internal/ai"
	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

type LLMReviewer struct {
	llm ai.LLMClient
}

func NewLLMReviewer(llm ai.LLMClient) *LLMReviewer {
	return &LLMReviewer{llm: llm}
}

// takes a git diff and returns the LLM's feedback
func (r *LLMReviewer) ReviewDiff(ctx context.Context, files []*gitdiff.File) (string, error) {
	prompt := ai.BuildPrompt(files, false)

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

func (r *LLMReviewer) ReviewDiffWithStructuredOutput(ctx context.Context, files []*gitdiff.File) (*CodeReview, error) {
	prompt := ai.BuildPrompt(files, true)

	req := ai.ChatRequest{
		Messages: []ai.Message{{Role: "user", Content: prompt}},
		ResponseFormat: &ai.ResponseFormat{
			Type: "json_schema",
			JSONSchema: &ai.JSONSchema{
				Name:   "code_review",
				Strict: true,
				Schema: GetAuraSchema(),
			},
		},
	}

	jsonResponse, err := r.llm.ChatStructured(ctx, req)
	if err != nil {
		return nil, err
	}

	var review CodeReview
	if err := json.Unmarshal([]byte(jsonResponse), &review); err != nil {
		return nil, fmt.Errorf("invalid structured LLM response: %w", err)
	}

	// Sort reviews by AuraLoss descending (highest first)
	slices.SortFunc(review.Reviews, func(a, b ReviewItem) int {
		return b.AuraLoss - a.AuraLoss
	})

	return &review, nil
}

func (r *LLMReviewer) ParseUnstructuredReview(input string) *CodeReview {
	var feedback CodeReview

	// Split by the horizontal rule
	for block := range strings.SplitSeq(input, "---") {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		item := ReviewItem{}

		for line := range strings.SplitSeq(block, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Simple key-value split by the first ":"
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			switch key {
			case "FILE":
				item.File = val
			case "LINE":
				fmt.Sscanf(val, "%d", &item.Line)
			case "TYPE":
				item.Type = val
			case "ISSUE":
				item.Issue = val
			case "SUGGESTION":
				item.Suggestion = val
			case "AURA_LOSS":
				fmt.Sscanf(val, "%d", &item.AuraLoss)
			}
		}

		if item.Issue != "" {
			feedback.Reviews = append(feedback.Reviews, item)
		}
	}

	feedback.Summary = fmt.Sprintf("Found %d issues via fallback review.", len(feedback.Reviews))
	return &feedback
}
