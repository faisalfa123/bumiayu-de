package retro

import (
	"strings"
	"testing"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
)

func TestNusantaraPanoramaAllThemes(t *testing.T) {
	themes := []PanoramaTheme{
		PanoramaAuto,
		PanoramaNusantara,
		PanoramaJawa,
		PanoramaSumatra,
		PanoramaBali,
		PanoramaKalimantan,
		PanoramaSulawesi,
		PanoramaPapua,
		PanoramaBumiayu,
		PanoramaBintang,
		PanoramaOff,
	}

	st := styles.GetTheme(styles.ThemeCyberArcade)
	now := time.Now()

	for _, theme := range themes {
		res := RenderNusantaraPanorama(theme, 10, 80, st, now)
		if theme == PanoramaOff {
			if res != "" {
				t.Errorf("Expected empty string for PanoramaOff, got %s", res)
			}
			continue
		}

		lines := strings.Split(res, "\n")
		if len(lines) != 3 {
			t.Errorf("Theme %s: expected 3 lines, got %d", theme, len(lines))
		}
	}
}

func TestResolvePanoramaThemeRotation(t *testing.T) {
	base := time.Unix(0, 0)
	t1, auto1 := ResolvePanoramaTheme(PanoramaAuto, base)
	if !auto1 {
		t.Errorf("Expected auto flag to be true")
	}

	// Advance 60 seconds
	t2, _ := ResolvePanoramaTheme(PanoramaAuto, base.Add(60*time.Second))
	if t1 == t2 {
		t.Errorf("Expected theme rotation after 60s, got %s both times", t1)
	}

	// Manual theme should not rotate
	manual, auto2 := ResolvePanoramaTheme(PanoramaJawa, base.Add(120*time.Second))
	if auto2 || manual != PanoramaJawa {
		t.Errorf("Expected manual theme Jawa, got %s (auto=%v)", manual, auto2)
	}
}
