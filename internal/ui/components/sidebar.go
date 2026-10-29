package components

import (
	"fmt"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// Workspace represents active desktop workspace
type Workspace int

const (
	WorkspaceDashboard Workspace = iota
	WorkspaceApps
	WorkspaceServices
	WorkspaceSecurity
	WorkspaceUsers
	WorkspaceFiles
	WorkspaceTasks
)

// RenderWorkspaceTabs renders the horizontal navigation tabs
func RenderWorkspaceTabs(active Workspace, lang i18n.Language, width int) string {
	t := i18n.Get(lang)

	tabs := []struct {
		ws   Workspace
		name string
		icon string
	}{
		{WorkspaceDashboard, t.TabDashboard, "📊"},
		{WorkspaceApps, t.TabApps, "📦"},
		{WorkspaceServices, t.TabServices, "⚙️"},
		{WorkspaceSecurity, t.TabSecurity, "🛡️"},
		{WorkspaceUsers, t.TabUsers, "👥"},
		{WorkspaceFiles, t.TabFiles, "📂"},
		{WorkspaceTasks, t.TabTasks, "⚡"},
	}

	var renderedTabs []string
	for _, tab := range tabs {
		label := fmt.Sprintf("%s %s", tab.icon, tab.name)
		if tab.ws == active {
			renderedTabs = append(renderedTabs, styles.ActiveTabStyle.Render(label))
		} else {
			renderedTabs = append(renderedTabs, styles.InactiveTabStyle.Render(label))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	border := strings.Repeat("─", width)
	return lipgloss.JoinVertical(lipgloss.Left, row, styles.SubtitleStyle.Render(border))
}
