package styles

import "github.com/charmbracelet/lipgloss"

type ThemeMode int

const (
	ThemeLight ThemeMode = iota // Terang (Default)
	ThemeDark                   // Gelap (Midnight)
	ThemeWarm                   // Hangat (Cream / Warm Ivory)
)

type ThemeStyles struct {
	Mode         ThemeMode
	BgColor      lipgloss.Color
	CardBg       lipgloss.Color
	TextColor    lipgloss.Color
	MutedColor   lipgloss.Color
	PrimaryColor lipgloss.Color
	AccentColor  lipgloss.Color
	BorderColor  lipgloss.Color
	DangerColor  lipgloss.Color
	SuccessColor lipgloss.Color

	// Pre-built styles
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	TopBar        lipgloss.Style
	Badge         lipgloss.Style
	Pill          lipgloss.Style
	PillWarning   lipgloss.Style
	PillDanger    lipgloss.Style
	ActiveTab     lipgloss.Style
	InactiveTab   lipgloss.Style
	Card          lipgloss.Style
	CardActive    lipgloss.Style
	Modal         lipgloss.Style
	RowSelected   lipgloss.Style
	RowNormal     lipgloss.Style
	Footer        lipgloss.Style
	ButtonActive  lipgloss.Style
	ButtonNormal  lipgloss.Style
}

// Current theme instance
func GetTheme(mode ThemeMode) ThemeStyles {
	ts := ThemeStyles{Mode: mode}

	switch mode {
	case ThemeLight:
		// WebUI Light Theme (Bersih, Kontras Lembut, Card Modern)
		ts.BgColor = lipgloss.Color("#F1F5F9")      // Slate 100
		ts.CardBg = lipgloss.Color("#FFFFFF")       // Pure White
		ts.TextColor = lipgloss.Color("#0F172A")    // Slate 900
		ts.MutedColor = lipgloss.Color("#64748B")   // Slate 500
		ts.PrimaryColor = lipgloss.Color("#0284C7") // Sky 600
		ts.AccentColor = lipgloss.Color("#3B82F6")  // Blue 500
		ts.BorderColor = lipgloss.Color("#CBD5E1")  // Slate 300
		ts.DangerColor = lipgloss.Color("#EF4444")  // Red 500
		ts.SuccessColor = lipgloss.Color("#10B981") // Emerald 500

	case ThemeWarm:
		// WebUI Warm Cream (Elegan & Nyaman di Mata)
		ts.BgColor = lipgloss.Color("#FEFCE8")      // Warm Cream
		ts.CardBg = lipgloss.Color("#FFFDF5")       // Soft Ivory
		ts.TextColor = lipgloss.Color("#292524")    // Stone 800
		ts.MutedColor = lipgloss.Color("#78716C")   // Stone 500
		ts.PrimaryColor = lipgloss.Color("#D97706") // Amber 600
		ts.AccentColor = lipgloss.Color("#059669")  // Emerald
		ts.BorderColor = lipgloss.Color("#D6D3D1")  // Stone 300
		ts.DangerColor = lipgloss.Color("#E11D48")  // Rose 600
		ts.SuccessColor = lipgloss.Color("#15803D") // Forest Green

	case ThemeDark:
		// Modern Dashboard Dark Mode (Layaknya Vercel / Pro Dashboard)
		ts.BgColor = lipgloss.Color("#0B0F17")      // Deep Space / Obsidian
		ts.CardBg = lipgloss.Color("#161D2B")       // Dark Slate Card
		ts.TextColor = lipgloss.Color("#F8FAFC")    // Crisp White
		ts.MutedColor = lipgloss.Color("#94A3B8")   // Slate Gray
		ts.PrimaryColor = lipgloss.Color("#10B981") // Jade Green / Emerald
		ts.AccentColor = lipgloss.Color("#38BDF8")  // Sky Blue
		ts.BorderColor = lipgloss.Color("#1E293B")  // Dark Border
		ts.DangerColor = lipgloss.Color("#F43F5E")  // Rose Red
		ts.SuccessColor = lipgloss.Color("#22C55E") // Bright Green
	}

	ts.Title = lipgloss.NewStyle().Bold(true).Foreground(ts.PrimaryColor)
	ts.Subtitle = lipgloss.NewStyle().Foreground(ts.MutedColor)

	ts.TopBar = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ts.BorderColor)

	ts.Badge = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.PrimaryColor)

	ts.Pill = lipgloss.NewStyle().
		Padding(0, 1)

	ts.PillWarning = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.AccentColor).
		Padding(0, 1)

	ts.PillDanger = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.DangerColor).
		Padding(0, 1)

	ts.ActiveTab = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.PrimaryColor).
		Padding(0, 2)

	ts.InactiveTab = lipgloss.NewStyle().
		Foreground(ts.MutedColor).
		Padding(0, 2)

	ts.Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ts.BorderColor).
		Padding(0, 1)

	ts.CardActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ts.PrimaryColor).
		Padding(0, 1)

	ts.Modal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ts.PrimaryColor).
		Padding(1, 2)

	ts.RowSelected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.AccentColor).
		Padding(0, 1)

	ts.RowNormal = lipgloss.NewStyle().
		Padding(0, 1)

	ts.Footer = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ts.BorderColor)

	ts.ButtonActive = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.PrimaryColor).
		Padding(0, 2)

	ts.ButtonNormal = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ts.BorderColor).
		Padding(0, 1)

	return ts
}
