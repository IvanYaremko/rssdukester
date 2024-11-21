package views

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#888888", Dark: "#666666"} // Grey
	text      = lipgloss.AdaptiveColor{Light: "#1A202C", Dark: "F8F9FA"}  // White
	highlight = lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#9D86E9"} // Purple
	special   = lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#61AFEF"} // Blue
	attention = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#98C379"} // Green
	warning   = lipgloss.AdaptiveColor{Light: "#E11D48", Dark: "#F87171"} // Red

	baseStyle = lipgloss.NewStyle().
			Foreground(text).
			PaddingTop(2).
			PaddingLeft(4)

	highlightStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true)

	subtleStyle = lipgloss.NewStyle().
			Foreground(subtle)

	errorStyle = lipgloss.NewStyle().
			Foreground(warning).
			Bold(true)

	attentionStyle = lipgloss.NewStyle().
			Foreground(attention)

	specialStyle = lipgloss.NewStyle().
			Foreground(special)

	itemNormalTitle = lipgloss.NewStyle().
			Width(width).
			Foreground(text)

	itemSelectedTitle = lipgloss.NewStyle().
				Width(width).
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(highlight).
				PaddingLeft(1)

	itemNormalDesc = lipgloss.NewStyle().
			Width(width).
			Foreground(text)

	itemSelectedDesc = lipgloss.NewStyle().
				Width(width).
				Foreground(highlight).
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(highlight).
				PaddingLeft(1)
)
