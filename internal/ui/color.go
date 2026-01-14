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
	ColorLoadingText      ColorName = "loading_text"
	ColorAI               ColorName = "ai"
	ColorAIResponseTag     ColorName = "ai_response_tags"
	ColorAIScroll         ColorName = "ai_scroll"
	ColorAuraLogo         ColorName = "aura_logo"
	ColorType             ColorName = "type"
	ColorIssue            ColorName = "issue"
	ColorSuggestion       ColorName = "suggestion"
	ColorAuraLoss         ColorName = "aura_loss"
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
	case ColorLoadingText:
		return lipgloss.Color("242")
	case ColorAI:
		return lipgloss.Color("205") // bright pink
	case ColorAIResponseTag:
		return lipgloss.Color("10") // bright green
	case ColorAIScroll:
		return lipgloss.Color("#89DCEB") // cyan

	// for preety print
	case ColorAuraLogo:
		return lipgloss.Color("205") // bright pink
	case ColorType:
		return lipgloss.Color("220") // bright yellow
	case ColorIssue:
		return lipgloss.Color("196") // bright red
	case ColorSuggestion:
		return lipgloss.Color("82") // bright yellow
	case ColorAuraLoss:
		return lipgloss.Color("214") // orange

	default:
		return lipgloss.Color("#CDD6F4") // default foreground
	}
}

var logoShimmerColors = []string{
    "93",  // bright magenta
    "129", // pink
    "165", // hot pink
    "201", // bright pink
    "199", // neon pink
    "196", // bright red
    "202", // orange-red
    "208", // orange
    "214", // gold
    "220", // bright yellow
    "226", // yellow
    "220",
    "214",
    "208",
    "202",
    "199",
    "196",
    "165",
    "129",
}
