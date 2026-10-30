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
func RenderDashboard(info system.HardwareInfo, catalog []system.AppCatalogItem, selectedApp int, activeCategory int, animFrame int, dialog AppDialogState, wm WMMode, wpType wallpaper.WallpaperType, wpPreset wallpaper.PresetScene, st styles.ThemeStyles, lang i18n.Language, width, height int) string {
	contentWidth := width - 4
	if contentWidth < 50 {
		contentWidth = 50
	}

	// 1. If Modal is Active (Install/Uninstall Confirmation or Log Output), render modal directly!
	if dialog.Active {
		return renderAppModal(dialog, st, contentWidth)
	}

	// Header Bar Desktop Elegan (Minimalis WebUI)
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		st.Title.Render("🏝️ BUMIAYU DESKTOP"),
		"  ",
		st.Subtitle.Render("Pusat Kendali Server Linux & ODROID-C2"),
	)

	// Hardware Telemetry Banner (Lengkap & Profesional)
	deviceType := "Linux Server (" + info.Arch + ")"
	if info.IsOdroidC2 {
		deviceType = "ODROID-C2 (Amlogic Quad-Core)"
	}
	ramBar := system.GetMemoryBar(info.MemUsagePct, 10)

	healthBadge := st.Badge.Render("🟢 Normal")
	if info.MemUsagePct > 80 || info.TemperatureC > 70 {
		healthBadge = st.PillDanger.Render("🔴 Beban Tinggi")
	}

	distro := system.DetectDistro()
	sysCol1 := fmt.Sprintf("💻 Perangkat: %s\n🌐 Hostname : %s | Distro: %s", deviceType, info.Hostname, distro)
	sysCol2 := fmt.Sprintf("🌡️ Suhu CPU: %.1f°C | Load: %.2f\n📊 Status  : %s", info.TemperatureC, info.Load1, healthBadge)
	sysCol3 := fmt.Sprintf("🧠 RAM : %d/%d MB (%s %.0f%%)\n💾 Swap: %d/%d MB", info.MemUsedMB, info.MemTotalMB, ramBar, info.MemUsagePct, info.SwapUsedMB, info.SwapTotalMB)

	colW := (contentWidth - 6) / 3
	infoBanner := st.Card.Width(contentWidth).Render(
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(colW).Render(sysCol1),
			"  ",
			lipgloss.NewStyle().Width(colW).Render(sysCol2),
			"  ",
			lipgloss.NewStyle().Width(colW).Render(sysCol3),
		),
	)

	// Category Filter Pill Bar (Desain Bersih & Elegan WebUI)
	var filterPills []string
	for idx, cat := range CategoryFilterList {
		if idx == activeCategory {
			filterPills = append(filterPills, st.ButtonActive.Render(cat))
		} else {
			filterPills = append(filterPills, st.ButtonNormal.Render(cat))
		}
	}
	filterBar := lipgloss.JoinHorizontal(lipgloss.Center, filterPills...)
	filterContainer := st.Card.Width(contentWidth).Render(
		lipgloss.JoinHorizontal(lipgloss.Center,
			st.Title.Render("📂 Filter: "),
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

		line := fmt.Sprintf("%-20s  %-15s  %s", app.Name, "["+app.Category+"]", badge)
		if i == selectedApp {
			listItems = append(listItems, st.RowSelected.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, st.RowNormal.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	// Limit visible app items to terminal height with smooth scrolling
	maxVisible := 6
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
			st.Title.Render(fmt.Sprintf("📦 Katalog Server (%d/%d):", len(filteredCatalog), len(catalog))),
			"",
			strings.Join(visibleList, "\n"),
		),
	)

	// Selected App Detail Panel (Window Modern & Cantik)
	var detailPanel string
	if selectedApp >= 0 && selectedApp < len(filteredCatalog) {
		app := filteredCatalog[selectedApp].Item
		statusText := "Belum Terpasang"
		actionBtn := st.ButtonActive.Render("👉 [Enter/I] Pasang")
		if app.Installed {
			statusText = "🟢 Siap Digunakan"
			actionBtn = st.PillDanger.Render("🗑️ [X] Hapus")
		}

		distroInstallMethod := "dietpi-software install " + app.DietPiID
		if distro == system.DistroArch {
			distroInstallMethod = "pacman -Sy " + app.ArchPkg
		} else if distro == system.DistroUbuntu || distro == system.DistroDebian {
			distroInstallMethod = "apt-get install -y " + app.PackageName
		}

		detailContent := fmt.Sprintf(
			"📌 %s [%s]\n📡 Port: %s | Status: %s\n⚙️ Command: %s\n\n📖 Deskripsi:\n%s\n\n⚡ Aksi:\n%s",
			app.Name, app.Category, app.Port, statusText, distroInstallMethod,
			app.Description,
			actionBtn,
		)
		detailPanel = st.CardActive.Width(contentWidth/2 - 2).Render(detailContent)
	}

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)

	// Living Pixel Art Animation Banner at Bottom (Desain Indah & Artistik WebUI)
	activeArtwork := wallpaper.GetAnimatedPaddyFarmer(animFrame)
	artBanner := st.Card.Width(contentWidth).Render(
		activeArtwork,
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		infoBanner,
		filterContainer,
		mainRow,
		artBanner,
	)
}
