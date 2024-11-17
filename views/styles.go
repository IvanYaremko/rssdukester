package views

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#888888", Dark: "#666666"} // Grey
	text      = lipgloss.AdaptiveColor{Light: "#1A202C", Dark: "F8F9FA"}  // White
	highlight = lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#8B5CF6"} // Purple
	special   = lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#3B82F6"} // Blue
	attention = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#10B981"} // Green
	warning   = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"} // Red

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
)
