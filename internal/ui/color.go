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
	case ColorLoadingText:
		return lipgloss.Color("242")
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

var logoShimmerColors = []string{
	"57",  // deep purple
	"63",  // violet
	"93",  // magenta
	"129", // pink
	"165", // hot pink
	"201", // bright pink
	"207", // light pink
	"199", // neon pink
	"163", // orange-pink
	"208", // bright orange
	"214", // yellow-orange
	"220", // bright yellow
	"226", // yellow
	"190", // lime
	"118", // green
	"82",  // bright green
	"45",  // cyan
	"39",  // bright cyan
	"33",  // blue
	"27",  // deep blue
	"93",  // magenta again for bounce
}
