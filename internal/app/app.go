package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/config"
	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/components"
	"github.com/bumiayu/bumiayu-de/internal/ui/help"
	"github.com/bumiayu/bumiayu-de/internal/ui/launcher"
	"github.com/bumiayu/bumiayu-de/internal/ui/retro"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/bumiayu/bumiayu-de/internal/ui/wm"
	"github.com/bumiayu/bumiayu-de/internal/web"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TickMsg triggers UI frame rendering
type TickMsg time.Time

// TelemetryMsg delivers background refreshed hardware data
type TelemetryMsg struct {
	Hardware  system.HardwareInfo
	Services  []system.ServiceItem
	Processes []system.ProcessItem
	Security  system.SecurityState
}

// AutoCleanMsg signals completion of periodic memory clean
type AutoCleanMsg string

// AppVersion is the current version displayed with BETA flag
const AppVersion = "v1.0.4"

// Model is the main state of Bumiayu TUI Desktop Environment
type Model struct {
	width  int
	height int

	ThemeMode     styles.ThemeMode
	Language      i18n.Language
	StatusMessage string
	AnimFrame     int
	LastActive    time.Time // For adaptive 24 FPS engine

	// Animated Background Runner Preset ("auto", "sawah", "gunung", "sakalibel", "senja", "malam", "curug", "off")
	AnimatedBG string
	// Background Dim / Opacity Level ("low", "med", "high")
	BGDimLevel string

	// FPS Diagnostics
	LastTickTime     time.Time
	CurrentFPS       float64
	FPSCategory      string
	TargetFPSConfig  int // 0: Adaptive 24FPS/Eco, 60: Ultra Smooth, 5: Low Eco

	// Auto Clean RAM Configuration
	AutoCleanMinutes int // 0: Off, 15, 30, 60, 360
	LastCleanTime    time.Time

	// Window Manager state (Tiling panes, focus, zoom, mode)
	WM wm.WindowManagerState

	// Unified Spotlight Launcher
	Spotlight launcher.SpotlightState

	// Arcade Help Modal state ('?' key)
	ShowHelp bool

	// Confirmation Safety Guard Dialog
	ShowConfirm   bool
	ConfirmPrompt string
	ConfirmAction func()

	// Live Telemetry caches
	Hardware        system.HardwareInfo
	Services        []system.ServiceItem
	Processes       []system.ProcessItem
	Ports           []system.PortAuditItem
	FirewallEnabled bool

	// Apps & Web server port
	Catalog       []system.AppCatalogItem
	CatalogCursor int
	WebServerPort int
}

// InitialModel initializes the state with default retro config
func InitialModel() Model {
	m := Model{
		ThemeMode:        styles.ThemeCyberArcade, // 16-bit Cyber-Arcade Default
		Language:         i18n.ID,
		LastActive:       time.Now(),
		LastTickTime:     time.Now(),
		CurrentFPS:       24.0,
		FPSCategory:      "24 FPS [AKTIF]",
		TargetFPSConfig:  24,
		AutoCleanMinutes: 0,
		AnimatedBG:       "auto",
		BGDimLevel:       "med",
		WebServerPort:    8080,
		Spotlight:        launcher.NewSpotlightState(),
		WM: wm.WindowManagerState{
			Mode:        wm.LayoutTiling2x2,
			FocusedPane: wm.PaneTelemetry,
			Zoomed:      false,
		},
	}

	cfg, err := config.LoadConfig()
	if err == nil {
		m.applyUserConfig(cfg)
	}

	if m.WebServerPort > 0 {
		web.StartWebServer(m.WebServerPort)
	}

	// Fetch initial data synchronously once
	m.fetchSystemDataSync()
	return m
}

