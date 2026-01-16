package ui

import (
	"fmt"
	"strings"

	"github.com/Mercury1565/Aura/internal/utils"
	"github.com/charmbracelet/lipgloss"
)

// generates the entire diff as one big string
func (m Model) renderDiffContent() string {
	if len(m.Error) > 0 {
		return m.renderError()
	}

	if m.ReviewData == nil {
		return m.renderLoading()
	}

	var doc strings.Builder
	feedback := m.ReviewData

	contentWidth := m.TerminalWidth - 2
	columnWidth := 4 * contentWidth / 10
	aiWindowWidth := contentWidth / 5

	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, true, false).
		BorderForeground(Color(ColorLineNumber)).
		Foreground(Color(ColorHeaderText)).
		Padding(1, 1).
		Align(lipgloss.Center).
		Width(m.TerminalWidth).
		Render(feedback.Summary)
	doc.WriteString(summaryBox + "\n")

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

			// Filter reviews that belong to this hunk
			var hunkFeedback strings.Builder
			for _, rev := range feedback.Reviews {
				line := int64(rev.Line)

				if rev.File == file.NewName &&
					line >= hunk.NewPosition &&
					line <= hunk.NewPosition+hunk.NewLines {

					var b strings.Builder

					labelAuraStyle := lipgloss.NewStyle().
						Foreground(Color(ColorIssue)).
						Bold(true).
						Align(lipgloss.Center).
						Width(aiWindowWidth - 4)

					labelTypeStyle := lipgloss.NewStyle().Foreground(Color(ColorAIResponseTag))
					labelIssueStyle := lipgloss.NewStyle().Foreground(Color(ColorAIResponseTag))
					labelSuggestionStyle := lipgloss.NewStyle().Foreground(Color(ColorAIResponseTag))
					typeValueStyle := lipgloss.NewStyle().Foreground(Color(ColorHeaderText))
					valueStyle := lipgloss.NewStyle().Foreground(Color(ColorAI)).Italic(true)

					b.WriteString(labelAuraStyle.Render(fmt.Sprintf("-%d %s", rev.AuraLoss, "AURA LOSS")))

					b.WriteString(labelTypeStyle.Render("\nLine: "))
					b.WriteString(valueStyle.Render(fmt.Sprintf("%d", rev.Line)))

					b.WriteString(labelTypeStyle.Render("\nType: "))
					b.WriteString(typeValueStyle.Render(fmt.Sprintf("%s", rev.Type)))

					b.WriteString(labelIssueStyle.Render("\nIssue: "))
					b.WriteString(valueStyle.Render(fmt.Sprintf("%s", rev.Issue)))

					b.WriteString(labelSuggestionStyle.Render("\nSuggestion: "))
					b.WriteString(valueStyle.Render(fmt.Sprintf("%s", rev.Suggestion)))

					// 3. Create the divider line
					divider := lipgloss.NewStyle().
						Foreground(Color(ColorLineNumber)).
						Faint(true).
						Render(strings.Repeat("─", aiWindowWidth-4))

					comment := lipgloss.NewStyle().
						Width(aiWindowWidth - 4).
						Render(b.String())

					hunkFeedback.WriteString(comment + divider + "\n")
				}
			}

			innerWidth := columnWidth - 2
			aiText := hunkFeedback.String()

			// Assemble the windows
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

			// determine the actual height of the code hunk
			leftHeight := lipgloss.Height(leftWindow)
			rightHeight := lipgloss.Height(rightWindow)
			maxCodeHeight := max(leftHeight, rightHeight)

			// prepare AI text with internal scrolling
			aiLines := strings.Split(aiText, "\n")
			displayLimit := maxCodeHeight - 2

			var finalAIText string
			if len(aiLines) > displayLimit {
				// Use the global offset
				start := m.AIScrollOffset
				if start > len(aiLines)-displayLimit {
					start = len(aiLines) - displayLimit
				}
				if start < 0 {
					start = 0
				}

				end := start + displayLimit
				finalAIText = strings.Join(aiLines[start:end], "\n")
				finalAIText += "\n" + lipgloss.NewStyle().Foreground(Color(ColorAIScroll)).Render("↓ more...")
			} else {
				finalAIText = aiText
			}

			// Force AI Window to match code height exactly
			aiWindow := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Color(ColorLineNumber)).
				Width(aiWindowWidth).
				Height(maxCodeHeight).
				Render(finalAIText)

			sideBySide := lipgloss.JoinHorizontal(lipgloss.Top, leftWindow, aiWindow, rightWindow)
			doc.WriteString(sideBySide + "\n")
		}
	}
	return doc.String()
}

func (m Model) renderError() string {
	var b strings.Builder

	for _, err := range m.Error {
		if err == nil {
			continue
		}
		b.WriteString("• ")
		b.WriteString(err.Error())
		b.WriteString("\n")
	}

	errorText := strings.TrimSpace(b.String())
	if errorText == "" {
		errorText = "Unknown error occurred."
	}

	errorBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Color(ColorRemoved)).
		Foreground(Color(ColorRemoved)).
		Padding(1, 2).
		Width(m.TerminalWidth).
		Render(
			"☠️ ☠️ AI FAILURE ON BOTH STRUCTURED AND UNSTRUCTURED OUTPUT ☠️ ☠️\n\n" + errorText,
		)

	return "\n" + errorBox
}

func (m Model) renderLoading() string {
	auraLogo := utils.AuraLogo
	colorIndex := m.Frame % len(logoShimmerColors)

	logo := lipgloss.NewStyle().
		Foreground(lipgloss.Color(logoShimmerColors[colorIndex])).
		Bold(true).
		Render(auraLogo)

	tagline := lipgloss.NewStyle().
		Foreground(Color(ColorLoadingText)).
		Render("Analyzing your code...")

	// Stack logo and tagline
	content := lipgloss.JoinVertical(lipgloss.Center, logo, tagline)

	// Perfect center placement
	return lipgloss.Place(
		m.TerminalWidth,
		m.TerminalHeight-3,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
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
	if m.TerminalWidth == 0 {
		return ""
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(Color(ColorAIResponseTag))).
		Align(lipgloss.Center).
		Width(m.TerminalWidth).
		Height(2)

	title := "[A]-[U]-[R]-[A]"

	instructions := lipgloss.NewStyle().
		Faint(true).
		Bold(false).
		Foreground(lipgloss.Color(Color(ColorAI))).
		Render("[q: quit | shift+↑/↓: scroll feedback]")

	return headerStyle.Render(title + "\n" + instructions)
}
