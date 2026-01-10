package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Mercury1565/Aura/internal/ai"
	"github.com/Mercury1565/Aura/internal/git"
	"github.com/Mercury1565/Aura/internal/reviewer"
)

func DryMode(cfg *ai.Config) {
	ctx := context.Background()
	modelName := cfg.ModelName

	llm, err := ai.NewGroqClient(modelName, cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := reviewer.NewLLMReviewer(llm)

	// Fetch raw diff from Git
	raw, err := git.GetStagedDiff(3)
	if err != nil {
		log.Fatalf("❌ Git Error: %v", err)
	}

	// Parse it into structured data
	files, err := git.ParseRawDiff(raw)
	if err != nil {
		log.Fatalf("❌ Parser Error: %v", err)
	}

	feedback, err := r.ReviewDiffWithStructuredOutput(ctx, files)
	if err != nil || feedback == nil || len(feedback.Reviews) == 0 {
		// Fallback
		fmt.Println("Falling back to unstructured review...")
		raw, fallbackErr := r.ReviewDiff(ctx, files)
		if fallbackErr != nil {
			log.Fatalf("❌ Fallback Error: %v", fallbackErr)
		}

		feedback = r.ParseUnstructuredReview(raw)
	}

	feedback.PrettyPrint()
}
