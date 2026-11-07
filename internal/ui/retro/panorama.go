package retro

import (
	"fmt"
	"strings"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// PanoramaTheme identifies cultural regional panorama themes
type PanoramaTheme string

const (
	PanoramaAuto       PanoramaTheme = "auto"
	PanoramaNusantara  PanoramaTheme = "nusantara"
	PanoramaJawa       PanoramaTheme = "jawa"
	PanoramaSumatra    PanoramaTheme = "sumatra"
	PanoramaBali       PanoramaTheme = "bali"
	PanoramaKalimantan PanoramaTheme = "kalimantan"
	PanoramaSulawesi   PanoramaTheme = "sulawesi"
	PanoramaPapua      PanoramaTheme = "papua"
	PanoramaBumiayu    PanoramaTheme = "bumiayu"
	PanoramaBintang    PanoramaTheme = "bintang"
	PanoramaOff        PanoramaTheme = "off"
)

// PanoramaData holds the 3-line ASCII motif and styling for a theme
type PanoramaData struct {
	ThemeName   string
	BadgeLabel  string
	SkyPattern  string
	MidPattern  string
	GndPattern  string
	SkyColor    lipgloss.Color
	MidColor    lipgloss.Color
	GndColor    lipgloss.Color
}

// ResolvePanoramaTheme resolves auto rotation every 60 seconds
func ResolvePanoramaTheme(t PanoramaTheme, now time.Time) (PanoramaTheme, bool) {
	if t != PanoramaAuto {
		return t, false
	}
	themes := []PanoramaTheme{
		PanoramaNusantara,
		PanoramaJawa,
		PanoramaSumatra,
		PanoramaBali,
		PanoramaKalimantan,
		PanoramaSulawesi,
		PanoramaPapua,
		PanoramaBumiayu,
		PanoramaBintang,
	}
	idx := (now.Unix() / 60) % int64(len(themes))
	return themes[idx], true
}

// GetPanoramaData returns the 3-tier ASCII scenery for each theme
func GetPanoramaData(t PanoramaTheme) PanoramaData {
	switch t {
	case PanoramaNusantara:
		return PanoramaData{
			ThemeName:  "Nusantara (Samudra & Pinisi)",
			BadgeLabel: "NUSANTARA",
			SkyPattern: "   ~~~   ~   ( )   ~~~   ~   ~   ~~~   ( )   ~~~   ~   ~   ~~~   ( )   ~~~  ",
			MidPattern: " /\\_/\\       |===|      /\\_/\\       |===|      /\\_/\\       |===|      /\\_/\\ ",
			GndPattern: "~~~~~~~\\~~~~~/~~~~~~~~~~~~~~~\\~~~~~/~~~~~~~~~~~~~~~\\~~~~~/~~~~~~~~~~~~~~~",
			SkyColor:   lipgloss.Color("#38BDF8"), // Sky Blue
			MidColor:   lipgloss.Color("#06B6D4"), // Ocean Cyan
			GndColor:   lipgloss.Color("#0284C7"), // Deep Sea
		}

	case PanoramaJawa:
		return PanoramaData{
			ThemeName:  "Jawa (Borobudur & Slamet)",
			BadgeLabel: "JAWA",
			SkyPattern: "      /\\_                ( )         _/\\_               /\\_            ( )  ",
			MidPattern: "  _  /   \\  _   _()_    /   \\       /|  |\\    _()_     /   \\  _   _()_      ",
			GndPattern: "_//\\\\_____//\\\\_/====\\_//=====\\\\_===//====\\\\_/====\\_===//=====\\\\_//====\\_===",
			SkyColor:   lipgloss.Color("#FBBF24"), // Gold Sun
			MidColor:   lipgloss.Color("#F59E0B"), // Temple Terracotta
			GndColor:   lipgloss.Color("#D97706"), // Fertile Soil
		}

	case PanoramaSumatra:
		return PanoramaData{
			ThemeName:  "Sumatra (Rumah Gadang & Toba)",
			BadgeLabel: "SUMATRA",
			SkyPattern: "  /\\  /\\        ~     ~        /\\  /\\         ~       /\\  /\\        ~       ",
			MidPattern: " /  \\/  \\    /~\\   __   __   /  \\/  \\     /~\\      /  \\/  \\    /~\\   __   ",
			GndPattern: "[========]==/    \\_/  \\_/  \\_[========]==/    \\=====[========]==/    \\_/  \\_",
			SkyColor:   lipgloss.Color("#F43F5E"), // Minang Crimson
			MidColor:   lipgloss.Color("#FB7185"), // Rose Gold
			GndColor:   lipgloss.Color("#BE123C"), // Lake Toba Blue-Red
		}

	case PanoramaBali:
		return PanoramaData{
			ThemeName:  "Bali (Candi Bentar & Pura)",
			BadgeLabel: "BALI",
			SkyPattern: "   )  (              ~  ~       )  (           ( )    )  (              ~   ",
			MidPattern: "  /|  |\\     _()_   / \\  / \\   /|  |\\     _()_ / \\   /|  |\\     _()_   / \\  ",
			GndPattern: "_/ |==| \\___/====\\_/   \\/   \\_/ |==| \\___/====\\___\\_/ |==| \\___/====\\_/   \\/",
			SkyColor:   lipgloss.Color("#FB923C"), // Coral Sunset
			MidColor:   lipgloss.Color("#F97316"), // Sandstone Temple
			GndColor:   lipgloss.Color("#EA580C"), // Dewata Waves
		}

	case PanoramaKalimantan:
		return PanoramaData{
			ThemeName:  "Kalimantan (Enggang & Mahakam)",
			BadgeLabel: "BORNEO",
			SkyPattern: "   \\__/         ~  ~   ~         \\__/         ~        \\__/         ~  ~    ",
			MidPattern: "  /____\\     /\\_/\\_/\\_/\\_/\\_    /____\\     /\\_/\\_/\\_  /____\\     /\\_/\\_/\\_  ",
			GndPattern: "===|  |====/~~~~~~~~~~~~~~~~~\\===|  |====/~~~~~~~~~~===|  |====/~~~~~~~~~~",
			SkyColor:   lipgloss.Color("#34D399"), // Canopy Mist
			MidColor:   lipgloss.Color("#10B981"), // Emerald Rainforest
			GndColor:   lipgloss.Color("#059669"), // River Mahakam
		}

	case PanoramaSulawesi:
		return PanoramaData{
			ThemeName:  "Sulawesi (Tongkonan & Pinisi)",
			BadgeLabel: "SULAWESI",
			SkyPattern: " \\______/        ~   ~          \\______/        ~     \\______/        ~     ",
			MidPattern: "  |    |     |====|     /|       |    |     |====|     |    |     |====|    ",
			GndPattern: "==|====|====/~~~~~~\\===/ |=======|====|====/~~~~~~\\===|====|====/~~~~~~\\===",
			SkyColor:   lipgloss.Color("#60A5FA"), // Coastal Sky
			MidColor:   lipgloss.Color("#3B82F6"), // Toraja Timber
			GndColor:   lipgloss.Color("#1D4ED8"), // Coral Reef
		}

	case PanoramaPapua:
		return PanoramaData{
			ThemeName:  "Papua (Honai & Cendrawasih)",
			BadgeLabel: "PAPUA",
			SkyPattern: "   (..)              ~           (..)           ~      (..)              ~  ",
			MidPattern: "  /====\\    ~~\\_/\\_/~   /|      /====\\    ~~\\_/\\_/~   /====\\    ~~\\_/\\_/~   ",
			GndPattern: "_/======\\__/~~~~~~~~~~\\/ |_____/======\\__/~~~~~~~~~~_/======\\__/~~~~~~~~~~\\/",
			SkyColor:   lipgloss.Color("#FACC15"), // Cendrawasih Feather
			MidColor:   lipgloss.Color("#EAB308"), // Honai Straw
			GndColor:   lipgloss.Color("#CA8A04"), // Highlands
		}

	case PanoramaBumiayu:
		return PanoramaData{
			ThemeName:  "Bumiayu (Sakalibel & Slamet)",
			BadgeLabel: "BUMIAYU",
			SkyPattern: "      /\\_                (o)         _/\\_               /\\_            (o)  ",
			MidPattern: " [oo]=|oo]==  ___...___...___... [oo]=|oo]==  ___...  [oo]=|oo]==  ___...   ",
			GndPattern: "============|===|===|===|===|======================|===|===|===|===|======",
			SkyColor:   lipgloss.Color("#2DD4BF"), // Mountain Teal
			MidColor:   lipgloss.Color("#14B8A6"), // Steam Locomotive
			GndColor:   lipgloss.Color("#0F766E"), // Sakalibel Viaduct
		}

	case PanoramaBintang:
		return PanoramaData{
			ThemeName:  "Bintang (Kosmik & Gugusan)",
			BadgeLabel: "BINTANG",
			SkyPattern: "    .     *        +         .          *        +        .     *        +  ",
			MidPattern: "  *    (   )    .     --===*      *          .    (   )     *      --===*   ",
			GndPattern: " .       +         *        .       +      *         .       +        *     ",
			SkyColor:   lipgloss.Color("#E2E8F0"), // Cosmic Starlight
			MidColor:   lipgloss.Color("#A855F7"), // Deep Nebula
			GndColor:   lipgloss.Color("#6366F1"), // Indigo Horizon
		}

	default:
		return GetPanoramaData(PanoramaNusantara)
	}
}

// RenderNusantaraPanorama renders a 3-line multi-tiered panorama running smoothly at 24 FPS
func RenderNusantaraPanorama(theme PanoramaTheme, frame int, width int, st styles.ThemeStyles, now time.Time) string {
	if theme == PanoramaOff || width < 25 {
		return ""
	}

	activeTheme, isAuto := ResolvePanoramaTheme(theme, now)
	data := GetPanoramaData(activeTheme)

	badgeWidth := 13
	sceneryWidth := width - badgeWidth
	if sceneryWidth < 10 {
		sceneryWidth = 10
	}

	// 24 FPS scrolling: 1 character shift every 2 ticks
	offset1 := (frame / 3)
	offset2 := (frame / 2)
	offset3 := (frame / 2)

	line1Scenery := scrollString(data.SkyPattern, offset1, sceneryWidth)
	line2Scenery := scrollString(data.MidPattern, offset2, sceneryWidth)
	line3Scenery := scrollString(data.GndPattern, offset3, sceneryWidth)

	autoTag := "MANUAL"
	if isAuto {
		autoTag = "AUTO"
	}

	// Badges on left (strictly 12 runes wide)
	badge1 := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(data.SkyColor).
		Render(" [NUSANTARA] ")

	badge2 := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(data.MidColor).
		Render(fmt.Sprintf(" %-11s", data.BadgeLabel))

	badge3 := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(data.GndColor).
		Render(fmt.Sprintf(" %-11s", autoTag))

	// Scenery styling
	s1 := lipgloss.NewStyle().Foreground(data.SkyColor).Render(" " + line1Scenery)
	s2 := lipgloss.NewStyle().Foreground(data.MidColor).Render(" " + line2Scenery)
	s3 := lipgloss.NewStyle().Foreground(data.GndColor).Render(" " + line3Scenery)

	row1 := badge1 + s1
	row2 := badge2 + s2
	row3 := badge3 + s3

	return row1 + "\n" + row2 + "\n" + row3
}

func scrollString(pattern string, offset int, targetWidth int) string {
	runes := []rune(pattern)
	pLen := len(runes)
	if pLen == 0 {
		return strings.Repeat(" ", targetWidth)
	}

	var sb strings.Builder
	for i := 0; i < targetWidth-1; i++ {
		idx := (offset + i) % pLen
		if idx < 0 {
			idx += pLen
		}
		sb.WriteRune(runes[idx])
	}
	return sb.String()
}
