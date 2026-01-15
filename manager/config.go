package manager

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the server configuration
// 配置文件结构体
type Config struct {
	HomeDir       string // User home directory
	SteamCMDDir   string // SteamCMD installation directory
	DSTInstallDir string // DST server installation directory
	ClusterDir    string // Cluster save directory
	BackupDir     string // Backup directory
}

// NewConfig creates a new configuration with default paths
// 创建新的默认配置
func NewConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return &Config{
		HomeDir:       home,
		SteamCMDDir:   filepath.Join(home, "steamcmd"),
		DSTInstallDir: filepath.Join(home, "dst"),
		ClusterDir:    filepath.Join(home, ".klei", "DoNotStarveTogether"),
		BackupDir:     filepath.Join(home, "dst-backups"),
	}
}

// EnsureDirs creates necessary directories
// 确保必要的目录存在
func (c *Config) EnsureDirs() error {
	dirs := []string{c.SteamCMDDir, c.DSTInstallDir, c.BackupDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	return nil
}
