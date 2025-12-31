package ui

import "github.com/charmbracelet/lipgloss"

type ColorName string

const (
	ColorAdded            ColorName = "primary"
	ColorRemoved          ColorName = "secondary"
	ColorHeaderText       ColorName = "header"
	ColorHeaderBackground ColorName = "header_background"
)

func Color(name ColorName) lipgloss.Color {
	switch name {

	case ColorAdded:
		return lipgloss.Color("10") // green
	case ColorRemoved:
		return lipgloss.Color("9") // red
	case ColorHeaderText:
		return lipgloss.Color("86") // cyan
	case ColorHeaderBackground:
		return lipgloss.Color("30") // light green

	default:
		return lipgloss.Color("7") // fallback
	}
}
