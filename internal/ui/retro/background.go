package retro

import (
	"fmt"
	"strings"
	"time"

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

// CanvasCell represents a single cell in the full-screen 2D background canvas
type CanvasCell struct {
	Char  rune
	FG    string // Full brightness TrueColor ANSI sequence (for outside panels)
	DimFG string // Translucent dimmed TrueColor ANSI sequence (for inside panels)
}

// VisualCell represents a parsed foreground cell with active ANSI styling
type VisualCell struct {
	Char rune
	ANSI string
}

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

func getDimColor(baseR, baseG, baseB uint8, dimLevel string) string {
	var factor float64
	switch dimLevel {
	case "low":
		factor = 0.35 // brighter, more visible
	case "high":
		factor = 0.12 // very dark/subtle
	default: // "med"
		factor = 0.22 // balanced medium
	}

	r := uint8(float64(baseR) * factor)
	g := uint8(float64(baseG) * factor)
	b := uint8(float64(baseB) * factor)

	if r < 12 {
		r = 12
	}
	if g < 15 {
		g = 15
	}
	if b < 24 {
		b = 24
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func hexToBgANSI(color lipgloss.Color) string {
	hex := string(color)
	if len(hex) == 7 && hex[0] == '#' {
		var r, g, b uint8
		_, _ = fmt.Sscanf(hex[1:], "%02x%02x%02x", &r, &g, &b)
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
	}
	return "\x1b[48;2;17;24;39m"
}

// makeCanvasCell constructs a CanvasCell with both vibrant and dimmed ANSI TrueColor
func makeCanvasCell(r rune, red, green, blue uint8, dimLevel string) CanvasCell {
	return CanvasCell{
		Char:  r,
		FG:    fmt.Sprintf("\x1b[38;2;%d;%d;%dm", red, green, blue),
		DimFG: getDimColor(red, green, blue, dimLevel),
	}
}

// Generate2DCanvas produces a full 2D screen grid of CanvasCells for the given preset
func Generate2DCanvas(bgType AnimatedBGType, frame, width, height int, dimLevel string) [][]CanvasCell {
	if width < 10 {
		width = 10
	}
	if height < 6 {
		height = 6
	}

	activeType, _ := ResolveActiveBG(bgType, time.Now())
	canvas := make([][]CanvasCell, height)
	for y := 0; y < height; y++ {
		canvas[y] = make([]CanvasCell, width)
		for x := 0; x < width; x++ {
			canvas[y][x] = CanvasCell{Char: ' '}
		}
	}

	skyH := (height * 35) / 100
	if skyH < 2 {
		skyH = 2
	}
	horizonH := (height * 70) / 100
	if horizonH <= skyH {
		horizonH = skyH + 2
	}

	switch activeType {
	case BGSawah:
		fillSawahCanvas(canvas, frame, width, height, skyH, horizonH, dimLevel)
	case BGGunung:
		fillGunungCanvas(canvas, frame, width, height, skyH, horizonH, dimLevel)
	case BGSakalibel:
		fillSakalibelCanvas(canvas, frame, width, height, skyH, horizonH, dimLevel)
	case BGSenja:
		fillSenjaCanvas(canvas, frame, width, height, skyH, horizonH, dimLevel)
	case BGMalam:
		fillMalamCanvas(canvas, frame, width, height, skyH, horizonH, dimLevel)
	case BGCurug:
		fillCurugCanvas(canvas, frame, width, height, skyH, horizonH, dimLevel)
	default:
		fillSawahCanvas(canvas, frame, width, height, skyH, horizonH, dimLevel)
	}

	return canvas
}

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

func overlayCanvasRunes(row []CanvasCell, runes []rune, startPos int, r, g, b uint8, dimLevel string) {
	rLen := len(runes)
	rowLen := len(row)
	for i := 0; i < rLen; i++ {
		idx := startPos + i
		if idx >= 0 && idx < rowLen {
			row[idx] = makeCanvasCell(runes[i], r, g, b, dimLevel)
		}
	}
}

// 1. SAWAH TERASERING 2D GENERATOR
func fillSawahCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int, dimLevel string) {
	// Sky zone (0 to skyH): Clouds, birds, morning sun
	for y := 0; y < skyH; y++ {
		var pattern string
		var red, green, blue uint8
		switch y % 3 {
		case 0:
			pattern = "          ☁️                        ~v~                ☁️           ~v~    ~v~          "
			red, green, blue = 147, 197, 253 // Light Sky Blue
		case 1:
			pattern = "   ~v~                 ☀️  Matahari Pagi                 ~v~                     ☁️     "
			red, green, blue = 253, 224, 71 // Sun Amber
		default:
			pattern = "                ☁️                                               ☁️                     "
			red, green, blue = 224, 242, 254 // Wispy Cloud White
		}
		runes := getScrolledRunes(pattern, (frame/5)+y*7, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], red, green, blue, dimLevel)
		}
	}

	// Horizon zone (skyH to horizonH): Hills, Saung gubuk, walking farmer & buffalo
	for y := skyH; y < horizonH; y++ {
		ridgePattern := "  /~~~~~~~~~~~~~~~~\\            /~~~~~~~~\\                    /~~~~~~~~~~~~~~~~\\         "
		runes := getScrolledRunes(ridgePattern, (frame/3)+y*11, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 99, 102, 241, dimLevel) // Mountain indigo
		}

		if y == horizonH-1 {
			// Saung and walking farmer & buffalo trail
			step := (frame / 6) % 2
			var farmerStr string
			if step == 0 {
				farmerStr = "[🛖 Saung]    (o_ /|\\ [===|   (o o)____/ @@"
			} else {
				farmerStr = "[🛖 Saung]    (o_ /|\\_ [===|  (- -)____/ @@"
			}
			actorPos := width / 4
			overlayCanvasRunes(c[y], []rune(farmerStr), actorPos, 253, 224, 71, dimLevel)
		}
	}

	// Ground zone (horizonH to height): Terraced rice terraces & irrigation stream
	for y := horizonH; y < height; y++ {
		var pattern string
		var red, green, blue uint8
		if y%2 == 0 {
			pattern = "🌾 ,,🌾,, ,,,🌾,,, _|||||_  🌾 ,,🌾,, ,,,🌾,,, _|||||_  🌾 ,,🌾,, ,,,🌾,,, _|||||_  "
			red, green, blue = 34, 197, 94 // Lush emerald rice shoots
		} else {
			pattern = "  ~ ~ ≈ ≈ ~ ~  Lembah Subur Slamet  ~ ~ ≈ ≈ ~ ~  Pematang Sawah  ~ ~ ≈ ≈ ~ ~  "
			red, green, blue = 56, 189, 248 // Irrigation water cyan
		}
		runes := getScrolledRunes(pattern, frame+y*5, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], red, green, blue, dimLevel)
		}
	}
}

