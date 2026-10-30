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

// CategoryFilterList holds all navigation categories
var CategoryFilterList = []string{
	"Semua",
	"Foto & Cloud",
	"Media",
	"VPN & Jaringan",
	"Web & Server",
	"Database",
	"Monitoring & Alat",
}

// RenderDashboard renders the All-in-One Desktop Launcher with live wallpaper art & direct app actions
func RenderDashboard(info system.HardwareInfo, catalog []system.AppCatalogItem, selectedApp int, activeCategory int, dialog AppDialogState, wm WMMode, wpType wallpaper.WallpaperType, wpPreset wallpaper.PresetScene, st styles.ThemeStyles, lang i18n.Language, width, height int) string {
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

	// Hardware Telemetry Banner (Lengkap & Profesional)
	deviceType := "Linux Server (" + info.Arch + ")"
	if info.IsOdroidC2 {
		deviceType = "ODROID-C2 (Amlogic S905 Quad-Core ARM64)"
	}
	ramBar := system.GetMemoryBar(info.MemUsagePct, 14)

	healthBadge := st.Badge.Render("🟢 Normal & Sehat")
	if info.MemUsagePct > 80 || info.TemperatureC > 70 {
		healthBadge = st.PillDanger.Render("🔴 Beban Tinggi!")
	}

	distro := system.DetectDistro()
	sysCol1 := fmt.Sprintf("💻 Perangkat: %s [%s]\n🌐 Hostname : %s | Uptime: %s", deviceType, distro, info.Hostname, info.UptimeStr)
	sysCol2 := fmt.Sprintf("🌡️ Suhu CPU: %.1f°C | Load Avg: %.2f, %.2f, %.2f\n📊 Status  : %s", info.TemperatureC, info.Load1, info.Load5, info.Load15, healthBadge)
	sysCol3 := fmt.Sprintf("🧠 RAM : %d/%d MB (%s %.1f%%)\n💾 Swap: %d/%d MB (%.1f%%)", info.MemUsedMB, info.MemTotalMB, ramBar, info.MemUsagePct, info.SwapUsedMB, info.SwapTotalMB, info.SwapUsagePct)

	colW := (contentWidth - 6) / 3
	infoBanner := st.Card.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("📊 MONITORING SISTEM & PERANGKAT KERAS REAL-TIME:"),
			"",
			lipgloss.JoinHorizontal(lipgloss.Top,
				lipgloss.NewStyle().Width(colW).Render(sysCol1),
				"  ",
				lipgloss.NewStyle().Width(colW).Render(sysCol2),
				"  ",
				lipgloss.NewStyle().Width(colW).Render(sysCol3),
			),
		),
	)

	// Category Filter Pill Bar
	var filterPills []string
	for idx, cat := range CategoryFilterList {
		if idx == activeCategory {
			filterPills = append(filterPills, st.ButtonActive.Render(fmt.Sprintf("[%s]", cat)))
		} else {
			filterPills = append(filterPills, st.ButtonNormal.Render(fmt.Sprintf(" %s ", cat)))
		}
	}
	filterBar := lipgloss.JoinHorizontal(lipgloss.Center, filterPills...)
	filterContainer := st.Card.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("📂 Filter Kategori [Tekan Tab / Shift+Tab]:"),
			filterBar,
		),
	)

	// Filter catalog based on active category
	type CatalogEntry struct {
		OriginalIndex int
		Item          system.AppCatalogItem
	}
	var filteredCatalog []CatalogEntry
	selectedCatName := CategoryFilterList[activeCategory]
	for i, app := range catalog {
		if activeCategory == 0 || app.Category == selectedCatName {
			filteredCatalog = append(filteredCatalog, CatalogEntry{OriginalIndex: i, Item: app})
		}
	}

	// App List Panel (With Category Badges)
	var listItems []string
	for i, entry := range filteredCatalog {
		app := entry.Item
		badge := st.Pill.Render("⚪ Belum Ada")
		if app.Installed {
			badge = st.Badge.Render("🟢 Terpasang")
		}

		line := fmt.Sprintf("%-22s  %-16s  %s", app.Name, "["+app.Category+"]", badge)
		if i == selectedApp {
			listItems = append(listItems, st.RowSelected.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, st.RowNormal.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	// Limit visible app items to terminal height with smooth scrolling
	maxVisible := 13
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
			st.Title.Render(fmt.Sprintf("📦 Katalog Aplikasi (%d/%d Ditampilkan):", len(filteredCatalog), len(catalog))),
			st.Subtitle.Render("Gunakan [↑ / ↓] pilih, [Enter/I] pasang, [X] hapus"),
			"",
			strings.Join(visibleList, "\n"),
		),
	)

	// Selected App Detail Panel
	var detailPanel string
	if selectedApp >= 0 && selectedApp < len(filteredCatalog) {
		app := filteredCatalog[selectedApp].Item
		statusText := "Belum Terpasang di Server"
		actionBtn := st.ButtonActive.Render("👉 Tekan [Enter] atau [I] untuk Pasang Sekarang")
		if app.Installed {
			statusText = "Sudah Terpasang & Siap Digunakan"
			actionBtn = st.PillDanger.Render("🗑️ Tekan [X] untuk Menghapus (Uninstall)")
		}

		distroInstallMethod := "dietpi-software install " + app.DietPiID
		if distro == system.DistroArch {
			distroInstallMethod = "pacman -Sy " + app.ArchPkg
		} else if distro == system.DistroUbuntu || distro == system.DistroDebian {
			distroInstallMethod = "apt-get install -y " + app.PackageName
		}

		thirdPartyBadge := ""
		if app.IsThirdParty {
			thirdPartyBadge = "\n🛡️ Tipe: Vendor Pihak Ketiga Resmi Terverifikasi"
		}

		detailContent := fmt.Sprintf(
			"📌 Aplikasi: %s\n🏷️ Kategori: %s\n📡 Port Layanan: %s\n⚙️ Perintah Sistem: %s\n📊 Status: %s%s\n\n📖 Deskripsi:\n%s\n\n💡 Rekomendasi Senior:\n%s\n\n⚡ Aksi Cepat:\n%s",
			app.Name, app.Category, app.Port, distroInstallMethod, statusText, thirdPartyBadge,
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
			st.RowNormal.Render("⌨️ Navigasi: [Tab/Shift+Tab] Ganti Kategori | [↑/↓] Pilih App | [Enter/I] Pasang | [X] Hapus | [1-3] Tab | [T] Tema | [M] WM | [W] Setup | [q] Keluar"),
		),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		sub,
		"",
		infoBanner,
		"",
		filterContainer,
		"",
		mainRow,
		"",
		artBanner,
	)
}
