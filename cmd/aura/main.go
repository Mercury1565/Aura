package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Mercury1565/Aura/internal/ai"
	"github.com/Mercury1565/Aura/internal/git"
	"github.com/Mercury1565/Aura/internal/reviewer"
	"github.com/joho/godotenv"
)

func aiTest() {
	ctx := context.Background()

	_ = godotenv.Load()
	modelName := os.Getenv("MODEL_NAME")

	if modelName == "" {
		modelName = "llama-3.3-70b-versatile"
	}

	llm, err := ai.NewGroqClient(modelName)
	if err != nil {
		log.Fatal(err)
	}

	r := reviewer.NewLLMReviewer(llm)

	diff := `
			--- a/main.go
			+++ b/main.go
			@@ -10,5 +10,6 @@ func main() {
			- apiKey := "AIza_Secret_Key_123"
			+ apiKey := os.Getenv("API_KEY")
			+ fmt.Println("Debugging here...")
		`

	feedback, err := r.ReviewDiff(ctx, diff)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- LLM Review ---")
	fmt.Println(feedback)
}

func gitTest() {
	gitDiff := git.BuildGitSummary()
	fmt.Println("--- Git Summary ---")
	fmt.Println(gitDiff)
}

func main() {
	aiTest()
	// gitTest()
}
