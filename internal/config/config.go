package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type UserConfig struct {
	ConfigVersion      string `json:"config_version"`
	Theme              string `json:"theme"`    // "light", "warm", "dark"
	WMModel            string `json:"wm_model"` // "floating", "tiling", "cardgrid"
	Language           string `json:"language"` // "id", "en"
	AutoStartSSH       bool   `json:"autostart_ssh"`
	FirstRunCompleted  bool   `json:"first_run_completed"`
	LastBackupTime     string `json:"last_backup_time"`
}

// GetConfigDir returns path to ~/.config/bumiayu
func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "bumiayu")
}

// ConfigPath returns main config file path
func ConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.json")
}

// BackupConfigPath returns backup file path
func BackupConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.backup.json")
}

// HasExistingSetup checks if previous setup or backup exists
func HasExistingSetup() bool {
	if _, err := os.Stat(BackupConfigPath()); err == nil {
		return true
	}
	if _, err := os.Stat(ConfigPath()); err == nil {
		return true
	}
	return false
}

// LoadConfig loads active config or default if missing
func LoadConfig() (UserConfig, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

// LoadBackupConfig loads backup config for recovery
func LoadBackupConfig() (UserConfig, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(BackupConfigPath())
	if err != nil {
		// Fallback to main config
		data, err = os.ReadFile(ConfigPath())
		if err != nil {
			return cfg, err
		}
	}

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

// SaveConfig saves config and creates backup
func SaveConfig(cfg UserConfig) error {
	dir := GetConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	cfg.LastBackupTime = time.Now().Format("2006-01-02 15:04:05")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Write main config
	if err := os.WriteFile(ConfigPath(), data, 0644); err != nil {
		return err
	}

	// Create backup mirror
	return os.WriteFile(BackupConfigPath(), data, 0644)
}

// DefaultConfig provides clean default setup
func DefaultConfig() UserConfig {
	return UserConfig{
		ConfigVersion:     "v0.9.0-beta.3",
		Theme:             "light",
		WMModel:           "floating",
		Language:          "id",
		AutoStartSSH:      false,
		FirstRunCompleted: false,
	}
}
