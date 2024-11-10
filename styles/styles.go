package styles

import "github.com/charmbracelet/lipgloss"

var (
	Subtle    = lipgloss.AdaptiveColor{Light: "#888888", Dark: "#666666"} // Grey
	Text      = lipgloss.AdaptiveColor{Light: "#1A202C", Dark: "F8F9FA"}  // White
	Highlight = lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#8B5CF6"} // Purple
	Special   = lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#3B82F6"} // Blue
	Attention = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#10B981"} // Green
	Warning   = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"} // Red

	BaseStyle = lipgloss.NewStyle().
			Foreground(Text).
			PaddingTop(2).
			PaddingLeft(4)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(Highlight).
			Bold(true)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	AttentionStyle = lipgloss.NewStyle().
			Foreground(Attention)
)
