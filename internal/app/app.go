package app

import (
	"time"

	"github.com/bumiayu/bumiayu-de/internal/config"
	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/components"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/bumiayu/bumiayu-de/internal/ui/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TickMsg triggers telemetry refresh every 2 seconds
type TickMsg time.Time

// AppActionMsg reports result of background install/uninstall
type AppActionMsg struct {
	Logs    string
	Success bool
}

// Model is the main state of Bumiayu TUI Desktop Environment
type Model struct {
	width  int
	height int

	ActiveWorkspace components.Workspace
	ThemeMode       styles.ThemeMode
	WMMode          views.WMMode
	Language        i18n.Language
	StatusMessage   string

	// Telemetry data
	Hardware system.HardwareInfo

	// 2. Apps data & cursor & dialog
	Catalog       []system.AppCatalogItem
	CatalogCursor int
	PackageCount  int
	AppDialog     views.AppDialogState

	// 3. Settings cursor
	SettingsCursor int

	// 4. First-run OOBE Wizard & Recovery State
	Wizard views.WizardState
}

// InitialModel initializes the state with config or opens wizard if first run
func InitialModel() Model {
	m := Model{
		ActiveWorkspace: components.WorkspaceDashboard,
		ThemeMode:       styles.ThemeLight,
		WMMode:          views.WMFloating,
		Language:        i18n.ID,
	}

	// Check configuration & first run
	cfg, err := config.LoadConfig()
	if err == nil && cfg.FirstRunCompleted {
		m.applyUserConfig(cfg)
		m.Wizard.Active = false
	} else {
		// Launch First-run OOBE Wizard
		m.Wizard.Active = true
		m.Wizard.HasBackup = config.HasExistingSetup()
		if m.Wizard.HasBackup {
			m.Wizard.CurrentStep = views.StepRecovery
		} else {
			m.Wizard.CurrentStep = views.StepLanguage
		}
		m.Wizard.DraftConfig = config.DefaultConfig()
	}

	m.refreshAllData()
	return m
}

func (m *Model) applyUserConfig(cfg config.UserConfig) {
	switch cfg.Theme {
	case "warm":
		m.ThemeMode = styles.ThemeWarm
	case "dark":
		m.ThemeMode = styles.ThemeDark
	default:
		m.ThemeMode = styles.ThemeLight
	}

	switch cfg.WMModel {
	case "tiling":
		m.WMMode = views.WMTiling
	case "cardgrid":
		m.WMMode = views.WMCardGrid
	default:
		m.WMMode = views.WMFloating
	}

	if cfg.Language == "en" {
		m.Language = i18n.EN
	} else {
		m.Language = i18n.ID
	}
}

func (m *Model) refreshAllData() {
	m.Hardware = system.GetHardwareInfo()
	m.Catalog = system.GetAppCatalog()
	m.PackageCount = system.GetInstalledPackageCount()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Init starts the Bubble Tea loop
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// Update handles keyboard and window events
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case TickMsg:
		m.Hardware = system.GetHardwareInfo()
		return m, tickCmd()

	case AppActionMsg:
		m.AppDialog.IsRunning = false
		m.AppDialog.Logs = msg.Logs
		m.AppDialog.Success = msg.Success
		m.refreshAllData()
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// 1. WIZARD / RECOVERY HANDLING
		if m.Wizard.Active {
			return m.handleWizardKeys(key)
		}

		// 2. APP INSTALL/UNINSTALL MODAL HANDLING
		if m.AppDialog.Active {
			return m.handleAppModalKeys(key)
		}

		// 3. GLOBAL DESKTOP SHORTCUTS
		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit

		// Quick Theme Toggle (T)
		case "t", "T":
			m.ThemeMode = (m.ThemeMode + 1) % 3
			themeName := "Terang"
			if m.ThemeMode == styles.ThemeWarm {
				themeName = "Hangat"
			} else if m.ThemeMode == styles.ThemeDark {
				themeName = "Gelap"
			}
			m.StatusMessage = "Tema diubah ke: " + themeName
			m.saveCurrentConfig()
			return m, nil

		// Open Setup Wizard anytime (W)
		case "w", "W":
			m.Wizard.Active = true
			m.Wizard.CurrentStep = views.StepLanguage
			m.Wizard.SelectedRow = 0
			m.Wizard.DraftConfig = config.DefaultConfig()
			return m, nil

		// Direct Tab Switching via numbers 1 - 4
		case "1":
			m.ActiveWorkspace = components.WorkspaceDashboard
			return m, nil
		case "2":
			m.ActiveWorkspace = components.WorkspaceApps
			return m, nil
		case "3":
			m.ActiveWorkspace = components.WorkspaceSettings
			return m, nil
		case "4":
			m.ActiveWorkspace = components.WorkspaceGuide
			return m, nil

		// Dynamic Arrow Tab Switching (Left / Right)
		case "left", "h":
			if m.ActiveWorkspace > 0 {
				m.ActiveWorkspace--
			} else {
				m.ActiveWorkspace = components.WorkspaceGuide
			}
			return m, nil

		case "right", "l":
			m.ActiveWorkspace = (m.ActiveWorkspace + 1) % 4
			return m, nil

		case "tab":
			m.ActiveWorkspace = (m.ActiveWorkspace + 1) % 4
			return m, nil

		case "shift+tab":
			if m.ActiveWorkspace > 0 {
				m.ActiveWorkspace--
			} else {
				m.ActiveWorkspace = components.WorkspaceGuide
			}
			return m, nil

		// Navigation inside active workspace
		case "up", "k":
			m.moveCursor(-1)
			return m, nil

		case "down", "j":
			m.moveCursor(1)
			return m, nil

		// Action Keys
		case "enter", "i", "I":
			if m.ActiveWorkspace == components.WorkspaceApps && len(m.Catalog) > 0 {
				selected := m.Catalog[m.CatalogCursor]
				action := "install"
				if selected.Installed {
					action = "uninstall"
				}
				m.AppDialog = views.AppDialogState{
					Active: true,
					Action: action,
					App:    selected,
				}
			} else if m.ActiveWorkspace == components.WorkspaceSettings {
				m.applySettingSelection()
			}
			return m, nil

		case "x", "X":
			if m.ActiveWorkspace == components.WorkspaceApps && len(m.Catalog) > 0 {
				selected := m.Catalog[m.CatalogCursor]
				if selected.Installed {
					m.AppDialog = views.AppDialogState{
						Active: true,
						Action: "uninstall",
						App:    selected,
					}
				} else {
					m.StatusMessage = "Aplikasi ini belum terpasang."
				}
			}
			return m, nil

		case " ":
			if m.ActiveWorkspace == components.WorkspaceSettings {
				m.applySettingSelection()
			}
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) handleWizardKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.Wizard.SelectedRow > 0 {
			m.Wizard.SelectedRow--
		}
	case "down", "j":
		m.Wizard.SelectedRow++
	case "enter", " ":
		switch m.Wizard.CurrentStep {
		case views.StepRecovery:
			if m.Wizard.SelectedRow == 0 {
				// Recover previous setup
				cfg, err := config.LoadBackupConfig()
				if err == nil {
					m.applyUserConfig(cfg)
					m.StatusMessage = "Pengaturan lama berhasil dipulihkan (Recovery)!"
				}
				m.Wizard.Active = false
				return m, nil
			}
			// Fresh Setup
			m.Wizard.CurrentStep = views.StepLanguage
			m.Wizard.SelectedRow = 0

		case views.StepLanguage:
			if m.Wizard.SelectedRow == 1 {
				m.Language = i18n.EN
				m.Wizard.DraftConfig.Language = "en"
			} else {
				m.Language = i18n.ID
				m.Wizard.DraftConfig.Language = "id"
			}
			m.Wizard.CurrentStep = views.StepTheme
			m.Wizard.SelectedRow = 0

		case views.StepTheme:
			switch m.Wizard.SelectedRow {
			case 1:
				m.ThemeMode = styles.ThemeWarm
				m.Wizard.DraftConfig.Theme = "warm"
			case 2:
				m.ThemeMode = styles.ThemeDark
				m.Wizard.DraftConfig.Theme = "dark"
			default:
				m.ThemeMode = styles.ThemeLight
				m.Wizard.DraftConfig.Theme = "light"
			}
			m.Wizard.CurrentStep = views.StepWM
			m.Wizard.SelectedRow = 0

		case views.StepWM:
			switch m.Wizard.SelectedRow {
			case 1:
				m.WMMode = views.WMTiling
				m.Wizard.DraftConfig.WMModel = "tiling"
			case 2:
				m.WMMode = views.WMCardGrid
				m.Wizard.DraftConfig.WMModel = "cardgrid"
			default:
				m.WMMode = views.WMFloating
				m.Wizard.DraftConfig.WMModel = "floating"
			}
			m.Wizard.CurrentStep = views.StepSSH
			m.Wizard.SelectedRow = 0

		case views.StepSSH:
			m.Wizard.DraftConfig.AutoStartSSH = (m.Wizard.SelectedRow == 0)
			m.Wizard.DraftConfig.FirstRunCompleted = true
			config.SaveConfig(m.Wizard.DraftConfig)
			m.Wizard.CurrentStep = views.StepComplete

		case views.StepComplete:
			m.Wizard.Active = false
			m.StatusMessage = "Selamat datang di Bumiayu TUI!"
		}
	}
	return m, nil
}

