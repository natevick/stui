package download

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/natevick/s3-tui/internal/download"
)

// Model is the download view model
type Model struct {
	progress    download.Progress
	progressBar progress.Model
	active      bool
	width       int
	height      int
}

// New creates a new download view
func New() Model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	return Model{
		progressBar: p,
	}
}

// SetSize sets the view size
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.progressBar.Width = width - 20
}

// SetProgress updates the download progress
func (m *Model) SetProgress(p download.Progress) {
	m.progress = p
	m.active = p.Status == download.StatusInProgress || p.Status == download.StatusPending
}

// IsActive returns true if a download is in progress
func (m Model) IsActive() bool {
	return m.active
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progress.FrameMsg:
		progressModel, cmd := m.progressBar.Update(msg)
		m.progressBar = progressModel.(progress.Model)
		return m, cmd
	}
	return m, nil
}

// View renders the view
func (m Model) View() string {
	if !m.active && m.progress.TotalFiles == 0 {
		return m.renderNoDownload()
	}

	var sb strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Padding(0, 1).
		Render("Downloads")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	// Status
	statusStyle := lipgloss.NewStyle().Padding(0, 1)
	switch m.progress.Status {
	case download.StatusInProgress:
		sb.WriteString(statusStyle.Foreground(lipgloss.Color("214")).Render("⏳ Downloading..."))
	case download.StatusCompleted:
		sb.WriteString(statusStyle.Foreground(lipgloss.Color("78")).Render("✓ Download complete"))
	case download.StatusFailed:
		sb.WriteString(statusStyle.Foreground(lipgloss.Color("196")).Render("✗ Download failed"))
	case download.StatusCancelled:
		sb.WriteString(statusStyle.Foreground(lipgloss.Color("240")).Render("⊘ Download cancelled"))
	}
	sb.WriteString("\n\n")

	// Overall progress
	percent := m.progress.PercentComplete() / 100
	sb.WriteString(lipgloss.NewStyle().Padding(0, 1).Render(m.progressBar.ViewAs(percent)))
	sb.WriteString("\n\n")

	// Stats
	statsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	stats := fmt.Sprintf("Files: %d/%d  •  %s / %s",
		m.progress.CompletedFiles,
		m.progress.TotalFiles,
		humanize.Bytes(uint64(m.progress.DownloadedBytes)),
		humanize.Bytes(uint64(m.progress.TotalBytes)),
	)
	sb.WriteString(statsStyle.Render(stats))
	sb.WriteString("\n")

	if m.progress.FailedFiles > 0 {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Padding(0, 1).
			Render(fmt.Sprintf("Failed: %d files", m.progress.FailedFiles)))
		sb.WriteString("\n")
	}

	// Current file
	if m.progress.CurrentFile != "" && m.progress.Status == download.StatusInProgress {
		sb.WriteString("\n")
		sb.WriteString(statsStyle.Render(fmt.Sprintf("Current: %s", truncatePath(m.progress.CurrentFile, m.width-20))))
	}

	// File list (last 10 files)
	if len(m.progress.Files) > 0 {
		sb.WriteString("\n\n")
		sb.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Padding(0, 1).
			Render("Recent files:"))
		sb.WriteString("\n")

		count := 0
		for _, fp := range m.progress.Files {
			if count >= 10 {
				break
			}

			var statusIcon string
			var style lipgloss.Style
			switch fp.Status {
			case download.StatusCompleted:
				statusIcon = "✓"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
			case download.StatusInProgress:
				statusIcon = "⏳"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
			case download.StatusFailed:
				statusIcon = "✗"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			case download.StatusCancelled:
				statusIcon = "⊘"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			default:
				statusIcon = "○"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			}

			line := fmt.Sprintf("  %s %s (%s)",
				statusIcon,
				truncatePath(fp.Key, m.width-30),
				humanize.Bytes(uint64(fp.Size)),
			)
			sb.WriteString(style.Render(line))
			sb.WriteString("\n")
			count++
		}

		if len(m.progress.Files) > 10 {
			sb.WriteString(statsStyle.Render(fmt.Sprintf("  ... and %d more files", len(m.progress.Files)-10)))
		}
	}

	// Help
	sb.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	if m.active {
		sb.WriteString(helpStyle.Render("Press Esc to cancel"))
	} else {
		sb.WriteString(helpStyle.Render("Press 1 to go to Buckets, 2 to go to Browser"))
	}

	return sb.String()
}

func (m Model) renderNoDownload() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("240"))

	return style.Render("No downloads in progress\n\nPress 'd' on a file or folder in the Browser to download")
}

func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}
