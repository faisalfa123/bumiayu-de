package retro

import (
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// MascotMood defines facial/body expression of the retro mascot
type MascotMood int

const (
	MoodNormal MascotMood = iota
	MoodHappy
	MoodWork
	MoodAlert
)

// GetAnimatedMascot returns an 8-bit retro pixel mascot sprite with 4-frame animation loop
func GetAnimatedMascot(frame int, mood MascotMood, st styles.ThemeStyles) string {
	// 24 FPS timing: change frame every 6 ticks (4 animation changes per second)
	f := (frame / 6) % 4

	cPrimary := lipgloss.NewStyle().Foreground(st.PrimaryColor).Bold(true)
	cAccent := lipgloss.NewStyle().Foreground(st.AccentColor).Bold(true)
	cDanger := lipgloss.NewStyle().Foreground(st.DangerColor).Bold(true)
	cMuted := lipgloss.NewStyle().Foreground(st.MutedColor)

	var antenna, eyes, mouth, feet string

	switch f {
	case 0:
		// Frame 0: Idle Standing, antenna blink left
		antenna = cAccent.Render(" ⚡ ")
		feet = cPrimary.Render(" /  \\ ")
	case 1:
		// Frame 1: Slight bounce down
		antenna = cPrimary.Render(" ✦ ")
		feet = cPrimary.Render(" |  | ")
	case 2:
		// Frame 2: Antenna blink right, arms wiggle
		antenna = cAccent.Render("  ⚡")
		feet = cPrimary.Render(" \\  / ")
	case 3:
		// Frame 3: Return bounce
		antenna = cMuted.Render(" · ")
		feet = cPrimary.Render(" |  | ")
	}

	switch mood {
	case MoodHappy:
		eyes = "^ _ ^"
		mouth = "▿"
	case MoodWork:
		eyes = "• ᴗ •"
		mouth = "—"
	case MoodAlert:
		eyes = "O _ O"
		mouth = "!"
		antenna = cDanger.Render(" ⚠️ ")
	default:
		// Normal mood with periodic blink
		if f == 1 {
			eyes = "- _ -"
		} else {
			eyes = "• _ •"
		}
		mouth = "o"
	}

	// Compact 3-line retro robot mascot
	line1 := "   " + antenna
	line2 := cPrimary.Render(" [ ") + cAccent.Render(eyes) + " " + cMuted.Render(mouth) + cPrimary.Render(" ]")
	line3 := "  " + feet

	return lipgloss.JoinVertical(lipgloss.Left, line1, line2, line3)
}

// GetBannerMascot returns a compact single-line badge with animated mascot icon
func GetBannerMascot(frame int, st styles.ThemeStyles) string {
	f := (frame / 6) % 4
	sprites := []string{
		"👾 (•_•)",
		"👾 (•ᴗ•)",
		"👾 (⌐■_■)",
		"👾 (^-^)",
	}
	sprite := sprites[f]
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(st.PrimaryColor).
		Padding(0, 1).
		Render(sprite + " BUMIAYU-DE")
}

// GetPixelIsland returns a retro 8-bit mountain & coconut tree landscape
func GetPixelIsland(frame int, st styles.ThemeStyles, width int) string {
	if width < 30 {
		return ""
	}
	sunY := (frame / 12) % 2
	sun := "☀️"
	if sunY == 1 {
		sun = "✨"
	}

	cSun := lipgloss.NewStyle().Foreground(st.PrimaryColor).Render(sun)
	cMountain := lipgloss.NewStyle().Foreground(st.AccentColor).Render("   /\\     /\\/\\   ")
	cTrees := lipgloss.NewStyle().Foreground(st.SuccessColor).Render("  🌴 /\\/\\ 🌴 /\\/\\🌴 ")
	cGround := lipgloss.NewStyle().Foreground(st.BorderColor).Render(strings.Repeat("━", 24))

	return lipgloss.JoinVertical(lipgloss.Center,
		"        "+cSun,
		cMountain,
		cTrees,
		cGround,
	)
}
