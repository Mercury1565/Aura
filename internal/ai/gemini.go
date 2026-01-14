package ai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(ctx context.Context, model string, c *Config) (*GeminiClient, error) {
	apiKey := "vbghui98ytgfvbnj8"
	// apiKey := c.GeminiAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiClient) Model() string {
	return g.model
}

func (g *GeminiClient) Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error) {
	out := make(chan ChatChunk)

	go func() {
		defer close(out)

		var prompt string
		for _, m := range req.Messages {
			prompt += fmt.Sprintf("%s: %s\n", m.Role, m.Content)
		}

		resp, err := g.client.Models.GenerateContent(
			ctx,
			g.model,
			genai.Text(prompt),
			nil,
		)
		if err != nil {
			out <- ChatChunk{Err: err, Done: true}
			return
		}

		if len(resp.Candidates) == 0 {
			out <- ChatChunk{Content: "", Done: true}
			return
		}

		for _, part := range resp.Candidates[0].Content.Parts {
			out <- ChatChunk{
				Content: fmt.Sprintf("%v", part),
			}
		}

		out <- ChatChunk{Done: true}
	}()

	return out, nil
}