func (m *Model) handleAppModalKeys(key string) (tea.Model, tea.Cmd) {
	if m.AppDialog.IsRunning {
		return m, nil
	}

	if m.AppDialog.Logs != "" {
		if key == "enter" || key == "esc" || key == "q" || key == " " {
			m.AppDialog.Active = false
			m.AppDialog.Logs = ""
			return m, nil
		}
	}

	if key == "y" || key == "Y" {
		m.AppDialog.IsRunning = true
		m.AppDialog.Logs = "Menjalankan proses " + m.AppDialog.Action + " " + m.AppDialog.App.PackageName + "..."
		action := m.AppDialog.Action
		pkg := m.AppDialog.App.PackageName
		return m, func() tea.Msg {
			out, err := system.ExecutePackageAction(action, pkg)
			success := err == nil
			return AppActionMsg{Logs: out, Success: success}
		}
	}

	if key == "n" || key == "N" || key == "esc" {
		m.AppDialog.Active = false
		return m, nil
	}
	return m, nil
}

func (m *Model) moveCursor(delta int) {
	switch m.ActiveWorkspace {
	case components.WorkspaceApps:
		m.CatalogCursor = clamp(m.CatalogCursor+delta, 0, len(m.Catalog)-1)
	case components.WorkspaceSettings:
		m.SettingsCursor = clamp(m.SettingsCursor+delta, 0, 5)
	}
}

