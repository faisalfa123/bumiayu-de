package views

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderDashboard renders the All-in-One Desktop Launcher with live hardware telemetry & direct app actions
func RenderDashboard(info system.HardwareInfo, catalog []system.AppCatalogItem, selectedApp int, dialog AppDialogState, wm WMMode, st styles.ThemeStyles, lang i18n.Language, width, height int) string {
	contentWidth := width - 6
	if contentWidth < 50 {
		contentWidth = 50
	}

	// 1. If Modal is Active (Install/Uninstall Confirmation or Log Output), render modal directly!
	if dialog.Active {
		return renderAppModal(dialog, st, contentWidth)
	}

	header := st.Title.Render("🏝️ BUMIAYU DESKTOP LAUNCHER")
	sub := st.Subtitle.Render("Pusat Kontrol Sistem Ringan & Peluncur Aplikasi Cepat untuk ODROID-C2 & Server Linux")

	// Hardware Telemetry Banner
	deviceType := "Linux Standard Server (" + info.Arch + ")"
	if info.IsOdroidC2 {
		deviceType = "ODROID-C2 (Amlogic S905 Quad-Core ARM64)"
	}
	ramBar := system.GetMemoryBar(info.MemUsagePct, 16)

	healthBadge := st.Badge.Render("🟢 Normal & Sehat")
	if info.MemUsagePct > 80 || info.TemperatureC > 70 {
		healthBadge = st.PillDanger.Render("🔴 Beban Tinggi!")
	}

	infoBanner := st.Card.Width(contentWidth).Render(
		fmt.Sprintf("💻 %s  |  Host: %s  |  Suhu: %.1f°C  |  RAM: %d/%d MB %s  |  Status: %s",
			deviceType, info.Hostname, info.TemperatureC, info.MemUsedMB, info.MemTotalMB, ramBar, healthBadge),
	)

	// App List Panel
	var listItems []string
	for i, app := range catalog {
		badge := st.Pill.Render("⚪ Belum Ada")
		if app.Installed {
			badge = st.Badge.Render("🟢 Terpasang")
		}

		line := fmt.Sprintf("%-18s  %-14s  %s", app.Name, "["+app.Category+"]", badge)
		if i == selectedApp {
			listItems = append(listItems, st.RowSelected.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, st.RowNormal.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	listPanel := st.Card.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("📦 Pilih Aplikasi untuk Dikelola:"),
			st.Subtitle.Render("Gunakan [↑ / ↓] untuk memilih, [Enter/I] pasang, [X] hapus"),
			"",
			strings.Join(listItems, "\n"),
		),
	)

	// Selected App Detail Panel
	var detailPanel string
	if selectedApp >= 0 && selectedApp < len(catalog) {
		app := catalog[selectedApp]
		statusText := "Belum Terpasang di Server"
		actionBtn := st.ButtonActive.Render("👉 Tekan [Enter] atau [I] untuk Pasang Sekarang")
		if app.Installed {
			statusText = "Sudah Terpasang & Siap Digunakan"
			actionBtn = st.PillDanger.Render("🗑️ Tekan [X] untuk Menghapus (Uninstall)")
		}

		detailContent := fmt.Sprintf(
			"📌 Aplikasi: %s\n🏷️ Kategori: %s\n📡 Port Layanan: %s\n📊 Status: %s\n\n📖 Deskripsi:\n%s\n\n💡 Rekomendasi Senior:\n%s\n\n⚡ Aksi Cepat:\n%s",
			app.Name, app.Category, app.Port, statusText,
			app.Description,
			app.Insight,
			actionBtn,
		)
		detailPanel = st.CardActive.Width(contentWidth/2 - 2).Render(detailContent)
	}

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)

	// Bottom action hint bar
	hintBar := st.Card.Width(contentWidth).Render(
		st.RowNormal.Render("⌨️ Langsung Aksi: [↑/↓] Pilih Aplikasi | [Enter/I] Pasang | [X] Hapus | [T] Ganti Tema | [M] Ganti Model Tampilan | [W] Wizard"),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		sub,
		"",
		infoBanner,
		"",
		mainRow,
		"",
		hintBar,
	)
}
