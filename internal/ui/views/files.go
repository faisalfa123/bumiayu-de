package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderFiles renders the file browser and config inspector
func RenderFiles(currentPath string, entries []system.FileEntry, selectedIndex int, previewText string, lang i18n.Language, width, height int) string {
	t := i18n.Get(lang)
	contentWidth := width - 6
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := styles.TitleStyle.Render("📂 Penjelajah Berkas & Inspektur Konfigurasi: " + currentPath)

	var listItems []string
	for i, f := range entries {
		icon := "📄"
		if f.IsDir {
			icon = "📁"
		}

		sizeStr := fmt.Sprintf("%d B", f.Size)
		if f.IsDir {
			sizeStr = "<DIR>"
		}

		line := fmt.Sprintf("%s %-22s  %-10s  %-8s", icon, f.Name, f.Permissions, sizeStr)
		if i == selectedIndex {
			listItems = append(listItems, styles.RowSelectedStyle.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, styles.RowNormalStyle.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	if len(listItems) == 0 {
		listItems = append(listItems, styles.SubtitleStyle.Render("Direktori kosong atau tidak ada akses."))
	}

	// Limit list items to visible height
	if len(listItems) > 16 {
		listItems = listItems[:16]
	}

	listPanel := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left, listItems...),
	)

	// Preview Panel
	previewTitle := styles.BadgeStyle.Render("👁️ Pratinjau Berkas")
	if previewText == "" {
		previewText = "Pilih berkas teks atau konfigurasi untuk melihat isinya."
	}
	previewPanel := styles.CardActiveStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			previewTitle,
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#cbd5e1")).Render(previewText),
		),
	)

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", previewPanel)

	// Senior Insight Box
	insightHeader := styles.InsightTitleStyle.Render(t.SeniorInsightTitle)
	insightBody := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0f2fe")).Render(t.InsightFiles)
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
