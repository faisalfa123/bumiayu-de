package retro

import (
	"strings"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// AnimatedBGType defines switchable animated background presets
type AnimatedBGType string

const (
	BGOff       AnimatedBGType = "off"
	BGAuto      AnimatedBGType = "auto"
	BGSawah     AnimatedBGType = "sawah"
	BGGunung    AnimatedBGType = "gunung"
	BGSenja     AnimatedBGType = "senja"
	BGMalam     AnimatedBGType = "malam"
	BGCurug     AnimatedBGType = "curug"
	BGSakalibel AnimatedBGType = "sakalibel"
)

// NextBG cycles to the next background preset for shortcut toggle [B]
func NextBG(current AnimatedBGType) AnimatedBGType {
	switch current {
	case BGAuto:
		return BGSawah
	case BGSawah:
		return BGGunung
	case BGGunung:
		return BGSenja
	case BGSenja:
		return BGSakalibel
	case BGSakalibel:
		return BGMalam
	case BGMalam:
		return BGCurug
	case BGCurug:
		return BGOff
	default:
		return BGAuto
	}
}

// BGName returns human-readable label with emoji
func BGName(bg AnimatedBGType) string {
	switch bg {
	case BGOff:
		return "Mati (Minimalis Gelap)"
	case BGAuto:
		return "Dinamis Jam Server (Auto)"
	case BGSawah:
		return "🌾 Sawah & Kerbau Berjalan"
	case BGGunung:
		return "🏔️ Gunung Slamet & Awan"
	case BGSenja:
		return "🏮 Lentera Senja Bumiayu"
	case BGSakalibel:
		return "🚂 Kereta Api Sakalibel"
	case BGMalam:
		return "🌌 Langit Malam & Meteor"
	case BGCurug:
		return "🌊 Curug Cipendok Air Terjun"
	default:
		return string(bg)
	}
}

// ResolveActiveBG resolves dynamic time for BGAuto based on current hour
func ResolveActiveBG(bgType AnimatedBGType, t time.Time) (AnimatedBGType, string) {
	if bgType != BGAuto {
		return bgType, ""
	}

	hour := t.Hour()
	switch {
	case hour >= 5 && hour < 11:
		return BGGunung, "Pagi: Fajar Slamet"
	case hour >= 11 && hour < 15:
		return BGSawah, "Siang: Sawah Subur"
	case hour >= 15 && hour < 18:
		return BGSenja, "Sore: Lentera Senja"
	case hour >= 18 && hour < 21:
		return BGSakalibel, "Malam: Sakalibel Ekspres"
	default:
		return BGMalam, "Larut: Langit Berbintang"
	}
}

// RenderAnimatedBackground renders a moving 24 FPS ASCII / pixel art panorama runner
func RenderAnimatedBackground(bgType AnimatedBGType, frame int, width int, st styles.ThemeStyles, maxLines int) string {
	if bgType == BGOff || width < 20 || maxLines < 2 {
		return ""
	}

	activeType, timeTag := ResolveActiveBG(bgType, time.Now())

	contentW := width - 4
	if contentW < 20 {
		contentW = 20
	}

	var title string
	var line1, line2, line3 string

	switch activeType {
	case BGSawah:
		title = "🌾 SAWAH TERASERING & KERBAU RUNNER"
		if timeTag != "" {
			title += " [" + timeTag + "]"
		}
		line1, line2, line3 = renderSawahLines(frame, contentW)

	case BGGunung:
		title = "🏔️ GUNUNG SLAMET & AWAN BERARAK"
		if timeTag != "" {
			title += " [" + timeTag + "]"
		}
		line1, line2, line3 = renderGunungLines(frame, contentW)

	case BGSakalibel:
		title = "🚂 KERETA API JEMBATAN SAKALIBEL BUMIAYU"
		if timeTag != "" {
			title += " [" + timeTag + "]"
		}
		line1, line2, line3 = renderSakalibelLines(frame, contentW)

	case BGSenja:
		title = "🏮 LENTERA SENJA & TEH POCI BUMIAYU"
		if timeTag != "" {
			title += " [" + timeTag + "]"
		}
		line1, line2, line3 = renderSenjaLines(frame, contentW)

	case BGMalam:
		title = "🌌 LANGIT MALAM BERBINTANG & METEOR"
		if timeTag != "" {
			title += " [" + timeTag + "]"
		}
		line1, line2, line3 = renderMalamLines(frame, contentW)

	case BGCurug:
		title = "🌊 CURUG CIPENDOK AIR TERJUN LEMBAH SLAMET"
		if timeTag != "" {
			title += " [" + timeTag + "]"
		}
		line1, line2, line3 = renderCurugLines(frame, contentW)

	default:
		title = "🌾 SAWAH TERASERING BUMIAYU"
		line1, line2, line3 = renderSawahLines(frame, contentW)
	}

	// Title banner line
	titleTag := lipgloss.NewStyle().
		Bold(true).
		Foreground(st.PrimaryColor).
		Render(" " + title + " • 24 FPS ")

	dividerLen := contentW - lipgloss.Width(titleTag)
	if dividerLen < 1 {
		dividerLen = 1
	}
	divider := lipgloss.NewStyle().
		Foreground(st.BorderColor).
		Render(strings.Repeat("─", dividerLen))
	headerLine := "  " + titleTag + divider

	if maxLines <= 2 {
		return lipgloss.JoinVertical(lipgloss.Left, headerLine, "  "+line2)
	}
	if maxLines == 3 {
		return lipgloss.JoinVertical(lipgloss.Left, headerLine, "  "+line1, "  "+line2)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		headerLine,
		"  "+line1,
		"  "+line2,
		"  "+line3,
	)
}

// getScrolledRunes safely creates a looped rune slice of exact targetWidth shifted by offset
func getScrolledRunes(pattern string, offset int, targetWidth int) []rune {
	pRunes := []rune(pattern)
	if len(pRunes) == 0 {
		return []rune(strings.Repeat(" ", targetWidth))
	}
	result := make([]rune, targetWidth)
	pLen := len(pRunes)
	normOffset := ((offset % pLen) + pLen) % pLen
	for i := 0; i < targetWidth; i++ {
		result[i] = pRunes[(normOffset+i)%pLen]
	}
	return result
}

// overlayRunes safely stamps overlay onto base runes at startPos
func overlayRunes(base []rune, overlay []rune, startPos int) []rune {
	bLen := len(base)
	oLen := len(overlay)
	if bLen == 0 || oLen == 0 {
		return base
	}

	for i := 0; i < oLen; i++ {
		idx := startPos + i
		if idx >= 0 && idx < bLen {
			base[idx] = overlay[i]
		}
	}
	return base
}

// 1. SAWAH TERASERING & KERBAU
func renderSawahLines(frame int, width int) (string, string, string) {
	// Sky: Slow parallax clouds and flying birds
	skyPattern := "          ☁️                        ~v~                ☁️           ~v~    ~v~          "
	skyRunes := getScrolledRunes(skyPattern, frame/5, width)
	line1 := lipgloss.NewStyle().Foreground(lipgloss.Color("#93C5FD")).Render(string(skyRunes))

	// Midground: Saung & walking farmer with buffalo
	midPattern := "   [🛖 Saung]             ,,,,,,                      [ Saung ]          ,,,,,,            "
	midRunes := getScrolledRunes(midPattern, frame/3, width)

	// Animate farmer walk cycle & buffalo head bob
	step := (frame / 6) % 2
	var farmerStr string
	if step == 0 {
		farmerStr = "(o_ /|\\ [===|   (o o)____/ @@"
	} else {
		farmerStr = "(o_ /|\\_ [===|  (- -)____/ @@"
	}
	actorPos := width / 4
	midRunes = overlayRunes(midRunes, []rune(farmerStr), actorPos)
	line2 := lipgloss.NewStyle().Foreground(lipgloss.Color("#FDE047")).Render(string(midRunes))

	// Foreground: Swaying lush green rice stalks and shimmering irrigation water
	forePattern := "🌾 ,,🌾,, ,,,🌾,,, _|||||_  ~ ~ ≈ ≈ ~ ~  🌾 ,,🌾,, ,,,🌾,,, _|||||_  ~ ~ ≈ ≈ ~ ~  "
	foreRunes := getScrolledRunes(forePattern, frame, width)
	line3 := lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")).Render(string(foreRunes))

	return line1, line2, line3
}

// 2. GUNUNG SLAMET & AWAN PARALLAX
func renderGunungLines(frame int, width int) (string, string, string) {
	// Peak with crater mist and drifting high clouds
	peakPattern := "              /▲▲▲\\ ^~^~^                 ( ☁️ )                 /▲▲▲\\ ^~^~^                 "
	peakRunes := getScrolledRunes(peakPattern, frame/6, width)
	line1 := lipgloss.NewStyle().Foreground(lipgloss.Color("#A5B4FC")).Render(string(peakRunes))

	// Slopes and dense pine forest foothills
	slopePattern := "       /═══════════════════\\     ▲▲▲▲  ▲▲▲  ▲▲▲▲       /═══════════════════\\     ▲▲▲▲  ▲▲▲   "
	slopeRunes := getScrolledRunes(slopePattern, frame/3, width)
	line2 := lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1")).Render(string(slopeRunes))

	// River valley: Kali Pedes & Kali Keruh with mountain wildflowers
	valleyPattern := "~~≈≈~~ Kali Pedes ~~≈≈~~  ❀ ✿ ❀  ~~≈≈~~ Kali Keruh ~~≈≈~~  ❀ ✿ ❀  ~~≈≈~~ Lembah Slamet ~~≈≈~~ "
	valleyRunes := getScrolledRunes(valleyPattern, frame, width)
	line3 := lipgloss.NewStyle().Foreground(lipgloss.Color("#06B6D4")).Render(string(valleyRunes))

	return line1, line2, line3
}

// 3. KERETA API SAKALIBEL EKSPRES
func renderSakalibelLines(frame int, width int) (string, string, string) {
	// Sky with steam puffs trailing from locomotive
	steamPattern := "        o   O   °   (awan)               o   O   °   (awan)               o   O   °         "
	steamRunes := getScrolledRunes(steamPattern, -frame/2, width)
	line1 := lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E8F0")).Render(string(steamRunes))

	// Express train moving rapidly across the bridge rails
	railPattern := "════════════════════════════════════════════════════════════════════════════════════════"
	railRunes := getScrolledRunes(railPattern, 0, width)
	trainStr := "[🚂 KAMANDAKA]═[🚃]═[🚃]═[🚃]═[🚃]"
	trainPos := ((frame * 2) % (width + len(trainStr))) - len(trainStr)
	railRunes = overlayRunes(railRunes, []rune(trainStr), trainPos)
	line2 := lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")).Render(string(railRunes))

	// Sakalibel Bridge tall arched pillars & Glagah river flowing beneath
	archPattern := "   ||  (◠◠)  ||        ||  (◠◠)  ||   ~~~≈≈≈ Glagah & Kalikeruh ~~~≈≈≈   ||  (◠◠)  ||        "
	archRunes := getScrolledRunes(archPattern, frame/4, width)
	line3 := lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Render(string(archRunes))

	return line1, line2, line3
}

// 4. LENTERA SENJA & TEH POCI
func renderSenjaLines(frame int, width int) (string, string, string) {
	// Sunset twilight sky with electrical wires & migrating birds
	skyPattern := "──────•─────────•───────•─────── ~v~   ~v~ ──────•─────────•───────•─────── ~v~   ~v~ "
	skyRunes := getScrolledRunes(skyPattern, frame/4, width)
	line1 := lipgloss.NewStyle().Foreground(lipgloss.Color("#F97316")).Render(string(skyRunes))

	// Hanging red lanterns with oscillating candle flicker
	lanternPattern := "     [🏮]               [🏮]               [🏮]               [🏮]               [🏮]     "
	lanternRunes := getScrolledRunes(lanternPattern, frame/3, width)
	line2 := lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Render(string(lanternRunes))

	// Angkringan street with steaming Teh Poci & Mendoan
	pociPattern := "♨️ Teh Poci Gula Batu Bumiayu ♨️  ░░▒▒▓▓  ♨️ Mendoan Anget ♨️  ░░▒▒▓▓  ♨️ Kopi Hangat ♨️    "
	pociRunes := getScrolledRunes(pociPattern, frame/2, width)
	line3 := lipgloss.NewStyle().Foreground(lipgloss.Color("#FBBF24")).Render(string(pociRunes))

	return line1, line2, line3
}

// 5. LANGIT MALAM BERBINTANG & METEOR
func renderMalamLines(frame int, width int) (string, string, string) {
	// Deep cosmic sky with twinkling stars and periodic meteor streak
	starPattern := "  ★     .       ✦       .    🌙    ★       .       ✦       .    ★       .       ✦       "
	starRunes := getScrolledRunes(starPattern, frame/8, width)

	// Meteor streak crossing sky every 64 frames
	meteorCycle := frame % 64
	if meteorCycle < 24 {
		meteorStr := "═══════★"
		meteorPos := meteorCycle * 4
		starRunes = overlayRunes(starRunes, []rune(meteorStr), meteorPos)
	}
	line1 := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E7FF")).Render(string(starRunes))

	// Nocturnal rolling hills & radio transmission tower with blinking red beacon
	beacon := "🔴"
	if (frame/12)%2 == 0 {
		beacon = "⚪"
	}
	towerStr := "/▲\\ (" + beacon + ")"
	hillPattern := "   .   ✦   /~~~~~~~~\\    ★    /~~~~~~~~~~~~~~~~\\    .   ✦   /~~~~~~~~\\   "
	hillRunes := getScrolledRunes(hillPattern, frame/4, width)
	hillRunes = overlayRunes(hillRunes, []rune(towerStr), width/3)
	line2 := lipgloss.NewStyle().Foreground(lipgloss.Color("#818CF8")).Render(string(hillRunes))

	// Fireflies dancing in the night meadow
	fireflyPattern := "  ✨   .   ✨   .   ✨   .   ✨  ,,.,,,,.  ✨   .   ✨   .   ✨  ,,.,,,,.  ✨   .   ✨   "
	fireflyRunes := getScrolledRunes(fireflyPattern, frame/2, width)
	line3 := lipgloss.NewStyle().Foreground(lipgloss.Color("#34D399")).Render(string(fireflyRunes))

	return line1, line2, line3
}

// 6. CURUG CIPENDOK AIR TERJUN
func renderCurugLines(frame int, width int) (string, string, string) {
	// Mountain rocky cliff plunge
	cliffPattern := "████\\ |||||||| /█████████\\ |||||||| /█████████\\ |||||||| /█████████\\ |||||||| /████"
	cliffRunes := getScrolledRunes(cliffPattern, frame/6, width)
	line1 := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B")).Render(string(cliffRunes))

	// Cascading vertical water streams
	fallCycle := frame % 3
	var fallPattern string
	switch fallCycle {
	case 0:
		fallPattern = "  ║ ░ ║ ▒ ║ ░ ║ ▒ ║  ░░▒▒  ║ ░ ║ ▒ ║ ░ ║ ▒ ║  ░░▒▒  ║ ░ ║ ▒ ║ ░ ║ ▒ ║  ░░▒▒  ║ ░ ║ ▒ ║ "
	case 1:
		fallPattern = "  ▒ ║ ░ ║ ▒ ║ ░ ║ ▒  ▒▒░░  ▒ ║ ░ ║ ▒ ║ ░ ║ ▒  ▒▒░░  ▒ ║ ░ ║ ▒ ║ ░ ║ ▒  ▒▒░░  ▒ ║ ░ ║ ▒ "
	default:
		fallPattern = "  ░ ║ ▒ ║ ░ ║ ▒ ║ ░  ░░▓▓  ░ ║ ▒ ║ ░ ║ ▒ ║ ░  ░░▓▓  ░ ║ ▒ ║ ░ ║ ▒ ║ ░  ░░▓▓  ░ ║ ▒ ║ ░ "
	}
	fallRunes := getScrolledRunes(fallPattern, 0, width)
	line2 := lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Render(string(fallRunes))

	// Swirling turquoise plunge pool rapids & spray
	poolPattern := "≈≈≈≈  Percikan Air Dingin Lembah Slamet  ≈≈≈≈ (●) (○) ≈≈≈≈  Curug Cipendok  ≈≈≈≈ (●) (○) "
	poolRunes := getScrolledRunes(poolPattern, frame, width)
	line3 := lipgloss.NewStyle().Foreground(lipgloss.Color("#0EA5E9")).Render(string(poolRunes))

	return line1, line2, line3
}
