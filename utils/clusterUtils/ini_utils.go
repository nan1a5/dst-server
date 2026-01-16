package clusterUtils

import (
	"dst-manager/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteDefaultConfig(clusterPath string, name string, desc string) error {
	if name == "" {
		name = "New DST Server"
	}
	if desc == "" {
		desc = "A new DST server"
	}
	// cluster.ini
	clusterIni := `[GAMEPLAY]
	game_mode = survival
	max_players = 6
	pvp = false
}
	pause_when_empty = true

	[NETWORK]
	cluster_name = ` + name + `
	cluster_description = ` + desc + `
	cluster_intention = cooperative

	[MISC]
	console_enabled = true
`
	if err := os.WriteFile(filepath.Join(clusterPath, "cluster.ini"), []byte(clusterIni), 0644); err != nil {
		return err
	}

	// Master/server.ini
	masterIni := `[NETWORK]
	server_port = 10999

	[SHARD]
	is_master = true

	[STEAM]
	authentication_port = 8766
	master_server_port = 27016
`
	if err := os.WriteFile(filepath.Join(clusterPath, "Master", "server.ini"), []byte(masterIni), 0644); err != nil {
		return err
	}

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
	if err := os.WriteFile(filepath.Join(clusterPath, "Caves", "server.ini"), []byte(cavesIni), 0644); err != nil {
		return err
	}

	// Create worldgenoverride.lua for Caves (essential for caves to work properly)
	cavesOverride := `return {
    override_enabled = true,
    preset = "DST_CAVE",
}
`
	if err := os.WriteFile(filepath.Join(clusterPath, "Caves", "worldgenoverride.lua"), []byte(cavesOverride), 0644); err != nil {
		return err
	}
	return nil
}

func createAdminList(clusterPath string) error {
	adminListPath := filepath.Join(clusterPath, "adminlist.txt")
	if err := os.WriteFile(adminListPath, []byte(""), 0644); err != nil {
		return err
	}
	return nil
}

func ReadAdminList(clusterPath string) ([]string, error) {
	adminListPath := filepath.Join(clusterPath, "adminlist.txt")
	if _, err := os.Stat(adminListPath); err != nil {
		if err := createAdminList(clusterPath); err != nil {
			return nil, fmt.Errorf("创建 adminlist 文件失败: %v", err)
		}
	}
	adminList, err := os.ReadFile(adminListPath)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(adminList), "\n"), nil
}

func AddAdmin(clusterPath string, username string) error {
	if username == "" {
		return errors.New("要添加的用户名不能为空")
	}
	if strings.Contains(username, "\n") {
		return errors.New("用户名不能包含换行符")
	}
	adminList, err := ReadAdminList(clusterPath)
	if err == nil {
		if strings.Contains(strings.Join(adminList, "\n"), username) {
			return errors.New("管理员已存在")
		}
	} else {
		return fmt.Errorf("读取 adminlist 文件失败: %v", err)
	}
	return writeAdminList(clusterPath, utils.MergeSlices(adminList, []string{username}))
}

func RemoveAdmin(clusterPath string, username string) error {
	if username == "" {
		return errors.New("要删除的用户名不能为空")
	}
	adminList, err := ReadAdminList(clusterPath)
	if err != nil {
		return fmt.Errorf("读取 adminlist 文件失败: %v", err)
	}

	return writeAdminList(clusterPath, utils.RemoveElement(adminList, username))
}

func writeAdminList(clusterPath string, adminList []string) error {
	adminListPath := filepath.Join(clusterPath, "adminlist.txt")
	oldAdminList, err := ReadAdminList(clusterPath)
	if err != nil {
		return fmt.Errorf("读取 adminlist 文件失败: %v", err)
	}
	adminList = utils.RemoveDuplicates(utils.MergeSlices(oldAdminList, adminList))
	adminListText := strings.Join(adminList, "\n")
	if err := os.WriteFile(adminListPath, []byte(adminListText), 0644); err != nil {
		return err
	}
	return nil
}
