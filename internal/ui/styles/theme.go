package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type ThemeMode int

const (
	ThemeCyberArcade ThemeMode = iota // 16-bit Cyber-Arcade Synthwave (Default)
	ThemeDOSNorton                    // Classic DOS / Norton Commander (Double-line blue/yellow)
	ThemeGameBoy                      // 8-bit Retro GameBoy 4-shade Green
	ThemeDark                         // Modern Midnight Obsidian
	ThemeWarm                         // Warm Cream / Amber
	ThemeMacLight                     // iPadOS Clean Light
	ThemeClassic                      // Classic Slate Light
)

type ThemeStyles struct {
	Mode         ThemeMode
	Name         string
	BgColor      lipgloss.Color
	CardBg       lipgloss.Color
	TextColor    lipgloss.Color
	MutedColor   lipgloss.Color
	PrimaryColor lipgloss.Color
	AccentColor  lipgloss.Color
	BorderColor  lipgloss.Color
	DangerColor  lipgloss.Color
	SuccessColor lipgloss.Color
	GlowColor    lipgloss.Color

	// Retro meter colors
	MeterHigh lipgloss.Color
	MeterMid  lipgloss.Color
	MeterLow  lipgloss.Color

	// Pre-built styles
	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	TopBar       lipgloss.Style
	Badge        lipgloss.Style
	Pill         lipgloss.Style
	PillWarning  lipgloss.Style
	PillDanger   lipgloss.Style
	ActiveTab    lipgloss.Style
	InactiveTab  lipgloss.Style
	Card         lipgloss.Style
	CardActive   lipgloss.Style
	Modal        lipgloss.Style
	RowSelected  lipgloss.Style
	RowNormal    lipgloss.Style
	Footer       lipgloss.Style
	ButtonActive lipgloss.Style
	ButtonNormal lipgloss.Style
}