func (m *Model) applyUserConfig(cfg config.UserConfig) {
	switch cfg.Theme {
	case "dos":
		m.ThemeMode = styles.ThemeDOSNorton
	case "gameboy":
		m.ThemeMode = styles.ThemeGameBoy
	case "dark":
		m.ThemeMode = styles.ThemeDark
	case "warm":
		m.ThemeMode = styles.ThemeWarm
	case "mac":
		m.ThemeMode = styles.ThemeMacLight
	default:
		m.ThemeMode = styles.ThemeCyberArcade
	}

	switch cfg.WMModel {
	case "master":
		m.WM.Mode = wm.LayoutMasterStack
	case "split":
		m.WM.Mode = wm.LayoutHorizontalSplit
	default:
		m.WM.Mode = wm.LayoutTiling2x2
	}

	if cfg.TargetFPS > 0 {
		m.TargetFPSConfig = cfg.TargetFPS
	}
	m.AutoCleanMinutes = cfg.AutoCleanMinutes

	if cfg.AnimatedBG != "" {
		m.AnimatedBG = cfg.AnimatedBG
	} else {
		m.AnimatedBG = "auto"
	}

	if cfg.BGDimLevel != "" {
		m.BGDimLevel = cfg.BGDimLevel
	} else {
		m.BGDimLevel = "med"
	}

	if cfg.Language == "en" {
		m.Language = i18n.EN
	} else {
		m.Language = i18n.ID
	}
}

func (m *Model) saveCurrentConfig() {
	themeStr := "cyber"
	switch m.ThemeMode {
	case styles.ThemeDOSNorton:
		themeStr = "dos"
	case styles.ThemeGameBoy:
		themeStr = "gameboy"
	case styles.ThemeDark:
		themeStr = "dark"
	case styles.ThemeWarm:
		themeStr = "warm"
	case styles.ThemeMacLight:
		themeStr = "mac"
	}

	wmStr := "tiling"
	switch m.WM.Mode {
	case wm.LayoutMasterStack:
		wmStr = "master"
	case wm.LayoutHorizontalSplit:
		wmStr = "split"
	}

	langStr := "id"
	if m.Language == i18n.EN {
		langStr = "en"
	}

	cfg := config.UserConfig{
		ConfigVersion:     AppVersion,
		Theme:             themeStr,
		WMModel:           wmStr,
		Language:          langStr,
		FirstRunCompleted: true,
		WebServerEnabled:  true,
		WebServerPort:     m.WebServerPort,
		TargetFPS:         m.TargetFPSConfig,
		AutoCleanMinutes:  m.AutoCleanMinutes,
		AnimatedBG:        m.AnimatedBG,
		BGDimLevel:        m.BGDimLevel,
	}
	_ = config.SaveConfig(cfg)
}

func (m *Model) fetchSystemDataSync() {
	m.Hardware = system.GetHardwareInfo()
	m.Services = system.GetServices()
	m.Processes = system.GetProcessList()
	sec := system.GetSecurityState()
	m.Ports = sec.OpenPorts
	m.FirewallEnabled = sec.UfwActive
	m.Catalog = system.GetAppCatalog()

	// Sync services & processes into Spotlight searchable catalog
	var svcNames, procNames []string
	for _, s := range m.Services {
		svcNames = append(svcNames, s.Name)
	}
	for i, p := range m.Processes {
		if i < 15 {
			procNames = append(procNames, fmt.Sprintf("%s (%d)", p.Name, p.PID))
		}
	}
	m.Spotlight.SyncDynamicItems(svcNames, procNames)
}

// Background goroutine for async telemetry without freezing 24 FPS UI
func pollTelemetryCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		hw := system.GetHardwareInfo()
		svcs := system.GetServices()
		procs := system.GetProcessList()
		sec := system.GetSecurityState()
		return TelemetryMsg{
			Hardware:  hw,
			Services:  svcs,
			Processes: procs,
			Security:  sec,
		}
	}
}

// Background timer for Auto Clean RAM
func scheduleAutoClean(minutes int) tea.Cmd {
	if minutes <= 0 {
		return nil
	}
	return func() tea.Msg {
		time.Sleep(time.Duration(minutes) * time.Minute)
		out, err := system.CleanRAMCache()
		if err != nil {
			return AutoCleanMsg("Gagal membersihkan RAM: " + err.Error())
		}
		return AutoCleanMsg(out)
	}
}

