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
	BGStars     AnimatedBGType = "stars"
	BGMatrix    AnimatedBGType = "matrix"
	BGSawah     AnimatedBGType = "sawah"
	BGGunung    AnimatedBGType = "gunung"
	BGSenja     AnimatedBGType = "senja"
	BGMalam     AnimatedBGType = "malam"
	BGCurug     AnimatedBGType = "curug"
	BGSakalibel AnimatedBGType = "sakalibel"
)

// CanvasCell represents a single background character on the desktop canvas
type CanvasCell struct {
	Char rune
	FG   string // TrueColor ANSI sequence
}

// VisualCell represents a parsed foreground cell with active ANSI styling
type VisualCell struct {
	Char rune
	ANSI string
}

// NextBG cycles through presets smoothly with key [B]
func NextBG(current AnimatedBGType) AnimatedBGType {
	switch current {
	case BGAuto:
		return BGStars
	case BGStars:
		return BGMatrix
	case BGMatrix:
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
		return "Mati (Hitam Polos)"
	case BGAuto:
		return "Dinamis Waktu Server (Auto)"
	case BGStars:
		return "Bintang Kosmik Hening"
	case BGMatrix:
		return "Digital Matrix Stream"
	case BGSawah:
		return "Sawah Lembah Slamet"
	case BGGunung:
		return "Gunung Slamet Megah"
	case BGSenja:
		return "Lentera Senja Bumiayu"
	case BGSakalibel:
		return "Kereta Api Sakalibel"
	case BGMalam:
		return "Langit Malam & Meteor"
	case BGCurug:
		return "Curug Cipendok Air Terjun"
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
		return BGStars, "Malam: Bintang Hening"
	}
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

// Generate2DCanvas produces a full 2D screen grid of CanvasCells for the desktop wallpaper
// Strictly uses single-width runes to prevent character wrapping
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
	case BGStars:
		fillStarsCanvas(canvas, frame, width, height)
	case BGMatrix:
		fillMatrixCanvas(canvas, frame, width, height)
	case BGSawah:
		fillSawahCanvas(canvas, frame, width, height, skyH, horizonH)
	case BGGunung:
		fillGunungCanvas(canvas, frame, width, height, skyH, horizonH)
	case BGSakalibel:
		fillSakalibelCanvas(canvas, frame, width, height, skyH, horizonH)
	case BGSenja:
		fillSenjaCanvas(canvas, frame, width, height, skyH, horizonH)
	case BGMalam:
		fillMalamCanvas(canvas, frame, width, height, skyH, horizonH)
	case BGCurug:
		fillCurugCanvas(canvas, frame, width, height, skyH, horizonH)
	default:
		fillStarsCanvas(canvas, frame, width, height)
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

func makeCell(r rune, red, green, blue uint8) CanvasCell {
	return CanvasCell{
		Char: r,
		FG:   fmt.Sprintf("\x1b[38;2;%d;%d;%dm", red, green, blue),
	}
}

// 1. STARS: Relaxing, sparse, serene cosmic stars
func fillStarsCanvas(c [][]CanvasCell, frame, width, height int) {
	starPatterns := []string{
		"  .                 *                   .                   +                 ·              ",
		"         *                     .                     *                   .                   ",
		"                   +                     ·                     *                     .       ",
		"  ·                     .                     +                     .                 *      ",
	}

	for y := 0; y < height; y++ {
		pIdx := (y + (frame / 48)) % len(starPatterns)
		runes := getScrolledRunes(starPatterns[pIdx], (frame/24)+(y*7), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 148, 163, 184) // Soft silver slate
			}
		}
	}
}

// 2. MATRIX: Soft digital stream
func fillMatrixCanvas(c [][]CanvasCell, frame, width, height int) {
	matrixPatterns := []string{
		"   1       :       0       .       1       *       0       :       .       1       *       ",
		"   0       .       1       *       0       :       1       .       0       *       :       ",
		"   :       1       .       0       *       1       :       0       .       1       0       ",
	}

	for y := 0; y < height; y++ {
		pIdx := (y + (frame / 16)) % len(matrixPatterns)
		runes := getScrolledRunes(matrixPatterns[pIdx], (frame/12)+(y*5), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 22, 101, 52) // Dark forest emerald green
			}
		}
	}
}