// GetTheme returns fully configured styles for requested theme
func GetTheme(mode ThemeMode) ThemeStyles {
	ts := ThemeStyles{Mode: mode}

	switch mode {
	case ThemeCyberArcade:
		// 16-Bit Cyber-Arcade Synthwave
		ts.Name = "16-Bit Cyber-Arcade"
		ts.BgColor = lipgloss.Color("#0A0E1A")      // Deep Cyber Navy
		ts.CardBg = lipgloss.Color("#111827")       // Dark Slate Void
		ts.TextColor = lipgloss.Color("#F3F4F6")    // Crisp White
		ts.MutedColor = lipgloss.Color("#9CA3AF")   // Silver Muted
		ts.PrimaryColor = lipgloss.Color("#00F5FF") // Electric Neon Cyan
		ts.AccentColor = lipgloss.Color("#E024A5")  // Neon Magenta / Violet
		ts.BorderColor = lipgloss.Color("#374151")  // Steel Gray Border
		ts.DangerColor = lipgloss.Color("#FF3366")  // Laser Red
		ts.SuccessColor = lipgloss.Color("#00FF88") // High-Voltage Neon Green
		ts.GlowColor = lipgloss.Color("#00F5FF")    // Cyan Glow

		ts.MeterHigh = lipgloss.Color("#FF3366")
		ts.MeterMid = lipgloss.Color("#FFCC00")
		ts.MeterLow = lipgloss.Color("#00FF88")

	case ThemeDOSNorton:
		// Classic DOS / Norton Commander 90s
		ts.Name = "Classic DOS / Norton"
		ts.BgColor = lipgloss.Color("#000084")      // Classic DOS Navy Blue
		ts.CardBg = lipgloss.Color("#0000AA")       // Pure Blue Panel
		ts.TextColor = lipgloss.Color("#FFFFFF")    // Pure Crisp White
		ts.MutedColor = lipgloss.Color("#00AAAA")   // Cyan Muted
		ts.PrimaryColor = lipgloss.Color("#FFFF55") // High-Intensity Yellow
		ts.AccentColor = lipgloss.Color("#55FFFF")  // Bright Cyan
		ts.BorderColor = lipgloss.Color("#AAAAAA")  // Gray Box
		ts.DangerColor = lipgloss.Color("#FF5555")  // Bright Red
		ts.SuccessColor = lipgloss.Color("#55FF55") // Bright Green
		ts.GlowColor = lipgloss.Color("#FFFF55")    // Bright Yellow Glow

		ts.MeterHigh = lipgloss.Color("#FF5555")
		ts.MeterMid = lipgloss.Color("#FFFF55")
		ts.MeterLow = lipgloss.Color("#55FF55")

	case ThemeGameBoy:
		// 8-Bit GameBoy 4-Shade Monochromatic Green
		ts.Name = "8-Bit GameBoy Green"
		ts.BgColor = lipgloss.Color("#0F380F")      // Deepest Green
		ts.CardBg = lipgloss.Color("#1F4B1F")       // Dark Forest
		ts.TextColor = lipgloss.Color("#9BBC0F")    // Brightest Yellow-Green
		ts.MutedColor = lipgloss.Color("#8BAC0F")   // Mid Light Green
		ts.PrimaryColor = lipgloss.Color("#9BBC0F") // Crisp Pixel Green
		ts.AccentColor = lipgloss.Color("#8BAC0F")  // Mid Accent
		ts.BorderColor = lipgloss.Color("#306230")  // Mid Forest
		ts.DangerColor = lipgloss.Color("#E76F51")  // Warm Brick
		ts.SuccessColor = lipgloss.Color("#9BBC0F") // GameBoy Green
		ts.GlowColor = lipgloss.Color("#9BBC0F")

		ts.MeterHigh = lipgloss.Color("#E76F51")
		ts.MeterMid = lipgloss.Color("#8BAC0F")
		ts.MeterLow = lipgloss.Color("#9BBC0F")

	case ThemeDark:
		// Modern Midnight Dark
		ts.Name = "Obsidian Dark"
		ts.BgColor = lipgloss.Color("#0B0F17")
		ts.CardBg = lipgloss.Color("#161D2B")
		ts.TextColor = lipgloss.Color("#F8FAFC")
		ts.MutedColor = lipgloss.Color("#94A3B8")
		ts.PrimaryColor = lipgloss.Color("#10B981")
		ts.AccentColor = lipgloss.Color("#38BDF8")
		ts.BorderColor = lipgloss.Color("#1E293B")
		ts.DangerColor = lipgloss.Color("#F43F5E")
		ts.SuccessColor = lipgloss.Color("#22C55E")
		ts.GlowColor = lipgloss.Color("#38BDF8")

		ts.MeterHigh = lipgloss.Color("#F43F5E")
		ts.MeterMid = lipgloss.Color("#FBBF24")
		ts.MeterLow = lipgloss.Color("#22C55E")

	case ThemeWarm:
		// Warm Cream / Ivory
		ts.Name = "Warm Ivory"
		ts.BgColor = lipgloss.Color("#FEFCE8")
		ts.CardBg = lipgloss.Color("#FFFDF5")
		ts.TextColor = lipgloss.Color("#292524")
		ts.MutedColor = lipgloss.Color("#78716C")
		ts.PrimaryColor = lipgloss.Color("#D97706")
		ts.AccentColor = lipgloss.Color("#059669")
		ts.BorderColor = lipgloss.Color("#D6D3D1")
		ts.DangerColor = lipgloss.Color("#E11D48")
		ts.SuccessColor = lipgloss.Color("#15803D")
		ts.GlowColor = lipgloss.Color("#D97706")

		ts.MeterHigh = lipgloss.Color("#E11D48")
		ts.MeterMid = lipgloss.Color("#D97706")
		ts.MeterLow = lipgloss.Color("#15803D")

	case ThemeMacLight:
		// Clean Apple / iPadOS Light
		ts.Name = "Clean Mac Light"
		ts.BgColor = lipgloss.Color("#E5E7EB")
		ts.CardBg = lipgloss.Color("#FFFFFF")
		ts.TextColor = lipgloss.Color("#0F172A")
		ts.MutedColor = lipgloss.Color("#475569")
		ts.PrimaryColor = lipgloss.Color("#2563EB")
		ts.AccentColor = lipgloss.Color("#0284C7")
		ts.BorderColor = lipgloss.Color("#CBD5E1")
		ts.DangerColor = lipgloss.Color("#DC2626")
		ts.SuccessColor = lipgloss.Color("#16A34A")
		ts.GlowColor = lipgloss.Color("#2563EB")

		ts.MeterHigh = lipgloss.Color("#DC2626")
		ts.MeterMid = lipgloss.Color("#F59E0B")
		ts.MeterLow = lipgloss.Color("#16A34A")

	default: // ThemeClassic
		ts.Name = "Classic Slate"
		ts.BgColor = lipgloss.Color("#F1F5F9")
		ts.CardBg = lipgloss.Color("#FFFFFF")
		ts.TextColor = lipgloss.Color("#0F172A")
		ts.MutedColor = lipgloss.Color("#64748B")
		ts.PrimaryColor = lipgloss.Color("#0284C7")
		ts.AccentColor = lipgloss.Color("#3B82F6")
		ts.BorderColor = lipgloss.Color("#CBD5E1")
		ts.DangerColor = lipgloss.Color("#EF4444")
		ts.SuccessColor = lipgloss.Color("#10B981")
		ts.GlowColor = lipgloss.Color("#0284C7")

		ts.MeterHigh = lipgloss.Color("#EF4444")
		ts.MeterMid = lipgloss.Color("#F59E0B")
		ts.MeterLow = lipgloss.Color("#10B981")
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
		Foreground(lipgloss.Color("#000000")).
		Background(ts.PrimaryColor)

	ts.Pill = lipgloss.NewStyle().
		Padding(0, 1)

	ts.PillWarning = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(ts.AccentColor).
		Padding(0, 1)

	ts.PillDanger = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(ts.DangerColor).
		Padding(0, 1)

	ts.ActiveTab = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(ts.PrimaryColor).
		Padding(0, 2)

	ts.InactiveTab = lipgloss.NewStyle().
		Foreground(ts.MutedColor).
		Padding(0, 2)

	// Box borders: For DOS mode we use double border, for Cyber we use rounded/thick
	cardBorder := lipgloss.RoundedBorder()
	if mode == ThemeDOSNorton {
		cardBorder = lipgloss.DoubleBorder()
	}

	ts.Card = lipgloss.NewStyle().
		Border(cardBorder).
		BorderForeground(ts.BorderColor).
		Padding(0, 1)

	ts.CardActive = lipgloss.NewStyle().
		Border(cardBorder).
		BorderForeground(ts.GlowColor).
		Bold(true).
		Padding(0, 1)

	modalBorder := lipgloss.ThickBorder()
	if mode == ThemeDOSNorton {
		modalBorder = lipgloss.DoubleBorder()
	}

	ts.Modal = lipgloss.NewStyle().
		Border(modalBorder).
		BorderForeground(ts.PrimaryColor).
		Background(ts.CardBg).
		Padding(1, 2)

	ts.RowSelected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(ts.PrimaryColor).
		Padding(0, 1)

	ts.RowNormal = lipgloss.NewStyle().
		Foreground(ts.TextColor).
		Padding(0, 1)

	ts.Footer = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ts.BorderColor)

	ts.ButtonActive = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(ts.PrimaryColor).
		Padding(0, 2)

	ts.ButtonNormal = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ts.BorderColor).
		Foreground(ts.TextColor).
		Padding(0, 1)

	return ts
}

// NextTheme cycles to next available theme
func NextTheme(current ThemeMode) ThemeMode {
	return (current + 1) % (ThemeClassic + 1)
}

// ThemeName returns human friendly string
func ThemeName(mode ThemeMode) string {
	switch mode {
	case ThemeCyberArcade:
		return "16-Bit Cyber-Arcade (Default)"
	case ThemeDOSNorton:
		return "Classic DOS / Norton (90s)"
	case ThemeGameBoy:
		return "8-Bit GameBoy Green"
	case ThemeDark:
		return "Obsidian Dark"
	case ThemeWarm:
		return "Warm Ivory"
	case ThemeMacLight:
		return "Clean Mac Light"
	default:
		return fmt.Sprintf("Theme #%d", mode)
	}
}
