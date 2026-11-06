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
