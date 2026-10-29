package views

import (
	"fmt"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderSecurity renders firewall rules and open ports audit
func RenderSecurity(sec system.SecurityState, lang i18n.Language, width, height int) string {
	t := i18n.Get(lang)
	contentWidth := width - 6
	if contentWidth < 40 {
		contentWidth = 40
	}

	header := styles.TitleStyle.Render("🛡️ Pusat Aturan Keamanan & Firewall Sistem (UFW)")

	// 1. Firewall status & rules panel
	ufwStatus := styles.PillDanger.Render("🔴 FIREWALL NONAKTIF (TIDAK AMAN)")
	if sec.UfwActive {
		ufwStatus = styles.BadgeStyle.Render("🟢 FIREWALL AKTIF (TERLINDUNGI)")
	}

	var rulesText string
	if len(sec.UfwRules) == 0 {
		rulesText = "Tidak ada aturan firewall khusus yang diterapkan."
	} else {
		for _, r := range sec.UfwRules {
			rulesText += fmt.Sprintf("[%d] %-12s %-12s %s\n", r.Index, r.Target, r.Action, r.From)
		}
	}

	firewallContent := fmt.Sprintf(
		"Status UFW: %s\n\nAturan Terdaftar:\n%s\n\n⚡ Aksi: Tekan [U] untuk Toggle Firewall ON / OFF",
		ufwStatus,
		rulesText,
	)

	firewallPanel := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			styles.TitleStyle.Render("Dinding Api (UFW)"),
			"",
			firewallContent,
		),
	)

	// 2. Open Ports Audit Panel
	var portsText string
	if len(sec.OpenPorts) == 0 {
		portsText = "Tidak ada port publik yang terdeteksi sedang mendengarkan."
	} else {
		for _, p := range sec.OpenPorts {
			riskBadge := "🟢 Lokal"
			if p.IsPublic {
				riskBadge = "⚠️ Publik"
			}
			portsText += fmt.Sprintf("• Port %-5d (%-4s) %s -> %s\n  %s\n", p.Port, p.Protocol, riskBadge, p.ProcessName, p.RiskNote)
		}
	}

	portsPanel := styles.CardStyle.Width(contentWidth/2 - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			styles.TitleStyle.Render("Audit Port & Koneksi Terbuka"),
			"",
			portsText,
		),
	)

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, firewallPanel, "  ", portsPanel)

	// Senior Insight Box
	insightHeader := styles.InsightTitleStyle.Render(t.SeniorInsightTitle)
	insightBody := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0f2fe")).Render(t.InsightSecurity)
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
