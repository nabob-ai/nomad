# nomad Production Server

> 🚀 **生产级 AI 应用服务器** - 完整的认证、监控、部署支持

---

## 📋 概览

nomad Server 是一个生产就绪的应用服务器层，提供：

- ✅ **认证授权**: API Key、JWT
- ✅ **速率限制**: 可配置的请求限制
- ✅ **CORS 支持**: 完整的跨域配置
- ✅ **结构化日志**: JSON 格式日志
- ✅ **健康检查**: Kubernetes 就绪探针
- ✅ **指标收集**: Prometheus 集成
- ✅ **Docker 支持**: 多阶段构建
- ✅ **Kubernetes**: 完整的 K8s 部署配置

---

## 🚀 快速开始

### 使用默认配置

```go
package main

import (
    "log"
    "github.com/nabob-ai/nomad/pkg/store"
    "github.com/nabob-ai/nomad/server"
)

func main() {
    // 创建存储
    st, _ := store.NewJSONStore(".data")

    // 创建依赖
    deps := &server.Dependencies{
        Store: st,
    }

    // 创建服务器（使用默认配置）
    srv, err := server.New(server.DefaultConfig(), deps)
    if err != nil {
        log.Fatal(err)
    }

    // 启动服务器
    srv.Start()
}
```

### 自定义配置

```go
config := &server.Config{
    Host: "0.0.0.0",
    Port: 8080,
    Mode: "production",

    // 认证配置
    Auth: server.AuthConfig{
        APIKey: server.APIKeyConfig{
            Enabled: true,
            HeaderName: "X-API-Key",
            Keys: []string{"your-secure-api-key"},
        },
    },

    // CORS 配置
    CORS: server.CORSConfig{
        Enabled: true,
        AllowOrigins: []string{"https://yourdomain.com"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    },

    // 速率限制
    RateLimit: server.RateLimitConfig{
        Enabled: true,
        RequestsPerIP: 1000,
        WindowSize: time.Minute,
    },
}

srv, _ := server.New(config, deps)
srv.Start()
```

---

## 🐳 Docker 部署

### 构建镜像

```bash
docker build -t agentsdk/server:latest -f server/deploy/docker/Dockerfile .
```

### 运行容器

```bash
docker run -p 8080:8080 \
  -e API_KEY=your-api-key \
  -e MODE=production \
  agentsdk/server:latest
```

### 使用 Docker Compose

```bash
cd server/deploy/docker
docker-compose up -d
```

---

## ☸️ Kubernetes 部署

### 应用配置

```bash
kubectl apply -f server/deploy/k8s/
```

### 检查状态

```bash
kubectl get pods -l app=agentsdk
kubectl get svc agentsdk-server
```

### 查看日志

```bash
kubectl logs -f deployment/agentsdk-server
```

### 扩容

```bash
kubectl scale deployment agentsdk-server --replicas=5
```

---

## 📝 环境变量

| 变量      | 描述                                  | 默认值          |
| --------- | ------------------------------------- | --------------- |
| `HOST`    | 服务器监听地址                        | `0.0.0.0`       |
| `PORT`    | 服务器端口                            | `8080`          |
| `MODE`    | 运行模式 (`development`/`production`) | `development`   |
| `API_KEY` | API 密钥                              | `dev-key-12345` |

---

## 🔐 认证

### API Key 认证

```bash
curl -H "X-API-Key: your-api-key" \
  http://localhost:8080/v1/agents
```

### JWT 认证

```bash
curl -H "Authorization: Bearer your-jwt-token" \
  http://localhost:8080/v1/agents
```

---

## 📊 监控

### 健康检查

```bash
curl http://localhost:8080/health
```

响应：

```json
{
  "status": "healthy",
  "checks": {
    "database": "ok",
    "timestamp": "2024-11-17T12:00:00Z"
  },
  "version": "2.0.0"
}
```

### Prometheus 指标

```bash
curl http://localhost:8080/metrics
```

---

## 🔧 配置选项

### CORS 配置

```go
CORS: server.CORSConfig{
    Enabled: true,
    AllowOrigins: []string{"https://app.example.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge: 86400,
}
```

### 速率限制配置

```go
RateLimit: server.RateLimitConfig{
    Enabled: true,
    RequestsPerIP: 100,        // 每个时间窗口的请求数
    WindowSize: time.Minute,    // 时间窗口大小
    BurstSize: 20,              // 突发容量
}
```

### 日志配置

```go
Logging: server.LoggingConfig{
    Level: "info",              // debug, info, warn, error
    Format: "json",             // json 或 text
    Output: "stdout",           // stdout 或文件路径
    Structured: true,           // 结构化日志
}
```

---

## 📡 API 端点

### Agent 管理

