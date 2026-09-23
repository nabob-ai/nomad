---
title: NomadOS
description: 统一运行时系统 - 多智能体协作的完整解决方案
navigation:
  icon: i-lucide-server
---

# NomadOS - 统一运行时系统

NomadOS 是 Nomad 框架的统一运行时系统，提供多智能体协作的完整解决方案。它管理所有 Agents、Stars、Workflows，并自动生成 REST API 端点，支持多种 Interface 类型。

## 🌟 核心特性

### 统一资源管理

- **Cosmos**: 智能体生命周期管理器，替代原有的 Pool 概念
- **Stars**: 多智能体协作单元，替代原有的 Room 概念
- **Workflows**: 工作流管理和执行
- **自动发现**: 自动注册和发现所有资源

### 多接口支持

- **HTTP Interface**: RESTful API 接口
- **A2A Interface**: Agent-to-Agent 通信接口
- **AGUI Interface**: 控制平面 UI 集成接口
- **插件化**: 支持自定义 Interface 扩展

### 自动 API 生成

- 为所有注册的 Agents 自动生成 REST 端点
- 为所有 Stars 自动生成协作管理 API
- 为所有 Workflows 自动生成执行 API
- 支持健康检查、指标监控、认证授权

## 🚀 快速开始

### 基本使用

```go
package main

import (
    "context"
    "log"

    "github.com/nabob-ai/nomad/pkg/nomados"
    "github.com/nabob-ai/nomad/pkg/cosmos"
)

func main() {
    // 创建依赖
    deps := createDependencies()

    // 创建 Cosmos
    cosmos := cosmos.New(&cosmos.Options{
        Dependencies: deps,
        MaxAgents:    10,
    })

    // 创建 NomadOS
    os, err := nomados.New(&nomados.Options{
        Name:   "MyNomadOS",
        Port:   8080,
        Cosmos: cosmos,
    })
    if err != nil {
        log.Fatal(err)
    }

    // 启动服务
    if err := os.Serve(); err != nil {
        log.Fatal(err)
    }
}
```

### 注册 Agent

```go
// 创建 Agent
agentConfig := &types.AgentConfig{
    AgentID:    "chat-agent",
    TemplateID: "chat-template",
    ModelConfig: &types.ModelConfig{
        Provider: "anthropic",
        Model:    "claude-sonnet-4-5",
        APIKey:   "your-api-key",
    },
}

ag, err := agent.Create(ctx, agentConfig, deps)
if err != nil {
    log.Fatal(err)
}

// 注册到 NomadOS
if err := os.RegisterAgent("chat-agent", ag); err != nil {
    log.Fatal(err)
}
```

### 创建 Stars 协作

```go
// 创建 Stars
starsInstance := stars.New(cosmos, "ChatTeam")

// 添加成员
if err := starsInstance.AddMember("leader", "agent-1", "leader"); err != nil {
    log.Fatal(err)
}
if err := starsInstance.AddMember("worker", "agent-2", "worker"); err != nil {
    log.Fatal(err)
}

// 注册到 NomadOS
if err := os.RegisterStars("chat-team", starsInstance); err != nil {
    log.Fatal(err)
}
```

## 📡 API 端点

NomadOS 自动生成以下 REST API 端点：

### Agent 管理

```
GET    /api/agents              # 列出所有 Agent
POST   /api/agents/{id}/run     # 运行指定 Agent
GET    /api/agents/{id}/status  # 获取 Agent 状态
```

### Stars 协作

```
GET    /api/stars                # 列出所有 Stars
POST   /api/stars/{id}/run       # 运行 Stars 协作
POST   /api/stars/{id}/join      # 加入 Stars
POST   /api/stars/{id}/leave     # 离开 Stars
GET    /api/stars/{id}/members   # 获取成员列表
```

### Workflow 执行

```
GET    /api/workflows             # 列出所有 Workflow
POST   /api/workflows/{id}/execute # 执行 Workflow
```

### 系统

```
GET    /health                    # 健康检查
GET    /metrics                   # Prometheus 指标
```

## 🎯 架构设计

### 核心组件

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP API      │    │   A2A Interface │    │   AGUI Interface│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────────────────────────────────────────────────────┐
│                    NomadOS Registry                             │
├─────────────────┬─────────────────┬─────────────────────────────┤
│   Agents        │    Stars        │      Workflows              │
└─────────────────┴─────────────────┴─────────────────────────────┘
         │                       │                       │
