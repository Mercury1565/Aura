package main

// Hello there

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Mercury1565/Aura/internal/ai"
	"github.com/Mercury1565/Aura/internal/git"
	"github.com/Mercury1565/Aura/internal/reviewer"
	"github.com/Mercury1565/Aura/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func aiTest() {
	ctx := context.Background()

	_ = godotenv.Load()
	modelName := os.Getenv("MODEL_NAME")

	llm, err := ai.NewGroqClient(modelName)
	if err != nil {
		log.Fatal(err)
	}

	r := reviewer.NewLLMReviewer(llm)

	// diff := `
	// 		--- a/main.go
	// 		+++ b/main.go
	// 		@@ -10,5 +10,6 @@ func main() {
	// 		- apiKey := "AIza_Secret_Key_123"
	// 		+ apiKey := os.Getenv("API_KEY")
	// 		+ fmt.Println("Debugging here...")
	// 	`

	// Fetch raw diff from Git
	raw, err := git.GetStagedDiff()
	if err != nil {
		log.Fatalf("❌ Git Error: %v", err)
	}

	// Parse it into structured data
	files, err := git.ParseRawDiff(raw)
	if err != nil {
		log.Fatalf("❌ Parser Error: %v", err)
	}

	// feedback, err := r.ReviewDiff(ctx, files)
	feedback, err := r.ReviewDiffWithStructuredOutput(ctx, files)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🚀🚀🚀--- LLM Review ---🚀🚀🚀")
	fmt.Println(feedback)
}

func gitTest() {
	gitDiff := git.BuildGitSummary()
	fmt.Println("--- Git Summary ---")
	fmt.Println(gitDiff)
}

func UITest() {
	// 1. Get the staged diff
	raw, err := git.GetStagedDiff()
	if err != nil {
		log.Fatalf("❌ Git Error: %v", err)
	}

	// 2. Parse it into structs
	files, err := git.ParseRawDiff(raw)
	if err != nil {
		log.Fatalf("❌ Parser Error: %v", err)
	}

	// 3. Initialize the UI Model with the files
	m := ui.InitialModel(files)

	// 4. Start the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen()) // WithAltScreen is the "IDE" feel
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	// aiTest()
	// gitTest()
	UITest()
}
