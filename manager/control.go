package manager

import (
	"dst-manager/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StartServer starts the DST server in a screen session
// 启动服务器
func (m *Manager) StartServer() error {
	m.Log("正在启动服务器，请稍候喵...")

	// Check if already running
	if m.IsRunning() {
		m.Log("服务器已经在运行了喵！不要重复启动哦~")
		return fmt.Errorf("服务器已经在运行了喵！不要重复启动哦~")
	}

	// Select Cluster
	cluster := m.SelectCluster("请选择要启动的存档喵:")
	if cluster == "" {
		return fmt.Errorf("请选择要启动的存档喵~")
	}
	m.Log("即将启动存档: %s", cluster)

	// Executable path
	// 64-bit executable is standard now
	binPath := filepath.Join(m.Config.DSTInstallDir, "bin64", "dontstarve_dedicated_server_nullrenderer_x64")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		// Fallback to 32-bit
		binPath = filepath.Join(m.Config.DSTInstallDir, "bin", "dontstarve_dedicated_server_nullrenderer")
	}

	// We launch Master and Caves separately
	// 启动 Master
	m.startShard(binPath, cluster, "Master")

	// 启动 Caves
	m.startShard(binPath, cluster, "Caves")

	m.Log("服务器启动指令已发送！可以用 screen -ls 查看后台进程喵~")
	return nil
}

func (m *Manager) startShard(binPath, clusterName, shardName string) {
	screenName := fmt.Sprintf("dst_%s", shardName)
	cmd := fmt.Sprintf("cd %s && %s -console -cluster %s -shard %s",
		filepath.Dir(binPath), binPath, clusterName, shardName)

	// Use screen to run in background
	// screen -dmS <name> <command>
	err := utils.RunCommand("screen", "-dmS", screenName, "bash", "-c", cmd)
	if err != nil {
		m.Log("启动 %s 失败了喵: %v", shardName, err)
	} else {
		m.Log("%s 世界启动成功！", shardName)
	}
}

// StopServer stops the DST server
// 停止服务器
func (m *Manager) StopServer() {
	m.Log("正在停止服务器，会保存存档喵...")

	m.stopShard("Master")
	m.stopShard("Caves")

	m.Log("服务器已停止，休息一下吧主人~")
}

func (m *Manager) StopMaster() {
	m.Log("正在停止地面服务器，会保存存档喵...")

	m.stopShard("Master")

	m.Log("地面服务器已停止，休息一下吧主人~")
}

func (m *Manager) StopCaves() {
	m.Log("正在停止洞穴服务器，会保存存档喵...")

	m.stopShard("Caves")

	m.Log("洞穴服务器已停止，休息一下吧主人~")
}

func (m *Manager) stopShard(shardName string) {
	screenName := fmt.Sprintf("dst_%s", shardName)
	// Send c_shutdown(true) to save and exit
	// screen -S <name> -p 0 -X stuff "c_shutdown(true)\n"
	cmd := "c_shutdown(true)\n"

	// Check if screen exists first
	out, _ := utils.RunCommandOutput("screen", "-ls")
	if !strings.Contains(out, screenName) {
		m.Log("%s 似乎没有在运行喵。", shardName)
		return
	}

	m.Log("正在向 %s 发送关闭指令...", shardName)
	utils.RunCommand("screen", "-S", screenName, "-p", "0", "-X", "stuff", cmd)

	// Wait a bit
	time.Sleep(3 * time.Second)
}

// IsRunning checks if the server is running
// 检查服务器是否运行
func (m *Manager) IsRunning() bool {
	out, err := utils.RunCommandOutput("screen", "-ls")
	if err != nil {
		return false
	}
	return strings.Contains(out, "dst_Master") || strings.Contains(out, "dst_Caves")
}
