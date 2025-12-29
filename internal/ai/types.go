package ai

import "context"

type Message struct {
	Role    string
	Content string
}

type ChatRequest struct {
	Messages    []Message
	Temperature float32
	Stream      bool
}

type ChatChunk struct {
	Content string
	Done    bool
	Err     error
}

type LLMClient interface {
	Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
	Model() string
}
