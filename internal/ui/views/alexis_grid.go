package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// AlexisIconItem represents an icon on the main Alexis Launcher
type AlexisIconItem struct {
	ID          string
	Title       string
	Subtitle    string
	Icon        string
	CircleColor string
	BadgeText   string
}

// GetAlexisIcons returns the list of primary launcher icons
func GetAlexisIcons() []AlexisIconItem {
	return []AlexisIconItem{
		{
			ID:          "dashboard",
			Title:       "Dasbor Utama",
			Subtitle:    "Monitoring Sistem",
			Icon:        "📊",
			CircleColor: "#38BDF8", // Pastel Sky Blue
			BadgeText:   "Real-time",
		},
		{
			ID:          "apps",
			Title:       "Toko Aplikasi",
			Subtitle:    "55+ App Server",
			Icon:        "📦",
			CircleColor: "#10B981", // Pastel Emerald
			BadgeText:   "Katalog",
		},
		{
			ID:          "immich",
			Title:       "Immich SBC",
			Subtitle:    "Optimizer ODROID-C2",
			Icon:        "📸",
			CircleColor: "#F43F5E", // Pastel Rose
			BadgeText:   "FaisalFa",
		},
		{
			ID:          "storage",
			Title:       "Penyimpanan",
			Subtitle:    "Disk & Berkas",
			Icon:        "📁",
			CircleColor: "#F59E0B", // Pastel Amber
			BadgeText:   "Storage",
		},
		{
			ID:          "network",
			Title:       "Jaringan & Port",
			Subtitle:    "VPN & Web Server",
			Icon:        "🌐",
			CircleColor: "#8B5CF6", // Pastel Purple
			BadgeText:   "Network",
		},
		{
			ID:          "security",
			Title:       "Keamanan Server",
			Subtitle:    "Firewall & Hardening",
			Icon:        "🛡️",
			CircleColor: "#EC4899", // Pastel Pink
			BadgeText:   "Secure",
		},
		{
			ID:          "settings",
			Title:       "Pengaturan",
			Subtitle:    "Tema & Window",
			Icon:        "⚙️",
			CircleColor: "#64748B", // Pastel Slate
			BadgeText:   "Config",
		},
		{
			ID:          "guide",
			Title:       "Petunjuk & Bantuan",
			Subtitle:    "Panduan & Update",
			Icon:        "📖",
			CircleColor: "#14B8A6", // Pastel Teal
			BadgeText:   "Bantuan",
		},
	}
}

// RenderAlexisGrid renders the Alexis Icons Launcher (Layar 1)
func RenderAlexisGrid(selectedIndex int, animFrame int, st styles.ThemeStyles, width, height int) string {
	items := GetAlexisIcons()

	// Outer container sizing
	containerWidth := width - 4
	if containerWidth > 110 {
		containerWidth = 110
	}
	if containerWidth < 70 {
		containerWidth = 70
	}

	// WebUI Title Bar
	topControls := lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).Render("● ● ●")
	headerTitle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8FAFC")).Render("Bumiayu DE  •  Alexis Launcher")
	headerBar := lipgloss.JoinHorizontal(lipgloss.Center,
		topControls,
		"    ",
		headerTitle,
	)

	subBar := lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).Render("Pilih aplikasi atau modul server menggunakan tombol panah [← / → / ↑ / ↓] lalu tekan [Enter]")

	// Animated indicator
	spinners := []string{"◐", "◓", "◑", "◒"}
	spinner := spinners[(animFrame/6)%len(spinners)]
	statusPill := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#38BDF8")).
		Render(fmt.Sprintf("%s Sistem Aktif & Siap", spinner))

	headerBox := lipgloss.JoinVertical(lipgloss.Center,
		headerBar,
		"",
		subBar,
		statusPill,
	)

	// Render Grid Cards (4 columns x 2 rows)
	cardWidth := (containerWidth - 14) / 4
	if cardWidth < 16 {
		cardWidth = 16
	}

	var renderedCards []string
	for i, it := range items {
		isSelected := (i == selectedIndex)

		// Circular Icon Badge Effect
		iconCircleStyle := lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color(it.CircleColor)).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8FAFC"))
		subStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))

		cardBorderColor := lipgloss.Color("#1E293B")
		cardBg := lipgloss.Color("#161D2B")
		cursorMark := "  "

		if isSelected {
			cardBorderColor = lipgloss.Color("#38BDF8") // Sky accent on selection
			cardBg = lipgloss.Color("#1E293B")
			cursorMark = "👉 "
			titleStyle = titleStyle.Foreground(lipgloss.Color("#38BDF8"))
		}

		cardContent := lipgloss.JoinVertical(lipgloss.Center,
			iconCircleStyle.Render(" "+it.Icon+" "),
			"",
			titleStyle.Render(cursorMark+it.Title),
			subStyle.Render(it.Subtitle),
		)

		card := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cardBorderColor).
			Background(cardBg).
			Width(cardWidth).
			Height(7).
			Align(lipgloss.Center, lipgloss.Center).
			Render(cardContent)

		renderedCards = append(renderedCards, card)
	}

	// 2 rows of 4 cards
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, renderedCards[0], "  ", renderedCards[1], "  ", renderedCards[2], "  ", renderedCards[3])
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, renderedCards[4], "  ", renderedCards[5], "  ", renderedCards[6], "  ", renderedCards[7])

	cardsGrid := lipgloss.JoinVertical(lipgloss.Center,
		row1,
		"",
		row2,
	)

	// Bottom WebUI navigation pills
	footerNav := lipgloss.JoinHorizontal(lipgloss.Center,
		st.ButtonActive.Render("[Enter] Buka Dasbor"),
		"    ",
		st.ButtonNormal.Render("[↑ / ↓ / ← / →] Navigasi"),
		"    ",
		st.ButtonNormal.Render("[T] Ganti Tema"),
		"    ",
		st.ButtonNormal.Render("[q] Keluar"),
	)

	webContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Background(lipgloss.Color("#0B0F17")).
		Padding(1, 2).
		Width(containerWidth).
		Render(
			lipgloss.JoinVertical(lipgloss.Center,
				headerBox,
				"",
				cardsGrid,
				"",
				footerNav,
			),
		)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, webContainer)
}