// Dynamic adaptive FPS engine ticker
func (m Model) tickAdaptive() tea.Cmd {
	var interval time.Duration

	switch m.TargetFPSConfig {
	case 60:
		// Ultra Smooth 60 FPS (~16.6ms)
		interval = time.Millisecond * 16
	case 5:
		// Static Low Eco Mode (5 FPS = 200ms)
		interval = time.Millisecond * 200
	default: // 24 FPS or Adaptive
		if !m.Spotlight.Active && !m.ShowHelp && !m.ShowConfirm && time.Since(m.LastActive) > 2500*time.Millisecond {
			interval = time.Millisecond * 500 // 2 FPS idle mode (MODE HEMAT)
		} else {
			interval = time.Millisecond * 41 // 24 FPS active
		}
	}

	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Init starts the Bubble Tea loop with smooth 24 FPS tick, background telemetry, and auto-clean
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.tickAdaptive(), pollTelemetryCmd())
	if m.AutoCleanMinutes > 0 {
		cmds = append(cmds, scheduleAutoClean(m.AutoCleanMinutes))
	}
	return tea.Batch(cmds...)
}

// Update handles keyboard, mouse, and animation events
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case TickMsg:
		m.AnimFrame++

		// Rolling FPS Calculation
		now := time.Now()
		delta := now.Sub(m.LastTickTime).Seconds()
		m.LastTickTime = now

		if delta > 0 {
			instantFPS := 1.0 / delta
			// Smooth weighted rolling average
			m.CurrentFPS = (m.CurrentFPS * 0.8) + (instantFPS * 0.2)
		}

		// Categorize FPS behavior
		isIdle := time.Since(m.LastActive) > 2500*time.Millisecond && !m.Spotlight.Active && !m.ShowHelp && !m.ShowConfirm
		if isIdle && m.TargetFPSConfig != 60 {
			m.FPSCategory = "MODE HEMAT"
		} else if delta > 0.8 {
			m.FPSCategory = "HANG / LAG"
		} else if m.TargetFPSConfig == 60 {
			m.FPSCategory = "ULTRA SMOOTH"
		} else {
			m.FPSCategory = "NORMAL 24FPS"
		}

		return m, m.tickAdaptive()

	case TelemetryMsg:
		m.Hardware = msg.Hardware
		m.Services = msg.Services
		m.Processes = msg.Processes
		m.Ports = msg.Security.OpenPorts
		m.FirewallEnabled = msg.Security.UfwActive
		return m, pollTelemetryCmd()

	case AutoCleanMsg:
		m.StatusMessage = "🧹 [Auto-Clean]: " + string(msg)
		m.LastCleanTime = time.Now()
		m.fetchSystemDataSync()
		if m.AutoCleanMinutes > 0 {
			return m, scheduleAutoClean(m.AutoCleanMinutes)
		}
		return m, nil

	case tea.KeyMsg:
		m.LastActive = time.Now() // Wake up adaptive engine to full 24 FPS
		key := msg.String()

		// 1. SAFETY CONFIRMATION DIALOG MODAL
		if m.ShowConfirm {
			switch key {
			case "y", "Y", "enter":
				m.ShowConfirm = false
				if m.ConfirmAction != nil {
					m.ConfirmAction()
				}
				return m, nil
			case "n", "N", "esc":
				m.ShowConfirm = false
				m.StatusMessage = "Tindakan dibatalkan."
				return m, nil
			default:
				return m, nil
			}
		}

		// 2. ARCADE HELP 'HOW TO PLAY' MODAL
		if m.ShowHelp {
			switch key {
			case "?", "esc", "enter", "q":
				m.ShowHelp = false
				return m, nil
			default:
				return m, nil
			}
		}

		// 3. UNIFIED SPOTLIGHT LAUNCHER HANDLING
		if m.Spotlight.Active {
			switch key {
			case "esc":
				m.Spotlight.Active = false
				return m, nil

			case "tab", "right":
				// Accept Ghost-Text Autocomplete!
				m.Spotlight.AcceptGhostText()
				return m, nil

			case "up":
				if m.Spotlight.Cursor > 0 {
					m.Spotlight.Cursor--
				}
				return m, nil

			case "down":
				if m.Spotlight.Cursor < len(m.Spotlight.Results)-1 {
					m.Spotlight.Cursor++
				}
				return m, nil

			case "enter":
				m.Spotlight.Active = false
				if len(m.Spotlight.Results) > 0 && m.Spotlight.Cursor < len(m.Spotlight.Results) {
					item := m.Spotlight.Results[m.Spotlight.Cursor]
					return m.executeLauncherItem(item)
				}
				return m, nil

			case "backspace":
				if len(m.Spotlight.Query) > 0 {
					m.Spotlight.Query = m.Spotlight.Query[:len(m.Spotlight.Query)-1]
					m.Spotlight.Filter()
				}
				return m, nil

			default:
				// Append printable single characters
				if len(key) == 1 {
					m.Spotlight.Query += key
					m.Spotlight.Filter()
				}
				return m, nil
			}
		}

		// 4. MAIN WINDOW MANAGER CONTROLS & NAVIGATION
		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit

		case " ":
			// Open Spotlight Launcher
			m.Spotlight.Active = true
			m.Spotlight.Query = ""
			m.Spotlight.Filter()
			return m, nil

		case "/":
			// Open Spotlight Launcher with query prefix
			m.Spotlight.Active = true
			m.Spotlight.Query = ""
			m.Spotlight.Filter()
			return m, nil

		case "?":
			// Open Arcade How to Play / Help Modal
			m.ShowHelp = true
			return m, nil

		case "tab":
			m.WM.NextPane()
			return m, nil

		case "shift+tab":
			m.WM.PrevPane()
			return m, nil

		case "f", "F":
			// Toggle Fullscreen / Zoom of focused pane
			m.WM.ToggleZoom()
			return m, nil

		case "1":
			m.WM.FocusedPane = wm.PaneTelemetry
			return m, nil

		case "2":
			m.WM.FocusedPane = wm.PaneServices
			return m, nil

		case "3":
			m.WM.FocusedPane = wm.PaneProcesses
			return m, nil

		case "4":
			m.WM.FocusedPane = wm.PaneNetwork
			return m, nil

		case "t", "T":
			// Quick Cycle Theme
			m.ThemeMode = styles.NextTheme(m.ThemeMode)
			m.StatusMessage = fmt.Sprintf("Tema aktif: %s", styles.ThemeName(m.ThemeMode))
			m.saveCurrentConfig()
			return m, nil

		case "b", "B":
			// Quick Cycle Animated Background
			next := retro.NextBG(retro.AnimatedBGType(m.AnimatedBG))
			m.AnimatedBG = string(next)
			m.StatusMessage = fmt.Sprintf("Background: %s", retro.BGName(next))
			m.saveCurrentConfig()
			return m, nil

		case "l", "L":
			// Quick Language Toggle
			if m.Language == i18n.ID {
				m.Language = i18n.EN
				m.StatusMessage = "Language set to English."
			} else {
				m.Language = i18n.ID
				m.StatusMessage = "Bahasa diubah ke Bahasa Indonesia."
			}
			m.saveCurrentConfig()
			return m, nil

		case "esc":
			if m.WM.Zoomed {
				m.WM.Zoomed = false
			}
			m.StatusMessage = ""
			return m, nil

		// PANE-SPECIFIC SCROLLING & SHORTCUTS
		case "up", "k":
			switch m.WM.FocusedPane {
			case wm.PaneServices:
				if m.WM.ServiceCursor > 0 {
					m.WM.ServiceCursor--
				}
			case wm.PaneProcesses:
				if m.WM.ProcessCursor > 0 {
					m.WM.ProcessCursor--
				}
			case wm.PaneNetwork:
				if m.WM.NetworkCursor > 0 {
					m.WM.NetworkCursor--
				}
			default:
				m.WM.PrevPane()
			}
			return m, nil

		case "down", "j":
			switch m.WM.FocusedPane {
			case wm.PaneServices:
				if m.WM.ServiceCursor < len(m.Services)-1 {
					m.WM.ServiceCursor++
				}
			case wm.PaneProcesses:
				if m.WM.ProcessCursor < len(m.Processes)-1 {
					m.WM.ProcessCursor++
				}
			case wm.PaneNetwork:
				if m.WM.NetworkCursor < len(m.Ports)-1 {
					m.WM.NetworkCursor++
				}
			default:
				m.WM.NextPane()
			}
			return m, nil

		case "left", "h":
			m.WM.PrevPane()
			return m, nil

		case "right":
			m.WM.NextPane()
			return m, nil

		case "s", "S":
			// Start Service
			if m.WM.FocusedPane == wm.PaneServices && len(m.Services) > 0 {
				svc := m.Services[m.WM.ServiceCursor]
				out, err := system.ControlService("start", svc.Name)
				if err != nil {
					m.StatusMessage = fmt.Sprintf("Gagal Start %s: %v", svc.Name, err)
				} else {
					m.StatusMessage = fmt.Sprintf("Berhasil Start %s: %s", svc.Name, out)
				}
				m.fetchSystemDataSync()
			}
			return m, nil

		case "r", "R":
			// Restart Service
			if m.WM.FocusedPane == wm.PaneServices && len(m.Services) > 0 {
				svc := m.Services[m.WM.ServiceCursor]
				out, err := system.ControlService("restart", svc.Name)
				if err != nil {
					m.StatusMessage = fmt.Sprintf("Gagal Restart %s: %v", svc.Name, err)
				} else {
					m.StatusMessage = fmt.Sprintf("Berhasil Restart %s: %s", svc.Name, out)
				}
				m.fetchSystemDataSync()
			}
			return m, nil

		// Stop Service with Safety Guard
		case "x", "X":
			if m.WM.FocusedPane == wm.PaneServices && len(m.Services) > 0 {
				svc := m.Services[m.WM.ServiceCursor]
				if strings.Contains(svc.Name, "ssh") || strings.Contains(svc.Name, "systemd") {
					m.promptConfirm(fmt.Sprintf("⚠️ BAHAYA: Menghentikan %s dapat memutus koneksi remote! Yakin Stop? (y/n)", svc.Name), func() {
						_, _ = system.ControlService("stop", svc.Name)
						m.StatusMessage = fmt.Sprintf("Layanan %s dihentikan.", svc.Name)
						m.fetchSystemDataSync()
					})
				} else {
					_, _ = system.ControlService("stop", svc.Name)
					m.StatusMessage = fmt.Sprintf("Layanan %s dihentikan.", svc.Name)
					m.fetchSystemDataSync()
				}
			}
			return m, nil

		case "K":
			// Kill process with Safety Guard
			if m.WM.FocusedPane == wm.PaneProcesses && len(m.Processes) > 0 {
				proc := m.Processes[m.WM.ProcessCursor]
				m.promptConfirm(fmt.Sprintf("⚠️ Yakin ingin mematikan proses [%s] (PID %d)? (y/n)", proc.Name, proc.PID), func() {
					err := system.KillProcess(proc.PID, true)
					if err != nil {
						m.StatusMessage = fmt.Sprintf("Gagal mematikan PID %d: %v", proc.PID, err)
					} else {
						m.StatusMessage = fmt.Sprintf("Proses %s (PID %d) berhasil dimatikan.", proc.Name, proc.PID)
					}
					m.fetchSystemDataSync()
				})
			}
			return m, nil

		case "u", "U":
			// Toggle UFW Firewall
			newStatus := !m.FirewallEnabled
			out, err := system.ToggleUFW(newStatus)
			if err != nil {
				m.StatusMessage = fmt.Sprintf("Gagal toggle firewall: %v", err)
			} else {
				m.FirewallEnabled = newStatus
				m.StatusMessage = fmt.Sprintf("Status Firewall: %s", out)
			}
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) promptConfirm(prompt string, action func()) {
	m.ShowConfirm = true
	m.ConfirmPrompt = prompt
	m.ConfirmAction = action
}

