package reviewer

func GetAuraSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"summary": map[string]any{"type": "string"},
			"reviews": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"file": map[string]any{"type": "string"},
						"line": map[string]any{"type": "integer"},
						"type": map[string]any{
							"type": "string",
							"enum": []string{"BUG", "SECURITY", "STYLE", "PERFORMANCE", "COMPLEXITY"},
						},
						"detail":     map[string]any{"type": "string"},
						"suggestion": map[string]any{"type": "string"},
						"aura_loss": map[string]any{
							"type":    "integer",
							"minimum": 0,
							"maximum": 10000,
						},
					},
					"required":             []string{"file", "line", "type", "detail", "suggestion", "aura_loss"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []string{"summary", "reviews"},
		"additionalProperties": false,
	}
}
