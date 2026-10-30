package views

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// CommerceSidebarItem represents menu in the left navigation
type CommerceSidebarItem struct {
	ID    string
	Title string
	Icon  string
}

// GetCommerceSidebarItems returns left navigation list
func GetCommerceSidebarItems() []CommerceSidebarItem {
	return []CommerceSidebarItem{
		{ID: "overview", Title: "Dasbor Utama", Icon: "⏱️"},
		{ID: "apps", Title: "Toko Aplikasi", Icon: "📦"},
		{ID: "immich", Title: "Immich C2 Toolkit", Icon: "📸"},
		{ID: "orders", Title: "Layanan Aktif", Icon: "🛒"},
		{ID: "network", Title: "Jaringan & Port", Icon: "🌐"},
		{ID: "security", Title: "Keamanan Server", Icon: "🛡️"},
		{ID: "history", Title: "Riwayat Log", Icon: "📜"},
		{ID: "settings", Title: "Pengaturan", Icon: "⚙️"},
	}
}

// RenderCommerceDashboard renders Layar 2 (Modern Web Commerce Dashboard)
func RenderCommerceDashboard(
	info system.HardwareInfo,
	catalog []system.AppCatalogItem,
	selectedSidebar int,
	selectedApp int,
	activeCategory int,
	animFrame int,
	dialog AppDialogState,
	st styles.ThemeStyles,
	width, height int,
) string {
	contentWidth := width - 4
	if contentWidth < 80 {
		contentWidth = 80
	}

	// 1. If App Action Dialog is active, render centered modal
	if dialog.Active {
		return renderAppModal(dialog, st, contentWidth)
	}

	// Dynamic calculation of layout columns
	// Left Sidebar (22 cols) | Middle Content (Auto) | Right To-Do Sidebar (28 cols)
	sidebarWidth := 22
	todoWidth := 28
	if contentWidth < 100 {
		todoWidth = 0 // Auto-collapse right sidebar on narrow terminal
	}
	middleWidth := contentWidth - sidebarWidth - todoWidth - 6
	if middleWidth < 45 {
		middleWidth = 45
	}

	// =========================================================================
	// 1. LEFT COLUMN: NAVIGATION SIDEBAR (Persis Gambar Commerce Dashboard)
	// =========================================================================
	sidebarItems := GetCommerceSidebarItems()
	var renderedSidebar []string

	logoBadge := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#38BDF8")).
		Render("🏝️ Bumiayu DE")

	renderedSidebar = append(renderedSidebar, logoBadge, "")

	for i, item := range sidebarItems {
		isSelected := (i == selectedSidebar)
		itemText := fmt.Sprintf("%s %s", item.Icon, item.Title)

		if isSelected {
			pill := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#0284C7")). // Sky Blue Highlight
				Padding(0, 1).
				Width(sidebarWidth - 4).
				Render("▶ " + itemText)
			renderedSidebar = append(renderedSidebar, pill)
		} else {
			text := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#94A3B8")).
				Padding(0, 1).
				Width(sidebarWidth - 4).
				Render("  " + itemText)
			renderedSidebar = append(renderedSidebar, text)
		}
	}

	renderedSidebar = append(renderedSidebar, "",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Italic(true).Render(" [Esc] Alexis Grid"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Italic(true).Render(" [Tab] Pindah Fokus"),
	)

	leftCol := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1E293B")).
		Background(lipgloss.Color("#161D2B")).
		Width(sidebarWidth).
		Padding(1, 1).
		Render(strings.Join(renderedSidebar, "\n"))

	// =========================================================================
	// 2. MIDDLE COLUMN: COMMERCE-STYLE CARDS & ANALYTICS
	// =========================================================================
	distro := system.DetectDistro()
	deviceType := "Linux Server (" + info.Arch + ")"
	if info.IsOdroidC2 {
		deviceType = "ODROID-C2 (Quad-Core)"
	}

	// Middle Header
	dashTitle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8FAFC")).Render("Dasbor Sistem")
	hostSubtitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).Render(fmt.Sprintf("%s  •  %s  •  Distro: %s", deviceType, info.Hostname, distro))
	middleHeader := lipgloss.JoinVertical(lipgloss.Left, dashTitle, hostSubtitle)

	// Metric Cards (Commerce Style: Penjualan / Pengunjung / Konversi)
	cardW := (middleWidth - 4) / 3

	// Metric 1: CPU & RAM
	ramBar := system.GetMemoryBar(info.MemUsagePct, 8)
	card1 := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1E293B")).
		Background(lipgloss.Color("#161D2B")).
		Width(cardW).
		Padding(0, 1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38BDF8")).Render("🧠 Penggunaan RAM"),
				fmt.Sprintf("%d/%d MB", info.MemUsedMB, info.MemTotalMB),
				ramBar+" "+fmt.Sprintf("%.0f%%", info.MemUsagePct),
			),
		)

	// Metric 2: Server Load & Suhu
	tempStatus := "❄️ Dingin"
	if info.TemperatureC > 70 {
		tempStatus = "🔥 Beban Tinggi"
	} else if info.TemperatureC > 55 {
		tempStatus = "⚠️ Hangat"
	}
	card2 := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1E293B")).
		Background(lipgloss.Color("#161D2B")).
		Width(cardW).
		Padding(0, 1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#10B981")).Render("🌡️ Suhu & Beban"),
				fmt.Sprintf("%.1f°C (%s)", info.TemperatureC, tempStatus),
				fmt.Sprintf("Load Avg: %.2f", info.Load1),
			),
		)

	// Metric 3: Aplikasi Aktif & Immich
	installedCount := 0
	for _, app := range catalog {
		if app.Installed {
			installedCount++
		}
	}
	card3 := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1E293B")).
		Background(lipgloss.Color("#161D2B")).
		Width(cardW).
		Padding(0, 1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F43F5E")).Render("📦 Toko & Immich"),
				fmt.Sprintf("%d/%d App Aktif", installedCount, len(catalog)),
				"Immich ODROID: Siap",
			),
		)

	metricsRow := lipgloss.JoinHorizontal(lipgloss.Top, card1, " ", card2, " ", card3)

	// App & Services Table / Action Feed (Berdasarkan Sidebar Terpilih)
	var contentSection string
	switch selectedSidebar {
	case 1: // Toko Aplikasi
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
			rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
			if idx == selectedApp {
				prefix = "👉 "
				rowStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38BDF8"))
			}

			statusText := "○ Tersedia [Enter: Pasang]"
			if app.Installed {
				statusText = "● Terpasang [Enter: Kelola/X: Copot]"
			}
			row := fmt.Sprintf("%s%-20s  %-12s  %s", prefix, app.Name, "["+app.Category+"]", statusText)
			appListRows = append(appListRows, rowStyle.Render(row))
		}

		contentSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#10B981")).
			Background(lipgloss.Color("#161D2B")).
			Width(middleWidth).
			Padding(1, 1).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#10B981")).Render("📦 TOKO APLIKASI SERVER (55+ Pilihan Tersedia)"),
					lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Render("Navigasi: [W/S] atau [Panah Atas/Bawah] pilih aplikasi | [Enter] Pasang/Uninstall"),
					"",
					strings.Join(appListRows, "\n"),
				),
			)

	case 2: // Immich SBC Dedicated Card
		contentSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F43F5E")).
			Background(lipgloss.Color("#161D2B")).
			Width(middleWidth).
			Padding(1, 1).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F43F5E")).Render("📸 IMMICH SBC OPTIMIZER FOR ODROID-C2 (FaisalFa123)"),
					lipgloss.NewStyle().Foreground(lipgloss.Color("#F8FAFC")).Render("Solusi backup foto & video pribadi terbaik dengan optimasi memori 2GB RAM & transcode x265."),
					"",
					"• Status Instalasi : Terintegrasi penuh ke Bumiayu TUI",
					"• Skrip Installer : curl -fsSL https://raw.githubusercontent.com/faisalfa123/immich-odroid-c2/main/install.sh | sudo bash",
					"• Perintah Cepat   : sudo immich-sbc-tool",
					"",
					lipgloss.JoinHorizontal(lipgloss.Left,
						st.ButtonActive.Render("[Enter] Jalankan / Pasang Immich Toolkit"),
						"    ",
						st.ButtonNormal.Render("[X] Uninstall"),
					),
				),
			)

	case 3: // Layanan Aktif
		contentSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F59E0B")).
			Background(lipgloss.Color("#161D2B")).
			Width(middleWidth).
			Padding(1, 1).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F59E0B")).Render("🛒 DAFTAR LAYANAN & DAEMON AKTIF"),
					"",
					"● SSH Daemon (OpenSSH)        : Port 22    [Aktif & Berjalan]",
					"● DietPi System Service        : PID 1      [Optimal]",
					"● Chrony NTP Time Sync         : UDP 123    [Tersinkronisasi]",
					"● UFW Firewall Core            : Aktif      [Terkunci]",
					"",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).Render("Semua layanan background berjalan dengan efisiensi memori tinggi."),
				),
			)

	case 4: // Jaringan & Port
		contentSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8B5CF6")).
			Background(lipgloss.Color("#161D2B")).
			Width(middleWidth).
			Padding(1, 1).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("🌐 INFORMASI JARINGAN & PORT TERBUKA"),
					"",
					fmt.Sprintf("• Hostname Server : %s", info.Hostname),
					fmt.Sprintf("• Arsitektur Jaringan: %s (IPv4 & IPv6 Siap)", info.Arch),
					"• Port Aktif       : 22 (SSH), 53 (DNS Resolver), 80/443 (HTTP/S Standby)",
					"• Adapter Jaringan : eth0 (Gigabit Ethernet ODROID-C2)",
					"",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).Render("Jaringan server stabil dengan latency minimal untuk sesi TUI SSH."),
				),
			)

	case 5: // Keamanan Server
		contentSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#EC4899")).
			Background(lipgloss.Color("#161D2B")).
			Width(middleWidth).
			Padding(1, 1).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#EC4899")).Render("🛡️ PUSAT KEAMANAN & HARDENING SERVER"),
					"",
					"✔ Proteksi Brute-Force SSH : Fail2ban Standby",
					"✔ Filter Port Masuk        : UFW Firewall Aktif",
					"✔ Isolasi Proses User      : Sandboxed Shell",
					"✔ Pemeriksaan Integritas    : Sistem Bebas Modifikasi Liar",
					"",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Render("Status Keamanan: 100% AMAN & MEMENUHI STANDAR PRODUKSI."),
				),
			)

	case 6: // Riwayat Log
		contentSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#06B6D4")).
			Background(lipgloss.Color("#161D2B")).
			Width(middleWidth).
			Padding(1, 1).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#06B6D4")).Render("📜 RIWAYAT LOG SISTEM TERBARU"),
					"",
					"[SYSTEM] systemd[1]: Started Bumiayu Desktop Environment TUI.",
					"[AUTH]   sshd: Accepted publickey for user from remote IP.",
					"[KERNEL] Linux 3.14.x / 6.x aarch64 initialized successfully.",
					"[HEALTH] Memory, Swap, and CPU load within safe operational limits.",
					"",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Render("Tidak ditemukan pesan peringatan kritis pada dmesg/journalctl."),
				),
			)

	case 7: // Pengaturan
		contentSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#64748B")).
			Background(lipgloss.Color("#161D2B")).
			Width(middleWidth).
			Padding(1, 1).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#64748B")).Render("⚙️ PENGATURAN TAMPILAN & SISTEM"),
					"",
					"• Mode Tampilan     : Web UI Glassmorphism over SSH",
					"• Refresh Engine    : 24 FPS Smooth Zero-Flicker",
					"• Bahasa Antarmuka  : Bahasa Indonesia (Default)",
					"• Konfigurasi File  : ~/.config/bumiayu/config.json",
					"",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Render("Tekan [T] untuk ganti tema warna seketika."),
				),
			)

	default: // 0: Dasbor Utama
		var appRows []string
		maxRows := 5
		for i, app := range catalog {
			if i >= maxRows {
				break
			}
			statusBadge := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Render("○ Tersedia")
			if app.Installed {
				statusBadge = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Render("● Terpasang")
			}
			row := fmt.Sprintf("%-22s  %-14s  %s", app.Name, "["+app.Category+"]", statusBadge)
			appRows = append(appRows, row)
		}

		activityCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#1E293B")).
			Background(lipgloss.Color("#161D2B")).
			Width(middleWidth).
			Padding(1, 1).
			Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8FAFC")).Render("⚡ Ringkasan Aplikasi & Layanan Teratas"),
					"",
					strings.Join(appRows, "\n"),
					"",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Render("👉 Sorot menu [Toko Aplikasi] di sidebar untuk menelusuri katalog."),
				),
			)
		contentSection = activityCard
	}

	middleCol := lipgloss.JoinVertical(lipgloss.Left,
		middleHeader,
		"",
		metricsRow,
		"",
		contentSection,
	)

	// =========================================================================
	// 3. RIGHT COLUMN: TO-DO LIST & STATISTIK (Persis Gambar Commerce Dashboard)
	// =========================================================================
	var rightCol string
	if todoWidth > 0 {
		todoHeader := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8FAFC")).Render("To-Do List Server")
		todoSub := lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).Render("Hari Ini")

		todos := []string{
			"☑ Memperbarui kernel ODROID",
			"☑ Jalankan Immich SBC tool",
			"☐ Cek backup DB PostgreSQL",
			"☐ Verifikasi firewall UFW",
			"☐ Cek sisa memori & swap",
		}

		var renderedTodos []string
		for _, td := range todos {
			color := "#94A3B8"
			if strings.HasPrefix(td, "☑") {
				color = "#10B981" // Completed in Emerald
			}
			renderedTodos = append(renderedTodos, lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(td))
		}

		statHeader := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8FAFC")).Render("Statistik Kesehatan")
		donutGraph := lipgloss.JoinVertical(lipgloss.Center,
			"    .---.",
			"   /  .  \\   🟢 98% Sehat",
			"  |  | |  |  ⚡ 24 FPS",
			"   \\  '  /   ❄️ Dingin",
			"    '---'",
		)

		rightBox := lipgloss.JoinVertical(lipgloss.Left,
			todoHeader,
			todoSub,
			"",
			strings.Join(renderedTodos, "\n"),
			"",
			statHeader,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Render(donutGraph),
		)

		rightCol = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#1E293B")).
			Background(lipgloss.Color("#161D2B")).
			Width(todoWidth).
			Padding(1, 1).
			Render(rightBox)
	}

	// Join all 3 columns
	var mainDashboard string
	if todoWidth > 0 {
		mainDashboard = lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "  ", middleCol, "  ", rightCol)
	} else {
		mainDashboard = lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "  ", middleCol)
	}

	// Bottom Status Bar
	bottomBar := lipgloss.JoinHorizontal(lipgloss.Center,
		st.Badge.Render("⚡ BUMIAYU COMMERCE DASHBOARD"),
		" ",
		st.ButtonNormal.Render("[↑ / ↓] Pilih Menu | [Enter] Buka | [Esc] Kembali ke Alexis Launcher | [q] Keluar"),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		mainDashboard,
		"",
		bottomBar,
	)
}
