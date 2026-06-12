package tui

import "charm.land/lipgloss/v2"

// Colors
const (
	purple   = "#7C3AED"
	white    = "#E0E0E0"
	gray     = "#9CA3AF"
	dimgray  = "#6B7280"
	darkgray = "#4B5563"
	green    = "#10B981"
	red      = "#EF4444"
)

// Container – wraps the entire view
var CardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(purple)).
		Padding(1, 2)

// Typography hierarchy
var TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(purple))

var HeadingStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(gray))

var BodyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(white))

var DimmedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(dimgray))

var HintStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(gray))

// List items
var ItemSelectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(purple)).
		Padding(0, 1)

var ItemNormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(white)).
		Padding(0, 1)

var ItemDimmedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(darkgray)).
		Padding(0, 1)

// Apply button
var ApplyBtnStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(purple)).
		Padding(0, 3)

var ApplyBtnDimmedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(darkgray)).
		Padding(0, 3)

// Checkbox toggle
var CheckOnStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(green))

var CheckOffStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(dimgray))

// Feedback
var ErrorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(red))

var SuccessStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(green))

// Confirm screen key-value rows
var RowLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(gray)).
		Width(16).
		Align(lipgloss.Right)

var ValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(white))