- `POST /v1/agents` - 创建 Agent
- `GET /v1/agents` - 列出所有 Agents
- `GET /v1/agents/:id` - 获取 Agent 详情
- `PATCH /v1/agents/:id` - 更新 Agent
- `DELETE /v1/agents/:id` - 删除 Agent
- `POST /v1/agents/:id/run` - 运行 Agent
- `POST /v1/agents/:id/send` - 发送消息给 Agent
- `GET /v1/agents/:id/status` - 获取 Agent 状态
- `GET /v1/agents/:id/stats` - Agent 统计
- `POST /v1/agents/:id/resume` - 恢复 Agent
- `POST /v1/agents/chat` - Agent 对话
- `POST /v1/agents/chat/stream` - 流式对话

### Pool 管理 (v0.13.0+)

- `POST /v1/pool/agents` - 在池中创建 Agent
- `GET /v1/pool/agents` - 列出池中所有 Agents
- `GET /v1/pool/agents/:id` - 获取池中 Agent
- `POST /v1/pool/agents/:id/resume` - 恢复池中 Agent
- `DELETE /v1/pool/agents/:id` - 从池中移除 Agent
- `GET /v1/pool/stats` - 池统计信息

### Room 管理 (v0.13.0+)

- `POST /v1/rooms` - 创建 Room
- `GET /v1/rooms` - 列出所有 Rooms
- `GET /v1/rooms/:id` - 获取 Room 详情
- `DELETE /v1/rooms/:id` - 删除 Room
- `POST /v1/rooms/:id/join` - 加入 Room
- `POST /v1/rooms/:id/leave` - 离开 Room
- `POST /v1/rooms/:id/say` - 在 Room 中发送消息
- `GET /v1/rooms/:id/members` - 获取 Room 成员
- `GET /v1/rooms/:id/history` - 获取 Room 历史消息

### Memory 管理

- `POST /v1/memory/working` - 创建工作记忆
- `GET /v1/memory/working` - 列出工作记忆
- `GET /v1/memory/working/:id` - 获取工作记忆
- `PATCH /v1/memory/working/:id` - 更新工作记忆
- `DELETE /v1/memory/working/:id` - 删除工作记忆
- `POST /v1/memory/working/clear` - 清空工作记忆
- `POST /v1/memory/semantic` - 创建语义记忆
- `POST /v1/memory/semantic/search` - 搜索语义记忆
- `GET /v1/memory/provenance/:id` - 获取溯源信息
- `POST /v1/memory/consolidate` - 记忆整合

### Session 管理

- `POST /v1/sessions` - 创建会话
- `GET /v1/sessions` - 列出所有会话
- `GET /v1/sessions/:id` - 获取会话详情
- `PATCH /v1/sessions/:id` - 更新会话
- `DELETE /v1/sessions/:id` - 删除会话
- `GET /v1/sessions/:id/messages` - 获取会话消息
- `GET /v1/sessions/:id/checkpoints` - 获取会话检查点
- `POST /v1/sessions/:id/resume` - 恢复会话
- `GET /v1/sessions/:id/stats` - 会话统计

### Workflow 管理

- `POST /v1/workflows` - 创建工作流
- `GET /v1/workflows` - 列出所有工作流
- `GET /v1/workflows/:id` - 获取工作流详情
- `PATCH /v1/workflows/:id` - 更新工作流
- `DELETE /v1/workflows/:id` - 删除工作流
- `POST /v1/workflows/:id/execute` - 执行工作流
- `POST /v1/workflows/:id/suspend` - 暂停工作流
- `POST /v1/workflows/:id/resume` - 恢复工作流
- `GET /v1/workflows/:id/executions` - 获取执行记录
- `GET /v1/workflows/:id/executions/:eid` - 获取执行详情

### Tool 管理

- `POST /v1/tools` - 创建工具
- `GET /v1/tools` - 列出所有工具
- `GET /v1/tools/:id` - 获取工具详情
- `PATCH /v1/tools/:id` - 更新工具
- `DELETE /v1/tools/:id` - 删除工具
- `POST /v1/tools/:id/execute` - 执行工具

### Middleware 管理

- `POST /v1/middlewares` - 创建中间件
- `GET /v1/middlewares` - 列出所有中间件
- `GET /v1/middlewares/:id` - 获取中间件详情
- `PATCH /v1/middlewares/:id` - 更新中间件
- `DELETE /v1/middlewares/:id` - 删除中间件
- `POST /v1/middlewares/:id/enable` - 启用中间件
- `POST /v1/middlewares/:id/disable` - 禁用中间件
- `POST /v1/middlewares/:id/reload` - 重新加载中间件
- `GET /v1/middlewares/:id/stats` - 中间件统计
- `GET /v1/middlewares/registry` - 列出注册表
- `POST /v1/middlewares/registry/:id/install` - 安装中间件
- `DELETE /v1/middlewares/registry/:id/uninstall` - 卸载中间件
- `GET /v1/middlewares/registry/:id/info` - 获取中间件信息
- `POST /v1/middlewares/registry/reload-all` - 重新加载所有

