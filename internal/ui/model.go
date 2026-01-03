package ui

import (
	"context"

	"github.com/Mercury1565/Aura/internal/reviewer"
	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Ctx            context.Context
	DiffFiles      []*gitdiff.File
	Reviewer       *reviewer.LLMReviewer
	ReviewData     *reviewer.CodeReview
	IsLoading      bool
	TerminalWidth  int
	TerminalHeight int
	Viewport       viewport.Model
	Ready          bool
	AIScrollOffset int
	Error          []error
}

func InitialModel(
	files []*gitdiff.File,
	reviewer *reviewer.LLMReviewer,
	ctx context.Context,
) Model {
	return Model{
		Ctx:       ctx,
		DiffFiles: files,
		Reviewer:  reviewer,
	}
}

func (m Model) Init() tea.Cmd {
	return m.FetchReviewCmd()
}
