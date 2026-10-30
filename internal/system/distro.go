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
	DistroArch    DistroType = "Arch Linux"
	DistroUbuntu  DistroType = "Ubuntu"
	DistroDebian  DistroType = "Debian"
	DistroFedora  DistroType = "Fedora"
	DistroGeneric DistroType = "Linux"
)

// DetectDistro determines the underlying Linux distribution
func DetectDistro() DistroType {
	if runtime.GOOS != "linux" {
		return DistroGeneric
	}

	// 1. Check DietPi (Filesystem markers or CLI)
	if _, err := os.Stat("/boot/dietpi"); err == nil {
		return DistroDietPi
	}
	if _, err := os.Stat("/boot/dietpi.txt"); err == nil {
		return DistroDietPi
	}
	if _, err := exec.LookPath("dietpi-software"); err == nil {
		return DistroDietPi
	}

	// 2. Check Arch Linux
	if _, err := os.Stat("/etc/arch-release"); err == nil {
		return DistroArch
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		return DistroArch
	}

	// 3. Check Fedora
	if _, err := os.Stat("/etc/fedora-release"); err == nil {
		return DistroFedora
	}

	// 4. Read /etc/os-release
	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		content := strings.ToLower(string(data))
		if strings.Contains(content, "dietpi") {
			return DistroDietPi
		}
		if strings.Contains(content, "arch") || strings.Contains(content, "manjaro") || strings.Contains(content, "endeavouros") {
			return DistroArch
		}
		if strings.Contains(content, "ubuntu") {
			return DistroUbuntu
		}
		if strings.Contains(content, "debian") || strings.Contains(content, "raspbian") {
			return DistroDebian
		}
		if strings.Contains(content, "fedora") || strings.Contains(content, "rhel") || strings.Contains(content, "rocky") {
			return DistroFedora
		}
	}

	return DistroGeneric
}
