package ai

import "context"

type Message struct {
	Role    string
	Content string
}

type ResponseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

type JSONSchema struct {
	Name   string         `json:"name"`
	Strict bool           `json:"strict"`
	Schema map[string]any `json:"schema"`
}

type ChatRequest struct {
	Messages       []Message
	Temperature    float32
	Stream         bool
	ResponseFormat *ResponseFormat // Optional: nil for regular text
}

type ChatChunk struct {
	Content string
	Done    bool
	Err     error
}

type LLMClient interface {
	Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
	ChatStructured(ctx context.Context, req ChatRequest) (string, error)
	Model() string
}

type SchemaWrapper struct {
	Bytes []byte
}

func (s SchemaWrapper) MarshalJSON() ([]byte, error) {
	return s.Bytes, nil
}