┌─────────────────────────────────────────────────────────────────┐
│                      Cosmos                                     │
│                  (Agent Manager)                                │
└─────────────────────────────────────────────────────────────────┘
```

### Interface 抽象

```go
type Interface interface {
    // 基本信息
    Name() string
    Type() InterfaceType

    // 生命周期管理
    Start(ctx context.Context, os *NomadOS) error
    Stop(ctx context.Context) error

    // 事件回调
    OnAgentRegistered(agent *agent.Agent) error
    OnStarsRegistered(stars *stars.Stars) error
    OnWorkflowRegistered(wf workflow.Agent) error
}
```

## 🌌 与 Cosmos 和 Stars 的关系

### Cosmos (宇宙)

- **职责**: Agent 生命周期管理
- **功能**: 创建、销毁、监控 Agent
- **类似**: Kubernetes 的 Pod Manager

### Stars (星座)

- **职责**: 多 Agent 协作管理
- **功能**: 编组、通信、协作调度
- **类似**: Kubernetes 的 Service/Deployment

### NomadOS (星系操作系统)

- **职责**: 统一运行时和 API 网关
- **功能**: 资源注册、API 生成、接口管理
- **类似**: Kubernetes API Server + Ingress Controller

## ⚙️ 配置选项

```go
type Options struct {
    // 基本配置
    Name        string          // NomadOS 名称
    Port        int             // HTTP 端口 (默认: 8080)
    Cosmos      *cosmos.Cosmos  // Cosmos 实例 (必需)

    // API 配置
    APIPrefix   string          // API 路径前缀 (默认: /api)

    // 功能开关
    AutoDiscover   bool         // 自动发现 (默认: true)
    EnableCORS     bool         // CORS 支持 (默认: true)
    EnableAuth     bool         // 认证授权 (默认: false)
    EnableMetrics  bool         // 指标监控 (默认: true)
    EnableHealth   bool         // 健康检查 (默认: true)
    EnableLogging  bool         // 日志记录 (默认: true)

    // 认证配置
    APIKey        string        // API 密钥

    // 日志配置
    LogLevel      string        // 日志级别 (默认: info)
}
```

## 🔧 扩展开发

### 自定义 Interface

```go
type MyInterface struct {
    BaseInterface
    // 自定义字段
}

func (i *MyInterface) Start(ctx context.Context, os *NomadOS) error {
    // 自定义启动逻辑
    return nil
}

func (i *MyInterface) OnAgentRegistered(agent *agent.Agent) error {
    // 自定义 Agent 注册处理
    return nil
}
```

### 添加 Interface

```go
// 添加 HTTP Interface (默认已包含)
httpIface := nomados.NewHTTPInterface()
os.AddInterface(httpIface)

// 添加 A2A Interface
a2aIface := nomados.NewA2AInterface()
os.AddInterface(a2aIface)

// 添加 AGUI Interface
aguiIface := nomados.NewAGUIInterface("/agui")
os.AddInterface(aguiIface)
```

## 📊 监控和观测

### 健康检查

```bash
curl http://localhost:8080/health
```

### Prometheus 指标

```bash
curl http://localhost:8080/metrics
```

### 日志输出

NomadOS 使用结构化日志，支持不同级别：

```
🌟 NomadOS 'MyNomadOS' is running on http://localhost:8080
[Agent Create] Total tools loaded: 5
[Stars Join] Agent 'worker-1' joined stars 'chat-team'
```

## 🛡️ 安全特性

### 认证授权

```go
// 启用认证
os, err := nomados.New(&nomados.Options{
    EnableAuth: true,
    APIKey:     "your-secret-api-key",
})
```

### CORS 支持

```go
// 启用 CORS (默认已启用)
os, err := nomados.New(&nomados.Options{
    EnableCORS: true,
})
```

## 📚 示例项目

查看 `examples/nomados/` 目录下的完整示例：

- `basic/`: 基本 NomadOS 使用
- `interfaces/`: 多种 Interface 使用示例
- `collaboration/`: Stars 协作示例

## 相关文档

- [Cosmos 详细文档](/multi-agent/cosmos)
- [Stars 详细文档](/multi-agent/stars)
- [多Agent示例](/examples/multi-agent)
