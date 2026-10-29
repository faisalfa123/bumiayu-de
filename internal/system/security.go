package system

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// PortAuditItem represents a listening network port
type PortAuditItem struct {
	Port        int
	Protocol    string
	BindIP      string
	ProcessName string
	IsPublic    bool
	RiskNote    string
}

// UfwRule represents a configured firewall rule
type UfwRule struct {
	Index    int
	Target   string // Port or Service
	Action   string // ALLOW, DENY
	From     string // Anywhere, specific IP
	Protocol string // tcp, udp, any
}

// SecurityState represents the full security posture
type SecurityState struct {
	UfwActive    bool
	UfwRules     []UfwRule
	OpenPorts    []PortAuditItem
	SshRootLogin bool // true if PermitRootLogin is yes
	Fail2banOn   bool
}

// GetSecurityState retrieves firewall status, open ports, and SSH configuration
func GetSecurityState() SecurityState {
	state := SecurityState{}

	if runtime.GOOS != "linux" {
		// Mock realistic security state for development preview
		state.UfwActive = true
		state.Fail2banOn = true
		state.SshRootLogin = false
		state.UfwRules = []UfwRule{
			{Index: 1, Target: "22/tcp", Action: "ALLOW IN", From: "Anywhere", Protocol: "tcp"},
			{Index: 2, Target: "80/tcp", Action: "ALLOW IN", From: "Anywhere", Protocol: "tcp"},
			{Index: 3, Target: "443/tcp", Action: "ALLOW IN", From: "Anywhere", Protocol: "tcp"},
		}
		state.OpenPorts = []PortAuditItem{
			{Port: 22, Protocol: "tcp", BindIP: "0.0.0.0", ProcessName: "sshd", IsPublic: true, RiskNote: "Akses Remote Server. Pastikan menggunakan autentikasi kunci SSH."},
			{Port: 80, Protocol: "tcp", BindIP: "0.0.0.0", ProcessName: "nginx", IsPublic: true, RiskNote: "HTTP Web Standar (Traffic tidak terenkripsi, gunakan redirect ke 443)."},
			{Port: 443, Protocol: "tcp", BindIP: "0.0.0.0", ProcessName: "nginx", IsPublic: true, RiskNote: "HTTPS Terenkripsi Aman."},
			{Port: 53, Protocol: "udp", BindIP: "127.0.0.1", ProcessName: "systemd-resolved", IsPublic: false, RiskNote: "DNS Resolver Lokal (Aman, hanya terikat pada localhost)."},
		}
		return state
	}

	// 1. Read UFW status
	cmd := exec.Command("ufw", "status", "numbered")
	out, err := cmd.Output()
	if err == nil {
		str := string(out)
		if strings.Contains(str, "Status: active") {
			state.UfwActive = true
		}
		scanner := bufio.NewScanner(strings.NewReader(str))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "[") {
				// Format: [ 1] 22/tcp  ALLOW IN  Anywhere
				parts := strings.Fields(line)
				if len(parts) >= 4 {
					idxStr := strings.Trim(parts[0], "[]")
					var idx int
					fmt.Sscanf(idxStr, "%d", &idx)
					state.UfwRules = append(state.UfwRules, UfwRule{
						Index:  idx,
						Target: parts[1],
						Action: parts[2] + " " + parts[3],
						From:   strings.Join(parts[4:], " "),
					})
				}
			}
		}
	}

	// 2. Read Open Ports via ss
	ssCmd := exec.Command("ss", "-tulnp")
	ssOut, err := ssCmd.Output()
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(ssOut)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Netid") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				proto := fields[0]
				localAddr := fields[4]
				proc := ""
				if len(fields) >= 7 {
					proc = fields[6]
				}

				lastColon := strings.LastIndex(localAddr, ":")
				if lastColon != -1 {
					portStr := localAddr[lastColon+1:]
					ipStr := localAddr[:lastColon]
					var port int
					fmt.Sscanf(portStr, "%d", &port)
					if port > 0 {
						isPublic := !strings.HasPrefix(ipStr, "127.") && ipStr != "::1" && ipStr != "[::1]"
						risk := "Port internal"
						if isPublic {
							risk = "Terbuka ke publik / LAN"
						}
						state.OpenPorts = append(state.OpenPorts, PortAuditItem{
							Port:        port,
							Protocol:    proto,
							BindIP:      ipStr,
							ProcessName: proc,
							IsPublic:    isPublic,
							RiskNote:    risk,
						})
					}
				}
			}
		}
	}

	return state
}

// ToggleUFW switches UFW firewall on or off
func ToggleUFW(enable bool) (string, error) {
	if runtime.GOOS != "linux" {
		return "Simulated: UFW toggle", nil
	}
	action := "disable"
	if enable {
		action = "enable"
	}
	cmd := exec.Command("ufw", "--force", action)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
