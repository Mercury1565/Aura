package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type GroqClient struct {
	client *openai.Client
	model  string
}

func NewGroqClient(model string) (*GroqClient, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY not set")
	}

	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = "https://api.groq.com/openai/v1"

	client := openai.NewClientWithConfig(cfg)

	return &GroqClient{
		client: client,
		model:  model,
	}, nil
}

func (g *GroqClient) Model() string {
	return g.model
}

func (g *GroqClient) Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error) {
	out := make(chan ChatChunk)

	go func() {
		defer close(out)

		msgs := make([]openai.ChatCompletionMessage, 0, len(req.Messages))
		for _, m := range req.Messages {
			msgs = append(msgs, openai.ChatCompletionMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}

		stream, err := g.client.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:       g.model,
				Messages:    msgs,
				Temperature: float32(req.Temperature),
				Stream:      true,
			},
		)
		if err != nil {
			out <- ChatChunk{Err: err, Done: true}
			return
		}
		defer stream.Close()

		for {
			resp, err := stream.Recv()
			if err != nil {
				out <- ChatChunk{Done: true}
				return
			}

			for _, choice := range resp.Choices {
				if choice.Delta.Content != "" {
					out <- ChatChunk{
						Content: choice.Delta.Content,
					}
				}
			}
		}
	}()

	return out, nil
}

func (g *GroqClient) ChatStructured(ctx context.Context, req ChatRequest) (string, error) {
	msgs := make([]openai.ChatCompletionMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	var respFormat *openai.ChatCompletionResponseFormat

	// SAFETY CHECK: proceed if ResponseFormat AND JSONSchema are provided
	if req.ResponseFormat != nil && req.ResponseFormat.JSONSchema != nil {
		schemaBytes, err := json.Marshal(req.ResponseFormat.JSONSchema.Schema)
		if err != nil {
			return "", fmt.Errorf("failed to marshal schema: %w", err)
		}

		// Initialize the struct pointer before assigning to its fields
		respFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatType(req.ResponseFormat.Type),
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   req.ResponseFormat.JSONSchema.Name,
				Strict: req.ResponseFormat.JSONSchema.Strict,
				Schema: SchemaWrapper{Bytes: schemaBytes},
			},
		}
	}

	resp, err := g.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:          g.model,
		Messages:       msgs,
		Temperature:    req.Temperature,
		ResponseFormat: respFormat, // nil (default) or our structured format
	})

	if err != nil {
		return "", fmt.Errorf("groq completion error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from groq")
	}

	return resp.Choices[0].Message.Content, nil
}
