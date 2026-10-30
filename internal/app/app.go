package app

import (
	"strings"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/config"
	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/components"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/bumiayu/bumiayu-de/internal/ui/views"
	"github.com/bumiayu/bumiayu-de/internal/ui/wallpaper"
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
	WallpaperType   wallpaper.WallpaperType
	WallpaperPreset wallpaper.PresetScene
	Language        i18n.Language
	StatusMessage   string

	// Telemetry data
	Hardware system.HardwareInfo

	// Apps data & cursor & dialog
	Catalog        []system.AppCatalogItem
	CatalogCursor  int
	ActiveCategory int
	PackageCount   int
	AppDialog      views.AppDialogState

	// Settings cursor & image input modal
	SettingsCursor int
	ShowImageInput bool
	ImageInputPath string
	AnimFrame      int

	// First-run OOBE Wizard & Recovery State
	Wizard views.WizardState

	// Screen Navigation Modes (WebUI Experience)
	CurrentScreen         ScreenMode
	AlexisCursor          int
	CommerceSidebarCursor int
}

// ScreenMode defines whether we are on Alexis Grid Launcher or Commerce Dashboard
type ScreenMode int

const (
	ScreenAlexisGrid ScreenMode = iota
	ScreenCommerceDashboard
)

// InitialModel initializes the state with config or opens wizard if first run
func InitialModel() Model {
	m := Model{
		ActiveWorkspace: components.WorkspaceDashboard,
		ThemeMode:       styles.ThemeLight,
		WMMode:          views.WMFloating,
		WallpaperType:   wallpaper.WallpaperDynamicTime,
		WallpaperPreset: wallpaper.SceneGunungSlamet,
		Language:        i18n.ID,
		ImageInputPath:  "/tmp/wallpaper.jpg",
	}

	cfg, err := config.LoadConfig()
	if err == nil && cfg.FirstRunCompleted {
		m.applyUserConfig(cfg)
		m.Wizard.Active = false
	} else {
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

	switch cfg.WallpaperType {
	case "custom":
		m.WallpaperType = wallpaper.WallpaperCustom
	case "preset":
		m.WallpaperType = wallpaper.WallpaperPreset
	default:
		m.WallpaperType = wallpaper.WallpaperDynamicTime
	}
	m.WallpaperPreset = wallpaper.PresetScene(cfg.WallpaperPreset)

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
	// 24 FPS = 1000ms / 24 ≈ 41.6ms
	return tea.Tick(time.Millisecond*42, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Init starts the Bubble Tea loop with smooth 24fps tick
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
		m.AnimFrame++
		// Polling hardware setiap ~1.5 detik (36 frame @ 24fps) agar CPU hemat daya
		if m.AnimFrame%36 == 0 {
			m.Hardware = system.GetHardwareInfo()
		}
		return m, tickCmd()

	case AppActionMsg:
		m.AppDialog.IsRunning = false
		m.AppDialog.Logs = msg.Logs
		m.AppDialog.Success = msg.Success
		m.refreshAllData()
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// 1. IMAGE PATH INPUT MODAL
		if m.ShowImageInput {
			switch key {
			case "enter":
				m.ShowImageInput = false
				path := strings.TrimSpace(m.ImageInputPath)
				if path != "" {
					_, err := wallpaper.ConvertImageToPixelArt(path, 60, 12)
					if err != nil {
						m.StatusMessage = "Gagal memproses gambar: " + err.Error()
					} else {
						m.WallpaperType = wallpaper.WallpaperCustom
						m.StatusMessage = "Gambar berhasil dikonversi menjadi pixel art & aktif!"
						m.saveCurrentConfig()
					}
				}
				return m, nil
			case "esc":
				m.ShowImageInput = false
				return m, nil
			case "backspace":
				if len(m.ImageInputPath) > 0 {
					m.ImageInputPath = m.ImageInputPath[:len(m.ImageInputPath)-1]
				}
				return m, nil
			default:
				if len(key) == 1 {
					m.ImageInputPath += key
				}
				return m, nil
			}
		}

		// 2. WIZARD / RECOVERY HANDLING
		if m.Wizard.Active {
			return m.handleWizardKeys(key)
		}

		// 3. APP INSTALL/UNINSTALL MODAL HANDLING
		if m.AppDialog.Active {
			return m.handleAppModalKeys(key)
		}

		// 4. SCREEN MODE NAVIGATION (Alexis Icons vs Commerce Dashboard)
		if m.CurrentScreen == ScreenAlexisGrid {
			switch key {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.AlexisCursor >= 4 {
					m.AlexisCursor -= 4
				}
				return m, nil
			case "down", "j":
				if m.AlexisCursor < 4 {
					m.AlexisCursor += 4
				}
				return m, nil
			case "left", "h":
				if m.AlexisCursor%4 > 0 {
					m.AlexisCursor--
				}
				return m, nil
			case "right", "l":
				if m.AlexisCursor%4 < 3 {
					m.AlexisCursor++
				}
				return m, nil
			case "enter", " ":
				// Enter opens Commerce Dashboard with relevant section
				m.CurrentScreen = ScreenCommerceDashboard
				switch m.AlexisCursor {
				case 0: // Dasbor Utama
					m.CommerceSidebarCursor = 0
				case 1: // Toko Aplikasi
					m.CommerceSidebarCursor = 1
				case 2: // Immich SBC
					m.CommerceSidebarCursor = 2
				case 3: // Penyimpanan
					m.CommerceSidebarCursor = 3
				case 4: // Jaringan
					m.CommerceSidebarCursor = 4
				case 5: // Keamanan
					m.CommerceSidebarCursor = 5
				case 6: // Pengaturan
					m.CommerceSidebarCursor = 7
				case 7: // Panduan
					m.CommerceSidebarCursor = 0
				}
				return m, nil
			case "t", "T":
				m.ThemeMode = (m.ThemeMode + 1) % 3
				return m, nil
			}
			return m, nil
		}

		// When on ScreenCommerceDashboard:
		if m.CurrentScreen == ScreenCommerceDashboard {
			switch key {
			case "esc", "backspace":
				m.CurrentScreen = ScreenAlexisGrid
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.CommerceSidebarCursor > 0 {
					m.CommerceSidebarCursor--
				}
				return m, nil
			case "down", "j":
				if m.CommerceSidebarCursor < 7 {
					m.CommerceSidebarCursor++
				}
				return m, nil
			case "left", "h":
				if m.CommerceSidebarCursor == 1 { // Toko Aplikasi
					if m.ActiveCategory > 0 {
						m.ActiveCategory--
					} else {
						m.ActiveCategory = len(views.CategoryFilterList) - 1
					}
					m.CatalogCursor = 0
				}
				return m, nil
			case "right", "l":
				if m.CommerceSidebarCursor == 1 { // Toko Aplikasi
					m.ActiveCategory = (m.ActiveCategory + 1) % len(views.CategoryFilterList)
					m.CatalogCursor = 0
				}
				return m, nil
			case "tab":
				if m.CommerceSidebarCursor == 1 {
					m.CatalogCursor = (m.CatalogCursor + 1) % len(m.getFilteredCatalog())
				}
				return m, nil
			case "enter", "i", "I":
				if m.CommerceSidebarCursor == 2 {
					// Quick install/manage Immich
					for _, item := range m.Catalog {
						if item.ID == "immich-sbc" {
							action := "install"
							if item.Installed {
								action = "uninstall"
							}
							m.AppDialog = views.AppDialogState{
								Active: true,
								Action: action,
								App:    item,
							}
							break
						}
					}
					return m, nil
				}
				if m.CommerceSidebarCursor == 1 {
					filtered := m.getFilteredCatalog()
					if len(filtered) > 0 && m.CatalogCursor < len(filtered) {
						selected := filtered[m.CatalogCursor]
						action := "install"
						if selected.Installed {
							action = "uninstall"
						}
						m.AppDialog = views.AppDialogState{
							Active: true,
							Action: action,
							App:    selected,
						}
					}
					return m, nil
				}
			case "x", "X":
				if m.CommerceSidebarCursor == 1 {
					filtered := m.getFilteredCatalog()
					if len(filtered) > 0 && m.CatalogCursor < len(filtered) {
						selected := filtered[m.CatalogCursor]
						if selected.Installed {
							m.AppDialog = views.AppDialogState{
								Active: true,
								Action: "uninstall",
								App:    selected,
							}
						}
					}
					return m, nil
				}
			}
		}

		// 5. GLOBAL DESKTOP SHORTCUTS (Fallback / Direct)
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

		// Quick Window Manager Mode Toggle (M)
		case "m", "M":
			m.WMMode = (m.WMMode + 1) % 3
			wmName := "Floating Window"
			if m.WMMode == views.WMTiling {
				wmName = "Tiling 2-Kolom"
			} else if m.WMMode == views.WMCardGrid {
				wmName = "Card Grid"
			}
			m.StatusMessage = "Model Window Manager diubah ke: " + wmName
			m.saveCurrentConfig()
			return m, nil

		// Open Image to Pixel Generator (G)
		case "g", "G":
			m.ShowImageInput = true
			return m, nil

		// Open Setup Wizard anytime (W)
		case "w", "W":
			m.Wizard.Active = true
			m.Wizard.CurrentStep = views.StepLanguage
			m.Wizard.SelectedRow = 0
			m.Wizard.DraftConfig = config.DefaultConfig()
			return m, nil

		// Direct Tab Switching via numbers 1 - 3
		case "1":
			m.ActiveWorkspace = components.WorkspaceDashboard
			return m, nil
		case "2":
			m.ActiveWorkspace = components.WorkspaceSettings
			return m, nil
		case "3":
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
			m.ActiveWorkspace = (m.ActiveWorkspace + 1) % 3
			return m, nil

		case "tab":
			if m.ActiveWorkspace == components.WorkspaceDashboard {
				m.ActiveCategory = (m.ActiveCategory + 1) % len(views.CategoryFilterList)
				m.CatalogCursor = 0
			} else {
				m.ActiveWorkspace = (m.ActiveWorkspace + 1) % 3
				return m, nil
			}
			return m, nil

		case "shift+tab":
			if m.ActiveWorkspace == components.WorkspaceDashboard {
				if m.ActiveCategory > 0 {
					m.ActiveCategory--
				} else {
					m.ActiveCategory = len(views.CategoryFilterList) - 1
				}
				m.CatalogCursor = 0
			} else {
				if m.ActiveWorkspace > 0 {
					m.ActiveWorkspace--
				} else {
					m.ActiveWorkspace = components.WorkspaceGuide
				}
				return m, nil
			}
			return m, nil

		// Navigation inside active workspace
		case "up", "k":
			m.moveCursor(-1)
			return m, nil

		case "down", "j":
			m.moveCursor(1)
			return m, nil

		// Action Keys: ENTER or I anywhere in Launcher to Install!
		case "enter", "i", "I":
			if m.ActiveWorkspace == components.WorkspaceDashboard {
				filtered := m.getFilteredCatalog()
				if len(filtered) > 0 && m.CatalogCursor < len(filtered) {
					selected := filtered[m.CatalogCursor]
					action := "install"
					if selected.Installed {
						action = "uninstall"
					}
					m.AppDialog = views.AppDialogState{
						Active: true,
						Action: action,
						App:    selected,
					}
				}
			} else if m.ActiveWorkspace == components.WorkspaceSettings {
				m.applySettingSelection()
			}
			return m, nil

		case "x", "X":
			if m.ActiveWorkspace == components.WorkspaceDashboard {
				filtered := m.getFilteredCatalog()
				if len(filtered) > 0 && m.CatalogCursor < len(filtered) {
					selected := filtered[m.CatalogCursor]
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
				cfg, err := config.LoadBackupConfig()
				if err == nil {
					m.applyUserConfig(cfg)
					m.StatusMessage = "Pengaturan lama berhasil dipulihkan (Recovery)!"
				}
				m.Wizard.Active = false
				return m, nil
			}
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
		m.AppDialog.Logs = "Sedang memproses " + m.AppDialog.Action + " " + m.AppDialog.App.Name + "..."
		action := m.AppDialog.Action
		item := m.AppDialog.App
		return m, func() tea.Msg {
			out, err := system.ExecutePackageAction(action, item)
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

// getFilteredCatalog returns applications matching active category
func (m *Model) getFilteredCatalog() []system.AppCatalogItem {
	if m.ActiveCategory == 0 {
		return m.Catalog
	}
	selectedCatName := views.CategoryFilterList[m.ActiveCategory]
	var list []system.AppCatalogItem
	for _, app := range m.Catalog {
		if app.Category == selectedCatName {
			list = append(list, app)
		}
	}
	return list
}

func (m *Model) moveCursor(delta int) {
	switch m.ActiveWorkspace {
	case components.WorkspaceDashboard:
		filtered := m.getFilteredCatalog()
		if len(filtered) > 0 {
			m.CatalogCursor = clamp(m.CatalogCursor+delta, 0, len(filtered)-1)
		} else {
			m.CatalogCursor = 0
		}
	case components.WorkspaceSettings:
		m.SettingsCursor = clamp(m.SettingsCursor+delta, 0, 11) // 0-2 themes, 3-5 WM, 6-11 Wallpaper
	}
}

func (m *Model) applySettingSelection() {
	if m.SettingsCursor < 3 {
		m.ThemeMode = styles.ThemeMode(m.SettingsCursor)
		m.StatusMessage = "Tema berhasil diterapkan!"
	} else if m.SettingsCursor < 6 {
		m.WMMode = views.WMMode(m.SettingsCursor - 3)
		m.StatusMessage = "Model Window Manager berhasil diubah!"
	} else {
		wpIdx := m.SettingsCursor - 6
		switch wpIdx {
		case 0, 1, 2, 3:
			m.WallpaperType = wallpaper.WallpaperPreset
			m.WallpaperPreset = wallpaper.PresetScene(wpIdx)
			m.StatusMessage = "Wallpaper Seni Bumiayu diterapkan!"
		case 4:
			m.WallpaperType = wallpaper.WallpaperDynamicTime
			m.StatusMessage = "Wallpaper Dinamis (Waktu Server) diterapkan!"
		case 5:
			m.ShowImageInput = true
		}
	}
	m.saveCurrentConfig()
}

func (m *Model) saveCurrentConfig() {
	cfg := config.UserConfig{
		ConfigVersion:     "v0.9.0-beta.13",
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
	switch m.WallpaperType {
	case wallpaper.WallpaperCustom:
		cfg.WallpaperType = "custom"
	case wallpaper.WallpaperPreset:
		cfg.WallpaperType = "preset"
	default:
		cfg.WallpaperType = "dynamic"
	}
	cfg.WallpaperPreset = int(m.WallpaperPreset)
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

	// Image input prompt modal
	if m.ShowImageInput {
		modalWidth := width - 12
		if modalWidth > 68 {
			modalWidth = 68
		}
		content := lipgloss.JoinVertical(lipgloss.Left,
			st.Badge.Render("🖼️ GENERATOR FOTO KE PIXEL ART BUMIAYU"),
			"",
			"Masukkan path file gambar (PNG / JPG) di server Anda:",
			st.Card.Render(m.ImageInputPath+"_"),
			"",
			st.Subtitle.Render("Contoh: /tmp/wallpaper.jpg atau gambar.png"),
			"",
			lipgloss.JoinHorizontal(lipgloss.Center,
				st.ButtonActive.Render("[Enter] Generate Pixel Art"),
				"    ",
				st.ButtonNormal.Render("[Esc] Batal"),
			),
		)
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, st.Modal.Width(modalWidth).Render(content))
	}

	// If Wizard is active, render centered Setup Wizard
	if m.Wizard.Active {
		return views.RenderWizard(m.Wizard, st, width, height)
	}

	// WEB UI SCREEN MODE RENDERING
	if m.CurrentScreen == ScreenAlexisGrid {
		content := views.RenderAlexisGrid(m.AlexisCursor, m.AnimFrame, st, width, height)
		lines := strings.Split(content, "\n")
		if len(lines) > height {
			lines = lines[:height]
		}
		return strings.Join(lines, "\n")
	}

	if m.CurrentScreen == ScreenCommerceDashboard {
		content := views.RenderCommerceDashboard(
			m.Hardware,
			m.Catalog,
			m.CommerceSidebarCursor,
			m.CatalogCursor,
			m.ActiveCategory,
			m.AnimFrame,
			m.AppDialog,
			st,
			width,
			height,
		)
		lines := strings.Split(content, "\n")
		if len(lines) > height {
			lines = lines[:height]
		}
		return strings.Join(lines, "\n")
	}

	// 1. Top Desktop Bar
	topBar := components.RenderTopBar(m.Hardware, st, m.Language, width)

	// 2. Workspace Navigation Tabs
	workspaceTabs := components.RenderWorkspaceTabs(m.ActiveWorkspace, st, width)

	// 3. Main Workspace View
	var currentView string
	switch m.ActiveWorkspace {
	case components.WorkspaceDashboard:
		currentView = views.RenderDashboard(m.Hardware, m.Catalog, m.CatalogCursor, m.ActiveCategory, m.AnimFrame, m.AppDialog, m.WMMode, m.WallpaperType, m.WallpaperPreset, st, m.Language, width, height)
	case components.WorkspaceSettings:
		currentView = views.RenderSettings(m.ThemeMode, m.WMMode, m.WallpaperType, m.WallpaperPreset, m.SettingsCursor, m.ImageInputPath, st, width, height)
	case components.WorkspaceGuide:
		currentView = views.RenderGuide(st, width, height)
	}

	// 4. Status Bar
	statusBar := components.RenderStatusBar(st, width, m.StatusMessage)

	fullContent := lipgloss.JoinVertical(lipgloss.Left,
		topBar,
		workspaceTabs,
		currentView,
		statusBar,
	)

	// Clamp to terminal height so no scrolling occurs inside alternate screen
	lines := strings.Split(fullContent, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	return strings.Join(lines, "\n")
}
