package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderDashboard renders system overview tailored to the active WM model and Theme
func RenderDashboard(info system.HardwareInfo, wm WMMode, st styles.ThemeStyles, lang i18n.Language, width, height int) string {
	contentWidth := width - 6
	if contentWidth < 50 {
		contentWidth = 50
	}

	header := st.Title.Render("🏝️ BUMIAYU DESKTOP ENVIRONMENT")
	sub := st.Subtitle.Render("Pusat Kontrol Sistem Ringan & Ramah untuk ODROID-C2 & Server Linux")

	// Telemetry data
	deviceType := "Linux Standard Server (" + info.Arch + ")"
	if info.IsOdroidC2 {
		deviceType = "ODROID-C2 (Amlogic S905 Quad-Core ARM64)"
	}

	ramBar := system.GetMemoryBar(info.MemUsagePct, 18)
	swapBar := system.GetMemoryBar(info.SwapUsagePct, 18)

	healthBadge := st.Badge.Render("🟢 Sistem Normal & Sehat")
	if info.MemUsagePct > 80 || info.TemperatureC > 70 {
		healthBadge = st.PillDanger.Render("🔴 Perhatian: Beban Tinggi!")
	}

	switch wm {
	case WMTiling:
		// Tiling: 2 Kolom Vertikal
		colLeft := st.Card.Width(contentWidth/2 - 2).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				st.Title.Render("💻 Perangkat Keras:"),
				"",
				fmt.Sprintf("• Model : %s\n• Host  : %s\n• Kernel: %s\n• Uptime: %s\n• Cores : %d Inti CPU\n• Status: %s",
					deviceType, info.Hostname, info.Kernel, info.UptimeStr, info.CPUCores, healthBadge),
			),
		)

		colRight := st.CardActive.Width(contentWidth/2 - 2).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				st.Title.Render("🧠 Penggunaan Memori & Suhu:"),
				"",
				fmt.Sprintf("• RAM  : %d / %d MB\n  %s\n\n• Swap : %d / %d MB\n  %s\n\n• Suhu : %.1f °C (Thermal Zone 0)\n• Load : %.2f, %.2f, %.2f",
					info.MemUsedMB, info.MemTotalMB, ramBar,
					info.SwapUsedMB, info.SwapTotalMB, swapBar,
					info.TemperatureC, info.Load1, info.Load5, info.Load15),
			),
		)

		mainContent := lipgloss.JoinHorizontal(lipgloss.Top, colLeft, "  ", colRight)
		return lipgloss.JoinVertical(lipgloss.Left, header, sub, "", mainContent)

	case WMCardGrid:
		// Card Grid: 3 Kartu Berjejer
		card1 := st.Card.Width(contentWidth/3 - 2).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				st.Title.Render("🏥 Kesehatan"),
				"",
				healthBadge,
				"",
				fmt.Sprintf("Host: %s\nUptime: %s", info.Hostname, info.UptimeStr),
			),
		)

		card2 := st.Card.Width(contentWidth/3 - 2).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				st.Title.Render("🧠 Memori"),
				"",
				fmt.Sprintf("%d / %d MB", info.MemUsedMB, info.MemTotalMB),
				ramBar,
				"",
				fmt.Sprintf("Swap: %d MB", info.SwapUsedMB),
			),
		)

		card3 := st.CardActive.Width(contentWidth/3 - 2).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				st.Title.Render("🌡️ Suhu & Beban"),
				"",
				fmt.Sprintf("Suhu: %.1f °C\nLoad: %.2f", info.TemperatureC, info.Load1),
				"",
				fmt.Sprintf("CPU: %d Cores", info.CPUCores),
			),
		)

		mainContent := lipgloss.JoinHorizontal(lipgloss.Top, card1, " ", card2, " ", card3)
		return lipgloss.JoinVertical(lipgloss.Left, header, sub, "", mainContent)

	default: // WMFloating
		// Floating / Center Window yang sangat anggun dan bersih
		windowWidth := contentWidth - 10
		if windowWidth > 80 {
			windowWidth = 80
		}

		cardContent := lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Center, st.Title.Render("💻 "+deviceType), "   ", healthBadge),
			st.Subtitle.Render("ODROID-C2 & Linux Server Control Center"),
			"",
			fmt.Sprintf("• Hostname     : %s\n• Kernel Linux : %s\n• Masa Aktif   : %s", info.Hostname, info.Kernel, info.UptimeStr),
			"",
			st.Title.Render("🧠 Penggunaan Sumber Daya:"),
			fmt.Sprintf("• Memori RAM   : %d MB / %d MB   %s", info.MemUsedMB, info.MemTotalMB, ramBar),
			fmt.Sprintf("• Ruang Swap   : %d MB / %d MB   %s", info.SwapUsedMB, info.SwapTotalMB, swapBar),
			fmt.Sprintf("• Suhu Inti    : %.1f °C  | Beban CPU: %.2f (1m), %.2f (5m)", info.TemperatureC, info.Load1, info.Load5),
			"",
			st.RowNormal.Render("💡 Tips Cepat: Tekan [2] untuk buka Toko Aplikasi, [3] Pengaturan Tampilan, [4] Petunjuk"),
		)

		floatingWindow := st.CardActive.Width(windowWidth).Render(cardContent)
		centered := lipgloss.Place(contentWidth, 16, lipgloss.Center, lipgloss.Center, floatingWindow)
		return lipgloss.JoinVertical(lipgloss.Left, header, sub, "", centered)
	}
}
