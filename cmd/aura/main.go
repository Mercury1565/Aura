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
var dryMode bool
var showPath bool
var showModel bool

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

	flag.Parse()
}

func main() {
	// Load Config
	cfg, err := ai.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Handle config set
	if len(os.Args) > 3 && os.Args[1] == "config" {
		cfg.HandleConfigSet(os.Args)
		return
	}

	// flags init
	flagsInit()

	if showPath {
		fmt.Printf("Config file used: %s\n", viper.ConfigFileUsed())
		return
	} else if showModel {
		if cfg.ModelName == "" {
			fmt.Println("Model name not set")
		} else {
			fmt.Printf("Model name: %s\n", cfg.ModelName)
		}
		return
	}

	// Validate existence of config variables
	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Core App features
	if dryMode {
		DryMode(cfg)
	} else {
		DefaultMode(cfg)
	}
}