// 2. GUNUNG SLAMET 2D GENERATOR
func fillGunungCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int, dimLevel string) {
	// Peak zone
	for y := 0; y < skyH; y++ {
		peakPattern := "              /▲▲▲\\ ^~^~^                 ( ☁️ )                 /▲▲▲\\ ^~^~^                 "
		runes := getScrolledRunes(peakPattern, (frame/6)+y*5, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 165, 180, 252, dimLevel)
		}
	}

	// Slopes and pine forests
	for y := skyH; y < horizonH; y++ {
		var pattern string
		if y%2 == 0 {
			pattern = "       /═══════════════════\\     ▲▲▲▲  ▲▲▲  ▲▲▲▲       /═══════════════════\\     ▲▲▲▲  ▲▲▲   "
		} else {
			pattern = "    ▲▲▲▲▲▲  ▲▲▲▲  ▲▲▲▲▲▲     /═════════════════\\    ▲▲▲▲▲▲  ▲▲▲▲  ▲▲▲▲▲▲     /═════════\\   "
		}
		runes := getScrolledRunes(pattern, (frame/3)+y*8, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 99, 102, 241, dimLevel)
		}
	}

	// Valley & streams
	for y := horizonH; y < height; y++ {
		valleyPattern := "~~≈≈~~ Kali Pedes ~~≈≈~~  ❀ ✿ ❀  ~~≈≈~~ Kali Keruh ~~≈≈~~  ❀ ✿ ❀  ~~≈≈~~ Lembah Slamet ~~≈≈~~ "
		runes := getScrolledRunes(valleyPattern, frame+y*4, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 6, 182, 212, dimLevel)
		}
	}
}

