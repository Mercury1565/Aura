package main

import (
	_ "embed"
	"flag"
	"fmt"
)

//go:embed help.txt
var helpContent string

func main() {
	dryMode := flag.Bool("d", false, "dry mode")
	flag.Usage = func() {
		fmt.Print(helpContent)
	}

	flag.Parse()

	if *dryMode {
		DryMode()
	} else {
		DefaultMode()
	}
}
