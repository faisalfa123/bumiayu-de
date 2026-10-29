package app

import (
	"fmt"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/components"
	"github.com/bumiayu/bumiayu-de/internal/ui/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TickMsg triggers telemetry refresh every 2 seconds
type TickMsg time.Time

// Model is the main state of Bumiayu TUI Desktop Environment
type Model struct {
	width  int
	height int

	ActiveWorkspace components.Workspace
	Language        i18n.Language
	StatusMessage   string

	// Telemetry data
	Hardware system.HardwareInfo

	// Workspaces Data & Cursors
	// 2. Apps
	Catalog          []system.AppCatalogItem
	CatalogCursor    int
	PackageCount     int

	// 3. Services
	Services         []system.ServiceItem
	ServicesCursor   int
	ShowingLogs      bool
	CurrentLogs      string

	// 4. Security
	Security         system.SecurityState

	// 5. Users
	Users            []system.UserAccount
	UsersCursor      int

	// 6. Files
	CurrentPath      string
	FileEntries      []system.FileEntry
	FilesCursor      int
	FilePreview      string

	// 7. Tasks
	Processes        []system.ProcessItem
	ProcessesCursor  int
}

// InitialModel initializes the state
func InitialModel() Model {
	m := Model{
		ActiveWorkspace: components.WorkspaceDashboard,
		Language:        i18n.ID,
		CurrentPath:     "/etc",
	}

	m.refreshAllData()
	return m
}

func (m *Model) refreshAllData() {
	m.Hardware = system.GetHardwareInfo()
	m.Catalog = system.GetAppCatalog()
	m.PackageCount = system.GetInstalledPackageCount()
	m.Services = system.GetServices()
	m.Security = system.GetSecurityState()
	m.Users = system.GetUserAccounts()
	m.Processes = system.GetProcessList()

	entries, err := system.ListDirectory(m.CurrentPath)
	if err == nil {
		m.FileEntries = entries
		if len(entries) > 0 && !entries[0].IsDir {
			preview, _ := system.ReadFilePreview(entries[0].Path, 15)
			m.FilePreview = preview
		}
	}
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
		// Periodic telemetry update
		m.Hardware = system.GetHardwareInfo()
		m.Processes = system.GetProcessList()
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// Switch language
		case "L", "l":
			if m.ActiveWorkspace != components.WorkspaceServices {
				if m.Language == i18n.ID {
					m.Language = i18n.EN
					m.StatusMessage = "Language switched to English"
				} else {
					m.Language = i18n.ID
					m.StatusMessage = "Bahasa diganti ke Bahasa Indonesia"
				}
				return m, nil
			}

		// Workspace navigation
		case "1":
			m.ActiveWorkspace = components.WorkspaceDashboard
			return m, nil
		case "2":
			m.ActiveWorkspace = components.WorkspaceApps
			return m, nil
		case "3":
			m.ActiveWorkspace = components.WorkspaceServices
			return m, nil
		case "4":
			m.ActiveWorkspace = components.WorkspaceSecurity
			return m, nil
		case "5":
			m.ActiveWorkspace = components.WorkspaceUsers
			return m, nil
		case "6":
			m.ActiveWorkspace = components.WorkspaceFiles
			return m, nil
		case "7":
			m.ActiveWorkspace = components.WorkspaceTasks
			return m, nil

		case "tab":
			m.ActiveWorkspace = (m.ActiveWorkspace + 1) % 7
			return m, nil

		case "shift+tab":
			if m.ActiveWorkspace == 0 {
				m.ActiveWorkspace = 6
			} else {
				m.ActiveWorkspace--
			}
			return m, nil

		// Navigation inside active workspace
		case "up", "k":
			m.moveCursor(-1)
			return m, nil

		case "down", "j":
			m.moveCursor(1)
			return m, nil

		case "esc":
			if m.ShowingLogs {
				m.ShowingLogs = false
				return m, nil
			}

		// Action keys
		case "s", "S":
			if m.ActiveWorkspace == components.WorkspaceServices && len(m.Services) > 0 {
				svc := m.Services[m.ServicesCursor].Name
				out, err := system.ControlService("start", svc)
				if err != nil {
					m.StatusMessage = "Gagal menyalakan: " + err.Error()
				} else {
					m.StatusMessage = "Service " + svc + " dinyalakan."
					_ = out
				}
				m.Services = system.GetServices()
			}

		case "t", "T":
			if m.ActiveWorkspace == components.WorkspaceServices && len(m.Services) > 0 {
				svc := m.Services[m.ServicesCursor].Name
				out, err := system.ControlService("stop", svc)
				if err != nil {
					m.StatusMessage = "Gagal menghentikan: " + err.Error()
				} else {
					m.StatusMessage = "Service " + svc + " dihentikan."
					_ = out
				}
				m.Services = system.GetServices()
			}

		case "r", "R":
			if m.ActiveWorkspace == components.WorkspaceServices && len(m.Services) > 0 {
				svc := m.Services[m.ServicesCursor].Name
				out, err := system.ControlService("restart", svc)
				if err != nil {
					m.StatusMessage = "Gagal restart: " + err.Error()
				} else {
					m.StatusMessage = "Service " + svc + " berhasil di-restart."
					_ = out
				}
				m.Services = system.GetServices()
			}

		case "u", "U":
			if m.ActiveWorkspace == components.WorkspaceSecurity {
				newState := !m.Security.UfwActive
				system.ToggleUFW(newState)
				m.Security = system.GetSecurityState()
				if newState {
					m.StatusMessage = "Firewall UFW Diaktifkan."
				} else {
					m.StatusMessage = "Firewall UFW Dinonaktifkan."
				}
			}

		case "g", "G":
			if m.ActiveWorkspace == components.WorkspaceUsers && len(m.Users) > 0 {
				user := m.Users[m.UsersCursor].Username
				system.ModifyUserSudoRole(user, true)
				m.StatusMessage = "Hak akses Sudo diberikan ke: " + user
				m.Users = system.GetUserAccounts()
			}

		case "v", "V":
			if m.ActiveWorkspace == components.WorkspaceUsers && len(m.Users) > 0 {
				user := m.Users[m.UsersCursor].Username
				system.ModifyUserSudoRole(user, false)
				m.StatusMessage = "Hak akses Sudo dicabut dari: " + user
				m.Users = system.GetUserAccounts()
			}

		case "K":
			if m.ActiveWorkspace == components.WorkspaceTasks && len(m.Processes) > 0 {
				pid := m.Processes[m.ProcessesCursor].PID
				name := m.Processes[m.ProcessesCursor].Name
				err := system.KillProcess(pid, false)
				if err != nil {
					m.StatusMessage = fmt.Sprintf("Gagal kill %s: %s", name, err.Error())
				} else {
					m.StatusMessage = fmt.Sprintf("Proses %s (PID %d) telah dihentikan (SIGTERM)", name, pid)
				}
				m.Processes = system.GetProcessList()
			}

		case "enter":
			if m.ActiveWorkspace == components.WorkspaceServices && len(m.Services) > 0 {
				svc := m.Services[m.ServicesCursor].Name
				m.CurrentLogs = system.GetServiceLogs(svc, 30)
				m.ShowingLogs = !m.ShowingLogs
			} else if m.ActiveWorkspace == components.WorkspaceFiles && len(m.FileEntries) > 0 {
				entry := m.FileEntries[m.FilesCursor]
				if entry.IsDir {
					m.CurrentPath = entry.Path
					entries, err := system.ListDirectory(m.CurrentPath)
					if err == nil {
						m.FileEntries = entries
						m.FilesCursor = 0
						m.FilePreview = ""
					}
				} else {
					preview, _ := system.ReadFilePreview(entry.Path, 18)
					m.FilePreview = preview
				}
			}
		}
	}

	return m, nil
}