// 3. KERETA API SAKALIBEL 2D GENERATOR
func fillSakalibelCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int, dimLevel string) {
	// Sky with steam puffs
	for y := 0; y < skyH; y++ {
		steamPattern := "        o   O   °   (awan)               o   O   °   (awan)               o   O   °         "
		runes := getScrolledRunes(steamPattern, (-frame/2)+y*9, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 226, 232, 240, dimLevel)
		}
	}

	// Rails and express train
	for y := skyH; y < horizonH; y++ {
		railPattern := "════════════════════════════════════════════════════════════════════════════════════════"
		runes := getScrolledRunes(railPattern, 0, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 148, 163, 184, dimLevel)
		}
		if y == horizonH-2 {
			trainStr := "[🚂 KAMANDAKA]═[🚃]═[🚃]═[🚃]═[🚃]"
			trainPos := ((frame * 2) % (width + len(trainStr))) - len(trainStr)
			overlayCanvasRunes(c[y], []rune(trainStr), trainPos, 245, 158, 11, dimLevel)
		}
	}

	// Arches & River Glagah
	for y := horizonH; y < height; y++ {
		archPattern := "   ||  (◠◠)  ||        ||  (◠◠)  ||   ~~~≈≈≈ Glagah & Kalikeruh ~~~≈≈≈   ||  (◠◠)  ||        "
		runes := getScrolledRunes(archPattern, (frame/4)+y*6, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 56, 189, 248, dimLevel)
		}
	}
}

// 4. LENTERA SENJA 2D GENERATOR
func fillSenjaCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int, dimLevel string) {
	for y := 0; y < skyH; y++ {
		skyPattern := "──────•─────────•───────•─────── ~v~   ~v~ ──────•─────────•───────•─────── ~v~   ~v~ "
		runes := getScrolledRunes(skyPattern, (frame/4)+y*7, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 249, 115, 22, dimLevel)
		}
	}

	for y := skyH; y < horizonH; y++ {
		lanternPattern := "     [🏮]               [🏮]               [🏮]               [🏮]               [🏮]     "
		runes := getScrolledRunes(lanternPattern, (frame/3)+y*5, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 239, 68, 68, dimLevel)
		}
	}

	for y := horizonH; y < height; y++ {
		pociPattern := "♨️ Teh Poci Gula Batu Bumiayu ♨️  ░░▒▒▓▓  ♨️ Mendoan Anget ♨️  ░░▒▒▓▓  ♨️ Kopi Hangat ♨️    "
		runes := getScrolledRunes(pociPattern, (frame/2)+y*4, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 251, 191, 36, dimLevel)
		}
	}
}

// 5. LANGIT MALAM BERBINTANG 2D GENERATOR
func fillMalamCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int, dimLevel string) {
	for y := 0; y < skyH; y++ {
		starPattern := "  ★     .       ✦       .    🌙    ★       .       ✦       .    ★       .       ✦       "
		runes := getScrolledRunes(starPattern, (frame/8)+y*13, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 224, 231, 255, dimLevel)
		}

		if y == 1 {
			meteorCycle := frame % 64
			if meteorCycle < 24 {
				meteorStr := "═══════★"
				meteorPos := meteorCycle * 4
				overlayCanvasRunes(c[y], []rune(meteorStr), meteorPos, 254, 240, 138, dimLevel)
			}
		}
	}

	for y := skyH; y < horizonH; y++ {
		hillPattern := "   .   ✦   /~~~~~~~~\\    ★    /~~~~~~~~~~~~~~~~\\    .   ✦   /~~~~~~~~\\   "
		runes := getScrolledRunes(hillPattern, (frame/4)+y*9, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 129, 140, 248, dimLevel)
		}
		if y == horizonH-1 {
			beacon := "🔴"
			if (frame/12)%2 == 0 {
				beacon = "⚪"
			}
			towerStr := "/▲\\ (" + beacon + ")"
			overlayCanvasRunes(c[y], []rune(towerStr), width/3, 239, 68, 68, dimLevel)
		}
	}

	for y := horizonH; y < height; y++ {
		fireflyPattern := "  ✨   .   ✨   .   ✨   .   ✨  ,,.,,,,.  ✨   .   ✨   .   ✨  ,,.,,,,.  ✨   .   ✨   "
		runes := getScrolledRunes(fireflyPattern, (frame/2)+y*5, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 52, 211, 153, dimLevel)
		}
	}
}

