package manager

import (
	"dst-manager/utils"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BackupCluster creates a backup of the cluster
// 备份存档
func (m *Manager) BackupCluster(filename string) {
	m.Log("开始备份存档，请稍候喵...")

	// Create backup dir
	if err := os.MkdirAll(m.Config.BackupDir, 0755); err != nil {
		m.Log("创建备份目录失败了喵: %v", err)
		return
	}

	backupPath := filepath.Join(m.Config.BackupDir, filename)

	// Target: ~/.klei/DoNotStarveTogether/Cluster_1
	// We backup the whole Cluster_1 folder
	clusterPath := filepath.Join(m.Config.ClusterDir, "Cluster_1")

	if _, err := os.Stat(clusterPath); os.IsNotExist(err) {
		m.Log("找不到存档目录喵: %s", clusterPath)
		return
	}

	// tar -czf <backup> -C <parent> Cluster_1
	parentDir := filepath.Dir(clusterPath)
	err := utils.RunCommand("tar", "-czf", backupPath, "-C", parentDir, "Cluster_1")
	if err != nil {
		m.Log("备份失败了喵: %v", err)
		return
	}

	m.Log("存档备份成功！文件保存在: %s", filename)
}

// ListBackups lists available backups and returns them
// 列出所有备份并返回列表
func (m *Manager) ListBackups() []string {
	entries, err := os.ReadDir(m.Config.BackupDir)
	if err != nil {
		m.Log("无法读取备份目录喵: %v", err)
		return nil
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tar.gz") {
			backups = append(backups, entry.Name())
		}
	}

	// Sort reverse (newest first)
	sort.Sort(sort.Reverse(sort.StringSlice(backups)))

	m.Log("找到以下备份文件喵:")
	for i, name := range backups {
		fmt.Printf("  [%d] %s\n", i+1, name)
	}
	return backups
}

// RestoreBackup restores a backup
// 恢复存档
func (m *Manager) RestoreBackup() {
	backups := m.ListBackups()
	if len(backups) == 0 {
		m.Log("没有找到备份文件喵~")
		return
	}

	input := utils.ReadInput("请输入要恢复的备份编号 (输入 0 取消): ")
	if input == "0" {
		return
	}

	var index int
	_, err := fmt.Sscanf(input, "%d", &index)
	if err != nil || index < 1 || index > len(backups) {
		m.Log("输入的编号不对喵~")
		return
	}

	selectedBackup := backups[index-1]
	backupPath := filepath.Join(m.Config.BackupDir, selectedBackup)

	if m.IsRunning() {
		m.Log("服务器正在运行喵！请先停止服务器再恢复存档~")
		return
	}

	m.Log("准备恢复备份: %s", selectedBackup)
	confirm := utils.ReadInput("这会覆盖当前的存档，确定要继续吗喵？(y/n): ")
	if strings.ToLower(confirm) != "y" {
		m.Log("操作已取消喵~")
		return
	}

	// Remove current Cluster_1
	clusterPath := filepath.Join(m.Config.ClusterDir, "Cluster_1")
	if err := os.RemoveAll(clusterPath); err != nil {
		m.Log("删除旧存档失败了喵: %v", err)
		return
	}

	// Restore
	// tar -xzf <backup> -C <parent>
	parentDir := filepath.Dir(clusterPath)
	err = utils.RunCommand("tar", "-xzf", backupPath, "-C", parentDir)
	if err != nil {
		m.Log("恢复存档失败了喵: %v", err)
	} else {
		m.Log("存档恢复成功啦！")
	}
}
