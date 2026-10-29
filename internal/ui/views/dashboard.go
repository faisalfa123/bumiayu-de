package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderDashboard renders the main system health and hardware overview
func RenderDashboard(info system.HardwareInfo, lang i18n.Language, width, height int) string {
	t := i18n.Get(lang)

	contentWidth := width - 6
	if contentWidth < 40 {
		contentWidth = 40
	}

	// 1. System Health Status Card
	healthTitle := styles.TitleStyle.Render("🏥 Status Kesehatan Sistem")
	healthBadge := styles.PillNormal.Render("🟢 Normal & Siaga")
	if info.MemUsagePct > 85 || info.TemperatureC > 75 {
		healthBadge = styles.PillDanger.Render("🔴 Kritis / Suhu Tinggi!")
	} else if info.MemUsagePct > 70 || info.TemperatureC > 60 {
		healthBadge = styles.PillWarning.Render("🟡 Waspada Beban Memori")
	}

	// 2. Hardware specs list
	deviceType := "Linux Standard Server (" + info.Arch + ")"
	if info.IsOdroidC2 {
		deviceType = "ODROID-C2 (Amlogic S905 - 64-bit Quad-Core ARM Cortex-A53)"
	}

	hwDetails := fmt.Sprintf(
		"• %s: %s\n• Kernel: %s\n• %s: %s\n• %s: %s\n• CPU Cores: %d Core(s)\n• %s: %.2f (1m), %.2f (5m), %.2f (15m)",
		t.Architecture, deviceType,
		info.Kernel,
		t.CurrentHost, info.Hostname,
		t.Uptime, info.UptimeStr,
		info.CPUCores,
		t.CpuLoad, info.Load1, info.Load5, info.Load15,
	)

	hwCard := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Center, healthTitle, "  ", healthBadge),
			"",
			hwDetails,
		),
	)

	// 3. Memory & Thermals Card
	memTitle := styles.TitleStyle.Render("🧠 Sumber Daya Memori & Suhu")
	ramBar := system.GetMemoryBar(info.MemUsagePct, 22)
	swapBar := system.GetMemoryBar(info.SwapUsagePct, 22)

	memDetails := fmt.Sprintf(
		"• RAM Fisik : %d MB / %d MB\n  %s\n\n• Swap Disk : %d MB / %d MB\n  %s\n\n• Suhu SoC  : %.1f °C\n• Status Fan: Pasif Heatsink (Thermal Zone 0)",
		info.MemUsedMB, info.MemTotalMB, ramBar,
		info.SwapUsedMB, info.SwapTotalMB, swapBar,
		info.TemperatureC,
	)

	memCard := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			memTitle,
			"",
			memDetails,
		),
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, hwCard, "  ", memCard)

	// 4. Senior Insight Box
	insightHeader := styles.InsightTitleStyle.Render(t.SeniorInsightTitle)
	insightBody := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0f2fe")).Render(t.InsightDashboard)
	insightBox := styles.InsightBoxStyle.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left, insightHeader, insightBody),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		topRow,
		"",
		insightBox,
	)
}
