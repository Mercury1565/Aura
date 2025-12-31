package ui

import (
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/charmbracelet/lipgloss"
)

func DiffSideBySide(fragment *gitdiff.TextFragment, width int) (string, string) {
	var leftLines, rightLines []string
	colWidth := (width / 2) - 1

	// Ensure we don't have borders/padding on the base styles that add height
	base := lipgloss.NewStyle().Padding(0).Margin(0)
	added := base.Foreground(Color(ColorAdded))
	removed := base.Foreground(Color(ColorRemoved))
	empty := base

	lines := fragment.Lines
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Clean the line content of any existing newlines
		cleanLine := strings.TrimRight(line.Line, "\r\n")

		switch line.Op {
		case gitdiff.OpContext:
			txt := truncate(cleanLine, colWidth)
			leftLines = append(leftLines, txt)
			rightLines = append(rightLines, txt)

		case gitdiff.OpDelete:
			if i+1 < len(lines) && lines[i+1].Op == gitdiff.OpAdd {
				leftLines = append(leftLines, removed.Render(truncate(cleanLine, colWidth)))
				rightLines = append(rightLines, added.Render(truncate(strings.TrimRight(lines[i+1].Line, "\r\n"), colWidth)))
				i++
			} else {
				leftLines = append(leftLines, removed.Render(truncate(cleanLine, colWidth)))
				rightLines = append(rightLines, empty.Render(strings.Repeat(" ", colWidth)))
			}

		case gitdiff.OpAdd:
			leftLines = append(leftLines, empty.Render(strings.Repeat(" ", colWidth)))
			rightLines = append(rightLines, added.Render(truncate(cleanLine, colWidth)))
		}
	}

	return strings.Join(leftLines, "\n"), strings.Join(rightLines, "\n")
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
