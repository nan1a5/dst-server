package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Config holds the server configuration
// 配置文件结构体
type Config struct {
	HomeDir       string
	SteamCMDDir   string
	DSTInstallDir string
	ClusterDir    string
	BackupDir     string
}

var (
	instance *Config
	once     sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}

		instance = &Config{
			HomeDir:       home,
			SteamCMDDir:   filepath.Join(home, "steamcmd"),
			DSTInstallDir: filepath.Join(home, "dst"),
			ClusterDir:    filepath.Join(home, ".klei", "DoNotStarveTogether"),
			BackupDir:     filepath.Join(home, "dst-backups"),
		}
	})
	return instance
}

func (c *Config) EnsureDirs() error {
	dirs := []string{c.SteamCMDDir, c.DSTInstallDir, c.BackupDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	return nil
}
