package wm

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/retro"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// PaneID identifies the 4 main tiled panels
type PaneID int

const (
	PaneTelemetry PaneID = iota
	PaneServices
	PaneProcesses
	PaneNetwork
)

// WindowManagerState holds layout, focused pane, zoom state, and scroll cursors
type WindowManagerState struct {
	FocusedPane PaneID
	Zoomed      bool

	// Scroll cursors for lists inside panes
	ServiceCursor int
	ProcessCursor int
	NetworkCursor int

	// Safe guard state
	ConfirmAction bool
	ConfirmPrompt string
	PendingAction func()
}

// NextPane moves focus to next pane
func (wm *WindowManagerState) NextPane() {
	wm.FocusedPane = (wm.FocusedPane + 1) % (PaneNetwork + 1)
}

// PrevPane moves focus to previous pane
func (wm *WindowManagerState) PrevPane() {
	if wm.FocusedPane == 0 {
		wm.FocusedPane = PaneNetwork
	} else {
		wm.FocusedPane--
	}
}

// ToggleZoom toggles fullscreen mode of the focused pane
func (wm *WindowManagerState) ToggleZoom() {
	wm.Zoomed = !wm.Zoomed
}

// RenderTiledLayout renders either the 2x2 tiled grid or the zoomed single pane
func RenderTiledLayout(
	wm WindowManagerState,
	info system.HardwareInfo,
	services []system.ServiceItem,
	processes []system.ProcessItem,
	ports []system.PortAuditItem,
	firewallEnabled bool,
	animFrame int,
	st styles.ThemeStyles,
	width, height int,
) string {
	contentWidth := width - 2
	if contentWidth < 40 {
		contentWidth = 40
	}
	contentHeight := height - 5 // Leave room for topbar, marquee, and footer
	if contentHeight < 15 {
		contentHeight = 15
	}

	// If Zoomed, render ONLY the focused pane with full size!
	if wm.Zoomed {
		switch wm.FocusedPane {
		case PaneTelemetry:
			return renderTelemetryPane(info, animFrame, st, contentWidth, contentHeight, true)
		case PaneServices:
			return renderServicesPane(services, wm.ServiceCursor, st, contentWidth, contentHeight, true)
		case PaneProcesses:
			return renderProcessesPane(processes, wm.ProcessCursor, st, contentWidth, contentHeight, true)
		case PaneNetwork:
			return renderNetworkPane(ports, firewallEnabled, wm.NetworkCursor, st, contentWidth, contentHeight, true)
		}
	}

	// 2x2 Grid Layout
	colW := (contentWidth - 2) / 2
	rowH := (contentHeight - 1) / 2
	if colW < 30 {
		colW = 30
	}
	if rowH < 8 {
		rowH = 8
	}

	p1 := renderTelemetryPane(info, animFrame, st, colW, rowH, wm.FocusedPane == PaneTelemetry)
	p2 := renderServicesPane(services, wm.ServiceCursor, st, colW, rowH, wm.FocusedPane == PaneServices)
	p3 := renderProcessesPane(processes, wm.ProcessCursor, st, colW, rowH, wm.FocusedPane == PaneProcesses)
	p4 := renderNetworkPane(ports, firewallEnabled, wm.NetworkCursor, st, colW, rowH, wm.FocusedPane == PaneNetwork)

	rowTop := lipgloss.JoinHorizontal(lipgloss.Top, p1, " ", p2)
	rowBottom := lipgloss.JoinHorizontal(lipgloss.Top, p3, " ", p4)

	return lipgloss.JoinVertical(lipgloss.Left, rowTop, rowBottom)
}

