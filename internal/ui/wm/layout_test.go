package wm

import (
	"testing"

	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
)

func TestRenderTiledLayoutSmallDimensions(t *testing.T) {
	st := styles.GetTheme(styles.ThemeCyberArcade)
	hw := system.GetHardwareInfo()
	svcs := system.GetServices()
	procs := system.GetProcessList()
	sec := system.GetSecurityState()

	wmState := WindowManagerState{
		FocusedPane: PaneTelemetry,
		Zoomed:      false,
	}

	// Test small widths where width-31 could previously become negative
	widths := []int{20, 30, 40, 50, 80, 120}
	heights := []int{10, 15, 20, 30}

	for _, w := range widths {
		for _, h := range heights {
			out := RenderTiledLayout(
				wmState,
				hw,
				svcs,
				procs,
				sec.OpenPorts,
				sec.UfwActive,
				0,
				st,
				w,
				h,
			)
			if len(out) == 0 {
				t.Errorf("Expected non-empty layout output for %dx%d", w, h)
			}
		}
	}
}
