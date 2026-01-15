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

	// Try to guess cluster name from backup filename
	// Format: backup_<cluster>_<timestamp>.tar.gz
	// e.g. backup_Cluster_1_20230101_120000.tar.gz
	var targetCluster string
	parts := strings.Split(selectedBackup, "_")
	if len(parts) >= 4 {
		// Attempt to reconstruct cluster name (Cluster_X)
		// Assuming cluster name starts with Cluster
		// backup, Cluster, 1, 2023..., tar.gz
		// This is weak parsing, but better than hardcoding.
		// Let's ask user instead.
	}

	targetCluster = m.SelectCluster("请选择要恢复到的存档位置 (这会覆盖该存档喵！):")
	if targetCluster == "" {
		return
	}

	if m.IsRunning() {
		m.Log("服务器正在运行喵！请先停止服务器再恢复存档~")
		return
	}

	m.Log("准备将备份 %s 恢复到 %s", selectedBackup, targetCluster)
	confirm := utils.ReadInput("这会完全覆盖目标存档，确定要继续吗喵？(y/n): ")
	if strings.ToLower(confirm) != "y" {
		m.Log("操作已取消喵~")
		return
	}

	// Remove current cluster dir
	clusterPath := filepath.Join(m.Config.ClusterDir, targetCluster)
	if err := os.RemoveAll(clusterPath); err != nil {
		m.Log("删除旧存档失败了喵: %v", err)
		return
	}

	// Restore
	// We need to be careful. The tar contains the top-level folder name (e.g. Cluster_1).
	// If we restore to Cluster_2, we might get Cluster_2/Cluster_1.
	// So we should extract to a temp dir, rename, then move.

	tempDir := filepath.Join(m.Config.BackupDir, "temp_restore")
	os.RemoveAll(tempDir)
	os.MkdirAll(tempDir, 0755)

	err = utils.RunCommand("tar", "-xzf", backupPath, "-C", tempDir)
	if err != nil {
		m.Log("解压备份失败了喵: %v", err)
		return
	}

	// Find the extracted folder
	entries, _ := os.ReadDir(tempDir)
	if len(entries) != 1 || !entries[0].IsDir() {
		m.Log("备份文件结构看起来怪怪的喵，无法自动恢复...")
		return
	}
	extractedName := entries[0].Name()
	extractedPath := filepath.Join(tempDir, extractedName)

	// Move to target
	err = os.Rename(extractedPath, clusterPath)
	if err != nil {
		m.Log("移动存档失败了喵: %v", err)
	} else {
		m.Log("存档恢复成功啦！已恢复到 %s", targetCluster)
	}
	os.RemoveAll(tempDir)
}
