package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/Mercury1565/Aura/internal/ai"
	"github.com/spf13/viper"
)

//go:embed __help__.txt
var helpContent string

// custom flags
var (
    dryMode   bool
    showPath  bool
    showModel bool
    reviewAll bool
    contextLines int
)

func flagsInit() {
	flag.Usage = func() {
		fmt.Print(helpContent)
	}

	flag.BoolVar(&dryMode, "d", false, "dry mode")
	flag.BoolVar(&dryMode, "dry", false, "dry mode")

	flag.BoolVar(&showPath, "w", false, "show config file path")
	flag.BoolVar(&showPath, "where", false, "show config file path")

	flag.BoolVar(&showModel, "m", false, "show model name")
	flag.BoolVar(&showModel, "model", false, "show model name")

	flag.BoolVar(&reviewAll, "a", true, "review all changes (including unstaged")
	flag.BoolVar(&reviewAll, "all", true, "review all changes (including unstaged")

	flag.IntVar(&contextLines, "c", 3, "number of context lines for the diff")
    flag.IntVar(&contextLines, "context", 3, "number of context lines for the diff")

	flag.Parse()
}

func main() {
	// Load Config
	cfg, err := ai.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Intercept config command early
    if len(os.Args) > 1 && os.Args[1] == "config" {
        if len(os.Args) > 3 {
            cfg.HandleConfigSet(os.Args)
        } else {
            fmt.Println("Usage: aura config <key> <value>")
        }
        return
    }

    flagsInit()

    // Handle information-only flags
    if showPath {
        fmt.Printf("Config file used: %s\n", viper.ConfigFileUsed())
        return
    }
    if showModel {
        fmt.Printf("Model name: %s\n", cfg.ModelName)
        return
    }

    if err := cfg.Validate(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // We choose the function once, then call it with the same params
    if dryMode {
        DryMode(cfg, contextLines, reviewAll)
    } else {
        DefaultMode(cfg, contextLines, reviewAll)
    }
}
