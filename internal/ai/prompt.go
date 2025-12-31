package ai

import (
	"fmt"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

const TestPrompt = "Review this git diff and find any glaring mistakes:\n\n%s"

// BuildPrompt creates the structured XML payload for the LLM
func BuildPrompt(files []*gitdiff.File) string {
	var sb strings.Builder

	sb.WriteString("<instruction>\n")
	sb.WriteString("Review the following staged changes for logic errors, security issues, and overall code cleanliness.\n")
	sb.WriteString("IMPORTANT: Return your response ONLY as a JSON object in the following format:\n")
	sb.WriteString(`{"summary": "string", "reviews": [{"file": "string", "line": int, "type": "BUG|AURA|STYLE", "message": "string"}]}` + "\n")
	sb.WriteString("Only comment on lines prefixed with '+'.\n")
	sb.WriteString("</instruction>\n\n")

	for _, file := range files {
		// Skip binary files or files with no text changes
		if file.IsBinary || len(file.TextFragments) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("<file name=\"%s\">\n", file.NewName))

		for _, fragment := range file.TextFragments {
			sb.WriteString(fmt.Sprintf("  <hunk start_line=\"%d\">\n", fragment.NewPosition))

			for _, line := range fragment.Lines {
				// line.String() preserves the ' ', '+', or '-' prefix
				sb.WriteString("    " + line.String())
			}

			sb.WriteString("  </hunk>\n")
		}

		sb.WriteString("</file>\n")
	}

	return sb.String()
}
