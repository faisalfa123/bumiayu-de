package views

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

type WMMode int

const (
	WMFloating WMMode = iota // Window Mengambang (Center Pop-up)
	WMTiling                 // Tiling 2 Kolom
	WMCardGrid               // Grid Kartu Minimalis
)

// RenderSettings renders the customization panel with live preview of WM styles
func RenderSettings(currentTheme styles.ThemeMode, currentWM WMMode, selectedRow int, st styles.ThemeStyles, width, height int) string {
	contentWidth := width - 6
	if contentWidth < 50 {
		contentWidth = 50
	}

	header := st.Title.Render("⚙️ PENGATURAN TAMPILAN & WINDOW MANAGER")
	sub := st.Subtitle.Render("Ubah gaya visual desktop dan palet warna sesuai kenyamanan mata Anda.")

	// Options List
	themeNames := []string{"☀️ Terang Bersih (Light Mode - Rekomendasi)", "🌾 Hangat & Lembut (Warm Ivory)", "🌙 Gelap Malam (Midnight Dark)"}
	wmNames := []string{"🪟 Window Mengambang (Floating WM)", "🔲 Tiling 2-Kolom (Split WM)", "📋 Kartu Grid Sederhana (Card Grid)"}

	// 1. Theme Selector Box
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
	themeBox := st.Card.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("🎨 Palet Warna / Tema:"),
			st.Subtitle.Render("Gunakan [1-3] atau panah untuk memilih"),
			"",
			strings.Join(themeItems, "\n"),
		),
	)

	// 2. WM Model Selector Box
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
	wmBox := st.Card.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("🖥️ Model Window Manager (WM):"),
			st.Subtitle.Render("Gunakan [4-6] atau panah untuk memilih"),
			"",
			strings.Join(wmItems, "\n"),
		),
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, themeBox, "  ", wmBox)

	// 3. LIVE PREVIEW BOX
	var previewDiagram string
	var previewDesc string
	switch currentWM {
	case WMFloating:
		previewDiagram = `
 ┌────────────────────────────────────────────────────────┐
 │ 🏝️ BUMIAYU    [Host: odroid-c2]   [RAM: 24%]   12:00  │
 ├────────────────────────────────────────────────────────┤
 │                                                        │
 │          ┌──────────────────────────────┐              │
 │          │  📦 TOKO APLIKASI            │              │
 │          ├──────────────────────────────┤              │
 │          │  ▶ Docker Engine  [Pasang]   │              │
 │          │    Nginx Server   [Hapus]    │              │
 │          │    WireGuard VPN  [Pasang]   │              │
 │          └──────────────────────────────┘              │
 │                                                        │
 ├────────────────────────────────────────────────────────┤
 │ [1] Desktop  [2] Toko App  [3] Pengaturan  [4] Panduan │
 └────────────────────────────────────────────────────────┘`
		previewDesc = "Model Floating: Jendela aplikasi muncul anggun di tengah desktop dengan taskbar di bawah. Tampilan sangat bersih, rapi, dan mudah dimengerti."

	case WMTiling:
		previewDiagram = `
 ┌────────────────────────────────────────────────────────┐
 │ 🏝️ BUMIAYU DESKTOP          [RAM: 480MB] [Suhu: 44.5°C] │
 ├─────────────────────────┬──────────────────────────────┤
 │  MENU WORKSPACE         │  DETAIL & AKSI APLIKASI      │
 │  ─────────────────      │  ──────────────────────      │
 │  ▶ [1] Toko Aplikasi    │  Docker Container Engine     │
 │    [2] Status Sistem    │  Status: Belum Terpasang     │
 │    [3] Pengaturan       │  Tekan [Enter] untuk Pasang  │
 │    [4] Petunjuk         │  RAM Terpakai: ~15 MB        │
 └─────────────────────────┴──────────────────────────────┘`
		previewDesc = "Model Tiling: Layar terbagi dalam 2 kolom terstruktur. Sisi kiri untuk menu cepat dan sisi kanan untuk layar kerja utama."

	case WMCardGrid:
		previewDiagram = `
 ┌────────────────────────────────────────────────────────┐
 │ 🏝️ BUMIAYU DESKTOP GRID                                │
 ├────────────────────────────┬───────────────────────────┤
 │ 🏥 STATUS SISTEM           │ 📦 APLIKASI CEPAT         │
 │ • CPU: 12%  • RAM: 480 MB  │ • Docker: Belum Pasang    │
 │ • Suhu: 44.5°C (Normal)    │ • Nginx:  Aktif Berjalan  │
 ├────────────────────────────┼───────────────────────────┤
 │ 🛡️ KEAMANAN & FIREWALL     │ ⚡ AKSI SISTEM             │
 │ • UFW: Aktif Terlindungi   │ • [I] Pasang Aplikasi Baru│
 └────────────────────────────┴───────────────────────────┘`
		previewDesc = "Model Card Grid: Tampilan kartu-kartu ringkas modern dengan informasi kunci yang mudah dicerna tanpa membuat bingung."
	}

	previewCard := st.CardActive.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Badge.Render("👁️ PRATINJAU LANGSUNG (LIVE PREVIEW)"),
			fmt.Sprintf("\n%s\n", previewDiagram),
			st.Subtitle.Render(previewDesc),
			"",
			st.RowNormal.Render("⚡ Petunjuk: Tekan [Space] atau [Enter] untuk menerapkan pilihan, atau tekan [T] untuk ganti tema cepat."),
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
