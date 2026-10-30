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
		// Tampilan Terang Bersih & Elegan
		ts.BgColor = lipgloss.Color("#F8FAFC")      // Slate 50
		ts.CardBg = lipgloss.Color("#FFFFFF")       // Pure White
		ts.TextColor = lipgloss.Color("#0F172A")    // Slate 900
		ts.MutedColor = lipgloss.Color("#64748B")   // Slate 500
		ts.PrimaryColor = lipgloss.Color("#0D9488") // Teal 600
		ts.AccentColor = lipgloss.Color("#2563EB")  // Royal Blue
		ts.BorderColor = lipgloss.Color("#CBD5E1")  // Slate 300
		ts.DangerColor = lipgloss.Color("#DC2626")  // Red 600
		ts.SuccessColor = lipgloss.Color("#16A34A") // Green 600

	case ThemeWarm:
		// Tampilan Hangat & Lembut di Mata
		ts.BgColor = lipgloss.Color("#FEFCE8")      // Warm Cream
		ts.CardBg = lipgloss.Color("#FFFDF5")       // Soft Ivory
		ts.TextColor = lipgloss.Color("#292524")    // Stone 800
		ts.MutedColor = lipgloss.Color("#78716C")   // Stone 500
		ts.PrimaryColor = lipgloss.Color("#D97706") // Amber 600
		ts.AccentColor = lipgloss.Color("#059669")  // Emerald
		ts.BorderColor = lipgloss.Color("#E7E5E4")  // Stone 300
		ts.DangerColor = lipgloss.Color("#E11D48")  // Rose 600
		ts.SuccessColor = lipgloss.Color("#15803D") // Forest Green

	case ThemeDark:
		// Tampilan Gelap (Midnight Blue)
		ts.BgColor = lipgloss.Color("#0F172A")      // Dark Slate
		ts.CardBg = lipgloss.Color("#1E293B")       // Slate Card
		ts.TextColor = lipgloss.Color("#F8FAFC")    // Crisp White
		ts.MutedColor = lipgloss.Color("#94A3B8")   // Slate Gray
		ts.PrimaryColor = lipgloss.Color("#10B981") // Jade Green
		ts.AccentColor = lipgloss.Color("#06B6D4")  // Cyan
		ts.BorderColor = lipgloss.Color("#334155")  // Slate Border
		ts.DangerColor = lipgloss.Color("#EF4444")  // Coral Red
		ts.SuccessColor = lipgloss.Color("#22C55E") // Bright Green
	}

	ts.Title = lipgloss.NewStyle().Bold(true).Foreground(ts.PrimaryColor)
	ts.Subtitle = lipgloss.NewStyle().Foreground(ts.MutedColor).Italic(true)

	ts.TopBar = lipgloss.NewStyle().
		Background(ts.CardBg).
		Foreground(ts.TextColor).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ts.BorderColor)

	ts.Badge = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.PrimaryColor)

	ts.Pill = lipgloss.NewStyle().
		Foreground(ts.TextColor).
		Background(ts.BgColor).
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
		Background(ts.CardBg).
		Padding(0, 2)

	ts.Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ts.BorderColor).
		Background(ts.CardBg).
		Foreground(ts.TextColor).
		Padding(1, 2)

	ts.CardActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ts.PrimaryColor).
		Background(ts.CardBg).
		Foreground(ts.TextColor).
		Padding(1, 2)

	ts.Modal = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(ts.PrimaryColor).
		Background(ts.CardBg).
		Foreground(ts.TextColor).
		Padding(1, 2)

	ts.RowSelected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.AccentColor).
		Padding(0, 1)

	ts.RowNormal = lipgloss.NewStyle().
		Foreground(ts.TextColor).
		Padding(0, 1)

	ts.Footer = lipgloss.NewStyle().
		Background(ts.CardBg).
		Foreground(ts.MutedColor).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ts.BorderColor)

	ts.ButtonActive = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.PrimaryColor).
		Padding(0, 2)

	ts.ButtonNormal = lipgloss.NewStyle().
		Foreground(ts.TextColor).
		Background(ts.BgColor).
		Border(lipgloss.NormalBorder()).
		BorderForeground(ts.BorderColor).
		Padding(0, 1)

	return ts
}
