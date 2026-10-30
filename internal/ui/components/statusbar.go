package components

import (
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderStatusBar renders the bottom shortcut and help bar
func RenderStatusBar(st styles.ThemeStyles, width int, statusMessage string) string {
	left := st.Badge.Render("⚡ BUMIAYU")
	shortcuts := st.ButtonNormal.Render("[↑/↓]: Pilih App | [Enter/I]: Pasang | [X]: Hapus | [1-3]/[←→]: Tab | [T]: Tema | [M]: WM | [W]: Wizard | [q]: Keluar")

	var msg string
	if statusMessage != "" {
		msg = st.PillWarning.Render("ℹ️ " + statusMessage)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Center, left, " ", shortcuts, " ", msg)
	return st.Footer.Width(width).Render(content)
}
