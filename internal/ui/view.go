package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// generates the entire diff as one big string
func (m Model) renderDiffContent() string {
	var doc strings.Builder

	contentWidth := m.TerminalWidth - 1
	leftWidth := contentWidth / 2
	rightWidth := contentWidth - leftWidth

	for _, file := range m.DiffFiles {
		header := lipgloss.NewStyle().
			Foreground(Color(ColorHeaderText)).
			Bold(true).
			Render(fmt.Sprintf(
				"\n 📂 %s %s",
				file.NewName,
				strings.Repeat(" ", contentWidth),
			))

		doc.WriteString(header + "\n")

		for _, hunk := range file.TextFragments {
			left, right := DiffSideBySide(hunk, m.TerminalWidth)

			leftWindow := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(leftWidth).Render(left)

			rightWindow := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(rightWidth).Render(right)

			sideBySide := lipgloss.JoinHorizontal(
				lipgloss.Top,
				leftWindow,
				rightWindow,
			)
			doc.WriteString(sideBySide + "\n")
		}
	}
	return doc.String()
}

func (m Model) View() string {
	if !m.Ready {
		return "Initializing Aura..."
	}

	// displays the viewport's current "window"
	return fmt.Sprintf("%s\n%s",
		m.headerView(),
		m.Viewport.View(),
	)
}

func (m Model) headerView() string {
	return lipgloss.NewStyle().Bold(true).Render(" AURA GIT DIFF VIEW (q to quit)")
}