func (m *Model) executeLauncherItem(item launcher.LauncherItem) (tea.Model, tea.Cmd) {
	cmd := item.Command

	switch {
	case cmd == "clean_ram":
		out, _ := system.CleanRAMCache()
		m.StatusMessage = "⚡ " + out
		m.LastCleanTime = time.Now()
		m.fetchSystemDataSync()

	case cmd == "audit_ports":
		m.WM.FocusedPane = wm.PaneNetwork
		m.StatusMessage = "🔍 Fokus ke audit Port & Jaringan."

	case cmd == "toggle_ufw":
		newStatus := !m.FirewallEnabled
		_, _ = system.ToggleUFW(newStatus)
		m.FirewallEnabled = newStatus
		m.StatusMessage = fmt.Sprintf("🛡️ Status Firewall UFW: %t", newStatus)

	case cmd == "kill_highest_oom":
		// Find highest OOM process
		var highestProc *system.ProcessItem
		for i := range m.Processes {
			if highestProc == nil || m.Processes[i].OOMScore > highestProc.OOMScore {
				highestProc = &m.Processes[i]
			}
		}
		if highestProc != nil {
			proc := *highestProc
			m.promptConfirm(fmt.Sprintf("⚠️ Bunuh proses OOM tertinggi [%s] (PID %d, Skor %d)? (y/n)", proc.Name, proc.PID, proc.OOMScore), func() {
				_ = system.KillProcess(proc.PID, true)
				m.StatusMessage = fmt.Sprintf("Proses %s (PID %d) berhasil dimatikan.", proc.Name, proc.PID)
				m.fetchSystemDataSync()
			})
		}

	case cmd == "next_theme":
		m.ThemeMode = styles.NextTheme(m.ThemeMode)
		m.StatusMessage = fmt.Sprintf("Tema aktif: %s", styles.ThemeName(m.ThemeMode))
		m.saveCurrentConfig()

	// 1. Layout Mode Switcher from Spotlight
	case cmd == "layout_tiling":
		m.WM.Mode = wm.LayoutTiling2x2
		m.StatusMessage = "Layout diubah ke: Tiling Grid 2x2"
		m.saveCurrentConfig()

	case cmd == "layout_master":
		m.WM.Mode = wm.LayoutMasterStack
		m.StatusMessage = "Layout diubah ke: Compact Master-Stack"
		m.saveCurrentConfig()

	case cmd == "layout_split":
		m.WM.Mode = wm.LayoutHorizontalSplit
		m.StatusMessage = "Layout diubah ke: Horizontal Split"
		m.saveCurrentConfig()

	// 2. FPS Diagnostics Configuration from Spotlight
	case cmd == "fps_adaptive":
		m.TargetFPSConfig = 24
		m.StatusMessage = "Mode FPS: Adaptif 24 FPS (Hemat daya 2 FPS saat diam)"
		m.saveCurrentConfig()

	case cmd == "fps_60":
		m.TargetFPSConfig = 60
		m.StatusMessage = "Mode FPS: Kunci 60 FPS Ultra Smooth"
		m.saveCurrentConfig()

	case cmd == "fps_eco":
		m.TargetFPSConfig = 5
		m.StatusMessage = "Mode FPS: Kunci 5 FPS Statis Hemat Daya"
		m.saveCurrentConfig()

	// 3. Auto Clean RAM Schedule Configuration from Spotlight
	case cmd == "autoclean_15m":
		m.AutoCleanMinutes = 15
		m.StatusMessage = "Auto Clean RAM dijadwalkan: Tiap 15 Menit"
		m.saveCurrentConfig()
		return *m, scheduleAutoClean(15)

	case cmd == "autoclean_30m":
		m.AutoCleanMinutes = 30
		m.StatusMessage = "Auto Clean RAM dijadwalkan: Tiap 30 Menit"
		m.saveCurrentConfig()
		return *m, scheduleAutoClean(30)

	case cmd == "autoclean_1h":
		m.AutoCleanMinutes = 60
		m.StatusMessage = "Auto Clean RAM dijadwalkan: Tiap 1 Jam"
		m.saveCurrentConfig()
		return *m, scheduleAutoClean(60)

	case cmd == "autoclean_off":
		m.AutoCleanMinutes = 0
		m.StatusMessage = "Auto Clean RAM dinonaktifkan (Manual saja via >clean)"
		m.saveCurrentConfig()

	// 4. Animated Background Configuration from Spotlight
	case cmd == "bg_auto":
		m.AnimatedBG = "auto"
		m.StatusMessage = "Background: Animasi Dinamis Sesuai Jam Server (Auto)"
		m.saveCurrentConfig()

	case cmd == "bg_sawah":
		m.AnimatedBG = "sawah"
		m.StatusMessage = "Background: 🌾 Sawah Terasering & Kerbau Berjalan (24 FPS)"
		m.saveCurrentConfig()

	case cmd == "bg_gunung":
		m.AnimatedBG = "gunung"
		m.StatusMessage = "Background: 🏔️ Gunung Slamet Megah & Awan Berarak (24 FPS)"
		m.saveCurrentConfig()

	case cmd == "bg_sakalibel":
		m.AnimatedBG = "sakalibel"
		m.StatusMessage = "Background: 🚂 Kereta Api Jembatan Sakalibel (24 FPS)"
		m.saveCurrentConfig()

	case cmd == "bg_senja":
		m.AnimatedBG = "senja"
		m.StatusMessage = "Background: 🏮 Lentera Senja & Teh Poci Bumiayu (24 FPS)"
		m.saveCurrentConfig()

	case cmd == "bg_malam":
		m.AnimatedBG = "malam"
		m.StatusMessage = "Background: 🌌 Langit Malam Berbintang & Meteor (24 FPS)"
		m.saveCurrentConfig()

	case cmd == "bg_curug":
		m.AnimatedBG = "curug"
		m.StatusMessage = "Background: 🌊 Curug Cipendok Air Terjun (24 FPS)"
		m.saveCurrentConfig()

	case cmd == "bg_off":
		m.AnimatedBG = "off"
		m.StatusMessage = "Background: Animasi dimatikan (Mode Minimalis Bersih)"
		m.saveCurrentConfig()

	case cmd == "bg_dim_low":
		m.BGDimLevel = "low"
		m.StatusMessage = "Opasitas Latar: Rendah (Transparan Terang)"
		m.saveCurrentConfig()

	case cmd == "bg_dim_med":
		m.BGDimLevel = "med"
		m.StatusMessage = "Opasitas Latar: Sedang (Seimbang / Default)"
		m.saveCurrentConfig()

	case cmd == "bg_dim_high":
		m.BGDimLevel = "high"
		m.StatusMessage = "Opasitas Latar: Tinggi (Sangat Redup / Fokus Teks)"
		m.saveCurrentConfig()

	case strings.HasPrefix(cmd, "jump_"):
		var idx int
		fmt.Sscanf(cmd, "jump_%d", &idx)
		m.WM.FocusedPane = wm.PaneID(idx)
		m.StatusMessage = fmt.Sprintf("Lompat ke Panel %d", idx+1)

	case strings.HasPrefix(cmd, "service:"):
		svcName := strings.TrimPrefix(cmd, "service:")
		m.WM.FocusedPane = wm.PaneServices
		for i, s := range m.Services {
			if s.Name == svcName {
				m.WM.ServiceCursor = i
				break
			}
		}
		m.StatusMessage = fmt.Sprintf("Layanan %s dipilih. Tekan [S] Start, [T] Stop, [R] Restart.", svcName)

	case strings.HasPrefix(cmd, "kill_process:"):
		procStr := strings.TrimPrefix(cmd, "kill_process:")
		var name string
		var pid int
		fmt.Sscanf(procStr, "%s (%d)", &name, &pid)
		if pid > 0 {
			m.promptConfirm(fmt.Sprintf("⚠️ Matikan proses %s (PID %d)? (y/n)", name, pid), func() {
				_ = system.KillProcess(pid, true)
				m.StatusMessage = fmt.Sprintf("Proses %s berhasil dimatikan.", name)
				m.fetchSystemDataSync()
			})
		}
	}

	return *m, nil
}

