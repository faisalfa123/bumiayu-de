package views

import (
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/bumiayu/bumiayu-de/internal/ui/wallpaper"
	"github.com/charmbracelet/lipgloss"
)

type WMMode int

const (
	WMFloating WMMode = iota // Window Mengambang (Center Pop-up)
	WMTiling                 // Tiling 2 Kolom
	WMCardGrid               // Grid Kartu Minimalis
)

// RenderSettings renders the customization panel with Wallpaper, Theme, and WM settings
func RenderSettings(currentTheme styles.ThemeMode, currentWM WMMode, wpType wallpaper.WallpaperType, wpPreset wallpaper.PresetScene, selectedRow int, imageInput string, st styles.ThemeStyles, width, height int) string {
	contentWidth := width - 6
	if contentWidth < 50 {
		contentWidth = 50
	}

	header := st.Title.Render("⚙️ PENGATURAN TAMPILAN, WINDOW MANAGER & WALLPAPER")
	sub := st.Subtitle.Render("Sesuaikan gaya visual, tata letak jendela, dan karya seni wallpaper terminal Anda.")

	// 1. Theme Selector
	themeNames := []string{"☀️ Terang Bersih (Light Mode)", "🌾 Hangat Lembut (Warm Ivory)", "🌙 Gelap Malam (Midnight Dark)"}
	var themeItems []string
	for i, name := range themeNames {
		cursor := "  "
		style := st.RowNormal
		if styles.ThemeMode(i) == currentTheme {
			name += "  [AKTIF]"
		}
		if selectedRow == i {
			cursor = "▶ "
			style = st.RowSelected
		}
		themeItems = append(themeItems, style.Render(cursor+name))
	}
	themeBox := st.Card.Width(contentWidth/3 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("🎨 Tema Warna:"),
			"",
			strings.Join(themeItems, "\n"),
		),
	)

	// 2. WM Mode Selector
	wmNames := []string{"🪟 Floating Window", "🔲 Tiling 2-Kolom", "📋 Card Grid"}
	var wmItems []string
	for i, name := range wmNames {
		cursor := "  "
		style := st.RowNormal
		if WMMode(i) == currentWM {
			name += "  [AKTIF]"
		}
		if selectedRow == i+3 {
			cursor = "▶ "
			style = st.RowSelected
		}
		wmItems = append(wmItems, style.Render(cursor+name))
	}
	wmBox := st.Card.Width(contentWidth/3 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("🖥️ Model Window Manager:"),
			"",
			strings.Join(wmItems, "\n"),
		),
	)

	// 3. Wallpaper Selector
	wpNames := []string{
		"🌾 Animasi Hidup: Sawah & Petani Bumiayu",
		"🌄 Seni Bumiayu: Gunung Slamet",
		"🏮 Seni Bumiayu: Lentera Senja",
		"🌌 Seni Bumiayu: Malam Bintang",
		"⏰ Dinamis Waktu (Day/Night)",
		"🖼️ Generator Foto ke Pixel Art",
	}
	var wpItems []string
	for i, name := range wpNames {
		cursor := "  "
		style := st.RowNormal
		if selectedRow == i+6 {
			cursor = "▶ "
			style = st.RowSelected
		}
		wpItems = append(wpItems, style.Render(cursor+name))
	}
	wpBox := st.Card.Width(contentWidth/3 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("🖼️ Mode Wallpaper:"),
			"",
			strings.Join(wpItems, "\n"),
		),
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, themeBox, " ", wmBox, " ", wpBox)

	// 4. Live Wallpaper Art Preview Box (Indah & Artistik)
	activeArtwork := wallpaper.GetAnimatedPaddyFarmer(0)
	if selectedRow >= 6 {
		wpIdx := selectedRow - 6
		switch wpIdx {
		case 1:
			activeArtwork = wallpaper.GetCurrentWallpaper(wallpaper.WallpaperPreset, wallpaper.SceneGunungSlamet)
		case 2:
			activeArtwork = wallpaper.GetCurrentWallpaper(wallpaper.WallpaperPreset, wallpaper.SceneLenteraSenja)
		case 3:
			activeArtwork = wallpaper.GetCurrentWallpaper(wallpaper.WallpaperPreset, wallpaper.SceneMalamBintang)
		case 4:
			activeArtwork = wallpaper.GetCurrentWallpaper(wallpaper.WallpaperDynamicTime, 0)
		case 5:
			activeArtwork = "🖼️ Tekan [Enter] atau [G] untuk memasukkan gambar JPG/PNG dan mengubahnya ke Pixel Art ANSI!"
		}
	}

	previewCard := st.CardActive.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Badge.Render("🖼️ PRATINJAU WALLPAPER TERMINAL AKTIF:"),
			"",
			activeArtwork,
			"",
			st.RowNormal.Render("⌨️ Tombol Navigasi: [↑/↓] Pilih Menu | [Enter/Space] Terapkan | [G] Generator Foto | [1-3] Pindah Tab"),
		),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		sub,
		"",
		topRow,
		"",
		previewCard,
	)
}
