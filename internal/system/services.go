package system

import (
	"bufio"
	"os/exec"
	"runtime"
	"strings"
)

// ServiceItem represents a systemd service unit
type ServiceItem struct {
	Name        string
	Description string
	ActiveState string // "active", "inactive", "failed"
	SubState    string // "running", "dead", "exited"
	Enabled     bool
	MemoryMB    float64
}

// GetServices fetches systemd services or mock services if non-linux
func GetServices() []ServiceItem {
	if runtime.GOOS != "linux" {
		return []ServiceItem{
			{Name: "ssh.service", Description: "OpenBSD Secure Shell server", ActiveState: "active", SubState: "running", Enabled: true, MemoryMB: 4.8},
			{Name: "nginx.service", Description: "A high performance web server", ActiveState: "active", SubState: "running", Enabled: true, MemoryMB: 6.2},
			{Name: "ufw.service", Description: "Uncomplicated firewall", ActiveState: "active", SubState: "exited", Enabled: true, MemoryMB: 1.1},
			{Name: "fail2ban.service", Description: "Fail2ban Service", ActiveState: "active", SubState: "running", Enabled: true, MemoryMB: 14.5},
			{Name: "cron.service", Description: "Regular background program processing daemon", ActiveState: "active", SubState: "running", Enabled: true, MemoryMB: 2.3},
			{Name: "docker.service", Description: "Docker Application Container Engine", ActiveState: "inactive", SubState: "dead", Enabled: false, MemoryMB: 0.0},
			{Name: "samba.service", Description: "Samba SMB Daemon", ActiveState: "inactive", SubState: "dead", Enabled: false, MemoryMB: 0.0},
			{Name: "wireguard.service", Description: "WireGuard VPN Tunnel", ActiveState: "inactive", SubState: "dead", Enabled: false, MemoryMB: 0.0},
			{Name: "systemd-timesyncd.service", Description: "Network Time Synchronization", ActiveState: "active", SubState: "running", Enabled: true, MemoryMB: 2.1},
			{Name: "cups.service", Description: "CUPS Scheduler (Print Server)", ActiveState: "failed", SubState: "failed", Enabled: false, MemoryMB: 0.0},
		}
	}

	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var services []ServiceItem
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			name := fields[0]
			activeState := fields[2]
			subState := fields[3]
			desc := ""
			if len(fields) > 4 {
				desc = strings.Join(fields[4:], " ")
			}

			// Filter to standard server services of interest
			services = append(services, ServiceItem{
				Name:        name,
				Description: desc,
				ActiveState: activeState,
				SubState:    subState,
				Enabled:     checkServiceEnabled(name),
			})
		}
	}

	return services
}

func checkServiceEnabled(name string) bool {
	cmd := exec.Command("systemctl", "is-enabled", name)
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out)) == "enabled"
}

// ControlService executes start, stop, restart, enable, disable
func ControlService(action, name string) (string, error) {
	if runtime.GOOS != "linux" {
		return "Simulated: " + action + " on " + name + " succeeded", nil
	}

	cmd := exec.Command("systemctl", action, name)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// GetServiceLogs fetches the last N lines of journalctl for a unit
func GetServiceLogs(name string, lines int) string {
	if runtime.GOOS != "linux" {
		return "[Simulated Logs for " + name + "]\nOct 29 20:15:01 odroid systemd[1]: Started " + name + ".\nOct 29 20:15:02 odroid " + name + "[1042]: Ready to accept client connections.\nOct 29 20:15:10 odroid " + name + "[1042]: Configuration loaded successfully. No errors detected."
	}

	cmd := exec.Command("journalctl", "-u", name, "-n", "30", "--no-pager")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Gagal memuat log: " + err.Error()
	}
	return string(out)
}
