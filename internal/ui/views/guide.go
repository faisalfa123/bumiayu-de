package views

import (
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderGuide renders user manual and practical cheat sheet
func RenderGuide(st styles.ThemeStyles, width, height int) string {
	contentWidth := width - 6
	if contentWidth < 50 {
		contentWidth = 50
	}

	header := st.Title.Render("📖 PETUNJUK PENGGUNAAN & PANDUAN LENGKAP BUMIAYU TUI")
	sub := st.Subtitle.Render("Dibuat agar siapa saja dapat mengelola server Linux dan ODROID-C2 dengan mudah dan percaya diri.")

	// Card 1: Navigasi
	navText := `• Tombol Angka [1 - 4]: Pindah tab secara instan.
• Tombol Panah [← / →]: Pindah tab ke kiri atau kanan.
• Tombol Panah [↑ / ↓] atau [j / k]: Memilih item dalam daftar.
• Tombol [Tab] / [Shift+Tab]: Berpindah fokus antar panel.
• Tombol [q] atau [Ctrl+C]: Keluar dari Bumiayu TUI dengan aman.`

	navCard := st.Card.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("🎮 1. Cara Bernavigasi"),
			"",
			navText,
		),
	)

	// Card 2: Install & Uninstall Aplikasi
	appText := `• Masuk ke tab [2] Toko Aplikasi.
• Pilih aplikasi (misal: Docker, Nginx, WireGuard, Samba).
• Tekan [Enter] atau [I] untuk Memasang (Install).
• Tekan [X] untuk Menghapus (Uninstall).
• Jendela konfirmasi akan muncul. Tekan [Y] untuk Ya atau [N] untuk Batal.
• Proses instalasi dan log akan tampil secara transparan.`

	appCard := st.Card.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("📦 2. Pasang & Hapus Aplikasi"),
			"",
			appText,
		),
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, navCard, "  ", appCard)

	// Card 3: Pengaturan Tampilan & Tips Hemat RAM
	tipsText := `• Mengubah Tampilan: Masuk ke tab [3] Pengaturan untuk memilih model tampilan Window Manager (Floating / Tiling / Card Grid) dan memilih tema (Terang / Hangat / Gelap).
• Ganti Tema Cepat: Tekan tombol [T] kapan saja untuk beralih antara Terang dan Gelap.
• Ganti Bahasa: Tekan tombol [L] kapan saja untuk beralih antara Bahasa Indonesia dan English.
• Pencegahan OOM (Out Of Memory) di ODROID-C2: Bumiayu TUI hanya memakan ~12 MB RAM. Pastikan server memiliki Swap minimal 1GB agar proses berat seperti Docker tidak menyebabkan sistem freeze.`

	tipsCard := st.CardActive.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("💡 3. Kustomisasi Tampilan & Tips Server"),
			"",
			tipsText,
		),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		sub,
		"",
		topRow,
		"",
		tipsCard,
	)
}
