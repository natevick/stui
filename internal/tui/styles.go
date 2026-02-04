package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles holds all the styling for the TUI
type Styles struct {
	// App chrome
	App       lipgloss.Style
	Header    lipgloss.Style
	StatusBar lipgloss.Style
	HelpBar   lipgloss.Style

	// Tabs
	Tab          lipgloss.Style
	ActiveTab    lipgloss.Style
	TabSeparator lipgloss.Style

	// List items
	Item         lipgloss.Style
	SelectedItem lipgloss.Style
	Folder       lipgloss.Style
	File         lipgloss.Style

	// Info display
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Info     lipgloss.Style
	Dim      lipgloss.Style

	// Progress
	Progress      lipgloss.Style
	ProgressBar   lipgloss.Style
	ProgressTrack lipgloss.Style

	// Messages
	Error   lipgloss.Style
	Success lipgloss.Style
	Warning lipgloss.Style

	// Prompt
	Prompt      lipgloss.Style
	PromptInput lipgloss.Style

	// Bookmark
	Bookmark lipgloss.Style
}

// Colors used throughout the app
var (
	ColorPrimary   = lipgloss.Color("39")  // Blue
	ColorSecondary = lipgloss.Color("213") // Pink
	ColorSuccess   = lipgloss.Color("78")  // Green
	ColorWarning   = lipgloss.Color("214") // Orange
	ColorError     = lipgloss.Color("196") // Red
	ColorDim       = lipgloss.Color("240") // Gray
	ColorBright    = lipgloss.Color("255") // White
	ColorFolder    = lipgloss.Color("226") // Yellow
	ColorFile      = lipgloss.Color("252") // Light gray
)

// DefaultStyles creates the default style set
func DefaultStyles() Styles {
	return Styles{
		App: lipgloss.NewStyle().
			Padding(0, 1),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(ColorDim).
			Padding(0, 1).
			MarginBottom(1),

		StatusBar: lipgloss.NewStyle().
			Foreground(ColorDim).
			Padding(0, 1),

		HelpBar: lipgloss.NewStyle().
			Foreground(ColorDim).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(ColorDim).
			Padding(0, 1),

		Tab: lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(ColorDim),

		ActiveTab: lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(ColorBright).
			Background(ColorPrimary).
			Bold(true),

		TabSeparator: lipgloss.NewStyle().
			Foreground(ColorDim),

		Item: lipgloss.NewStyle().
			Padding(0, 1),

		SelectedItem: lipgloss.NewStyle().
			Padding(0, 1).
			Background(ColorPrimary).
			Foreground(ColorBright),

		Folder: lipgloss.NewStyle().
			Foreground(ColorFolder).
			Bold(true),

		File: lipgloss.NewStyle().
			Foreground(ColorFile),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary),

		Subtitle: lipgloss.NewStyle().
			Foreground(ColorSecondary),

		Info: lipgloss.NewStyle().
			Foreground(ColorDim),

		Dim: lipgloss.NewStyle().
			Foreground(ColorDim),

		Progress: lipgloss.NewStyle().
			Padding(0, 1),

		ProgressBar: lipgloss.NewStyle().
			Foreground(ColorSuccess),

		ProgressTrack: lipgloss.NewStyle().
			Foreground(ColorDim),

		Error: lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true),

		Success: lipgloss.NewStyle().
			Foreground(ColorSuccess),

		Warning: lipgloss.NewStyle().
			Foreground(ColorWarning),

		Prompt: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary),

		PromptInput: lipgloss.NewStyle().
			Foreground(ColorBright),

		Bookmark: lipgloss.NewStyle().
			Foreground(ColorSecondary),
	}
}
