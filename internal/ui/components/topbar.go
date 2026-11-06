package components

import (
	"fmt"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/retro"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderTopBar renders the desktop status bar with animated 8-bit retro mascot and telemetry
func RenderTopBar(info system.HardwareInfo, animFrame int, st styles.ThemeStyles, lang i18n.Language, width int) string {
	// Left: Animated Mascot Badge
	mascotBadge := retro.GetBannerMascot(animFrame, st)

	// Center: Telemetry Pills
	var tempPill string
	if info.TemperatureC > 70 {
		tempPill = st.PillDanger.Render(fmt.Sprintf("🔥 %.1f°C", info.TemperatureC))
	} else if info.TemperatureC > 55 {
		tempPill = st.PillWarning.Render(fmt.Sprintf("⚠️ %.1f°C", info.TemperatureC))
	} else {
		tempPill = st.ButtonNormal.Render(fmt.Sprintf("❄️ %.1f°C", info.TemperatureC))
	}

	ramText := fmt.Sprintf("RAM: %d/%d MB (%.0f%%)", info.MemUsedMB, info.MemTotalMB, info.MemUsagePct)
	var ramPill string
	if info.MemUsagePct > 85 {
		ramPill = st.PillDanger.Render(ramText)
	} else if info.MemUsagePct > 70 {
		ramPill = st.PillWarning.Render(ramText)
	} else {
		ramPill = st.ButtonNormal.Render(ramText)
	}

	cpuPill := st.ButtonNormal.Render(fmt.Sprintf("Load: %.2f", info.Load1))
	themePill := lipgloss.NewStyle().
		Foreground(st.PrimaryColor).
		Bold(true).
		Render(fmt.Sprintf("[%s]", st.Name))

	clockPill := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(st.AccentColor).
		Padding(0, 1).
		Render(time.Now().Format("15:04:05"))

	left := lipgloss.JoinHorizontal(lipgloss.Center, mascotBadge, " ", themePill)
	center := lipgloss.JoinHorizontal(lipgloss.Center, tempPill, " ", ramPill, " ", cpuPill)
	right := clockPill

	totalUsed := lipgloss.Width(left) + lipgloss.Width(center) + lipgloss.Width(right)
	spacing := (width - totalUsed) / 2
	if spacing < 1 {
		spacing = 1
	}

	spaces := lipgloss.NewStyle().Width(spacing).Render("")
	barContent := lipgloss.JoinHorizontal(lipgloss.Center, left, spaces, center, spaces, right)

	return st.TopBar.Width(width).Render(barContent)
}
