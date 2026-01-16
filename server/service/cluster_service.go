package service

import (
	"dst-manager/config"
	"dst-manager/utils/clusterUtils"

	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Config struct {
	Gameplay GameplayConfig `ini:"GAMEPLAY"`
	Network  NetworkConfig  `ini:"NETWORK"`
	Misc     MiscConfig     `ini:"MISC"`
	Shard    ShardConfig    `ini:"SHARD"`
}

type GameplayConfig struct {
	GameMode       string `ini:"game_mode"`
	MaxPlayers     int    `ini:"max_players"`
	PVP            bool   `ini:"pvp"`
	PauseWhenEmpty bool   `ini:"pause_when_empty"`
}

type NetworkConfig struct {
	LanOnlyCluster     bool   `ini:"lan_only_cluster"`
	ClusterPassword    string `ini:"cluster_password"`
	ClusterDescription string `ini:"cluster_description"`
	ClusterName        string `ini:"cluster_name"`
	OfflineCluster     bool   `ini:"offline_cluster"`
	ClusterLanguage    string `ini:"cluster_language"`
	ClusterCloudID     string `ini:"cluster_cloud_id"`
}

type MiscConfig struct {
	ConsoleEnabled bool `ini:"console_enabled"`
}

type ShardConfig struct {
	ShardEnabled bool   `ini:"shard_enabled"`
	BindIP       string `ini:"bind_ip"`
	MasterIP     string `ini:"master_ip"`
	MasterPort   int    `ini:"master_port"`
	ClusterKey   string `ini:"cluster_key"`
}

type WorldPreset struct {
	Desc                string `json:"desc"`
	HideMinimap         bool   `json:"hideminimap"`
	ID                  string `json:"id"`
	Location            string `json:"location"`
	MaxPlaylistPosition int    `json:"max_playlist_position"`
	MinPlaylistPosition int    `json:"min_playlist_position"`
	Name                string `json:"name"`
	Playstyle           string `json:"playstyle"`
	Version             int    `json:"version"`

	// 世界设置
	Overrides map[string]any `json:"overrides"`

	RandomSetPieces   []string `json:"random_set_pieces"`
	RequiredPrefabs   []string `json:"required_prefabs"`
	RequiredSetpieces []string `json:"required_setpieces"`

	SettingsDesc string `json:"settings_desc"`
	SettingsID   string `json:"settings_id"`
	SettingsName string `json:"settings_name"`

	WorldgenDesc string `json:"worldgen_desc"`
	WorldgenID   string `json:"worldgen_id"`
	WorldgenName string `json:"worldgen_name"`
}

type ClusterService interface {
	ListClusters() ([]string, error)
	CreateCluster(clusterName string, clusterToken string) error
	DeleteCluster(clusterName string) error
	RenameCluster(clusterName string, newName string) error
	GetAdminList(clusterName string) ([]string, error)
	AddAdmin(clusterName string, username string) error
	RemoveAdmin(clusterName string, username string) error
	GetBlackList(clusterName string) ([]string, error)
	SetBlackList(clusterName string, blackList []string) error
	SetToken(clusterName string, token string) error
	LoadConfig(clusterName string) (*Config, error)
	SetConfig(clusterName string, config *Config) error
	SetModOverride(clusterName string, modOverride string) error
	GetModOverride(clusterName string) (string, error)
	GetLevelOverride(clusterName string, levelName string) (*WorldPreset, error)
	SetLevelOverride(clusterName string, levelName string, override *WorldPreset) error
	GetServerLog(clusterName string, levelName string) ([]string, error)
	GetServerChatLog(clusterName string, levelName string) ([]string, error)
}

type clusterService struct {
	Config *config.Config
}

func NewClusterService() ClusterService {
	return &clusterService{
		Config: config.NewConfig(),
	}
}

func (c *clusterService) ListClusters() ([]string, error) {
	entries, err := os.ReadDir(c.Config.ClusterDir)
	if err != nil {
		// Try to create if not exists
		os.MkdirAll(c.Config.ClusterDir, 0755)
		return []string{}, err
	}

	var clusters []string
	for _, entry := range entries {
		if entry.IsDir() {
			clusters = append(clusters, entry.Name())
		}
	}
	sort.Strings(clusters)
	return clusters, nil
}

func (c *clusterService) CreateCluster(clusterName string, clusterToken string) error {
	if clusterName == "" {
		return errors.New("存档名不得为空")
	}

	// Validate name (simple check)
	if strings.Contains(clusterName, "/") || strings.Contains(clusterName, "\\") || strings.Contains(clusterName, " ") {
		return errors.New("存档名包含非法字符或空格，请使用字母数字下划线")
	}

	clusterPath := filepath.Join(c.Config.ClusterDir, clusterName)
	if _, err := os.Stat(clusterPath); err == nil {
		return errors.New("存档名已存在了")
	}

	if clusterToken == "" {
		return errors.New("token 不能为空！")
	}

	// Create directories
	dirs := []string{
		clusterPath,
		filepath.Join(clusterPath, "Master"),
		filepath.Join(clusterPath, "Caves"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建存档目录失败: %v", err)
		}
	}

	// Write cluster_token.txt
	if err := os.WriteFile(filepath.Join(clusterPath, "cluster_token.txt"), []byte(clusterToken), 0644); err != nil {
		return fmt.Errorf("写入 cluster_token 文件失败: %v", err)
	}

	// Write default configs
	if err := clusterUtils.WriteDefaultConfig(clusterPath, clusterName, ""); err != nil {
		return fmt.Errorf("写入默认配置失败: %v", err)
	}

	return nil
}

func (c *clusterService) DeleteCluster(clusterName string) error {
	// TODO 检查服务器是否运行
	if clusterName == "" {
		return errors.New("要删除的存档名不能为空")
	}

	if err := os.RemoveAll(filepath.Join(c.Config.ClusterDir, clusterName)); err != nil {
		return fmt.Errorf("删除存档目录失败: %v", err)
	}

	return nil
}

func (c *clusterService) RenameCluster(clusterName string, newName string) error {
	// TODO 检查服务器是否运行
	return nil
}

func (c *clusterService) GetAdminList(clusterName string) ([]string, error) {
	if clusterName == "" {
		return nil, errors.New("要获取的存档名不能为空")
	}

	clusterPath := filepath.Join(c.Config.ClusterDir, clusterName)
	if _, err := os.Stat(clusterPath); err != nil {
		return nil, fmt.Errorf("存档目录不存在: %v", err)
	}

	adminListPath := filepath.Join(clusterPath, "adminlist.txt")

	adminList, err := clusterUtils.ReadAdminList(adminListPath)
	if err != nil {
		return nil, fmt.Errorf("读取 adminlist 文件失败: %v", err)
	}

	return adminList, nil
}

func (c *clusterService) AddAdmin(clusterName string, username string) error {
	if clusterName == "" {
		return errors.New("存档名不能为空")
	}

	clusterPath := filepath.Join(c.Config.ClusterDir, clusterName)
	if _, err := os.Stat(clusterPath); err != nil {
		return fmt.Errorf("存档目录不存在: %v", err)
	}

	if err := clusterUtils.AddAdmin(clusterPath, username); err != nil {
		return fmt.Errorf("添加管理员失败: %v", err)
	}

	return nil
}

func (c *clusterService) RemoveAdmin(clusterName string, username string) error {
	if clusterName == "" {
		return errors.New("存档名不能为空")
	}

	clusterPath := filepath.Join(c.Config.ClusterDir, clusterName)
	if _, err := os.Stat(clusterPath); err != nil {
		return fmt.Errorf("存档目录不存在: %v", err)
	}

	if err := clusterUtils.RemoveAdmin(clusterPath, username); err != nil {
		return fmt.Errorf("删除管理员失败: %v", err)
	}

	return nil
}

func (c *clusterService) GetBlackList(clusterName string) ([]string, error) {
	return nil, nil
}

func (c *clusterService) SetBlackList(clusterName string, blackList []string) error {
	return nil
}

func (c *clusterService) SetToken(clusterName string, token string) error {
	return nil
}

func (c *clusterService) LoadConfig(clusterName string) (*Config, error) {
	return nil, nil
}

func (c *clusterService) SetConfig(clusterName string, config *Config) error {
	return nil
}
func (c *clusterService) SetModOverride(clusterName string, modOverride string) error {
	return nil
}

func (c *clusterService) GetModOverride(clusterName string) (string, error) {
	return "", nil
}
func (c *clusterService) GetLevelOverride(clusterName string, levelName string) (*WorldPreset, error) {
	return nil, nil
}

func (c *clusterService) SetLevelOverride(clusterName string, levelName string, override *WorldPreset) error {
	return nil
}

func (c *clusterService) GetServerLog(clusterName string, levelName string) ([]string, error) {
	return nil, nil
}

func (c *clusterService) GetServerChatLog(clusterName string, levelName string) ([]string, error) {
	return nil, nil
}
