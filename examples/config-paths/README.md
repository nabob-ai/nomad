# Cross-Platform Paths 示例

本示例演示 Nomad 的跨平台路径管理系统，自动适配不同操作系统的标准路径约定。

## 功能特点

- 🖥️ 支持 macOS、Linux、Windows
- 📁 遵循各平台标准路径约定
- 🔧 便捷的文件路径生成方法
- 🏗️ 自动创建目录结构

## 运行示例

```bash
go run ./examples/config-paths/
```

## 平台路径约定

### macOS

```
配置: ~/Library/Application Support/nomad/
数据: ~/Library/Application Support/nomad/
缓存: ~/Library/Caches/nomad/
日志: ~/Library/Logs/nomad/
```

### Linux (XDG)

```
配置: ~/.config/nomad/
数据: ~/.local/share/nomad/
缓存: ~/.cache/nomad/
日志: ~/.local/share/nomad/logs/
```

### Windows

```
配置: %APPDATA%\nomad\
数据: %LOCALAPPDATA%\nomad\
缓存: %LOCALAPPDATA%\nomad\cache\
日志: %LOCALAPPDATA%\nomad\logs\
```

## 使用方式

### 1. 获取标准目录

```go
import "github.com/nabob-ai/nomad/pkg/config"

// 获取各类目录路径
configDir := config.ConfigDir()  // 配置目录
dataDir := config.DataDir()      // 数据目录
cacheDir := config.CacheDir()    // 缓存目录
logDir := config.LogDir()        // 日志目录
```

### 2. 获取文件路径

```go
// 配置文件
settingsPath := config.ConfigFile("settings.yaml")
// macOS: ~/Library/Application Support/nomad/settings.yaml

// 数据文件
dataPath := config.DataFile("data.json")
// macOS: ~/Library/Application Support/nomad/data.json

// 数据库文件
dbPath := config.DatabaseFile("sessions.db")
// macOS: ~/Library/Application Support/nomad/sessions.db

// 日志文件
logPath := config.LogFile("app.log")
// macOS: ~/Library/Logs/nomad/app.log

// 缓存文件
cachePath := config.CacheFile("temp.cache")
// macOS: ~/Library/Caches/nomad/temp.cache
```

### 3. 自定义应用名

```go
// 为不同应用创建路径管理器
paths := config.NewPaths("myapp")

configDir := paths.ConfigDir()
// macOS: ~/Library/Application Support/myapp/

dataDir := paths.DataDir()
// macOS: ~/Library/Application Support/myapp/
```

### 4. 确保目录存在

```go
// 确保目录存在（递归创建）
err := config.EnsureDir("/path/to/nested/directory")
if err != nil {
    log.Fatal(err)
}

// 常见用法：确保配置目录存在
err = config.EnsureDir(config.ConfigDir())
```

## 实际应用示例

### SQLite 数据库

```go
import (
    "github.com/nabob-ai/nomad/pkg/config"
    "github.com/nabob-ai/nomad/pkg/session/sqlite"
)

// 使用标准数据目录存储数据库
dbPath := config.DatabaseFile("sessions.db")

// 确保目录存在
config.EnsureDir(filepath.Dir(dbPath))

// 创建 SQLite 服务
service, err := sqlite.New(dbPath)
```

### 日志文件

```go
import (
    "github.com/nabob-ai/nomad/pkg/config"
    "log"
    "os"
)

// 使用标准日志目录
logPath := config.LogFile("app.log")

// 确保目录存在
config.EnsureDir(config.LogDir())

// 创建日志文件
logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
if err != nil {
    log.Fatal(err)
}
log.SetOutput(logFile)
```

### 配置文件

```go
import (
    "github.com/nabob-ai/nomad/pkg/config"
    "gopkg.in/yaml.v3"
    "os"
)

type AppConfig struct {
    APIKey string `yaml:"api_key"`
    Model  string `yaml:"model"`
}

// 加载配置
configPath := config.ConfigFile("config.yaml")
data, err := os.ReadFile(configPath)
if err != nil {
    // 使用默认配置
}

var cfg AppConfig
yaml.Unmarshal(data, &cfg)
```

## 环境变量覆盖

支持通过环境变量覆盖默认路径：

```bash
# Linux XDG 规范
export XDG_CONFIG_HOME=~/.myconfig
export XDG_DATA_HOME=~/.mydata
export XDG_CACHE_HOME=~/.mycache

# 运行应用
go run main.go
```

## 相关示例

- [session-sqlite](../session-sqlite/) - SQLite 会话存储
- [permission](../permission/) - 权限系统（配置持久化）
- [desktop](../desktop/) - 桌面应用集成