### Telemetry 遥测

- `POST /v1/telemetry/metrics` - 记录指标
- `GET /v1/telemetry/metrics` - 列出指标
- `POST /v1/telemetry/traces` - 记录追踪
- `POST /v1/telemetry/traces/query` - 查询追踪
- `POST /v1/telemetry/logs` - 记录日志
- `POST /v1/telemetry/logs/query` - 查询日志

### Eval 评估

- `POST /v1/eval/text` - 文本评估
- `POST /v1/eval/session` - 会话评估
- `POST /v1/eval/batch` - 批量评估
- `POST /v1/eval/custom` - 自定义评估
- `GET /v1/eval/evals` - 列出评估
- `GET /v1/eval/evals/:id` - 获取评估详情
- `DELETE /v1/eval/evals/:id` - 删除评估
- `POST /v1/eval/benchmarks` - 创建基准测试
- `GET /v1/eval/benchmarks` - 列出基准测试
- `GET /v1/eval/benchmarks/:id` - 获取基准测试
- `DELETE /v1/eval/benchmarks/:id` - 删除基准测试
- `POST /v1/eval/benchmarks/:id/run` - 运行基准测试

### MCP 管理

- `POST /v1/mcp/servers` - 创建 MCP 服务器
- `GET /v1/mcp/servers` - 列出 MCP 服务器
- `GET /v1/mcp/servers/:id` - 获取 MCP 服务器
- `PATCH /v1/mcp/servers/:id` - 更新 MCP 服务器
- `DELETE /v1/mcp/servers/:id` - 删除 MCP 服务器
- `POST /v1/mcp/servers/:id/connect` - 连接 MCP 服务器
- `POST /v1/mcp/servers/:id/disconnect` - 断开 MCP 服务器

### System 系统

- `GET /v1/system/config` - 列出配置
- `GET /v1/system/config/:key` - 获取配置
- `PUT /v1/system/config/:key` - 更新配置
- `DELETE /v1/system/config/:key` - 删除配置
- `GET /v1/system/info` - 系统信息
- `GET /v1/system/health` - 健康检查
- `GET /v1/system/stats` - 系统统计
- `POST /v1/system/reload` - 重新加载
- `POST /v1/system/gc` - 垃圾回收
- `POST /v1/system/backup` - 备份

### 可观测性

- `GET /health` - 健康检查
- `GET /metrics` - Prometheus 指标

完整 API 文档请参考: [API Reference](../../docs/content/14.api-reference/)

---

## 🏗️ 架构

```
┌─────────────────────────────────────┐
│   Client SDKs                        │
│   - client-js, React, AI SDK         │
└────────────┬────────────────────────┘
             │ HTTP/WebSocket
┌────────────▼────────────────────────┐
│   server/ (生产级应用服务器)         │
│   ├── 认证授权                       │
│   ├── 速率限制                       │
│   ├── CORS 处理                      │
│   ├── 结构化日志                     │
│   ├── 健康检查                       │
│   └── 指标收集                       │
└────────────┬────────────────────────┘
             │ 纯 Go 接口
┌────────────▼────────────────────────┐
│   pkg/ (核心 SDK)                   │
│   - Agent, Memory, Workflow...       │
└─────────────────────────────────────┘
```

---

## 🔄 与 cmd/agentsdk 的对比

| 特性         | cmd/agentsdk | server/             |
| ------------ | ------------ | ------------------- |
| **定位**     | 演示/开发    | 生产部署            |
| **认证**     | ❌           | ✅ API Key + JWT    |
| **速率限制** | ❌           | ✅                  |
| **CORS**     | 基础         | 完整配置            |
| **日志**     | 简单         | 结构化              |
| **监控**     | ❌           | ✅ Health + Metrics |
| **部署**     | 手动         | Docker + K8s        |
| **生产就绪** | ❌           | ✅                  |

---

## 🛠️ 开发

### 本地运行

```bash
go run ./cmd/nomad-server
```

### 构建

```bash
go build -o nomad-server ./cmd/nomad-server
```

### 测试

```bash
go test ./server/...
```

---

## 📚 相关文档

- [架构设计](../SERVER_ARCHITECTURE.md) - 完整架构文档
- [核心 SDK](../docs/content/16.architecture/2.core-sdk.md) - pkg/ 设计
- [HTTP 层](../docs/content/16.architecture/3.http-layer.md) - 原 cmd/ 设计
- [客户端 SDK](../docs/content/16.architecture/4.client-sdk.md) - client-sdks 设计

---
