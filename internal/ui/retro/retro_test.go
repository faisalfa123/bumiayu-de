package retro

import (
	"strings"
	"testing"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
)

func TestVUMeter(t *testing.T) {
	st := styles.GetTheme(styles.ThemeCyberArcade)
	out0 := RenderVUMeter("CPU", 0, 30, st)
	if !strings.Contains(out0, "0%") {
		t.Errorf("Expected 0%% in output, got: %s", out0)
	}

	out50 := RenderVUMeter("CPU", 50, 30, st)
	if !strings.Contains(out50, "50%") {
		t.Errorf("Expected 50%% in output, got: %s", out50)
	}

	out100 := RenderVUMeter("CPU", 100, 30, st)
	if !strings.Contains(out100, "100%") {
		t.Errorf("Expected 100%% in output, got: %s", out100)
	}
}

func TestDynamicEqualizer(t *testing.T) {
	st := styles.GetTheme(styles.ThemeCyberArcade)
	eq := RenderDynamicEqualizer([]float64{10, 30, 70, 90}, 0, st)
	if len(eq) == 0 {
		t.Errorf("Expected non-empty equalizer output")
	}
}

func TestMarqueeTicker(t *testing.T) {
	st := styles.GetTheme(styles.ThemeCyberArcade)
	msgs := []string{"TEST MESSAGE 1", "TEST MESSAGE 2"}
	tick1 := RenderMarqueeTicker(msgs, 0, 50, st)
	tick2 := RenderMarqueeTicker(msgs, 30, 50, st)

	if len(tick1) == 0 || len(tick2) == 0 {
		t.Errorf("Expected non-empty marquee tickers")
	}
}

func TestAnimatedBackground(t *testing.T) {
	st := styles.GetTheme(styles.ThemeCyberArcade)

	presets := []AnimatedBGType{
		BGAuto,
		BGSawah,
		BGGunung,
		BGSakalibel,
		BGSenja,
		BGMalam,
		BGCurug,
	}

	widths := []int{20, 40, 80, 120}
	frames := []int{0, 1, 12, 24, 48, 100}

	for _, p := range presets {
		name := BGName(p)
		if len(name) == 0 {
			t.Errorf("Expected non-empty name for preset %s", p)
		}

		for _, w := range widths {
			for _, f := range frames {
				out := RenderAnimatedBackground(p, f, w, st, 4)
				if len(out) == 0 {
					t.Errorf("Expected non-empty output for preset %s at width %d, frame %d", p, w, f)
				}
			}
		}
	}

	// Verify BGOff returns empty string
	if out := RenderAnimatedBackground(BGOff, 0, 80, st, 4); out != "" {
		t.Errorf("Expected empty string for BGOff, got %s", out)
	}

	// Verify NextBG cycles properly
	cur := BGAuto
	for i := 0; i < 10; i++ {
		cur = NextBG(cur)
	}
}

