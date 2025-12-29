package ai

import (
	"context"
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
