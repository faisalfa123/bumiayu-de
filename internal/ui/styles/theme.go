package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Color Palette
	ColorPrimary   = lipgloss.Color("#10B981") // Jade Green (Bumiayu Emerald)
	ColorSecondary = lipgloss.Color("#06B6D4") // Cyan Horizon
	ColorAccent    = lipgloss.Color("#F59E0B") // Amber Glow
	ColorDanger    = lipgloss.Color("#EF4444") // Coral Red
	ColorSuccess   = lipgloss.Color("#22C55E") // Bright Green
	ColorBgDark    = lipgloss.Color("#0F172A") // Deep Slate Background
	ColorBgCard    = lipgloss.Color("#1E293B") // Card Background
	ColorText      = lipgloss.Color("#F8FAFC") // Crisp White
	ColorMuted     = lipgloss.Color("#94A3B8") // Slate Gray
	ColorBorder    = lipgloss.Color("#334155") // Subtle Border
	ColorHighlight = lipgloss.Color("#38BDF8") // Sky Blue

	// Typography & Borders
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	// Top Bar Component
	TopBarStyle = lipgloss.NewStyle().
			Background(ColorBgCard).
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorBorder)

	BadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Foreground(ColorBgDark).
			Background(ColorPrimary)

	PillNormal = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(ColorBgDark).
			Padding(0, 1)

	PillWarning = lipgloss.NewStyle().
			Foreground(ColorBgDark).
			Background(ColorAccent).
			Bold(true).
			Padding(0, 1)

	PillDanger = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(ColorDanger).
			Bold(true).
			Padding(0, 1)

	// Workspace Tabs
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBgDark).
			Background(ColorPrimary).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Background(ColorBgCard).
			Padding(0, 2)

	// Main Content Panel Box
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Background(ColorBgCard).
			Padding(1, 2)

	CardActiveStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Background(ColorBgCard).
			Padding(1, 2)

	// Senior Insight Box
	InsightBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Background(lipgloss.Color("#082f49")).
			Foreground(lipgloss.Color("#e0f2fe")).
			Padding(0, 1)

	InsightTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorHighlight)

	// Tables and Lists
	RowSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorBgDark).
				Background(ColorSecondary).
				Padding(0, 1)

	RowNormalStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Padding(0, 1)

	// Footer Bar
	FooterStyle = lipgloss.NewStyle().
			Background(ColorBgCard).
			Foreground(ColorMuted).
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorBorder)
)
