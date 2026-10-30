package app

import (
	"time"

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
}

// InitialModel initializes the state with default clean Light Theme
func InitialModel() Model {
	m := Model{
		ActiveWorkspace: components.WorkspaceDashboard,
		ThemeMode:       styles.ThemeLight, // Default: Tampilan Cerah Bersih & Simpel
		WMMode:          views.WMFloating,  // Default: Window Mengambang Minimalis
		Language:        i18n.ID,
	}

	m.refreshAllData()
	return m
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

		// If App Modal is Active, handle modal shortcuts
		if m.AppDialog.Active {
			if m.AppDialog.IsRunning {
				return m, nil // Wait for completion
			}

			if m.AppDialog.Logs != "" {
				// Finished view: any key closes modal
				if key == "enter" || key == "esc" || key == "q" || key == " " {
					m.AppDialog.Active = false
					m.AppDialog.Logs = ""
					return m, nil
				}
			}

			// Confirmation view
			if key == "y" || key == "Y" {
				m.AppDialog.IsRunning = true
				m.AppDialog.Logs = "Menjalankan perintah " + m.AppDialog.Action + " " + m.AppDialog.App.PackageName + "..."
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

		// Global navigation shortcuts
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

func (m *Model) moveCursor(delta int) {
	switch m.ActiveWorkspace {
	case components.WorkspaceApps:
		m.CatalogCursor = clamp(m.CatalogCursor+delta, 0, len(m.Catalog)-1)
	case components.WorkspaceSettings:
		m.SettingsCursor = clamp(m.SettingsCursor+delta, 0, 5) // 0-2 themes, 3-5 WM models
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

// View renders the complete Desktop Environment
func (m Model) View() string {
	width := m.width
	if width < 80 {
		width = 80
	}
	height := m.height
	if height < 24 {
		height = 24
	}

	// 1. Get active theme styling
	st := styles.GetTheme(m.ThemeMode)

	// 2. Top Desktop Bar
	topBar := components.RenderTopBar(m.Hardware, st, m.Language, width)

	// 3. Workspace Navigation Tabs
	workspaceTabs := components.RenderWorkspaceTabs(m.ActiveWorkspace, st, width)

	// 4. Main Workspace View
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

	// 5. Status Bar
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
