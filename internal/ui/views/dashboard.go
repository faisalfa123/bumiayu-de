package views

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/bumiayu/bumiayu-de/internal/ui/wallpaper"
	"github.com/charmbracelet/lipgloss"
)

// RenderDashboard renders the All-in-One Desktop Launcher with live wallpaper art & direct app actions
func RenderDashboard(info system.HardwareInfo, catalog []system.AppCatalogItem, selectedApp int, dialog AppDialogState, wm WMMode, wpType wallpaper.WallpaperType, wpPreset wallpaper.PresetScene, st styles.ThemeStyles, lang i18n.Language, width, height int) string {
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

	distro := system.DetectDistro()
	infoBanner := st.Card.Width(contentWidth).Render(
		fmt.Sprintf("💻 %s [%s]  |  Host: %s  |  Suhu: %.1f°C  |  RAM: %d/%d MB %s  |  Status: %s",
			deviceType, distro, info.Hostname, info.TemperatureC, info.MemUsedMB, info.MemTotalMB, ramBar, healthBadge),
	)

	// App List Panel (With Category Badges)
	var listItems []string
	for i, app := range catalog {
		badge := st.Pill.Render("⚪ Belum Ada")
		if app.Installed {
			badge = st.Badge.Render("🟢 Terpasang")
		}

		line := fmt.Sprintf("%-20s  %-18s  %s", app.Name, "["+app.Category+"]", badge)
		if i == selectedApp {
			listItems = append(listItems, st.RowSelected.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, st.RowNormal.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	// Limit visible app items to terminal height with smooth scrolling
	maxVisible := 14
	startIdx := 0
	if selectedApp >= maxVisible {
		startIdx = selectedApp - maxVisible + 1
	}
	endIdx := startIdx + maxVisible
	if endIdx > len(listItems) {
		endIdx = len(listItems)
	}
	visibleList := listItems[startIdx:endIdx]

	listPanel := st.Card.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render(fmt.Sprintf("📦 Katalog Aplikasi (%d Aplikasi Tersedia):", len(catalog))),
			st.Subtitle.Render("Gunakan [↑ / ↓] untuk memilih, [Enter/I] pasang, [X] hapus"),
			"",
			strings.Join(visibleList, "\n"),
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

		thirdPartyBadge := ""
		if app.IsThirdParty {
			thirdPartyBadge = "\n🛡️ Tipe: Vendor Pihak Ketiga Terpercaya (Official Verified)"
		}

		detailContent := fmt.Sprintf(
			"📌 Aplikasi: %s\n🏷️ Kategori: %s\n📡 Port Layanan: %s\n📊 Status: %s%s\n\n📖 Deskripsi:\n%s\n\n💡 Rekomendasi Senior:\n%s\n\n⚡ Aksi Cepat:\n%s",
			app.Name, app.Category, app.Port, statusText, thirdPartyBadge,
			app.Description,
			app.Insight,
			actionBtn,
		)
		detailPanel = st.CardActive.Width(contentWidth/2 - 2).Render(detailContent)
	}

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)

	// Mini Wallpaper Artwork Banner at Bottom
	activeArtwork := wallpaper.GetCurrentWallpaper(wpType, wpPreset)
	artBanner := st.Card.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			activeArtwork,
			"",
			st.RowNormal.Render("⌨️ Aksi Cepat: [↑/↓] Pilih App | [Enter/I] Pasang | [X] Hapus | [T] Tema | [M] Model WM | [W] Wallpaper/Wizard | [q] Keluar"),
		),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		sub,
		"",
		infoBanner,
		"",
		mainRow,
		"",
		artBanner,
	)
}
