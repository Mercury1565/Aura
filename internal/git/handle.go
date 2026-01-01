package git

import (
	"fmt"
	"log"
)

// Build a git diff summary string of what the staged changes look like
func BuildGitSummary(contextWidth int) string {
	response := ""

	// Fetch raw diff from Git
	raw, err := GetStagedDiff(contextWidth)
	if err != nil {
		log.Fatalf("❌ Git Error: %v", err)
	}

	// Parse it into structured data
	files, err := ParseRawDiff(raw)
	if err != nil {
		log.Fatalf("❌ Parser Error: %v", err)
	}

	// Append results
	response += fmt.Sprintf("🔍 Aura found %d changed files in your staging area:\n", len(files))
	for _, f := range files {
		status := "modified"
		if f.IsNew {
			status = "new"
		}
		response += fmt.Sprintf("  - 🚀🚀 [%s] %s\n", status, f.NewName)

		// TextFragments directly from the file
		for i, fragment := range f.TextFragments {
			response += fmt.Sprintf("    Hunk #%d: +%d lines, -%d lines (Starts at line %d) \n",
				i+1, fragment.NewLines, fragment.OldLines, fragment.NewPosition+1)

			// append the actual changes
			for _, line := range fragment.Lines {
				response += fmt.Sprintf("      %s", line.String())
			}
		}
	}

	return response
}
