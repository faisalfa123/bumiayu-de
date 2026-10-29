package components

import (
	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderStatusBar renders the bottom shortcut and help bar
func RenderStatusBar(lang i18n.Language, width int, statusMessage string) string {
	t := i18n.Get(lang)

	left := styles.BadgeStyle.Render("⚡ BUMIAYU")
	shortcuts := styles.PillNormal.Render(t.HelpNav + " | " + t.HelpLangToggle + " | " + t.HelpQuit)

	var msg string
	if statusMessage != "" {
		msg = styles.PillWarning.Render("ℹ️ " + statusMessage)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Center, left, " ", shortcuts, " ", msg)
	return styles.FooterStyle.Width(width).Render(content)
}
