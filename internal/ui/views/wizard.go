package views

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/config"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

type WizardStep int

const (
	StepRecovery WizardStep = iota
	StepLanguage
	StepTheme
	StepWM
	StepSSH
	StepComplete
)

type WizardState struct {
	Active        bool
	CurrentStep   WizardStep
	SelectedRow   int
	HasBackup     bool
	DraftConfig   config.UserConfig
	StatusMessage string
}

// RenderWizard renders the first-run OOBE setup wizard and recovery modal
func RenderWizard(w WizardState, st styles.ThemeStyles, width, height int) string {
	contentWidth := width - 6
	if contentWidth < 50 {
		contentWidth = 50
	}
	modalWidth := contentWidth - 10
	if modalWidth > 76 {
		modalWidth = 76
	}

	var stepTitle string
	var stepBody string

	switch w.CurrentStep {
	case StepRecovery:
		stepTitle = "🔄 TERDETEKSI PENGATURAN SEBELUMNYA (RECOVERY)"
		options := []string{
			"1. Pulihkan Pengaturan Lama (Recovery: Tema, WM & Konfigurasi Anda)",
			"2. Mulai Konfigurasi Baru (Fresh Setup dari Awal)",
		}
		var rendered []string
		for i, opt := range options {
			if i == w.SelectedRow {
				rendered = append(rendered, st.RowSelected.Render("▶ "+opt))
			} else {
				rendered = append(rendered, st.RowNormal.Render("  "+opt))
			}
		}

		stepBody = lipgloss.JoinVertical(lipgloss.Left,
			st.Subtitle.Render("Sistem mendeteksi ada cadangan konfigurasi sebelumnya di ~/.config/bumiayu/."),
			"",
			strings.Join(rendered, "\n\n"),
			"",
			st.RowNormal.Render("Gunakan [↑ / ↓] lalu tekan [Enter] untuk memilih."),
		)

	case StepLanguage:
		stepTitle = "🌐 LANGKAH 1: PILIH BAHASA / LANGUAGE"
		options := []string{
			"1. Bahasa Indonesia (Bahasa Resmi & Panduan Lokal)",
			"2. English (International Standard)",
		}
		var rendered []string
		for i, opt := range options {
			if i == w.SelectedRow {
				rendered = append(rendered, st.RowSelected.Render("▶ "+opt))
			} else {
				rendered = append(rendered, st.RowNormal.Render("  "+opt))
			}
		}
		stepBody = lipgloss.JoinVertical(lipgloss.Left,
			st.Subtitle.Render("Pilih bahasa antarmuka yang paling nyaman untuk Anda:"),
			"",
			strings.Join(rendered, "\n\n"),
			"",
			st.RowNormal.Render("Tekan [Enter] untuk melanjutkan ke langkah berikutnya."),
		)

	case StepTheme:
		stepTitle = "🎨 LANGKAH 2: PILIH TEMA TAMPILAN (PALET WARNA)"
		options := []string{
			"1. ☀️ Terang Bersih (Light Mode - Sangat Nyaman di Terminal)",
			"2. 🌾 Hangat Lembut (Warm Ivory / Cream)",
			"3. 🌙 Gelap Malam (Midnight Dark)",
		}
		var rendered []string
		for i, opt := range options {
			if i == w.SelectedRow {
				rendered = append(rendered, st.RowSelected.Render("▶ "+opt))
			} else {
				rendered = append(rendered, st.RowNormal.Render("  "+opt))
			}
		}
		stepBody = lipgloss.JoinVertical(lipgloss.Left,
			st.Subtitle.Render("Pilih gaya kecerahan antarmuka (Bisa diubah kapan saja via tombol [T]):"),
			"",
			strings.Join(rendered, "\n\n"),
			"",
			st.RowNormal.Render("Tekan [Enter] untuk memilih dan lanjut."),
		)

	case StepWM:
		stepTitle = "🖥️ LANGKAH 3: PILIH MODEL WINDOW MANAGER (WM)"
		options := []string{
			"1. 🪟 Window Mengambang (Floating WM: Pop-up di tengah desktop)",
			"2. 🔲 Tiling 2-Kolom (Split WM: Menu di kiri, aksi di kanan)",
			"3. 📋 Kartu Grid Sederhana (Card Grid: Kartu ringkas berdampingan)",
		}
		var rendered []string
		for i, opt := range options {
			if i == w.SelectedRow {
				rendered = append(rendered, st.RowSelected.Render("▶ "+opt))
			} else {
				rendered = append(rendered, st.RowNormal.Render("  "+opt))
			}
		}
		stepBody = lipgloss.JoinVertical(lipgloss.Left,
			st.Subtitle.Render("Pilih tata letak (layout) kerja yang paling intuitif bagi Anda:"),
			"",
			strings.Join(rendered, "\n\n"),
			"",
			st.RowNormal.Render("Tekan [Enter] untuk memilih model WM."),
		)

	case StepSSH:
		stepTitle = "⚡ LANGKAH 4: INTEGRASI AUTO-START SAAT SSH"
		options := []string{
			"1. Ya, Buka Bumiayu TUI otomatis setiap kali login SSH",
			"2. Tidak, saya akan membukanya manual dengan mengetik 'bumiayu'",
		}
		var rendered []string
		for i, opt := range options {
			if i == w.SelectedRow {
				rendered = append(rendered, st.RowSelected.Render("▶ "+opt))
			} else {
				rendered = append(rendered, st.RowNormal.Render("  "+opt))
			}
		}
		stepBody = lipgloss.JoinVertical(lipgloss.Left,
			st.Subtitle.Render("Apakah Anda ingin Bumiayu TUI menyapa otomatis saat Anda SSH ke ODROID/server?"),
			"",
			strings.Join(rendered, "\n\n"),
			"",
			st.RowNormal.Render("Tekan [Enter] untuk menyimpan preferensi."),
		)

	case StepComplete:
		stepTitle = "🎉 SETUP SELESAI & SIAP DIGUNAKAN!"
		summary := fmt.Sprintf(
			"• Bahasa Terpilih    : %s\n• Tema Tampilan      : %s\n• Window Manager     : %s\n• Auto-Start SSH     : %t\n• File Konfigurasi   : ~/.config/bumiayu/config.json",
			w.DraftConfig.Language,
			w.DraftConfig.Theme,
			w.DraftConfig.WMModel,
			w.DraftConfig.AutoStartSSH,
		)
		stepBody = lipgloss.JoinVertical(lipgloss.Left,
			st.Subtitle.Render("Konfigurasi Anda telah disimpan dengan aman dan cadangan dibuat."),
			"",
			st.Card.Render(summary),
			"",
			st.ButtonActive.Render("👉 Tekan [Enter] atau [Spasi] untuk Memulai Desktop Bumiayu"),
		)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		st.Badge.Render("🏝️ BUMIAYU TUI - WIZARD SETUP AWAL"),
		"",
		st.Title.Render(stepTitle),
		"",
		stepBody,
	)

	return lipgloss.Place(contentWidth, 20, lipgloss.Center, lipgloss.Center, st.Modal.Width(modalWidth).Render(content))
}