// 6. CURUG CIPENDOK 2D GENERATOR
func fillCurugCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int, dimLevel string) {
	for y := 0; y < skyH; y++ {
		cliffPattern := "████\\ |||||||| /█████████\\ |||||||| /█████████\\ |||||||| /█████████\\ |||||||| /████"
		runes := getScrolledRunes(cliffPattern, (frame/6)+y*6, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 100, 116, 139, dimLevel)
		}
	}

	for y := skyH; y < horizonH; y++ {
		fallCycle := (frame + y) % 3
		var fallPattern string
		switch fallCycle {
		case 0:
			fallPattern = "  ║ ░ ║ ▒ ║ ░ ║ ▒ ║  ░░▒▒  ║ ░ ║ ▒ ║ ░ ║ ▒ ║  ░░▒▒  ║ ░ ║ ▒ ║ ░ ║ ▒ ║  ░░▒▒  ║ ░ ║ ▒ ║ "
		case 1:
			fallPattern = "  ▒ ║ ░ ║ ▒ ║ ░ ║ ▒  ▒▒░░  ▒ ║ ░ ║ ▒ ║ ░ ║ ▒  ▒▒░░  ▒ ║ ░ ║ ▒ ║ ░ ║ ▒  ▒▒░░  ▒ ║ ░ ║ ▒ "
		default:
			fallPattern = "  ░ ║ ▒ ║ ░ ║ ▒ ║ ░  ░░▓▓  ░ ║ ▒ ║ ░ ║ ▒ ║ ░  ░░▓▓  ░ ║ ▒ ║ ░ ║ ▒ ║ ░  ░░▓▓  ░ ║ ▒ ║ ░ "
		}
		runes := getScrolledRunes(fallPattern, 0, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 56, 189, 248, dimLevel)
		}
	}

	for y := horizonH; y < height; y++ {
		poolPattern := "≈≈≈≈  Percikan Air Dingin Lembah Slamet  ≈≈≈≈ (●) (○) ≈≈≈≈  Curug Cipendok  ≈≈≈≈ (●) (○) "
		runes := getScrolledRunes(poolPattern, frame+y*4, width)
		for x := 0; x < width; x++ {
			c[y][x] = makeCanvasCell(runes[x], 14, 165, 233, dimLevel)
		}
	}
}

// parseANSILine parses one line of ANSI text into visual cells with active ANSI codes
func parseANSILine(line string, targetWidth int) []VisualCell {
	cells := make([]VisualCell, targetWidth)
	for i := range cells {
		cells[i] = VisualCell{Char: ' '}
	}

	runes := []rune(line)
	col := 0
	currentANSI := ""
	inEscape := false
	var escapeBuf []rune

	for i := 0; i < len(runes) && col < targetWidth; i++ {
		r := runes[i]
		if r == 0x1b {
			inEscape = true
			escapeBuf = []rune{r}
			continue
		}
		if inEscape {
			escapeBuf = append(escapeBuf, r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
				currentANSI += string(escapeBuf)
				escapeBuf = nil
			}
			continue
		}

		cells[col] = VisualCell{
			Char: r,
			ANSI: currentANSI,
		}
		col++
	}

	return cells
}

// CompositeCanvas composites foreground Window Manager text over the full-screen 2D background canvas
func CompositeCanvas(canvas [][]CanvasCell, fgText string, width, height int, dimLevel string, cardBg lipgloss.Color) string {
	if len(canvas) == 0 || width < 10 || height < 4 {
		return fgText
	}

	lines := strings.Split(fgText, "\n")
	cardBgANSI := hexToBgANSI(cardBg)

	var sb strings.Builder

	for y := 0; y < height; y++ {
		var fgLine string
		if y < len(lines) {
			fgLine = lines[y]
		}

		vCells := parseANSILine(fgLine, width)
		var bgRow []CanvasCell
		if y < len(canvas) {
			bgRow = canvas[y]
		}

		inPanel := false

		for x := 0; x < width; x++ {
			vCell := vCells[x]
			var bgCell CanvasCell
			if x < len(bgRow) {
				bgCell = bgRow[x]
			} else {
				bgCell = CanvasCell{Char: ' '}
			}

			// Toggle panel state on vertical borders
			if vCell.Char == '│' || vCell.Char == '║' {
				inPanel = !inPanel
			}

			if vCell.Char == ' ' {
				// Empty space: reveal animated background!
				if inPanel {
					// Inside panel: render translucent dimmed background character with card background!
					if bgCell.Char != 0 && bgCell.Char != ' ' {
						sb.WriteString(bgCell.DimFG + cardBgANSI + string(bgCell.Char) + "\x1b[0m")
					} else {
						sb.WriteString(cardBgANSI + " \x1b[0m")
					}
				} else {
					// Outside panel (gutter/margin): render 100% full vibrant color!
					if bgCell.Char != 0 && bgCell.Char != ' ' {
						sb.WriteString(bgCell.FG + string(bgCell.Char) + "\x1b[0m")
					} else {
						sb.WriteString(" ")
					}
				}
			} else {
				// Non-space: active foreground text, numbers, meters, borders!
				if inPanel {
					sb.WriteString(vCell.ANSI + cardBgANSI + string(vCell.Char) + "\x1b[0m")
				} else {
					sb.WriteString(vCell.ANSI + string(vCell.Char) + "\x1b[0m")
				}
			}
		}

		if y < height-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
