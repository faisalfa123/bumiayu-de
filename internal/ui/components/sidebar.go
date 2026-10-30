package components

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// Workspace represents active desktop workspace
type Workspace int

const (
	WorkspaceDashboard Workspace = iota
	WorkspaceSettings
	WorkspaceGuide
)

// RenderWorkspaceTabs renders horizontal navigation tabs with active theme
func RenderWorkspaceTabs(active Workspace, st styles.ThemeStyles, width int) string {
	tabs := []struct {
		ws   Workspace
		name string
		num  string
	}{
		{WorkspaceDashboard, "Desktop & Aplikasi (All-in-One)", "1"},
		{WorkspaceSettings, "Pengaturan Tampilan & WM", "2"},
		{WorkspaceGuide, "Petunjuk Penggunaan", "3"},
	}

	var renderedTabs []string
	for _, tab := range tabs {
		label := fmt.Sprintf("[%s] %s", tab.num, tab.name)
		if tab.ws == active {
			renderedTabs = append(renderedTabs, st.ActiveTab.Render(label))
		} else {
			renderedTabs = append(renderedTabs, st.InactiveTab.Render(label))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	return row
}
