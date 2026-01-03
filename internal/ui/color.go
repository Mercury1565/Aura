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
	ColorAIResponeTag     ColorName = "ai_response_tags"
	ColorAIScroll         ColorName = "ai_scroll"
)

func Color(name ColorName) lipgloss.Color {
	switch name {

	case ColorAdded:
		return lipgloss.Color("#A6E3A1") // soft green
	case ColorAddedBG:
		return lipgloss.Color("22") // dark green background
	case ColorRemoved:
		return lipgloss.Color("#F38BA8") // soft red
	case ColorRemovedBG:
		return lipgloss.Color("52") // dark red background
	case ColorHeaderText:
		return lipgloss.Color("#89DCEB") // light cyan
	case ColorHeaderBackground:
		return lipgloss.Color("#313244") // dark gray
	case ColorLineNumber:
		return lipgloss.Color("#585B70") // muted gray
	case ColorAI:
		return lipgloss.Color("205") // bright pink
	case ColorAIResponeTag:
		return lipgloss.Color("10") // bright green
	case ColorAIScroll:
		return lipgloss.Color("#89DCEB") // cyan

	default:
		return lipgloss.Color("#CDD6F4") // default foreground
	}
}
