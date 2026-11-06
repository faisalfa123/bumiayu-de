package retro

import (
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderMarqueeTicker renders a smooth 24 FPS scrolling marquee line across terminal width
func RenderMarqueeTicker(messages []string, frame int, width int, st styles.ThemeStyles) string {
	if width < 20 {
		return ""
	}

	if len(messages) == 0 {
		messages = []string{
			"⚡ BUMIAYU TUI WINDOW MANAGER",
			"ADAPTIVE 24 FPS ENGINE READY",
			"PRESS [SPACE] FOR UNIFIED LAUNCHER",
			"PRESS [?] FOR ARCADE CONTROLS HELP",
			"LINUX SYSTEM MANAGEMENT ACTIVE",
		}
	}

	fullText := strings.Join(messages, "  ✦  ") + "  ✦  "
	runes := []rune(fullText)
	textLen := len(runes)
	if textLen == 0 {
		return ""
	}

	// 24 FPS ticker: shift 1 character every 3 ticks
	offset := (frame / 3) % textLen

	// Construct visible slice
	var sb strings.Builder
	for i := 0; i < width-4; i++ {
		idx := (offset + i) % textLen
		sb.WriteRune(runes[idx])
	}

	prefix := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(st.AccentColor).
		Render(" 📡 TICKER ")

	tickerContent := lipgloss.NewStyle().
		Foreground(st.PrimaryColor).
		Render(" " + sb.String())

	return lipgloss.JoinHorizontal(lipgloss.Center, prefix, tickerContent)
}
