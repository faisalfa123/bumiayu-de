package components

import (
	"fmt"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderTopBar renders the desktop status bar using active theme styles
func RenderTopBar(info system.HardwareInfo, st styles.ThemeStyles, lang i18n.Language, width int) string {
	t := i18n.Get(lang)

	// Left: App Name & Host
	hostBadge := st.Badge.Render("🏝️ " + t.AppTitle)
	hostInfo := lipgloss.NewStyle().Foreground(st.TextColor).Bold(true).Render(fmt.Sprintf("%s (%s)", info.Hostname, info.Arch))

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
	clockPill := st.ButtonNormal.Render(time.Now().Format("15:04:05"))

	left := lipgloss.JoinHorizontal(lipgloss.Center, hostBadge, " ", hostInfo)
	center := lipgloss.JoinHorizontal(lipgloss.Center, tempPill, " ", ramPill, " ", cpuPill)
	right := lipgloss.JoinHorizontal(lipgloss.Center, clockPill)

	totalUsed := lipgloss.Width(left) + lipgloss.Width(center) + lipgloss.Width(right)
	spacing := (width - totalUsed) / 2
	if spacing < 1 {
		spacing = 1
	}

	spaces := lipgloss.NewStyle().Width(spacing).Render("")
	barContent := lipgloss.JoinHorizontal(lipgloss.Center, left, spaces, center, spaces, right)

	return st.TopBar.Width(width).Render(barContent)
}
