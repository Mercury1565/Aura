package ui

import (
	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	DiffFiles      []*gitdiff.File
	TerminalWidth  int
	TerminalHeight int
	Viewport       viewport.Model
	Ready          bool
}

func InitialModel(files []*gitdiff.File) Model {
	return Model{DiffFiles: files}
}

func (m Model) Init() tea.Cmd {
	return nil
}
