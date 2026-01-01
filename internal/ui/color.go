package ui

import "github.com/charmbracelet/lipgloss"

type ColorName string

const (
	ColorAdded            ColorName = "primary"
	ColorAddedBG          ColorName = "primary_background"
	ColorRemoved          ColorName = "secondary"
	ColorRemovedBG        ColorName = "secondary_background"
	ColorHeaderText       ColorName = "header"
	ColorHeaderBackground ColorName = "header_background"
	ColorLineNumber       ColorName = "line_number"
	ColorAI               ColorName = "ai"
)

func Color(name ColorName) lipgloss.Color {
	switch name {

	case ColorAdded:
		return lipgloss.Color("#A6E3A1")
	case ColorAddedBG:
		return lipgloss.Color("22")
	case ColorRemoved:
		return lipgloss.Color("#F38BA8")
	case ColorRemovedBG:
		return lipgloss.Color("52")
	case ColorHeaderText:
		return lipgloss.Color("#89DCEB")
	case ColorHeaderBackground:
		return lipgloss.Color("#313244")
	case ColorLineNumber:
		return lipgloss.Color("#585B70")
	case ColorAI:
		return lipgloss.Color("205")

	default:
		return lipgloss.Color("#CDD6F4")
	}
}
