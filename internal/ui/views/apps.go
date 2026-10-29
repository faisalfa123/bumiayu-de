package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderApps renders the package catalog and 1-click app hub
func RenderApps(catalog []system.AppCatalogItem, selectedIndex int, pkgCount int, lang i18n.Language, width, height int) string {
	t := i18n.Get(lang)
	contentWidth := width - 6
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := fmt.Sprintf("📦 Katalog Aplikasi & Paket Sistem (Terdeteksi %d paket terpasang via APT/dpkg)", pkgCount)
	headerView := styles.TitleStyle.Render(header)

	var listItems []string
	for i, app := range catalog {
		statusBadge := styles.PillNormal.Render("⚪ Belum Terpasang")
		if app.Installed {
			statusBadge = styles.BadgeStyle.Render("🟢 Terpasang")
		}

		line := fmt.Sprintf("%-20s  %-18s  %s", app.Name, "["+app.Category+"]", statusBadge)
		if i == selectedIndex {
			listItems = append(listItems, styles.RowSelectedStyle.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, styles.RowNormalStyle.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	listPanel := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left, listItems...),
	)

	// Detail Panel for currently selected app
	var detailPanel string
	if selectedIndex >= 0 && selectedIndex < len(catalog) {
		app := catalog[selectedIndex]
		desc := app.Description
		if lang == i18n.EN && app.DescriptionEN != "" {
			desc = app.DescriptionEN
		}

		statusText := t.NotInstalled
		actionText := t.ActionInstall
		if app.Installed {
			statusText = t.Installed
			actionText = t.ActionRemove
		}

		detailContent := fmt.Sprintf(
			"📌 Aplikasi: %s\n🏷️ Kategori: %s\n📡 Port Layanan: %s\n🏷️ Nama Paket: %s\n📊 Status: %s\n\n📖 Deskripsi:\n%s\n\n💡 Rekomendasi Senior:\n%s\n\n⚡ Aksi Tersedia: %s",
			app.Name, app.Category, app.Port, app.PackageName, statusText,
			desc,
			app.Insight,
			actionText,
		)

		detailPanel = styles.CardActiveStyle.Width(contentWidth/2 - 2).Render(detailContent)
	}

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)

	// Senior Insight Box
	insightHeader := styles.InsightTitleStyle.Render(t.SeniorInsightTitle)
	insightBody := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0f2fe")).Render(t.InsightApps)
	insightBox := styles.InsightBoxStyle.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left, insightHeader, insightBody),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerView,
		"",
		mainRow,
		"",
		insightBox,
	)
}