// 3. SAWAH: Rice terraces & gentle valley breeze (single-width runes)
func fillSawahCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int) {
	// Sky: Gentle drifting clouds
	for y := 0; y < skyH; y++ {
		pattern := "         (~~~)                       .                 (~~~)               *             "
		runes := getScrolledRunes(pattern, (frame/24)+(y*11), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 100, 116, 139)
			}
		}
	}

	// Horizon: Distant hills & saung
	for y := skyH; y < horizonH; y++ {
		var pattern string
		if y == horizonH-1 {
			pattern = "   /~~~~~~~~\\       [saung]       /~~~~~~~~~~~~~\\              /~~~~~~~~\\                "
		} else {
			pattern = "  /~~~~~~~~~~\\                   /~~~~~~~~~~~~~~~\\            /~~~~~~~~~~\\               "
		}
		runes := getScrolledRunes(pattern, (frame/16)+(y*9), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 71, 85, 105)
			}
		}
	}

	// Ground: Rice terrace ripples
	for y := horizonH; y < height; y++ {
		pattern := "  ~ . ~ , ~ . ~ _ - _ - _ ~ . ~ , ~ . ~ _ - _ - _ ~ . ~ , ~ . ~ _ - _ - _  "
		runes := getScrolledRunes(pattern, (frame/8)+(y*5), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 22, 101, 52)
			}
		}
	}
}

// 4. GUNUNG SLAMET: Mountain silhouette & pine lines
func fillGunungCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int) {
	for y := 0; y < skyH; y++ {
		pattern := "            /^^^\\  ^~^                 .                 /^^^\\  ^~^                 .    "
		runes := getScrolledRunes(pattern, (frame/24)+(y*7), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 99, 102, 241)
			}
		}
	}

	for y := skyH; y < horizonH; y++ {
		pattern := "      /===============\\      ^^^^   ^^   ^^^^       /===============\\      ^^^^   ^^     "
		runes := getScrolledRunes(pattern, (frame/16)+(y*8), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 79, 70, 229)
			}
		}
	}

	for y := horizonH; y < height; y++ {
		pattern := "  ~.~.~ Kali Pedes ~.~.~       *       ~.~.~ Kali Keruh ~.~.~       *       ~.~.~        "
		runes := getScrolledRunes(pattern, (frame/8)+(y*6), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 14, 116, 144)
			}
		}
	}
}

// 5. SAKALIBEL: Railway bridge & train
func fillSakalibelCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int) {
	for y := 0; y < skyH; y++ {
		pattern := "        o   O   o   O   .               o   O   o   O   .               o   O   o   O    "
		runes := getScrolledRunes(pattern, (-frame/16)+(y*8), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 148, 163, 184)
			}
		}
	}

	for y := skyH; y < horizonH; y++ {
		var pattern string
		if y == horizonH-2 {
			pattern = "==[KAMANDAKA]====[==]=[==]=[==]=========================================================="
		} else {
			pattern = "========================================================================================="
		}
		runes := getScrolledRunes(pattern, (frame/6)+(y*4), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 180, 83, 9)
			}
		}
	}

	for y := horizonH; y < height; y++ {
		pattern := "   ||   (===)   ||        ||   (===)   ||   ~.~.~ Kalikeruh ~.~.~   ||   (===)   ||      "
		runes := getScrolledRunes(pattern, (frame/12)+(y*5), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 14, 116, 144)
			}
		}
	}
}

// 6. SENJA: Lanterns & dusk sky
func fillSenjaCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int) {
	for y := 0; y < skyH; y++ {
		pattern := "-------.-------*-------.-------  ~v~  -------.-------*-------.-------  ~v~  -------.-----"
		runes := getScrolledRunes(pattern, (frame/20)+(y*6), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 194, 65, 12)
			}
		}
	}

	for y := skyH; y < horizonH; y++ {
		pattern := "     [+]               [+]               [+]               [+]               [+]         "
		runes := getScrolledRunes(pattern, (frame/14)+(y*5), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 185, 28, 28)
			}
		}
	}

	for y := horizonH; y < height; y++ {
		pattern := "  ~.~ Teh Poci Bumiayu ~.~     - - -     ~.~ Mendoan ~.~     - - -     ~.~ Teh Poci ~.~  "
		runes := getScrolledRunes(pattern, (frame/10)+(y*4), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 161, 98, 7)
			}
		}
	}
}

