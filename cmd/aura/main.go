package main

import (
	"context"
	_ "embed"
	"flag"
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

//go:embed help.txt
var helpContent string

func detachedModeTest() {
	ctx := context.Background()

	_ = godotenv.Load()
	modelName := os.Getenv("MODEL_NAME")

	llm, err := ai.NewGroqClient(modelName)
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

func gitTest() {
	gitDiff := git.BuildGitSummary(3)
	fmt.Println("--- Git Summary ---")
	fmt.Println(gitDiff)
}

func UITest() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Get the staged diff
	raw, err := git.GetStagedDiff(5)
	if err != nil {
		log.Fatalf("❌ Git Error: %v", err)
	}

	// 2. Parse it into structs
	files, err := git.ParseRawDiff(raw)
	if err != nil {
		log.Fatalf("❌ Parser Error: %v", err)
	}

	_ = godotenv.Load()
	modelName := os.Getenv("MODEL_NAME")

	llm, err := ai.NewGroqClient(modelName)
	if err != nil {
		log.Fatal(err)
	}

	r := reviewer.NewLLMReviewer(llm)

	// 3. Initialize the UI Model with the files
	m := ui.InitialModel(files, r, ctx)

	// 4. Start the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	dryMode := flag.Bool("d", false, "dry mode")
	flag.Usage = func() {
		fmt.Print(helpContent)
	}

	flag.Parse()

	if *dryMode {
		detachedModeTest()
	} else {
		UITest()
	}
}
