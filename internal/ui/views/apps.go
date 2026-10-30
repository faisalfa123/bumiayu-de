package views

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// AppDialogState holds state for install/uninstall confirmation and execution
type AppDialogState struct {
	Active     bool
	Action     string // "install" or "uninstall"
	App        system.AppCatalogItem
	IsRunning  bool
	Logs       string
	Success    bool
}

// RenderApps renders the simplified, beautiful app store with install/uninstall modal
func RenderApps(catalog []system.AppCatalogItem, selectedIndex int, pkgCount int, dialog AppDialogState, st styles.ThemeStyles, lang i18n.Language, width, height int) string {
	contentWidth := width - 6
	if contentWidth < 50 {
		contentWidth = 50
	}

	header := st.Title.Render("📦 TOKO APLIKASI SERVER (APP STORE)")
	sub := st.Subtitle.Render(fmt.Sprintf("Kelola aplikasi server siap pakai dengan mudah (Terdeteksi %d paket sistem aktif).", pkgCount))

	// If Modal is Active, Render Center Modal
	if dialog.Active {
		return renderAppModal(dialog, st, contentWidth)
	}

	// 1. App List
	var listItems []string
	for i, app := range catalog {
		statusBadge := st.Pill.Render("⚪ Belum Ada")
		if app.Installed {
			statusBadge = st.Badge.Render("🟢 Terpasang")
		}

		line := fmt.Sprintf("%-18s  %-16s  %s", app.Name, "["+app.Category+"]", statusBadge)
		if i == selectedIndex {
			listItems = append(listItems, st.RowSelected.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, st.RowNormal.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	listPanel := st.Card.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			st.Title.Render("Daftar Aplikasi Populer:"),
			st.Subtitle.Render("Gunakan [↑ / ↓] untuk memilih, [Enter/I] pasang, [X] hapus"),
			"",
			strings.Join(listItems, "\n"),
		),
	)

	// 2. Selected App Detail Panel
	var detailPanel string
	if selectedIndex >= 0 && selectedIndex < len(catalog) {
		app := catalog[selectedIndex]
		statusText := "Belum Terpasang di Sistem"
		actionHint := "[I] atau [Enter] untuk Memasang Aplikasi Ini"
		if app.Installed {
			statusText = "Sudah Terpasang & Siap Digunakan"
			actionHint = "[X] untuk Menghapus (Uninstall) Aplikasi Ini"
		}

		detailContent := fmt.Sprintf(
			"📌 Nama: %s\n🏷️ Kategori: %s\n📡 Port Akses: %s\n📊 Status: %s\n\n📖 Deskripsi Singkat:\n%s\n\n💡 Rekomendasi Ramah Pemula:\n%s\n\n⚡ Aksi Tersedia:\n%s",
			app.Name, app.Category, app.Port, statusText,
			app.Description,
			app.Insight,
			actionHint,
		)
		detailPanel = st.CardActive.Width(contentWidth/2 - 2).Render(detailContent)
	}

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)

	// Bottom quick hint bar
	hintBar := st.Card.Width(contentWidth).Render(
		st.RowNormal.Render("⌨️ Pintasan Cepat: Tekan [Enter] atau [I] untuk Pasang | Tekan [X] untuk Hapus | [1-4] atau [←/→] Pindah Tab"),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		sub,
		"",
		mainRow,
		"",
		hintBar,
	)
}

func renderAppModal(dialog AppDialogState, st styles.ThemeStyles, contentWidth int) string {
	modalWidth := contentWidth - 10
	if modalWidth > 70 {
		modalWidth = 70
	}

	actionTitle := "PASANG APLIKASI"
	actionVerb := "memasang"
	btnConfirm := "[Y] Ya, Pasang Sekarang"
	if dialog.Action == "uninstall" {
		actionTitle = "HAPUS APLIKASI"
		actionVerb = "menghapus"
		btnConfirm = "[Y] Ya, Hapus Sekarang"
	}

	if dialog.IsRunning {
		// Loading / Execution Log view
		content := lipgloss.JoinVertical(lipgloss.Left,
			st.Badge.Render("⏳ SEDANG MEMPROSES..."),
			"",
			fmt.Sprintf("Sedang %s %s di server...", actionVerb, dialog.App.Name),
			"",
			st.Card.Width(modalWidth-4).Render(dialog.Logs),
			"",
			st.Subtitle.Render("Mohon tunggu sebentar hingga proses selesai..."),
		)
		return lipgloss.Place(contentWidth, 18, lipgloss.Center, lipgloss.Center, st.Modal.Width(modalWidth).Render(content))
	}

	if dialog.Logs != "" {
		// Completed view
		statusBadge := st.Badge.Render("✅ SELESAI!")
		if !dialog.Success {
			statusBadge = st.PillDanger.Render("❌ GAGAL")
		}
		content := lipgloss.JoinVertical(lipgloss.Left,
			statusBadge,
			"",
			fmt.Sprintf("Proses %s %s telah selesai.", actionVerb, dialog.App.Name),
			"",
			st.Card.Width(modalWidth-4).Render(dialog.Logs),
			"",
			st.ButtonActive.Render("Tekan [Enter] atau [Esc] untuk Menutup"),
		)
		return lipgloss.Place(contentWidth, 18, lipgloss.Center, lipgloss.Center, st.Modal.Width(modalWidth).Render(content))
	}

	// Confirmation prompt view
	content := lipgloss.JoinVertical(lipgloss.Left,
		st.Badge.Render("⚠️ KONFIRMASI "+actionTitle),
		"",
		fmt.Sprintf("Apakah Anda yakin ingin %s:", actionVerb),
		st.Title.Render("👉 "+dialog.App.Name),
		st.Subtitle.Render("Paket: "+dialog.App.PackageName+" | Kategori: "+dialog.App.Category),
		"",
		dialog.App.Insight,
		"",
		lipgloss.JoinHorizontal(lipgloss.Center,
			st.ButtonActive.Render(btnConfirm),
			"    ",
			st.ButtonNormal.Render("[N] Batal / Esc"),
		),
	)

	return lipgloss.Place(contentWidth, 18, lipgloss.Center, lipgloss.Center, st.Modal.Width(modalWidth).Render(content))
}
