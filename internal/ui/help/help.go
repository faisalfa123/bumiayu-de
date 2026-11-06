package help

import (
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderQuickBar renders the persistent shortcuts help bar at the bottom
func RenderQuickBar(st styles.ThemeStyles, width int, statusMessage string, zoomActive bool) string {
	zoomStatus := "[f]: Zoom"
	if zoomActive {
		zoomStatus = "[f/Esc]: Unzoom"
	}

	btnKey := func(key, desc string) string {
		k := lipgloss.NewStyle().Bold(true).Foreground(st.PrimaryColor).Render("[" + key + "]")
		d := lipgloss.NewStyle().Foreground(st.TextColor).Render(" " + desc)
		return k + d
	}

	keys := []string{
		btnKey("Space", "Launcher"),
		btnKey("Tab/hjkl", "Nav"),
		btnKey(zoomStatus, ""),
		btnKey("1-4", "Panels"),
		btnKey("T", "Tema"),
		btnKey("?", "Bantuan"),
		btnKey("q", "Keluar"),
	}

	content := strings.Join(keys, "  ")

	if statusMessage != "" {
		notif := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(st.AccentColor).
			Padding(0, 1).
			Render("ℹ " + statusMessage)
		content = lipgloss.JoinHorizontal(lipgloss.Center, notif, "  ", content)
	}

	return st.Footer.Width(width).Render(content)
}

// RenderArcadeHelpModal renders a full retro arcade style "HOW TO PLAY / SYSTEM CONTROLS" modal
func RenderArcadeHelpModal(st styles.ThemeStyles, width, height int) string {
	modalWidth := 74
	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(st.PrimaryColor).
		Padding(0, 2).
		Render("🕹️  BUMIAYU TUI — HOW TO PLAY & CONTROLS  🕹️")

	secHeader := func(s string) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(st.AccentColor).
			Render("▶ " + s)
	}

	row := func(key, desc string) string {
		k := lipgloss.NewStyle().
			Bold(true).
			Foreground(st.PrimaryColor).
			Width(22).
			Render("  [" + key + "]")
		d := lipgloss.NewStyle().
			Foreground(st.TextColor).
			Render(desc)
		return k + " " + d
	}

	sec1 := lipgloss.JoinVertical(lipgloss.Left,
		secHeader("NAVIGASI & WINDOW MANAGER"),
		row("Tab / Shift+Tab", "Pindah fokus antar panel jendela"),
		row("h / j / k / l", "Navigasi arah (Vim-style)"),
		row("↑ / ↓ / ← / →", "Navigasi arah kursor panel aktif"),
		row("f", "Toggle Fullscreen / Zoom panel terpilih"),
		row("1 / 2 / 3 / 4", "Langsung lompat ke Panel 1 s/d 4"),
	)

	sec2 := lipgloss.JoinVertical(lipgloss.Left,
		secHeader("UNIFIED SPOTLIGHT LAUNCHER"),
		row("Space / /", "Buka bar pencarian cepat (Spotlight)"),
		row("Tab / →", "Terima tebakan teks (Ghost-Text Autocomplete)"),
		row("Enter", "Eksekusi perintah / buka aksi terpilih"),
		row(">clean", "Macro cepat: Bersihkan Memory Cache RAM"),
		row(">ports", "Macro cepat: Audit port terbuka publik"),
		row(">ufw", "Macro cepat: Nyala/Matikan Dinding Api"),
	)

	sec3 := lipgloss.JoinVertical(lipgloss.Left,
		secHeader("KONTROL SISTEM & ESTETIKA"),
		row("S / T / R", "Start / Stop / Restart Layanan Systemd"),
		row("K", "Hentikan Paksa Proses (Kill Process SIGTERM)"),
		row("U", "Nyalakan / Matikan Firewall UFW"),
		row("T", "Ganti Tema (Cyber-Arcade ⇄ DOS ⇄ GameBoy)"),
		row("L", "Ganti Bahasa (Indonesia ⇄ English)"),
		row("? / Esc", "Tutup jendela bantuan ini / Batalkan"),
	)

	footerHint := lipgloss.NewStyle().
		Italic(true).
		Foreground(st.MutedColor).
		Render("Tekan [?] atau [Esc] untuk kembali ke layar Window Manager")

	body := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		sec1,
		"",
		sec2,
		"",
		sec3,
		"",
		footerHint,
	)

	return st.Modal.Width(modalWidth).Render(body)
}