// 7. MALAM: Night sky & constellations
func fillMalamCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int) {
	for y := 0; y < skyH; y++ {
		pattern := "  *     .       +       .    )    *       .       +       .    *       .       +         "
		runes := getScrolledRunes(pattern, (frame/24)+(y*7), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 148, 163, 184)
			}
		}
	}

	for y := skyH; y < horizonH; y++ {
		pattern := "   .   +   /~~~~~~~~\\    *    /~~~~~~~~~~~~~~~~\\    .   +   /~~~~~~~~\\   "
		runes := getScrolledRunes(pattern, (frame/16)+(y*5), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 71, 85, 105)
			}
		}
	}

	for y := horizonH; y < height; y++ {
		pattern := "  .   *   .   +   .   *   .   +   ~.~.~.~.~   .   *   .   +   .   *   ~.~.~.~.~   .   *  "
		runes := getScrolledRunes(pattern, (frame/12)+(y*4), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 51, 65, 85)
			}
		}
	}
}

// 8. CURUG: Waterfall streams
func fillCurugCanvas(c [][]CanvasCell, frame, width, height, skyH, horizonH int) {
	for y := 0; y < skyH; y++ {
		pattern := "|||||   .   |||||||||||   .   |||||||||||   .   |||||||||||   .   |||||||||||   .   |||||"
		runes := getScrolledRunes(pattern, (frame/24)+(y*6), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 71, 85, 105)
			}
		}
	}

	for y := skyH; y < horizonH; y++ {
		pattern := "  | : . | : . |     | : . | : . |     | : . | : . |     | : . | : . |     | : . | : . |  "
		runes := getScrolledRunes(pattern, (frame/10)+(y*5), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 14, 116, 144)
			}
		}
	}

	for y := horizonH; y < height; y++ {
		pattern := "  ~.~.~ Curug Cipendok ~.~.~   (o)   ~.~.~ Lembah Slamet ~.~.~   (o)   ~.~.~ Curug ~.~.~ "
		runes := getScrolledRunes(pattern, (frame/8)+(y*4), width)
		for x := 0; x < width; x++ {
			r := runes[x]
			if r != ' ' {
				c[y][x] = makeCell(r, 8, 145, 178)
			}
		}
	}
}

// parseANSILine parses one line of ANSI text into visual cells with active ANSI codes
// strictly handles single vs multi-width runes to avoid terminal wrapping bugs
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

// CompositeCanvas composites foreground Window Manager text over the full-screen 2D background canvas.
// RULE 1: Inside the windows, content is 100% clean and solid (cardBg). ZERO particles in text/tables!
// RULE 2: Background ONLY renders on the desktop canvas (margins, headers, gutters between windows).
// RULE 3: Strict width bounds ensure ZERO line wrapping or terminal stutter!
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
				sb.WriteString(vCell.ANSI + cardBgANSI + string(vCell.Char) + "\x1b[0m")
				continue
			}

			if inPanel {
				// INSIDE WINDOW: 100% CLEAN SOLID CARD BACKGROUND!
				// No background particles, no stars, no noise in text or tables!
				if vCell.Char == ' ' {
					sb.WriteString(cardBgANSI + " \x1b[0m")
				} else {
					sb.WriteString(vCell.ANSI + cardBgANSI + string(vCell.Char) + "\x1b[0m")
				}
			} else {
				// OUTSIDE WINDOW (Desktop canvas, margins, gaps between windows):
				// Render the soothing wallpaper animation!
				if vCell.Char == ' ' {
					if bgCell.Char != 0 && bgCell.Char != ' ' {
						sb.WriteString(bgCell.FG + string(bgCell.Char) + "\x1b[0m")
					} else {
						sb.WriteString(" ")
					}
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
