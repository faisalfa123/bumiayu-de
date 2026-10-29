package system

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// HardwareInfo represents core machine and system telemetry
type HardwareInfo struct {
	Hostname     string
	OSName       string
	Kernel       string
	Arch         string
	CPUModel     string
	CPUCores     int
	CPUUsagePct  float64
	MemTotalMB   uint64
	MemUsedMB    uint64
	MemFreeMB    uint64
	MemUsagePct  float64
	SwapTotalMB  uint64
	SwapUsedMB   uint64
	SwapUsagePct float64
	TemperatureC float64
	IsOdroidC2   bool
	UptimeStr    string
	Load1        float64
	Load5        float64
	Load15       float64
}

// GetHardwareInfo fetches telemetry from /proc and /sys or mock if on non-Linux
func GetHardwareInfo() HardwareInfo {
	info := HardwareInfo{
		OSName:   runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUCores: runtime.NumCPU(),
	}

	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	} else {
		info.Hostname = "odroid-c2"
	}

	if runtime.GOOS == "linux" {
		readLinuxHardware(&info)
	} else {
		// Mock telemetry for development / cross-platform testing
		info.OSName = "Linux (Bumiayu Sim / " + runtime.GOOS + ")"
		info.Kernel = "6.1.0-odroid-arm64"
		info.CPUModel = "Amlogic S905 (Cortex-A53 quad-core)"
		info.IsOdroidC2 = true
		info.MemTotalMB = 2048
		info.MemUsedMB = 480
		info.MemFreeMB = 1568
		info.MemUsagePct = (float64(info.MemUsedMB) / float64(info.MemTotalMB)) * 100
		info.SwapTotalMB = 1024
		info.SwapUsedMB = 64
		info.SwapUsagePct = 6.25
		info.TemperatureC = 44.5
		info.CPUUsagePct = 12.4
		info.UptimeStr = "3d 14h 22m"
		info.Load1 = 0.35
		info.Load5 = 0.28
		info.Load15 = 0.19
	}

	return info
}

func readLinuxHardware(info *HardwareInfo) {
	// 1. Read Uptime
	if data, err := os.ReadFile("/proc/uptime"); err == nil {
		parts := strings.Fields(string(data))
		if len(parts) > 0 {
			if upSec, err := strconv.ParseFloat(parts[0], 64); err == nil {
				d := time.Duration(upSec) * time.Second
				days := int(d.Hours()) / 24
				hours := int(d.Hours()) % 24
				mins := int(d.Minutes()) % 60
				if days > 0 {
					info.UptimeStr = fmt.Sprintf("%dd %dh %dm", days, hours, mins)
				} else {
					info.UptimeStr = fmt.Sprintf("%dh %dm", hours, mins)
				}
			}
		}
	}

	// 2. Read Load Average
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			info.Load1, _ = strconv.ParseFloat(fields[0], 64)
			info.Load5, _ = strconv.ParseFloat(fields[1], 64)
			info.Load15, _ = strconv.ParseFloat(fields[2], 64)
		}
	}

	// 3. Read Memory Info from /proc/meminfo
	if f, err := os.Open("/proc/meminfo"); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		var memTotal, memFree, memAvailable, buffers, cached uint64
		var swapTotal, swapFree uint64

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			val, _ := strconv.ParseUint(fields[1], 10, 64)
			switch fields[0] {
			case "MemTotal:":
				memTotal = val / 1024
			case "MemFree:":
				memFree = val / 1024
			case "MemAvailable:":
				memAvailable = val / 1024
			case "Buffers:":
				buffers = val / 1024
			case "Cached:":
				cached = val / 1024
			case "SwapTotal:":
				swapTotal = val / 1024
			case "SwapFree:":
				swapFree = val / 1024
			}
		}

		info.MemTotalMB = memTotal
		if memAvailable > 0 {
			info.MemUsedMB = memTotal - memAvailable
			info.MemFreeMB = memAvailable
		} else {
			info.MemUsedMB = memTotal - (memFree + buffers + cached)
			info.MemFreeMB = memFree + buffers + cached
		}

		if info.MemTotalMB > 0 {
			info.MemUsagePct = (float64(info.MemUsedMB) / float64(info.MemTotalMB)) * 100
		}

		info.SwapTotalMB = swapTotal
		info.SwapUsedMB = swapTotal - swapFree
		if swapTotal > 0 {
			info.SwapUsagePct = (float64(info.SwapUsedMB) / float64(swapTotal)) * 100
		}
	}

	// 4. Read CPU and Model from /proc/cpuinfo
	if f, err := os.Open("/proc/cpuinfo"); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "model name") || strings.HasPrefix(line, "Hardware") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					info.CPUModel = strings.TrimSpace(parts[1])
					if strings.Contains(strings.ToLower(info.CPUModel), "odroid") || strings.Contains(strings.ToLower(info.CPUModel), "amlogic") {
						info.IsOdroidC2 = true
					}
					break
				}
			}
		}
	}

	// 5. Read Thermal Sensor (ODROID-C2 specific or generic Linux thermal zone)
	thermalPaths := []string{
		"/sys/devices/virtual/thermal/thermal_zone0/temp",
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/hwmon/hwmon0/temp1_input",
	}

	for _, p := range thermalPaths {
		if data, err := os.ReadFile(p); err == nil {
			raw := strings.TrimSpace(string(data))
			if tempVal, err := strconv.ParseFloat(raw, 64); err == nil {
				if tempVal > 1000 {
					info.TemperatureC = tempVal / 1000.0
				} else {
					info.TemperatureC = tempVal
				}
				break
			}
		}
	}

	// 6. Read Kernel version
	if data, err := os.ReadFile("/proc/version"); err == nil {
		parts := strings.Fields(string(data))
		if len(parts) >= 3 {
			info.Kernel = parts[2]
		}
	}
}
