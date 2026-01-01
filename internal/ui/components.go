package ui

import (
	"fmt"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/charmbracelet/lipgloss"
)

func DiffSideBySide(fragment *gitdiff.TextFragment, width int) (string, string, string, string) {
	var lNums, lLines, rNums, rLines []string

	// Reserve 4 chars for line numbers + 1 for padding
	numWidth := 5
	colWidth := (width / 2) - numWidth - 1

	styleNum := lipgloss.NewStyle().
		Foreground(Color(ColorLineNumber))

	added := lipgloss.NewStyle().
		Foreground(Color(ColorAdded)).
		Background(Color(ColorAddedBG))

	removed := lipgloss.NewStyle().
		Foreground(Color(ColorRemoved)).
		Background(Color(ColorRemovedBG))

	// Initialize counters from the fragment header
	curLeft := fragment.OldPosition
	curRight := fragment.NewPosition

	lines := fragment.Lines
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		cleanLine := strings.TrimRight(line.Line, "\r\n")

		switch line.Op {
		case gitdiff.OpContext:
			lNums = append(lNums, styleNum.Render(fmt.Sprintf("%4d ", curLeft)))
			lLines = append(lLines, truncate(cleanLine, colWidth))
			rNums = append(rNums, styleNum.Render(fmt.Sprintf("%4d ", curRight)))
			rLines = append(rLines, truncate(cleanLine, colWidth))
			curLeft++
			curRight++

		case gitdiff.OpDelete:
			// Left side gets a number and red text
			lNums = append(lNums, styleNum.Render(fmt.Sprintf("%4d ", curLeft)))
			lLines = append(lLines, removed.Render(truncate(cleanLine, colWidth)))
			curLeft++

			// Check if this is a "Change" (Delete followed by Add)
			if i+1 < len(lines) && lines[i+1].Op == gitdiff.OpAdd {
				nextClean := strings.TrimRight(lines[i+1].Line, "\r\n")
				rNums = append(rNums, styleNum.Render(fmt.Sprintf("%4d ", curRight)))
				rLines = append(rLines, added.Render(truncate(nextClean, colWidth)))
				curRight++
				i++ // Skip the next line
			} else {
				// Pure deletion: right side is empty/ghosted
				rNums = append(rNums, strings.Repeat(" ", numWidth))
				rLines = append(rLines, strings.Repeat(" ", colWidth))
			}

		case gitdiff.OpAdd:
			// Pure addition: left side is empty/ghosted
			lNums = append(lNums, strings.Repeat(" ", numWidth))
			lLines = append(lLines, strings.Repeat(" ", colWidth))

			rNums = append(rNums, styleNum.Render(fmt.Sprintf("%4d ", curRight)))
			rLines = append(rLines, added.Render(truncate(cleanLine, colWidth)))
			curRight++
		}
	}

	return strings.Join(lNums, "\n"), strings.Join(lLines, "\n"),
		strings.Join(rNums, "\n"), strings.Join(rLines, "\n")
}

func truncate(s string, w int) string {
	s = strings.ReplaceAll(s, "\t", "    ")

	// If the window is incredibly small, just return an empty string or dots
	if w <= 0 {
		return ""
	}
	if w <= 3 {
		return strings.Repeat(".", w)
	}

	if len(s) > w {
		return s[:w-3] + "..."
	}

	// Pad the string to the full column width to keep backgrounds consistent
	return s + strings.Repeat(" ", w-len(s))
}