// View renders the terminal window manager interface
func (m Model) View() string {
	width := m.width
	height := m.height
	if width < 80 {
		width = 80
	}
	if height < 24 {
		height = 24
	}

	st := styles.GetTheme(m.ThemeMode)

	// 1. CONFIRMATION SAFETY GUARD MODAL OVERLAY
	if m.ShowConfirm {
		promptRender := lipgloss.NewStyle().Bold(true).Foreground(st.DangerColor).Render(m.ConfirmPrompt)
		hint := lipgloss.NewStyle().Foreground(st.TextColor).Render("Tekan [Y / Enter] untuk Lanjut | [N / Esc] untuk Batal")
		box := st.Modal.Width(60).Render(lipgloss.JoinVertical(lipgloss.Center, promptRender, "", hint))
		placed := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
		if m.AnimatedBG != "off" {
			canvas := retro.Generate2DCanvas(retro.AnimatedBGType(m.AnimatedBG), m.AnimFrame, width, height, m.BGDimLevel)
			return retro.CompositeCanvas(canvas, placed, width, height, m.BGDimLevel, st.CardBg)
		}
		return placed
	}

	// 2. ARCADE HELP 'HOW TO PLAY' MODAL OVERLAY ('?' key)
	if m.ShowHelp {
		modal := help.RenderArcadeHelpModal(st, width, height)
		placed := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal)
		if m.AnimatedBG != "off" {
			canvas := retro.Generate2DCanvas(retro.AnimatedBGType(m.AnimatedBG), m.AnimFrame, width, height, m.BGDimLevel)
			return retro.CompositeCanvas(canvas, placed, width, height, m.BGDimLevel, st.CardBg)
		}
		return placed
	}

	// 3. UNIFIED SPOTLIGHT LAUNCHER MODAL OVERLAY ('Space' or '/' key)
	if m.Spotlight.Active {
		spotlightModal := launcher.RenderSpotlight(m.Spotlight, st, width, height)
		placed := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, spotlightModal)
		if m.AnimatedBG != "off" {
			canvas := retro.Generate2DCanvas(retro.AnimatedBGType(m.AnimatedBG), m.AnimFrame, width, height, m.BGDimLevel)
			return retro.CompositeCanvas(canvas, placed, width, height, m.BGDimLevel, st.CardBg)
		}
		return placed
	}

	// 4. MAIN WINDOW MANAGER INTERFACE
	// 4a. Topbar: 8-Bit Mascot + Telemetry pills + Live Rolling FPS Diagnostic + BETA Flag & Version
	autoCleanStatus := ""
	if m.AutoCleanMinutes > 0 {
		autoCleanStatus = fmt.Sprintf("Clean:%dm", m.AutoCleanMinutes)
	}
	topBar := components.RenderTopBar(m.Hardware, m.AnimFrame, st, m.Language, width, m.CurrentFPS, m.FPSCategory, AppVersion, autoCleanStatus)

	// 4b. Running 24 FPS Status Marquee Ticker
	tickerMsgs := []string{
		"⚡ BUMIAYU TUI WINDOW MANAGER",
		fmt.Sprintf("HOST: %s (%s)", m.Hardware.Hostname, m.Hardware.Arch),
		fmt.Sprintf("TEMA: %s", st.Name),
		fmt.Sprintf("BG: %s [%s]", retro.BGName(retro.AnimatedBGType(m.AnimatedBG)), strings.ToUpper(m.BGDimLevel)),
		fmt.Sprintf("FPS: %.0f [%s]", m.CurrentFPS, m.FPSCategory),
		"TEKAN [SPACE] UNTUK SPOTLIGHT LAUNCHER",
		"TEKAN [?] UNTUK BANTUAN ARCADE CONTROLS",
		fmt.Sprintf("CPU LOAD: %.2f | RAM: %d MB", m.Hardware.Load1, m.Hardware.MemUsedMB),
	}
	tickerLine := retro.RenderMarqueeTicker(tickerMsgs, m.AnimFrame, width, st)

	// 4c. Main Workspace: Tiled Layout, Master-Stack, or Zoomed Pane (Uses full height!)
	tiledLayout := wm.RenderTiledLayout(
		m.WM,
		m.Hardware,
		m.Services,
		m.Processes,
		m.Ports,
		m.FirewallEnabled,
		m.AnimFrame,
		st,
		width,
		height,
	)

	// 4d. Persistent Bottom Quickbar with Shortcuts Help
	quickBar := help.RenderQuickBar(st, width, m.StatusMessage, m.WM.Zoomed)

	fullContent := lipgloss.JoinVertical(lipgloss.Left,
		topBar,
		tickerLine,
		tiledLayout,
		quickBar,
	)

	// Clamp to terminal height so no scrolling occurs inside alternate screen buffer
	lines := strings.Split(fullContent, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	cleanContent := strings.Join(lines, "\n")

	// 5. FULL-SCREEN 2D COMPOSITOR WITH TRANSLUCENT TINT
	if m.AnimatedBG != "off" {
		canvas := retro.Generate2DCanvas(retro.AnimatedBGType(m.AnimatedBG), m.AnimFrame, width, height, m.BGDimLevel)
		return retro.CompositeCanvas(canvas, cleanContent, width, height, m.BGDimLevel, st.CardBg)
	}

	return cleanContent
}
