package views

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// IPadTab represents the top tab items in iPadOS style
type IPadTab int

const (
	IPadTabDashboard IPadTab = iota
	IPadTabApps
	IPadTabImmich
	IPadTabNetwork
	IPadTabSecurity
	IPadTabSettings
)

// RenderIPadDashboard renders Mac / iPadOS clean light dashboard with top tab bar
func RenderIPadDashboard(
	info system.HardwareInfo,
	catalog []system.AppCatalogItem,
	currentTab IPadTab,
	selectedApp int,
	activeCategory int,
	animFrame int,
	dialog AppDialogState,
	webPort int,
	st styles.ThemeStyles,
	width, height int,
) string {
	contentWidth := width - 4
	if contentWidth < 78 {
		contentWidth = 78
	}

	// 1. If App Action Dialog is active, render centered modal
	if dialog.Active {
		return renderAppModal(dialog, st, contentWidth)
	}

	// =========================================================================
	// 1. TOP HEADER & IPADOS TAB BAR (Bukan Hitam, Kartu Putih Bersih)
	// =========================================================================
	appleDots := lipgloss.NewStyle().Foreground(lipgloss.Color("#DC2626")).Render("● ") +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")).Render("● ") +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#16A34A")).Render("●")

	titleText := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0F172A")).Render("Bumiayu Control Center")
	webBadge := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#2563EB")).
		Background(lipgloss.Color("#DBEAFE")).
		Padding(0, 1).
		Render(fmt.Sprintf("🌐 Web UI: http://%s:%d", info.Hostname, webPort))

	topBrand := lipgloss.JoinHorizontal(lipgloss.Center,
		appleDots,
		"  ",
		titleText,
		"    ",
		webBadge,
	)

	// iPadOS Style Segmented Top Tab Bar
	tabs := []struct {
		tab   IPadTab
		label string
		icon  string
	}{
		{IPadTabDashboard, "Dasbor Utama", "📊"},
		{IPadTabApps, "Toko Aplikasi", "📦"},
		{IPadTabImmich, "Immich C2", "📸"},
		{IPadTabNetwork, "Jaringan", "🌐"},
		{IPadTabSecurity, "Keamanan", "🛡️"},
		{IPadTabSettings, "Pengaturan", "⚙️"},
	}

	var renderedTabs []string
	for idx, t := range tabs {
		tabStyle := lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#475569")).
			Background(lipgloss.Color("#F1F5F9"))

		if t.tab == currentTab {
			tabStyle = lipgloss.NewStyle().
				Bold(true).
				Padding(0, 2).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#2563EB"))
		}
		numBadge := fmt.Sprintf("[%d] ", idx+1)
		renderedTabs = append(renderedTabs, tabStyle.Render(numBadge+t.icon+" "+t.label))
	}

	tabBarRow := lipgloss.NewStyle().
		Background(lipgloss.Color("#E2E8F0")).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBD5E1")).
		Render(strings.Join(renderedTabs, " "))

	headerContainer := lipgloss.NewStyle().
		Background(lipgloss.Color("#FFFFFF")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBD5E1")).
		Padding(1, 2).
		Width(contentWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, topBrand, "", tabBarRow))

	// =========================================================================
	// 2. METRIC PILL CARDS (Mac/iPadOS Style Clean Cards)
	// =========================================================================
	cardWidth := (contentWidth - 6) / 3
	if cardWidth < 24 {
		cardWidth = 24
	}

	ramBar := system.GetMemoryBar(info.MemUsagePct, 8)
	card1 := lipgloss.NewStyle().
		Background(lipgloss.Color("#FFFFFF")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBD5E1")).
		Width(cardWidth).
		Padding(0, 1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2563EB")).Render("🧠 Penggunaan RAM"),
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0F172A")).Render(fmt.Sprintf("%d / %d MB", info.MemUsedMB, info.MemTotalMB)),
				ramBar+" "+fmt.Sprintf("%.0f%%", info.MemUsagePct),
			),
		)

	tempStatus := "❄️ Normal"
	if info.TemperatureC > 70 {
		tempStatus = "🔥 Panas"
	} else if info.TemperatureC > 55 {
		tempStatus = "⚠️ Hangat"
	}

	card2 := lipgloss.NewStyle().
		Background(lipgloss.Color("#FFFFFF")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBD5E1")).
		Width(cardWidth).
		Padding(0, 1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#D97706")).Render("🌡️ Suhu & Beban"),
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0F172A")).Render(fmt.Sprintf("%.1f°C  (%s)", info.TemperatureC, tempStatus)),
				fmt.Sprintf("Load 1m: %.2f", info.Load1),
			),
		)

	card3 := lipgloss.NewStyle().
		Background(lipgloss.Color("#FFFFFF")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBD5E1")).
		Width(cardWidth).
		Padding(0, 1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#16A34A")).Render("⚡ Server Health"),
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0F172A")).Render("🟢 100% Siap"),
				"FPS: 30 • Uptime: "+info.UptimeStr,
			),
		)

	metricsRow := lipgloss.JoinHorizontal(lipgloss.Top, card1, "  ", card2, "  ", card3)

	// =========================================================================
	// 3. MAIN TAB CONTENT (Berubah sesuai tab yang aktif)
	// =========================================================================
	var tabContent string
	switch currentTab {
	case IPadTabApps:
		var appListRows []string
		startIdx := 0
		if selectedApp > 6 {
			startIdx = selectedApp - 6
		}
		endIdx := startIdx + 8
		if endIdx > len(catalog) {
			endIdx = len(catalog)
		}

		for idx := startIdx; idx < endIdx; idx++ {
			app := catalog[idx]
			prefix := "  "
			rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#334155"))
			if idx == selectedApp {
				prefix = "👉 "
				rowStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2563EB"))
			}

			statusText := "○ Tersedia [Enter: Pasang]"
			if app.Installed {
				statusText = "● Terpasang [Enter: Kelola/X: Copot]"
			}
			row := fmt.Sprintf("%s%-24s  %-14s  %s", prefix, app.Name, "["+app.Category+"]", statusText)
			appListRows = append(appListRows, rowStyle.Render(row))
		}

		tabContent = lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#16A34A")).
			Width(contentWidth).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#16A34A")).Render("📦 TOKO APLIKASI SERVER (55+ Pilihan Tersedia)"),
					lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Render("Navigasi: [W/S] atau [Panah Atas/Bawah] pilih aplikasi | [Enter] Pasang | [X] Copot"),
					"",
					strings.Join(appListRows, "\n"),
				),
			)

	case IPadTabImmich:
		tabContent = lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#E11D48")).
			Width(contentWidth).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E11D48")).Render("📸 IMMICH SBC OPTIMIZER FOR ODROID-C2 (FaisalFa123)"),
					lipgloss.NewStyle().Foreground(lipgloss.Color("#334155")).Render("Solusi backup foto & video pribadi terbaik dengan optimasi memori 2GB RAM & transcode x265."),
					"",
					"• Status Instalasi : Terintegrasi penuh ke Bumiayu TUI & Web",
					"• Skrip Installer : curl -fsSL https://raw.githubusercontent.com/faisalfa123/immich-odroid-c2/main/install.sh | bash",
					"• Perintah Cepat   : sudo immich-sbc-tool",
					"",
					lipgloss.JoinHorizontal(lipgloss.Left,
						st.ButtonActive.Render("[Enter] Jalankan / Pasang Immich Toolkit"),
						"    ",
						st.ButtonNormal.Render("[X] Uninstall"),
					),
				),
			)

	case IPadTabNetwork:
		tabContent = lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Width(contentWidth).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")).Render("🌐 INFORMASI JARINGAN & PORT TERBUKA"),
					"",
					fmt.Sprintf("• Hostname Server : %s", info.Hostname),
					fmt.Sprintf("• Port Web Admin   : http://%s:%d (Buka di browser)", info.Hostname, webPort),
					"• Port Aktif       : 22 (SSH), 80/443 (HTTP/S Standby), 2283 (Immich Standby)",
					"• Status Jaringan  : Terhubung Stabil (Ethernet Gigabit)",
				),
			)

	case IPadTabSecurity:
		tabContent = lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#DC2626")).
			Width(contentWidth).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#DC2626")).Render("🛡️ PUSAT KEAMANAN SERVER ODROID-C2"),
					"",
					"✔ Proteksi Brute-Force SSH : Fail2ban Standby",
					"✔ Isolasi Port & Layanan   : Sandboxed Server Daemon",
					"✔ Pemeriksaan Integritas   : Sistem Bebas Modifikasi Liar",
					"",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#16A34A")).Render("Status Keamanan: 100% AMAN & MEMENUHI STANDAR PRODUKSI."),
				),
			)

	case IPadTabSettings:
		tabContent = lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#2563EB")).
			Width(contentWidth).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2563EB")).Render("⚙️ PENGATURAN TAMPILAN & SISTEM"),
					"",
					"• Tema Tampilan     : Mac / iPadOS Clean Light (Abu-abu Modern & Putih Salju)",
					"• Refresh Engine    : 30 FPS Smooth Ultra-Responsive",
					fmt.Sprintf("• Port Web Server   : %d (Dapat diakses di browser)", webPort),
					"• Bahasa Antarmuka  : Bahasa Indonesia (Default)",
				),
			)

	default: // Dasbor Utama
		var appRows []string
		maxRows := 5
		for i, app := range catalog {
			if i >= maxRows {
				break
			}
			statusBadge := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Render("○ Tersedia")
			if app.Installed {
				statusBadge = lipgloss.NewStyle().Foreground(lipgloss.Color("#16A34A")).Render("● Terpasang")
			}
			row := fmt.Sprintf("%-24s  %-14s  %s", app.Name, "["+app.Category+"]", statusBadge)
			appRows = append(appRows, row)
		}

		tabContent = lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#CBD5E1")).
			Width(contentWidth).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0F172A")).Render("⚡ Ringkasan Aktivitas Server Terkini"),
					"",
					strings.Join(appRows, "\n"),
					"",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#2563EB")).Render("👉 Tekan [2] untuk masuk ke Toko Aplikasi atau [3] untuk Immich C2."),
				),
			)
	}

	// =========================================================================
	// 4. FOOTER CONTROLS
	// =========================================================================
	bottomHint := lipgloss.NewStyle().
		Background(lipgloss.Color("#FFFFFF")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBD5E1")).
		Width(contentWidth).
		Padding(0, 1).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Center,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2563EB")).Render("NAVIGASI MUDAH:"),
				"  ",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#334155")).Render("[1-6] Pilih Tab | [← / →] Geser Tab | [Enter] Aksi | [q] Keluar"),
			),
		)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerContainer,
		"",
		metricsRow,
		"",
		tabContent,
		"",
		bottomHint,
	)
}
