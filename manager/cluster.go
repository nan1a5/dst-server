package manager

import (
	"dst-manager/utils"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListClusters returns a list of available cluster names
// 列出所有存档
func (m *Manager) ListClusters() []string {
	entries, err := os.ReadDir(m.Config.ClusterDir)
	if err != nil {
		// Try to create if not exists
		os.MkdirAll(m.Config.ClusterDir, 0755)
		return []string{}
	}

	var clusters []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if it looks like a cluster (has cluster.ini or cluster_token.txt)
			// Or just assume all dirs are clusters
			clusters = append(clusters, entry.Name())
		}
	}
	sort.Strings(clusters)
	return clusters
}

// SelectCluster allows the user to select a cluster
// 选择存档
func (m *Manager) SelectCluster(prompt string) string {
	clusters := m.ListClusters()
	if len(clusters) == 0 {
		m.Log("没有找到任何存档喵~ 请先创建一个吧！")
		return ""
	}

	m.Log("%s", prompt)
	for i, name := range clusters {
		fmt.Printf("  [%d] %s\n", i+1, name)
	}

	input := utils.ReadInput("请输入编号 (输入 0 取消): ")
	if input == "0" {
		return ""
	}

	var index int
	_, err := fmt.Sscanf(input, "%d", &index)
	if err != nil || index < 1 || index > len(clusters) {
		m.Log("输入的编号不对喵~")
		return ""
	}

	return clusters[index-1]
}

// CreateCluster creates a new cluster
// 创建新存档
func (m *Manager) CreateCluster() {
	m.Log("开始创建新存档喵...")

	name := utils.ReadInput("请输入存档目录名 (例如 Cluster_2): ")
	if name == "" {
		m.Log("存档名不能为空喵！")
		return
	}
	
	// Validate name (simple check)
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, " ") {
		m.Log("存档名包含非法字符或空格喵，请使用字母数字下划线~")
		return
	}

	clusterPath := filepath.Join(m.Config.ClusterDir, name)
	if _, err := os.Stat(clusterPath); err == nil {
		m.Log("这个存档名已经存在了喵！")
		return
	}

	token := utils.ReadInput("请输入 Cluster Token (从 Klei 官网获取): ")
	if token == "" {
		m.Log("Token 不能为空喵！没有 Token 服务器没法启动哦~")
		return
	}

	// Create directories
	dirs := []string{
		clusterPath,
		filepath.Join(clusterPath, "Master"),
		filepath.Join(clusterPath, "Caves"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			m.Log("创建目录失败了喵: %v", err)
			return
		}
	}

	// Write cluster_token.txt
	os.WriteFile(filepath.Join(clusterPath, "cluster_token.txt"), []byte(token), 0644)

	// Write default configs
	m.writeDefaultConfigs(clusterPath)

	m.Log("存档 %s 创建成功啦！快去启动试试吧喵~", name)
}

// writeDefaultConfigs writes default .ini files
func (m *Manager) writeDefaultConfigs(clusterPath string) {
	// cluster.ini
	clusterIni := `[GAMEPLAY]
game_mode = survival
max_players = 6
pvp = false
pause_when_empty = true

[NETWORK]
cluster_name = New Go DST Server
cluster_description = Created by DST Manager
cluster_intention = cooperative

[MISC]
console_enabled = true
`
	os.WriteFile(filepath.Join(clusterPath, "cluster.ini"), []byte(clusterIni), 0644)

	// Master/server.ini
	masterIni := `[NETWORK]
server_port = 10999

[SHARD]
is_master = true

[STEAM]
authentication_port = 8766
master_server_port = 27016
`
	os.WriteFile(filepath.Join(clusterPath, "Master", "server.ini"), []byte(masterIni), 0644)

	// Caves/server.ini
	cavesIni := `[NETWORK]
server_port = 10998

[SHARD]
is_master = false
name = Caves

[STEAM]
authentication_port = 8765
master_server_port = 27015

[WORLD]
id = 2
`
	os.WriteFile(filepath.Join(clusterPath, "Caves", "server.ini"), []byte(cavesIni), 0644)
	
	// Create worldgenoverride.lua for Caves (essential for caves to work properly)
	cavesOverride := `return {
    override_enabled = true,
    preset = "DST_CAVE",
}
`
	os.WriteFile(filepath.Join(clusterPath, "Caves", "worldgenoverride.lua"), []byte(cavesOverride), 0644)
}

// DeleteCluster deletes a cluster
// 删除存档
func (m *Manager) DeleteCluster() {
	cluster := m.SelectCluster("请选择要删除的存档:")
	if cluster == "" {
		return
	}

	confirm := utils.ReadInput(fmt.Sprintf("真的要删除存档 %s 吗？删除后找不回来的喵！(y/n): ", cluster))
	if strings.ToLower(confirm) != "y" {
		m.Log("吓死宝宝了，还好没删喵~")
		return
	}

	err := os.RemoveAll(filepath.Join(m.Config.ClusterDir, cluster))
	if err != nil {
		m.Log("删除失败了喵: %v", err)
	} else {
		m.Log("存档 %s 已经变成蝴蝶飞走了喵...", cluster)
	}
}

// ManageClusters menu
// 存档管理菜单
func (m *Manager) ManageClusters() {
	for {
		fmt.Println("\n============== 存档管理 ==============")
		fmt.Println("  1. 创建新存档")
		fmt.Println("  2. 删除存档")
		fmt.Println("  3. 查看存档列表")
		fmt.Println("  0. 返回主菜单")
		fmt.Println("======================================")

		choice := utils.ReadInput("请输入选项数字喵: ")
		switch choice {
		case "1":
			m.CreateCluster()
		case "2":
			m.DeleteCluster()
		case "3":
			clusters := m.ListClusters()
			m.Log("当前共有 %d 个存档喵:", len(clusters))
			for i, c := range clusters {
				fmt.Printf("  [%d] %s\n", i+1, c)
			}
		case "0":
			return
		default:
			m.Log("听不懂喵~")
		}
	}
}
