package app

import (
	"testing"

	"github.com/bumiayu/bumiayu-de/internal/ui/launcher"
	"github.com/bumiayu/bumiayu-de/internal/ui/wm"
)

func TestAppViewModelRenderingWithAnimatedBG(t *testing.T) {
	m := InitialModel()

	presets := []string{
		"auto",
		"sawah",
		"gunung",
		"sakalibel",
		"senja",
		"malam",
		"curug",
		"off",
	}

	widths := []int{40, 80, 120}
	heights := []int{15, 24, 30, 45}

	for _, p := range presets {
		m.AnimatedBG = p
		for _, w := range widths {
			for _, h := range heights {
				m.width = w
				m.height = h

				// Test normal 2x2 tiling
				m.WM.Mode = wm.LayoutTiling2x2
				m.WM.Zoomed = false
				out := m.View()
				if len(out) == 0 {
					t.Fatalf("Expected non-empty View output for preset %s at %dx%d", p, w, h)
				}

				// Test Master-Stack layout
				m.WM.Mode = wm.LayoutMasterStack
				out = m.View()
				if len(out) == 0 {
					t.Fatalf("Expected non-empty MasterStack output for preset %s at %dx%d", p, w, h)
				}

				// Test Horizontal Split layout
				m.WM.Mode = wm.LayoutHorizontalSplit
				out = m.View()
				if len(out) == 0 {
					t.Fatalf("Expected non-empty HorizontalSplit output for preset %s at %dx%d", p, w, h)
				}

				// Test Zoomed mode
				m.WM.Zoomed = true
				out = m.View()
				if len(out) == 0 {
					t.Fatalf("Expected non-empty Zoomed output for preset %s at %dx%d", p, w, h)
				}
				m.WM.Zoomed = false
			}
		}
	}
}

func TestExecuteLauncherBackgroundCommands(t *testing.T) {
	m := InitialModel()

	tests := []struct {
		cmd      string
		expected string
	}{
		{"bg_auto", "auto"},
		{"bg_sawah", "sawah"},
		{"bg_gunung", "gunung"},
		{"bg_sakalibel", "sakalibel"},
		{"bg_senja", "senja"},
		{"bg_malam", "malam"},
		{"bg_curug", "curug"},
		{"bg_off", "off"},
	}

	for _, tc := range tests {
		item := launcher.LauncherItem{Command: tc.cmd}
		_, _ = m.executeLauncherItem(item)
		if m.AnimatedBG != tc.expected {
			t.Errorf("Expected AnimatedBG to be %s after command %s, got %s", tc.expected, tc.cmd, m.AnimatedBG)
		}
	}
}
