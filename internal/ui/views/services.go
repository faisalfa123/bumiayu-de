package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderServices renders the systemd service controller
func RenderServices(services []system.ServiceItem, selectedIndex int, showLogs bool, logOutput string, lang i18n.Language, width, height int) string {
	t := i18n.Get(lang)
	contentWidth := width - 6
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := styles.TitleStyle.Render("⚙️ Manajer Layanan Latar Belakang (Systemd Units)")

	if showLogs {
		// Log Viewer Modal
		logTitle := styles.BadgeStyle.Render("📜 Live Journalctl Logs: " + services[selectedIndex].Name)
		logContent := styles.CardActiveStyle.Width(contentWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				logTitle,
				"",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#cbd5e1")).Render(logOutput),
				"",
				styles.SubtitleStyle.Render("Tekan [Esc] atau [L] untuk kembali ke daftar layanan"),
			),
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, "", logContent)
	}

	var listItems []string
	for i, s := range services {
		var stateBadge string
		switch s.ActiveState {
		case "active":
			stateBadge = styles.BadgeStyle.Render("🟢 RUNNING")
		case "failed":
			stateBadge = styles.PillDanger.Render("🔴 FAILED")
		default:
			stateBadge = styles.PillNormal.Render("⚪ STOPPED")
		}

		autostart := "[Manual]"
		if s.Enabled {
			autostart = "[AutoStart]"
		}

		line := fmt.Sprintf("%-26s  %s  %-12s", s.Name, stateBadge, autostart)
		if i == selectedIndex {
			listItems = append(listItems, styles.RowSelectedStyle.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, styles.RowNormalStyle.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	listPanel := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left, listItems...),
	)

	// Detail & Action Panel
	var detailPanel string
	if selectedIndex >= 0 && selectedIndex < len(services) {
		s := services[selectedIndex]
		detailContent := fmt.Sprintf(
			"📌 Layanan: %s\n📖 Deskripsi: %s\n\n📊 Status: %s (%s)\n🚀 Autostart: %t\n💾 Memori Digunakan: %.1f MB\n\n⚡ Tombol Kontrol Cepat:\n  [S] Nyalakan (Start)\n  [T] Hentikan (Stop)\n  [R] Mulai Ulang (Restart)\n  [E] Toggle Autostart (Enable/Disable)\n  [L] Lihat Cuplikan Log (Journalctl)",
			s.Name, s.Description,
			s.ActiveState, s.SubState,
			s.Enabled,
			s.MemoryMB,
		)
		detailPanel = styles.CardActiveStyle.Width(contentWidth/2 - 2).Render(detailContent)
	}

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)

	// Senior Insight Box
	insightHeader := styles.InsightTitleStyle.Render(t.SeniorInsightTitle)
	insightBody := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0f2fe")).Render(t.InsightServices)
	insightBox := styles.InsightBoxStyle.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left, insightHeader, insightBody),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		mainRow,
		"",
		insightBox,
	)
}
