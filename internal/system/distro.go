package system

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type DistroType string

const (
	DistroDietPi  DistroType = "DietPi"
	DistroUbuntu  DistroType = "Ubuntu"
	DistroDebian  DistroType = "Debian"
	DistroGeneric DistroType = "Linux"
)

// DetectDistro determines the underlying Linux distribution
func DetectDistro() DistroType {
	if runtime.GOOS != "linux" {
		return DistroGeneric
	}

	// 1. Check DietPi
	if _, err := os.Stat("/boot/dietpi"); err == nil {
		return DistroDietPi
	}
	if _, err := exec.LookPath("dietpi-software"); err == nil {
		return DistroDietPi
	}

	// 2. Read /etc/os-release
	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		content := strings.ToLower(string(data))
		if strings.Contains(content, "dietpi") {
			return DistroDietPi
		}
		if strings.Contains(content, "ubuntu") {
			return DistroUbuntu
		}
		if strings.Contains(content, "debian") {
			return DistroDebian
		}
	}

	return DistroGeneric
}
