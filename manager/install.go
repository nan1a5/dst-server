package manager

import (
	"dst-manager/utils"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// InstallDependencies installs necessary system dependencies
// 安装必要的系统依赖
func (m *Manager) InstallDependencies() {
	if runtime.GOOS != "linux" {
		m.Log("主人，咱喵检测到不是Linux系统，跳过依赖安装步骤哦~")
		return
	}

	m.Log("正在检查并安装系统依赖，可能需要主人输入密码呢...")
	// Ubuntu/Debian dependencies for SteamCMD (32-bit) and DST
	deps := []string{"lib32gcc-s1", "lib32stdc++6", "libcurl4-gnutls-dev:i386", "screen"}
	
	// Add architecture if needed
	utils.RunCommand("sudo", "dpkg", "--add-architecture", "i386")
	utils.RunCommand("sudo", "apt-get", "update")
	
	args := append([]string{"apt-get", "install", "-y"}, deps...)
	if err := utils.RunCommand("sudo", args...); err != nil {
		m.Log("依赖安装出错了喵: %v", err)
	} else {
		m.Log("系统依赖安装完成啦！")
	}
}

// InstallSteamCMD downloads and installs SteamCMD
// 下载并安装 SteamCMD
func (m *Manager) InstallSteamCMD() error {
	steamPath := filepath.Join(m.Config.SteamCMDDir, "steamcmd.sh")
	if _, err := os.Stat(steamPath); err == nil {
		m.Log("SteamCMD 已经安装过了喵~")
		return nil
	}

	m.Log("开始下载 SteamCMD...")
	if err := m.Config.EnsureDirs(); err != nil {
		return err
	}

	// Download steamcmd_linux.tar.gz
	tarUrl := "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
	tarPath := filepath.Join(m.Config.SteamCMDDir, "steamcmd_linux.tar.gz")
	
	// We use curl or wget. Let's assume curl exists or use Go's http client.
	// For simplicity, let's use system commands if available, checking utils.
	if utils.CheckCommandExists("curl") {
		utils.RunCommand("curl", "-o", tarPath, tarUrl)
	} else if utils.CheckCommandExists("wget") {
		utils.RunCommand("wget", "-O", tarPath, tarUrl)
	} else {
		return fmt.Errorf("没有找到 curl 或 wget，咱喵没法下载文件呢")
	}

	// Extract
	m.Log("正在解压 SteamCMD...")
	if err := utils.RunCommand("tar", "-xvzf", tarPath, "-C", m.Config.SteamCMDDir); err != nil {
		return err
	}

	// Remove tar
	os.Remove(tarPath)

	m.Log("SteamCMD 安装成功啦！")
	return nil
}

// InstallDST installs or updates the DST server
// 安装或更新 DST 服务端
func (m *Manager) InstallDST() error {
	m.Log("准备安装/更新 饥荒联机版服务端...")
	
	steamCmdPath := filepath.Join(m.Config.SteamCMDDir, "steamcmd.sh")
	installDir := m.Config.DSTInstallDir

	// Retry loop
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			m.Log("安装失败了，正在尝试第 %d 次重试喵...", i+1)
		}

		// cmd: ./steamcmd.sh +force_install_dir <path> +login anonymous +app_update 343050 validate +quit
		// Note: Sometimes running login first separately helps
		args := []string{
			"+force_install_dir", installDir,
			"+login", "anonymous",
			"+app_update", "343050", "validate",
			"+quit",
		}

		err := utils.RunCommand(steamCmdPath, args...)
		if err == nil {
			m.Log("饥荒联机版服务端安装/更新完成！可以开始冒险了喵！")
			return nil
		}
		
		m.Log("安装 DST 服务端出错: %v", err)
	}

	m.Log("呜呜...重试了 %d 次还是失败了喵。", maxRetries)
	m.Log("如果错误是 'Missing configuration'，请确认您的服务器架构是否为 x86/amd64。")
	return fmt.Errorf("安装失败")
}