func getPaneBox(title string, badge string, content string, isFocused bool, st styles.ThemeStyles, width, height int) string {
	borderStyle := st.Card
	if isFocused {
		borderStyle = st.CardActive
	}

	badgeRender := ""
	if badge != "" {
		badgeRender = " " + lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(st.AccentColor).
			Padding(0, 1).
			Render(badge)
	}

	headerColor := st.TextColor
	if isFocused {
		headerColor = st.PrimaryColor
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(headerColor).
		Render(title) + badgeRender

	inner := lipgloss.JoinVertical(lipgloss.Left,
		header,
		lipgloss.NewStyle().Foreground(st.BorderColor).Render(strings.Repeat("─", width-6)),
		content,
	)

	return borderStyle.
		Width(width).
		Height(height).
		Render(inner)
}

// 1. TELEMETRY PANE
func renderTelemetryPane(info system.HardwareInfo, frame int, st styles.ThemeStyles, width, height int, isFocused bool) string {
	meterW := width - 12

	cpuBar := retro.RenderVUMeter("CPU", info.CPUUsagePct, meterW, st)
	ramBar := retro.RenderVUMeter("RAM", info.MemUsagePct, meterW, st)
	tempBar := retro.RenderVUMeter("TEMP", (info.TemperatureC/100.0)*100.0, meterW, st)

	heart := retro.RenderHeartbeat(frame, info.MemUsagePct < 85 && info.TemperatureC < 75, st)
	eq := retro.RenderDynamicEqualizer([]float64{info.CPUUsagePct, info.MemUsagePct, 50, 65, info.Load1 * 20, 40}, frame, st)

	devName := "Linux Server"
	if info.IsOdroidC2 {
		devName = "ODROID-C2 (Amlogic S905)"
	}

	meta1 := fmt.Sprintf("💻 %s | Arch: %s", devName, info.Arch)
	meta2 := fmt.Sprintf("🌐 Host: %-12s | Uptime: %s", info.Hostname, info.UptimeStr)
	meta3 := fmt.Sprintf("📊 Load: %.2f %.2f | RAM: %d/%d MB", info.Load1, info.Load5, info.MemUsedMB, info.MemTotalMB)

	content := lipgloss.JoinVertical(lipgloss.Left,
		cpuBar,
		ramBar,
		tempBar,
		"",
		lipgloss.JoinHorizontal(lipgloss.Left, "EQ 24FPS: ", eq, "  ", heart),
		"",
		meta1,
		meta2,
		meta3,
	)

	return getPaneBox("📊 [1] TELEMETRI HARDWARE", "LIVE", content, isFocused, st, width, height)
}

// 2. SERVICES PANE
func renderServicesPane(services []system.ServiceItem, cursor int, st styles.ThemeStyles, width, height int, isFocused bool) string {
	if len(services) == 0 {
		return getPaneBox("⚙️ [2] LAYANAN SYSTEMD", "", "Memuat layanan...", isFocused, st, width, height)
	}

	maxRows := height - 6
	if maxRows < 3 {
		maxRows = 3
	}

	startIdx := 0
	if cursor >= maxRows {
		startIdx = cursor - maxRows + 1
	}
	endIdx := startIdx + maxRows
	if endIdx > len(services) {
		endIdx = len(services)
	}

	var rows []string
	for i := startIdx; i < endIdx; i++ {
		svc := services[i]
		isSelected := isFocused && (i == cursor)

		var stateBadge string
		if svc.ActiveState == "active" {
			stateBadge = lipgloss.NewStyle().Foreground(st.SuccessColor).Render("● RUN")
		} else if svc.ActiveState == "failed" {
			stateBadge = lipgloss.NewStyle().Foreground(st.DangerColor).Render("✖ FAIL")
		} else {
			stateBadge = lipgloss.NewStyle().Foreground(st.MutedColor).Render("○ STOP")
		}

		name := svc.Name
		if len(name) > width-24 {
			name = name[:width-27] + "..."
		}

		line := fmt.Sprintf("%-6s %-18s %4.1fMB", stateBadge, name, svc.MemoryMB)

		if isSelected {
			rows = append(rows, lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(st.PrimaryColor).
				Render("► "+line))
		} else {
			rows = append(rows, "  "+line)
		}
	}

	hint := lipgloss.NewStyle().Foreground(st.MutedColor).Render("[S]: Start | [T]: Stop | [R]: Restart")
	body := lipgloss.JoinVertical(lipgloss.Left,
		strings.Join(rows, "\n"),
		"",
		hint,
	)

	return getPaneBox("⚙️ [2] LAYANAN SYSTEMD", fmt.Sprintf("%d Unit", len(services)), body, isFocused, st, width, height)
}

// 3. PROCESSES PANE
func renderProcessesPane(processes []system.ProcessItem, cursor int, st styles.ThemeStyles, width, height int, isFocused bool) string {
	if len(processes) == 0 {
		return getPaneBox("⚡ [3] TASK MANAGER (OOM)", "", "Memuat daftar proses...", isFocused, st, width, height)
	}

	maxRows := height - 6
	if maxRows < 3 {
		maxRows = 3
	}

	startIdx := 0
	if cursor >= maxRows {
		startIdx = cursor - maxRows + 1
	}
	endIdx := startIdx + maxRows
	if endIdx > len(processes) {
		endIdx = len(processes)
	}

	var rows []string
	for i := startIdx; i < endIdx; i++ {
		p := processes[i]
		isSelected := isFocused && (i == cursor)

		var oomBadge string
		if p.OOMScore > 300 {
			oomBadge = lipgloss.NewStyle().Foreground(st.DangerColor).Bold(true).Render(fmt.Sprintf("OOM:%-3d!", p.OOMScore))
		} else {
			oomBadge = lipgloss.NewStyle().Foreground(st.MutedColor).Render(fmt.Sprintf("OOM:%-3d", p.OOMScore))
		}

		name := p.Name
		if len(name) > width-28 {
			name = name[:width-31] + "..."
		}

		line := fmt.Sprintf("%-5d %-12s %5.1fMB %s", p.PID, name, p.MemMB, oomBadge)

		if isSelected {
			rows = append(rows, lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(st.PrimaryColor).
				Render("► "+line))
		} else {
			rows = append(rows, "  "+line)
		}
	}

	hint := lipgloss.NewStyle().Foreground(st.MutedColor).Render("[K]: Kill Process SIGTERM")
	body := lipgloss.JoinVertical(lipgloss.Left,
		strings.Join(rows, "\n"),
		"",
		hint,
	)

	return getPaneBox("⚡ [3] TASK MANAGER (OOM)", fmt.Sprintf("%d Pros", len(processes)), body, isFocused, st, width, height)
}

// 4. NETWORK & PORTS PANE
func renderNetworkPane(ports []system.PortAuditItem, ufwActive bool, cursor int, st styles.ThemeStyles, width, height int, isFocused bool) string {
	var ufwPill string
	if ufwActive {
		ufwPill = lipgloss.NewStyle().Foreground(st.SuccessColor).Bold(true).Render("🛡️ UFW: AKTIF (Aman)")
	} else {
		ufwPill = lipgloss.NewStyle().Foreground(st.DangerColor).Bold(true).Render("⚠️ UFW: MATI (Terbuka)")
	}

	var rows []string
	maxRows := height - 7
	if maxRows < 2 {
		maxRows = 2
	}

	startIdx := 0
	if cursor >= maxRows {
		startIdx = cursor - maxRows + 1
	}
	endIdx := startIdx + maxRows
	if endIdx > len(ports) {
		endIdx = len(ports)
	}

	for i := startIdx; i < endIdx; i++ {
		p := ports[i]
		isSelected := isFocused && (i == cursor)

		var pubBadge string
		if p.IsPublic {
			pubBadge = lipgloss.NewStyle().Foreground(st.DangerColor).Render("PUBLIK")
		} else {
			pubBadge = lipgloss.NewStyle().Foreground(st.SuccessColor).Render("LOKAL ")
		}

		line := fmt.Sprintf("%-5d %-4s %-6s %s", p.Port, p.Protocol, pubBadge, p.ProcessName)

		if isSelected {
			rows = append(rows, lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(st.PrimaryColor).
				Render("► "+line))
		} else {
			rows = append(rows, "  "+line)
		}
	}

	hint := lipgloss.NewStyle().Foreground(st.MutedColor).Render("[U]: Toggle Firewall UFW")
	body := lipgloss.JoinVertical(lipgloss.Left,
		ufwPill,
		"",
		strings.Join(rows, "\n"),
		"",
		hint,
	)

	return getPaneBox("🛡️ [4] JARINGAN & PORT", fmt.Sprintf("%d Port", len(ports)), body, isFocused, st, width, height)
}
