package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var sb strings.Builder

	// Header with tabs
	sb.WriteString(m.renderHeader())
	sb.WriteString("\n")

	// Main content
	content := m.renderContent()
	sb.WriteString(content)

	// Prompt overlay
	if m.showPrompt {
		return m.renderWithPrompt(sb.String())
	}

	// Help overlay
	if m.showHelp {
		return m.renderWithHelp(sb.String())
	}

	// Status bar
	sb.WriteString("\n")
	sb.WriteString(m.renderStatusBar())

	return m.styles.App.Render(sb.String())
}

func (m Model) renderHeader() string {
	// Show simpler header when picking profile
	if m.activeView == ViewProfiles {
		title := m.styles.Title.Render("S3 TUI")
		subtitle := m.styles.Dim.Render("Select an AWS profile to continue")
		header := lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", subtitle)
		return m.styles.Header.Width(m.width - 2).Render(header)
	}

	tabs := []struct {
		name   string
		view   ViewType
		hotkey string
	}{
		{"Buckets", ViewBuckets, "1"},
		{"Browser", ViewBrowser, "2"},
		{"Bookmarks", ViewBookmarks, "3"},
	}

	var tabStrings []string
	for _, tab := range tabs {
		var style lipgloss.Style
		if m.activeView == tab.view {
			style = m.styles.ActiveTab
		} else {
			style = m.styles.Tab
		}
		tabStrings = append(tabStrings, style.Render(fmt.Sprintf("%s [%s]", tab.name, tab.hotkey)))
	}

	// Add download tab if active
	if m.downloadView.IsActive() || m.activeView == ViewDownload {
		var style lipgloss.Style
		if m.activeView == ViewDownload {
			style = m.styles.ActiveTab
		} else {
			style = m.styles.Tab.Foreground(ColorWarning)
		}
		tabStrings = append(tabStrings, style.Render("⏬ Downloads"))
	}

	tabLine := strings.Join(tabStrings, m.styles.TabSeparator.Render(" │ "))

	// Title
	title := m.styles.Title.Render("S3 TUI")

	// Profile info
	profile := m.styles.Dim.Render(fmt.Sprintf("Profile: %s", m.profileDisplay()))

	// Combine title, tabs, and profile
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		"  ",
		tabLine,
		"  ",
		profile,
	)

	return m.styles.Header.Width(m.width - 2).Render(header)
}

func (m Model) profileDisplay() string {
	if m.profile != "" {
		return m.profile
	}
	return "default"
}

func (m Model) renderContent() string {
	// Calculate content area
	contentHeight := m.height - 6 // header + status bar

	var content string
	switch m.activeView {
	case ViewProfiles:
		content = m.profilesView.View()
	case ViewBuckets:
		content = m.bucketsView.View()
	case ViewBrowser:
		content = m.browserView.View()
	case ViewDownload:
		content = m.downloadView.View()
	case ViewBookmarks:
		content = m.bookmarksView.View()
	default:
		content = "Unknown view"
	}

	// Ensure content fills the available space
	style := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(contentHeight)

	return style.Render(content)
}

func (m Model) renderStatusBar() string {
	// Left side: status message or error
	var leftContent string
	if m.errorMsg != "" {
		leftContent = m.styles.Error.Render("Error: " + m.errorMsg)
	} else if m.statusMsg != "" {
		leftContent = m.styles.Success.Render(m.statusMsg)
	} else {
		leftContent = m.renderContextualHelp()
	}

	// Right side: key hints
	rightContent := m.styles.Dim.Render("? help • q quit")

	// Calculate spacing
	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(rightContent)
	spacerWidth := m.width - leftWidth - rightWidth - 4

	if spacerWidth < 0 {
		spacerWidth = 1
	}

	spacer := strings.Repeat(" ", spacerWidth)

	return m.styles.HelpBar.Width(m.width - 2).Render(
		leftContent + spacer + rightContent,
	)
}

func (m Model) renderContextualHelp() string {
	switch m.activeView {
	case ViewProfiles:
		return m.styles.Dim.Render("↑↓ navigate • enter select profile • / filter")
	case ViewBuckets:
		return m.styles.Dim.Render("↑↓ navigate • enter select • / filter • ←→ tabs")
	case ViewBrowser:
		return m.styles.Dim.Render("↑↓ navigate • space select • enter open • d download • ←→ tabs")
	case ViewDownload:
		if m.downloadView.IsActive() {
			return m.styles.Dim.Render("esc cancel")
		}
		return m.styles.Dim.Render("←→ switch tabs")
	case ViewBookmarks:
		return m.styles.Dim.Render("↑↓ navigate • enter go to • x delete • ←→ tabs")
	default:
		return ""
	}
}

func (m Model) renderWithPrompt(base string) string {
	// Create prompt box
	promptStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Width(50)

	// Input with cursor
	input := m.promptInput
	cursor := "█"
	if m.promptCursor < len(input) {
		input = input[:m.promptCursor] + cursor + input[m.promptCursor:]
	} else {
		input = input + cursor
	}

	promptContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Title.Render(m.promptText),
		"",
		m.styles.PromptInput.Render(input),
		"",
		m.styles.Dim.Render("Enter to confirm • Esc to cancel"),
	)

	prompt := promptStyle.Render(promptContent)

	// Use lipgloss.Place to center the prompt over a background
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		prompt,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}

func (m Model) renderWithHelp(base string) string {
	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Width(60)

	helpContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Title.Render("Keyboard Shortcuts"),
		"",
		m.styles.Subtitle.Render("Navigation"),
		"  ↑/k, ↓/j    Move up/down",
		"  Enter       Open folder",
		"  Backspace   Go back",
		"  PgUp/PgDn   Page up/down",
		"",
		m.styles.Subtitle.Render("Views"),
		"  ←/→         Switch tabs",
		"  Tab         Next tab",
		"  Shift+Tab   Previous tab",
		"  1/2/3       Jump to tab",
		"",
		m.styles.Subtitle.Render("Selection & Actions"),
		"  Space       Select/deselect item",
		"  d           Download selected (or current)",
		"  s           Sync prefix to local",
		"  b           Add bookmark",
		"  r           Refresh",
		"  /           Filter list",
		"",
		m.styles.Subtitle.Render("General"),
		"  ?           Toggle this help",
		"  Esc         Cancel / Close",
		"  q           Quit",
		"",
		m.styles.Dim.Render("Press Esc or ? to close"),
	)

	help := helpStyle.Render(helpContent)

	// Use lipgloss.Place to center the help
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		help,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}

