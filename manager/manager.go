package manager

import (
	"dst-manager/config"
	"fmt"
)

// Manager handles the DST server operations
// 管理器结构体
type Manager struct {
	Config *config.Config
}

// NewManager creates a new Manager instance
// 创建新的管理器实例
func NewManager() *Manager {
	return &Manager{
		Config: config.NewConfig(),
	}
}

// Log prints a formatted message with a cute prefix
// 打印带萌系前缀的日志
func (m *Manager) Log(format string, a ...interface{}) {
	fmt.Printf("[小花酱] "+format+"\n", a...)
}
