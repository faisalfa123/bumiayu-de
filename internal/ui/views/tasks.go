package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderTasks renders the process manager and Anti-OOM monitor
func RenderTasks(procs []system.ProcessItem, selectedIndex int, lang i18n.Language, width, height int) string {
	t := i18n.Get(lang)
	contentWidth := width - 6
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := styles.TitleStyle.Render("⚡ Pengawas Proses Sistem & Anti-OOM (Out Of Memory) Monitor")

	var listItems []string
	for i, p := range procs {
		var oomBadge string
		if p.OOMScore > 300 {
			oomBadge = styles.PillDanger.Render(fmt.Sprintf("OOM: %d", p.OOMScore))
		} else if p.OOMScore > 100 {
			oomBadge = styles.PillWarning.Render(fmt.Sprintf("OOM: %d", p.OOMScore))
		} else {
			oomBadge = styles.PillNormal.Render(fmt.Sprintf("OOM: %d", p.OOMScore))
		}

		line := fmt.Sprintf("%-6d %-18s  %6.1f MB  %s", p.PID, p.Name, p.MemMB, oomBadge)
		if i == selectedIndex {
			listItems = append(listItems, styles.RowSelectedStyle.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, styles.RowNormalStyle.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	if len(listItems) > 16 {
		listItems = listItems[:16]
	}

	listPanel := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left, listItems...),
	)

	// Detail & Kill Panel
	var detailPanel string
	if selectedIndex >= 0 && selectedIndex < len(procs) {
		p := procs[selectedIndex]
		oomWarning := "Risiko OOM rendah. Proses ini aman berjalan normal."
		if p.OOMScore > 200 {
			oomWarning = "⚠️ PERINGATAN: Nilai OOM Score tinggi! Jika server kehabisan RAM, Linux Kernel akan mematikan proses ini secara paksa terlebih dahulu."
		}

		detailContent := fmt.Sprintf(
			"⚡ Proses: %s (PID: %d)\n👤 Pemilik: %s | Status: %s\n💾 Memori Fisik (RSS): %.1f MB\n🧠 OOM Score Kernel: %d / 1000\n\n🛡️ Analisis Risiko:\n%s\n\n⚡ Aksi Kontrol Cepat:\n  [K] Hentikan Paksa (Kill Process - SIGTERM)\n  [F] Hentikan Instan (Force Kill - SIGKILL)",
			p.Name, p.PID,
			p.User, p.State,
			p.MemMB,
			p.OOMScore,
			oomWarning,
		)
		detailPanel = styles.CardActiveStyle.Width(contentWidth/2 - 2).Render(detailContent)
	}

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)

	// Senior Insight Box
	insightHeader := styles.InsightTitleStyle.Render(t.SeniorInsightTitle)
	insightBody := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0f2fe")).Render(t.InsightTasks)
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
