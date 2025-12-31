package ai

import (
	"fmt"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

// BuildPrompt creates the structured XML payload for the LLM
func BuildPrompt(files []*gitdiff.File, isOutputStructured bool) string {
	var sb strings.Builder

	if isOutputStructured {
		sb.WriteString("<instruction>\n")
		sb.WriteString("Review the following staged changes for logic errors, security issues, and overall code cleanliness.\n")
		sb.WriteString("Only comment on lines prefixed with '+'.\n")
		sb.WriteString("</instruction>\n\n")
	} else {
		sb.WriteString("<instruction>\n")
		sb.WriteString("Review the following staged changes for logic, security, and 'aura' (cleanliness).\n")
		sb.WriteString("Only comment on lines prefixed with '+'.\n\n")

		sb.WriteString("### OUTPUT FORMAT\n")
		sb.WriteString("For every issue found, use exactly this structure:\n")
		sb.WriteString("ISSUE: [Short title]\n")
		sb.WriteString("FILE: [Filename]\n")
		sb.WriteString("LINE: [Line number]\n")
		sb.WriteString("TYPE: [BUG|SECURITY|STYLE|PERFORMANCE|COMPLEXITY]\n")
		sb.WriteString("DETAIL: [What is wrong]\n")
		sb.WriteString("SUGGESTION: [fix suggestion]\n")
		sb.WriteString("AURA_LOSS: [0-10000]\n")
		sb.WriteString("---\n\n")

		sb.WriteString("### EXAMPLE RESPONSE\n")
		sb.WriteString("ISSUE: Unhandled Nil Pointer\n")
		sb.WriteString("FILE: internal/git/client.go\n")
		sb.WriteString("LINE: 42\n")
		sb.WriteString("TYPE: BUG\n")
		sb.WriteString("DETAIL: The 'resp' variable is used before checking if 'err' is nil.\n")
		sb.WriteString("SUGGESTION: if err != nil { return err }\n")
		sb.WriteString("AURA_LOSS: 800\n")
		sb.WriteString("---\n")
		sb.WriteString("</instruction>\n\n")
	}

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
