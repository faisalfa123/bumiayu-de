package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderUsers renders user accounts and sudoer role permissions
func RenderUsers(users []system.UserAccount, selectedIndex int, lang i18n.Language, width, height int) string {
	t := i18n.Get(lang)
	contentWidth := width - 6
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := styles.TitleStyle.Render("👥 Manajemen Pengguna, Role & Hak Akses Sudo")

	var listItems []string
	for i, u := range users {
		roleBadge := styles.PillNormal.Render("Standard User")
		if u.IsSudoer {
			roleBadge = styles.BadgeStyle.Render("👑 SUDO / ADMIN")
		}

		keyBadge := "🔑 No SSH Key"
		if u.HasSSHKey {
			keyBadge = fmt.Sprintf("🔑 %d Key(s)", u.SSHKeyCount)
		}

		line := fmt.Sprintf("%-16s  %s  %s", u.Username, roleBadge, keyBadge)
		if i == selectedIndex {
			listItems = append(listItems, styles.RowSelectedStyle.Width(contentWidth/2-4).Render("▶ "+line))
		} else {
			listItems = append(listItems, styles.RowNormalStyle.Width(contentWidth/2-4).Render("  "+line))
		}
	}

	listPanel := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left, listItems...),
	)

	// Detail & Role Modification Panel
	var detailPanel string
	if selectedIndex >= 0 && selectedIndex < len(users) {
		u := users[selectedIndex]
		roleDesc := "Pengguna biasa, tidak memiliki hak mengubah konfigurasi inti server."
		if u.IsSudoer {
			roleDesc = "Administrator penuh (Sudoer). Dapat menjalankan perintah 'sudo' dan mengubah sistem."
		}

		detailContent := fmt.Sprintf(
			"👤 Username: %s (UID: %d, GID: %d)\n📁 Home Directory: %s\n💻 Shell: %s\n🔑 Kunci SSH: %d Kunci Terdaftar\n\n🛡️ Tingkat Hak Akses:\n%s\n\n⚡ Aksi Kontrol Role:\n  [G] Berikan Hak Akses Admin (Grant Sudo)\n  [V] Cabut Hak Akses Admin (Revoke Sudo)",
			u.Username, u.UID, u.GID,
			u.HomeDir,
			u.Shell,
			u.SSHKeyCount,
			roleDesc,
		)
		detailPanel = styles.CardActiveStyle.Width(contentWidth/2 - 2).Render(detailContent)
	}

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", detailPanel)

	// Senior Insight Box
	insightHeader := styles.InsightTitleStyle.Render(t.SeniorInsightTitle)
	insightBody := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0f2fe")).Render(t.InsightUsers)
	insightBox := styles.InsightBoxStyle.Width(contentWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left, insightHeader, insightBody),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		mainRow,
		"",
		insightBox,
	)
}
