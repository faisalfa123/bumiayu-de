package retro

import (
	"fmt"
	"math"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

var eqBlocks = []rune{' ', ' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// RenderVUMeter renders a high-contrast 16-bit arcade audio VU-meter style bar
func RenderVUMeter(label string, valuePct float64, width int, st styles.ThemeStyles) string {
	if width < 15 {
		width = 15
	}
	if valuePct < 0 {
		valuePct = 0
	}
	if valuePct > 100 {
		valuePct = 100
	}

	barWidth := width - len(label) - 8 // account for label and " 100% "
	if barWidth < 5 {
		barWidth = 5
	}

	filledLen := int(math.Round((valuePct / 100.0) * float64(barWidth)))
	if filledLen > barWidth {
		filledLen = barWidth
	}

	var sb strings.Builder
	for i := 0; i < barWidth; i++ {
		pct := (float64(i) / float64(barWidth)) * 100.0
		if i < filledLen {
			var col lipgloss.Color
			if pct > 80 {
				col = st.MeterHigh
			} else if pct > 60 {
				col = st.MeterMid
			} else {
				col = st.MeterLow
			}
			sb.WriteString(lipgloss.NewStyle().Foreground(col).Render("█"))
		} else {
			sb.WriteString(lipgloss.NewStyle().Foreground(st.BorderColor).Render("░"))
		}
	}

	// Dynamic text color for value percentage
	valColor := st.SuccessColor
	if valuePct > 80 {
		valColor = st.DangerColor
	} else if valuePct > 60 {
		valColor = st.AccentColor
	}

	valStr := lipgloss.NewStyle().Foreground(valColor).Bold(true).Render(fmt.Sprintf("%3.0f%%", valuePct))
	lblStr := lipgloss.NewStyle().Foreground(st.TextColor).Bold(true).Render(label)

	return fmt.Sprintf("%-6s [%s] %s", lblStr, sb.String(), valStr)
}

// RenderDynamicEqualizer generates a live 8-bar retro stereo equalizer at 24 FPS
func RenderDynamicEqualizer(baseVals []float64, frame int, st styles.ThemeStyles) string {
	if len(baseVals) == 0 {
		baseVals = []float64{30, 45, 60, 40, 75, 50, 35, 65}
	}

	var bars []string
	for i, val := range baseVals {
		// Calculate subtle 24 FPS oscillation
		t := float64(frame)*0.15 + float64(i)*0.8
		osc := math.Sin(t) * 12.0
		curr := val + osc
		if curr < 5 {
			curr = 5
		}
		if curr > 100 {
			curr = 100
		}

		// Map 0-100 to block index 0-8
		idx := int((curr / 100.0) * 8.0)
		if idx < 0 {
			idx = 0
		}
		if idx > 8 {
			idx = 8
		}

		char := string(eqBlocks[idx])
		var col lipgloss.Color
		if curr > 80 {
			col = st.MeterHigh
		} else if curr > 50 {
			col = st.MeterMid
		} else {
			col = st.MeterLow
		}

		bars = append(bars, lipgloss.NewStyle().Foreground(col).Bold(true).Render(char))
	}

	return strings.Join(bars, " ")
}

// RenderHeartbeat returns a 24 FPS pulsating retro heart or power indicator
func RenderHeartbeat(frame int, healthGood bool, st styles.ThemeStyles) string {
	f := (frame / 3) % 4
	if !healthGood {
		symbols := []string{"💔 [!!]", "❤️ [ALERT]", "💔 [CRIT]", "❤️ [HIGH]"}
		return lipgloss.NewStyle().Foreground(st.DangerColor).Bold(true).Render(symbols[f])
	}

	symbols := []string{"♥ [OK]", "❤️ [PULSE]", "💚 [ACTIVE]", "♥ [STABLE]"}
	col := st.SuccessColor
	if f == 1 {
		col = st.PrimaryColor
	}
	return lipgloss.NewStyle().Foreground(col).Bold(true).Render(symbols[f])
}
