package components

import (
	"fmt"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderTopBar renders the desktop status bar
func RenderTopBar(info system.HardwareInfo, lang i18n.Language, width int) string {
	t := i18n.Get(lang)

	// Left: App Name & Host
	hostBadge := styles.BadgeStyle.Render("🏝️ " + t.AppTitle)
	hostInfo := styles.PillNormal.Render(fmt.Sprintf("%s: %s (%s)", t.CurrentHost, info.Hostname, info.Arch))

	// Center: Telemetry Pills
	// 1. Temperature pill with dynamic color
	var tempPill string
	if info.TemperatureC > 70 {
		tempPill = styles.PillDanger.Render(fmt.Sprintf("🔥 %.1f°C", info.TemperatureC))
	} else if info.TemperatureC > 55 {
		tempPill = styles.PillWarning.Render(fmt.Sprintf("⚠️ %.1f°C", info.TemperatureC))
	} else {
		tempPill = styles.PillNormal.Render(fmt.Sprintf("❄️ %.1f°C", info.TemperatureC))
	}

	// 2. RAM pill
	ramText := fmt.Sprintf("RAM: %d/%d MB (%.0f%%)", info.MemUsedMB, info.MemTotalMB, info.MemUsagePct)
	var ramPill string
	if info.MemUsagePct > 85 {
		ramPill = styles.PillDanger.Render(ramText)
	} else if info.MemUsagePct > 70 {
		ramPill = styles.PillWarning.Render(ramText)
	} else {
		ramPill = styles.PillNormal.Render(ramText)
	}

	// 3. CPU / Load
	cpuPill := styles.PillNormal.Render(fmt.Sprintf("Load: %.2f %.2f", info.Load1, info.Load5))

	// Right: Clock
	clockPill := styles.PillNormal.Render(time.Now().Format("15:04:05"))

	left := lipgloss.JoinHorizontal(lipgloss.Center, hostBadge, " ", hostInfo)
	center := lipgloss.JoinHorizontal(lipgloss.Center, tempPill, " ", ramPill, " ", cpuPill)
	right := lipgloss.JoinHorizontal(lipgloss.Center, clockPill)

	// Calculate space distribution
	totalUsed := lipgloss.Width(left) + lipgloss.Width(center) + lipgloss.Width(right)
	spacing := (width - totalUsed) / 2
	if spacing < 1 {
		spacing = 1
	}

	spaces := lipgloss.NewStyle().Width(spacing).Render("")
	barContent := lipgloss.JoinHorizontal(lipgloss.Center, left, spaces, center, spaces, right)

	return styles.TopBarStyle.Width(width).Render(barContent)
}
