package views

import (
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderGuide renders user manual, update instructions, and uninstaller guide
func RenderGuide(st styles.ThemeStyles, width, height int) string {
	contentWidth := width - 4
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
	appText := `• Masuk ke tab [1] Desktop & Aplikasi.
• Tersedia Immich SBC Optimizer resmi (ODROID-C2 & Server) terintegrasi!
• Pilih aplikasi (Immich, Docker, Nginx, WireGuard, Samba, dll).
• Tekan [Enter] atau [I] untuk Memasang (Install).
• Tekan [X] untuk Menghapus (Uninstall).
• Jendela konfirmasi akan muncul. Tekan [Y] untuk Ya atau [N] untuk Batal.
• Tool Immich SBC dapat dibuka langsung dengan: sudo immich-sbc-tool.`

	appCard := st.Card.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("📦 2. Pasang & Hapus Aplikasi"),
			"",
			appText,
		),
	)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, navCard, "  ", appCard)

	// Card 3: Update & Uninstall Bumiayu TUI
	updateUninstallText := `🔄 CARA UPDATE / INSTAL ULANG:
• Cukup jalankan kembali perintah install 1-baris. Sistem otomatis mendeteksi versi terpasang dan memperbarui binary ke versi rilis terbaru tanpa merusak setting Anda.

🗑️ CARA UNINSTALL (MENGHAPUS BERSIH):
• Cukup ketik perintah berikut di terminal server:
    sudo bumiayu-uninstall
• Atau via curl:
    curl -fsSL https://raw.githubusercontent.com/faisalfa123/bumiayu-de/main/scripts/uninstall.sh | sudo bash`

	manageCard := st.CardActive.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("🛠️ 3. Pembaruan Otomatis (Update) & Cara Uninstall"),
			"",
			updateUninstallText,
		),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		sub,
		"",
		topRow,
		"",
		manageCard,
	)
}
