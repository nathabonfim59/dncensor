package tui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			Padding(0, 1)

	DimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Padding(0, 1)

	HintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EF4444"))

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#10B981"))

	CheckboxStyle = lipgloss.NewStyle().
			Padding(0, 1)

	CheckboxChecked = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10B981")).
				Padding(0, 1)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2)

	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#9CA3AF"))

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0"))
)
