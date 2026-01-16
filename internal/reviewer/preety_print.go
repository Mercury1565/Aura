package reviewer

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	SummaryColor    = lipgloss.Color("252") // bright white
	TypeColor       = lipgloss.Color("220") // bright yellow
	IssueColor      = lipgloss.Color("214") // orange
	SuggestionColor = lipgloss.Color("82")  // green
	AuraLossColor   = lipgloss.Color("196") // bright red
)

func (r *CodeReview) PrettyPrint() {
	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(SummaryColor)).
		Italic(true)

	typeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(TypeColor))
	issueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(IssueColor)).Bold(true)
	suggestionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(SuggestionColor)).Bold(true)
	auraLossStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(AuraLossColor)).Bold(true)

	fmt.Println(summaryStyle.Render(r.Summary) + "\n")

	for _, rev := range r.Reviews {
		fmt.Printf("📂 %s (Line %d)\n", rev.File, rev.Line)
		fmt.Printf("   %s %s\n", typeStyle.Render("TYPE:"), rev.Type)
		fmt.Printf("   %s %s\n", issueStyle.Render("ISSUE:"), rev.Issue)
		fmt.Printf("   %s %s\n", suggestionStyle.Render("SUGGEST:"), rev.Suggestion)
		fmt.Printf("   %s -%d points\n\n", auraLossStyle.Render("AURA LOSS:"), rev.AuraLoss)
	}
}
