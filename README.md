# DST Server Manager (Go Version)

这是为 Ubuntu 系统编写的饥荒联机版（DST）服务器管理工具。
它可以帮助你快速安装、更新、启动和备份 DST 服务器。

## 功能特性

*   **自动安装依赖**: 自动检测并安装 SteamCMD 和 DST 所需的系统库 (lib32gcc-s1 等)。
*   **一键更新**: 支持更新 SteamCMD 和 DST 服务端。
*   **进程管理**: 使用 `screen` 在后台运行服务器 (Master 和 Caves)。
*   **备份管理**: 支持一键备份存档到 tar.gz 文件，并支持恢复。
*   **简单易用**: 交互式数字菜单。

## 使用说明

### 1. 编译或下载

如果你已经在 Ubuntu 上安装了 Go 环境：

```bash
git clone <repository_url>
cd dst-server
go build -o dst-manager main.go
```

或者在 Windows 上交叉编译 (目标为 Linux):

```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o dst-manager main.go
```

### 2. 运行

将编译好的 `dst-manager` 上传到你的 Ubuntu 服务器，然后运行：

```bash
chmod +x dst-manager
./dst-manager
```

### 3. 首次配置

1.  运行程序后，**首先选择 `1`** 进行环境安装。
    *   这会安装系统依赖（可能需要 sudo 密码）。
    *   下载并安装 SteamCMD。
    *   下载并安装 DST 服务端。
2.  安装完成后，你需要准备服务器存档配置。
    *   默认存档路径: `~/.klei/DoNotStarveTogether/Cluster_1/`
    *   请确保你已经获取了 `cluster_token.txt` 并配置了 `cluster.ini`。
    *   你可以从本地电脑生成的存档直接上传到该目录。
3.  配置好存档后，选择 `2` 启动服务器。

### 4. 目录结构

*   `~/steamcmd`: SteamCMD 安装目录
*   `~/dst-server`: DST 服务端安装目录
*   `~/.klei/DoNotStarveTogether`: 存档目录
*   `~/dst-backups`: 备份文件存放目录

## 注意事项

*   本工具依赖 `screen` 来管理后台进程。
*   请确保你的服务器有足够的内存 (建议至少 2GB)。
*   恢复存档功能会覆盖当前的 `Cluster_1`，请谨慎操作。

## 许可证

MIT
