package system

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// ProcessItem represents a running process
type ProcessItem struct {
	PID         int
	User        string
	Name        string
	CPUUsagePct float64
	MemMB       float64
	MemUsagePct float64
	OOMScore    int // 0 - 1000 from /proc/[pid]/oom_score
	State       string
}

// GetProcessList retrieves running processes sorted by memory usage
func GetProcessList() []ProcessItem {
	if runtime.GOOS != "linux" {
		return []ProcessItem{
			{PID: 1042, User: "www-data", Name: "nginx: worker", CPUUsagePct: 0.8, MemMB: 18.5, MemUsagePct: 0.9, OOMScore: 50, State: "S"},
			{PID: 1205, User: "root", Name: "dockerd", CPUUsagePct: 2.1, MemMB: 84.0, MemUsagePct: 4.1, OOMScore: 200, State: "S"},
			{PID: 1350, User: "root", Name: "containerd", CPUUsagePct: 1.5, MemMB: 42.0, MemUsagePct: 2.0, OOMScore: 120, State: "S"},
			{PID: 980, User: "root", Name: "fail2ban-server", CPUUsagePct: 0.2, MemMB: 14.5, MemUsagePct: 0.7, OOMScore: 30, State: "S"},
			{PID: 850, User: "root", Name: "sshd: /usr/sbin/sshd", CPUUsagePct: 0.1, MemMB: 6.8, MemUsagePct: 0.3, OOMScore: 10, State: "S"},
			{PID: 1420, User: "odroid", Name: "bumiayu-de", CPUUsagePct: 1.2, MemMB: 12.4, MemUsagePct: 0.6, OOMScore: 5, State: "R"},
			{PID: 512, User: "systemd-resolve", Name: "systemd-resolved", CPUUsagePct: 0.1, MemMB: 8.2, MemUsagePct: 0.4, OOMScore: 0, State: "S"},
		}
	}

	var list []ProcessItem
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return list
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue // Not a PID folder
		}

		item := readProcStat(pid)
		if item.PID > 0 {
			list = append(list, item)
		}
	}

	// Sort by MemMB descending
	sort.Slice(list, func(i, j int) bool {
		return list[i].MemMB > list[j].MemMB
	})

	if len(list) > 50 {
		list = list[:50] // Limit to top 50 to avoid any overhead
	}

	return list
}

func readProcStat(pid int) ProcessItem {
	item := ProcessItem{PID: pid}

	// 1. Read /proc/[pid]/comm
	commPath := filepath.Join("/proc", strconv.Itoa(pid), "comm")
	if data, err := os.ReadFile(commPath); err == nil {
		item.Name = strings.TrimSpace(string(data))
	} else {
		return item
	}

	// 2. Read /proc/[pid]/oom_score
	oomPath := filepath.Join("/proc", strconv.Itoa(pid), "oom_score")
	if data, err := os.ReadFile(oomPath); err == nil {
		score, _ := strconv.Atoi(strings.TrimSpace(string(data)))
		item.OOMScore = score
	}

	// 3. Read /proc/[pid]/status for VmRSS and State
	statusPath := filepath.Join("/proc", strconv.Itoa(pid), "status")
	if f, err := os.Open(statusPath); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "State:") {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					item.State = parts[1]
				}
			} else if strings.HasPrefix(line, "VmRSS:") {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					kb, _ := strconv.ParseFloat(parts[1], 64)
					item.MemMB = kb / 1024.0
				}
			}
		}
	}

	return item
}

// KillProcess terminates a process by PID
func KillProcess(pid int, force bool) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

// GetMemoryBar renders a visual ASCII gauge bar
func GetMemoryBar(pct float64, width int) string {
	fillWidth := int((pct / 100.0) * float64(width))
	if fillWidth > width {
		fillWidth = width
	}
	if fillWidth < 0 {
		fillWidth = 0
	}
	emptyWidth := width - fillWidth

	return fmt.Sprintf("[%s%s] %.1f%%", strings.Repeat("█", fillWidth), strings.Repeat("░", emptyWidth), pct)
}
