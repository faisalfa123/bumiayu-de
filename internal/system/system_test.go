package system_test

import (
	"testing"

	"github.com/bumiayu/bumiayu-de/internal/i18n"
	"github.com/bumiayu/bumiayu-de/internal/system"
)

func TestHardwareInfo(t *testing.T) {
	info := system.GetHardwareInfo()

	if info.Hostname == "" {
		t.Errorf("Hostname should not be empty")
	}

	if info.CPUCores <= 0 {
		t.Errorf("CPUCores should be positive, got %d", info.CPUCores)
	}

	if info.MemTotalMB == 0 {
		t.Errorf("MemTotalMB should be greater than 0")
	}

	if info.MemUsagePct < 0 || info.MemUsagePct > 100 {
		t.Errorf("MemUsagePct out of bounds: %f", info.MemUsagePct)
	}
}

func TestAppCatalog(t *testing.T) {
	catalog := system.GetAppCatalog()

	if len(catalog) < 5 {
		t.Errorf("Expected at least 5 curated server apps, got %d", len(catalog))
	}

	foundDocker := false
	for _, item := range catalog {
		if item.ID == "docker" {
			foundDocker = true
			if item.Category == "" || item.Description == "" {
				t.Errorf("Docker catalog entry missing metadata")
			}
		}
	}

	if !foundDocker {
		t.Errorf("Expected docker in catalog")
	}
}

func TestMemoryBar(t *testing.T) {
	bar0 := system.GetMemoryBar(0, 10)
	if bar0 != "[░░░░░░░░░░] 0.0%" {
		t.Errorf("Unexpected bar0: %s", bar0)
	}

	bar50 := system.GetMemoryBar(50, 10)
	if bar50 != "[█████░░░░░] 50.0%" {
		t.Errorf("Unexpected bar50: %s", bar50)
	}

	bar100 := system.GetMemoryBar(100, 10)
	if bar100 != "[██████████] 100.0%" {
		t.Errorf("Unexpected bar100: %s", bar100)
	}
}

func TestI18nTranslations(t *testing.T) {
	id := i18n.Get(i18n.ID)
	if id.AppTitle != "BUMIAYU TUI" {
		t.Errorf("Expected BUMIAYU TUI, got %s", id.AppTitle)
	}
	if id.TabDashboard == "" || id.InsightDashboard == "" {
		t.Errorf("Indonesian translation strings missing")
	}

	en := i18n.Get(i18n.EN)
	if en.AppTitle != "BUMIAYU TUI" {
		t.Errorf("Expected BUMIAYU TUI in EN, got %s", en.AppTitle)
	}
	if en.TabDashboard == "" || en.InsightDashboard == "" {
		t.Errorf("English translation strings missing")
	}
}
