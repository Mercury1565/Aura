package ui

import (
	"github.com/Mercury1565/Aura/internal/reviewer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case []error:
		m.Error = msg
		m.IsLoading = false
		m.Viewport.SetContent(m.renderDiffContent())
		return m, nil

	case *reviewer.CodeReview:
		m.ReviewData = msg
		m.IsLoading = false
		m.Viewport.SetContent(m.renderDiffContent())
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// Horizontal scroll for AI columns or internal AI shifting
		case "shift+down":
			m.AIScrollOffset++
			m.Viewport.SetContent(m.renderDiffContent())
			return m, nil

		case "shift+up":
			if m.AIScrollOffset > 0 {
				m.AIScrollOffset--
				m.Viewport.SetContent(m.renderDiffContent())
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		m.TerminalHeight = msg.Height

		// Adjust heights
		headerHeight := 2
		verticalMarginHeight := headerHeight

		if !m.Ready {
			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.Viewport.YPosition = headerHeight
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - verticalMarginHeight
		}

		// Re-render content on resize so column math stays correct
		m.Viewport.SetContent(m.renderDiffContent())
	}

	// Handle standard vertical scrolling for the whole page
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) FetchReviewCmd() tea.Cmd {
	return func() tea.Msg {
		// Try structured first
		feedback, err := m.Reviewer.ReviewDiffWithStructuredOutput(m.Ctx, m.DiffFiles)

		if err != nil || feedback == nil || len(feedback.Reviews) == 0 {
			// Fallback
			raw, fallbackErr := m.Reviewer.ReviewDiff(m.Ctx, m.DiffFiles)

			if fallbackErr != nil {
				return []error{err, fallbackErr}
			}

			feedback = m.Reviewer.ParseUnstructuredReview(raw)
		}
		return feedback
	}
}
