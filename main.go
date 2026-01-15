package main

import (
	"dst-manager/manager"
	"dst-manager/utils"
	"fmt"
	"os"
	"time"
)

func main() {
	// Catch potential panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n!!! 小花酱晕倒了喵 (Panic) !!!\n错误信息: %v\n", r)
		}
	}()

	mgr := manager.NewManager()

	fmt.Println("========================================")
	mgr.Log("欢迎使用饥荒联机版服务器管理助手喵！")
	mgr.Log("我是小花酱，会帮主人管理服务器哦~")
	fmt.Println("========================================")

	for {
		printMenu()
		choice := utils.ReadInput("请输入选项数字喵: ")

		switch choice {
		case "1":
			mgr.InstallDependencies()
			if err := mgr.InstallSteamCMD(); err != nil {
				mgr.Log("SteamCMD 安装失败了喵: %v", err)
				continue
			}
			mgr.InstallDST()
		case "2":
			mgr.StartServer()
		case "3":
			mgr.StopServer()
		case "4":
			mgr.StopServer()
			mgr.StartServer()
		case "5":
			// mgr.BackupCluster()
			name := utils.ReadInput("请输入备份文件名喵: ")
			if name == "" {
				timestamp := time.Now().Format("20070831_162739")
				name = fmt.Sprintf("backup_%s.tar.gz", timestamp)
			}
			mgr.BackupCluster(name)
		case "6":
			mgr.ListBackups()
		case "7":
			mgr.RestoreBackup()
		case "8":
			mgr.ManageClusters()
		case "0":
			mgr.Log("好的喵，小花酱先退下了，主人要注意休息哦~")
			os.Exit(0)
		default:
			mgr.Log("看不懂这个指令喵，请重新输入~")
		}

		fmt.Println("\n按回车键继续喵...")
		utils.ReadInput("")
	}
}

func printMenu() {
	fmt.Println("\n============== 功能菜单 ==============")
	fmt.Println("  1. 安装/更新环境 (依赖+SteamCMD+DST)")
	fmt.Println("  2. 启动服务器")
	fmt.Println("  3. 停止服务器")
	fmt.Println("  4. 重启服务器")
	fmt.Println("  5. 备份存档")
	fmt.Println("  6. 查看备份列表")
	fmt.Println("  7. 恢复存档")
	fmt.Println("  8. 存档管理")
	fmt.Println("  0. 退出")
	fmt.Println("======================================")
}
