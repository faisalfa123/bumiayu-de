package components

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// Workspace represents active desktop workspace
type Workspace int

const (
	WorkspaceDashboard Workspace = iota
	WorkspaceApps
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
		{WorkspaceDashboard, "Desktop & Status", "1"},
		{WorkspaceApps, "Toko Aplikasi", "2"},
		{WorkspaceSettings, "Pengaturan Tampilan", "3"},
		{WorkspaceGuide, "Petunjuk Penggunaan", "4"},
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
	border := strings.Repeat("─", width)
	return lipgloss.JoinVertical(lipgloss.Left, row, st.Subtitle.Render(border))
}
