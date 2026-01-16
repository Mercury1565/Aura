package ai

import (
	"fmt"
	"strings"

	"github.com/Mercury1565/Aura/internal/utils"
	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

func BuildPrompt(files []*gitdiff.File, isOutputStructured bool) string {
	var sb strings.Builder

	if isOutputStructured {
		sb.WriteString("<instruction>\n")
		sb.WriteString(utils.BaseInstruction)
		sb.WriteString("\n")
		sb.WriteString("</instruction>\n\n")
	} else {
		sb.WriteString("<instruction>\n")
		sb.WriteString(utils.BaseInstruction)
		sb.WriteString("\n\n")
		sb.WriteString(`### OUTPUT FORMAT
For every issue found, use exactly this structure:
ISSUE: [Short title]
FILE: [Filename]
LINE: [Line number]
TYPE: [BUG|SECURITY|STYLE|PERFORMANCE|COMPLEXITY]
DETAIL: [What is wrong]
SUGGESTION: [fix suggestion]
AURA_LOSS: [0-10000]

### EXAMPLE RESPONSE
ISSUE: Unhandled Nil Pointer
FILE: internal/git/client.go
LINE: 42
TYPE: BUG
DETAIL: The 'resp' variable is used before checking if 'err' is nil.
SUGGESTION: if err != nil { return err }
AURA_LOSS: 800`)
		sb.WriteString("</instruction>")
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
