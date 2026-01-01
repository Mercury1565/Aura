package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// generates the entire diff as one big string
func (m Model) renderDiffContent() string {
	var doc strings.Builder

	contentWidth := m.TerminalWidth - 2
	columnWidth := 4 * contentWidth / 10
	aiWindowWidth := contentWidth / 5

	// Get feedback once per render
	feedback, err := m.Reviewer.ReviewDiffWithStructuredOutput(m.Ctx, m.DiffFiles)
	if err != nil {
		return "AI Review failed: " + err.Error()
	}

	// AI Summary Box (Displayed once for the entire diff)
	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(Color(ColorLineNumber)).
		Padding(1, 2).
		Width(m.TerminalWidth).
		Render(fmt.Sprintf("AI SUMMARY: %s", feedback.Summary))
	doc.WriteString(summaryBox + "\n\n")

	for _, file := range m.DiffFiles {
		// File Header
		header := lipgloss.NewStyle().
			Foreground(Color(ColorHeaderText)).
			Background(Color(ColorHeaderBackground)).
			Bold(true).
			Width(m.TerminalWidth).
			Render(fmt.Sprintf(" 📂 %s", file.NewName))
		doc.WriteString(header + "\n")

		for _, hunk := range file.TextFragments {
			lNums, lLines, rNums, rLines := DiffSideBySide(hunk, m.TerminalWidth)

			// Filter reviews that belong to THIS hunk
			var hunkFeedback strings.Builder
			for _, rev := range feedback.Reviews {
				line := int64(rev.Line)

				if rev.File == file.NewName &&
					line >= hunk.NewPosition &&
					line <= hunk.NewPosition+hunk.NewLines {

					comment := lipgloss.NewStyle().
						Foreground(Color(ColorAI)).
						Italic(true).
						Width(aiWindowWidth - 4).
						Render(fmt.Sprintf("📍 Line %d: %s\n  🚀%s", rev.Line, rev.Detail, rev.Suggestion))
					hunkFeedback.WriteString(comment + "\n\n")
				}
			}

			// Assemble the windows
			innerWidth := columnWidth - 2
			aiText := hunkFeedback.String()

			leftContent := lipgloss.NewStyle().
				Width(innerWidth).
				Render(
					lipgloss.JoinHorizontal(lipgloss.Top, lNums, lLines),
				)
			rightContent := lipgloss.NewStyle().
				Width(innerWidth).
				Render(
					lipgloss.JoinHorizontal(lipgloss.Top, rNums, rLines),
				)

			leftWindow := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(columnWidth).
				Render(leftContent)

			rightWindow := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(columnWidth).
				Render(rightContent)

			aiWindow := lipgloss.NewStyle().
				Padding(2, 1).
				BorderForeground(Color(ColorLineNumber)).
				Width(aiWindowWidth).
				Render(aiText)

			sideBySide := lipgloss.JoinHorizontal(lipgloss.Top, leftWindow, aiWindow, rightWindow)
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
