package components

import (
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderStatusBar renders the bottom shortcut and help bar
func RenderStatusBar(st styles.ThemeStyles, width int, statusMessage string) string {
	left := st.Badge.Render("⚡ BUMIAYU")
	shortcuts := st.Pill.Render("[1-4] atau [←/→]: Pindah Tab | [↑/↓]: Pilih | [Enter/I]: Pasang | [X]: Hapus | [T]: Tema | [q]: Keluar")

	var msg string
	if statusMessage != "" {
		msg = st.PillWarning.Render("ℹ️ " + statusMessage)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Center, left, " ", shortcuts, " ", msg)
	return st.Footer.Width(width).Render(content)
}