func (m *Model) applySettingSelection() {
	if m.SettingsCursor < 3 {
		m.ThemeMode = styles.ThemeMode(m.SettingsCursor)
		m.StatusMessage = "Tema berhasil diterapkan!"
	} else {
		m.WMMode = views.WMMode(m.SettingsCursor - 3)
		m.StatusMessage = "Model Window Manager berhasil diubah!"
	}
	m.saveCurrentConfig()
}

func (m *Model) saveCurrentConfig() {
	cfg := config.UserConfig{
		ConfigVersion:     "v0.9.0-beta.3",
		FirstRunCompleted: true,
		AutoStartSSH:      false,
	}
	switch m.ThemeMode {
	case styles.ThemeWarm:
		cfg.Theme = "warm"
	case styles.ThemeDark:
		cfg.Theme = "dark"
	default:
		cfg.Theme = "light"
	}
	switch m.WMMode {
	case views.WMTiling:
		cfg.WMModel = "tiling"
	case views.WMCardGrid:
		cfg.WMModel = "cardgrid"
	default:
		cfg.WMModel = "floating"
	}
	if m.Language == i18n.EN {
		cfg.Language = "en"
	} else {
		cfg.Language = "id"
	}
	config.SaveConfig(cfg)
}

func clamp(v, min, max int) int {
	if max < min {
		return min
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// View renders the complete Desktop Environment or Wizard modal
func (m Model) View() string {
	width := m.width
	if width < 80 {
		width = 80
	}
	height := m.height
	if height < 24 {
		height = 24
	}

	st := styles.GetTheme(m.ThemeMode)

	// If Wizard is active, render centered Setup Wizard
	if m.Wizard.Active {
		return views.RenderWizard(m.Wizard, st, width, height)
	}

	// 1. Top Desktop Bar
	topBar := components.RenderTopBar(m.Hardware, st, m.Language, width)

	// 2. Workspace Navigation Tabs
	workspaceTabs := components.RenderWorkspaceTabs(m.ActiveWorkspace, st, width)

	// 3. Main Workspace View
	var currentView string
	switch m.ActiveWorkspace {
	case components.WorkspaceDashboard:
		currentView = views.RenderDashboard(m.Hardware, m.WMMode, st, m.Language, width, height)
	case components.WorkspaceApps:
		currentView = views.RenderApps(m.Catalog, m.CatalogCursor, m.PackageCount, m.AppDialog, st, m.Language, width, height)
	case components.WorkspaceSettings:
		currentView = views.RenderSettings(m.ThemeMode, m.WMMode, m.SettingsCursor, st, width, height)
	case components.WorkspaceGuide:
		currentView = views.RenderGuide(st, width, height)
	}

	// 4. Status Bar
	statusBar := components.RenderStatusBar(st, width, m.StatusMessage)

	return lipgloss.JoinVertical(lipgloss.Left,
		topBar,
		workspaceTabs,
		"",
		currentView,
		"",
		statusBar,
	)
}
