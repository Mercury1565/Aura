package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var AuraLogoColor = lipgloss.Color("205") // bright pink

func LogoPrettyPrint() {
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(AuraLogoColor)).
		Bold(true)

	fmt.Println(logoStyle.Render(AuraLogo))

	// Code review in progress message
	reviewStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("242")).
		Bold(true).
		Italic(true)

	fmt.Println(reviewStyle.Render("Code review in progress...\n"))
}