func (m *Model) moveCursor(delta int) {
	switch m.ActiveWorkspace {
	case components.WorkspaceApps:
		m.CatalogCursor = clamp(m.CatalogCursor+delta, 0, len(m.Catalog)-1)
	case components.WorkspaceServices:
		m.ServicesCursor = clamp(m.ServicesCursor+delta, 0, len(m.Services)-1)
	case components.WorkspaceUsers:
		m.UsersCursor = clamp(m.UsersCursor+delta, 0, len(m.Users)-1)
	case components.WorkspaceFiles:
		m.FilesCursor = clamp(m.FilesCursor+delta, 0, len(m.FileEntries)-1)
		if len(m.FileEntries) > 0 && !m.FileEntries[m.FilesCursor].IsDir {
			preview, _ := system.ReadFilePreview(m.FileEntries[m.FilesCursor].Path, 18)
			m.FilePreview = preview
		}
	case components.WorkspaceTasks:
		m.ProcessesCursor = clamp(m.ProcessesCursor+delta, 0, len(m.Processes)-1)
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

	// 1. Top Desktop Bar
	topBar := components.RenderTopBar(m.Hardware, m.Language, width)

	// 2. Workspace Tabs
	workspaceTabs := components.RenderWorkspaceTabs(m.ActiveWorkspace, m.Language, width)

	// 3. Main Workspace View
	var currentView string
	switch m.ActiveWorkspace {
	case components.WorkspaceDashboard:
		currentView = views.RenderDashboard(m.Hardware, m.Language, width, height)
	case components.WorkspaceApps:
		currentView = views.RenderApps(m.Catalog, m.CatalogCursor, m.PackageCount, m.Language, width, height)
	case components.WorkspaceServices:
		currentView = views.RenderServices(m.Services, m.ServicesCursor, m.ShowingLogs, m.CurrentLogs, m.Language, width, height)
	case components.WorkspaceSecurity:
		currentView = views.RenderSecurity(m.Security, m.Language, width, height)
	case components.WorkspaceUsers:
		currentView = views.RenderUsers(m.Users, m.UsersCursor, m.Language, width, height)
	case components.WorkspaceFiles:
		currentView = views.RenderFiles(m.CurrentPath, m.FileEntries, m.FilesCursor, m.FilePreview, m.Language, width, height)
	case components.WorkspaceTasks:
		currentView = views.RenderTasks(m.Processes, m.ProcessesCursor, m.Language, width, height)
	}

	// 4. Status Bar
	statusBar := components.RenderStatusBar(m.Language, width, m.StatusMessage)

	return lipgloss.JoinVertical(lipgloss.Left,
		topBar,
		workspaceTabs,
		currentView,
		"",
		statusBar,
	)
}
