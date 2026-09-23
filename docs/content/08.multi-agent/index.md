---
title: 多Agent系统总览
description: 构建多个 Agent 协作的复杂系统
navigation: false
---

# 多Agent系统总览

Nomad 提供完整的多 Agent 系统架构，从单 Agent 内部的任务委派，到多 Agent 协作，再到统一运行时管理。

## 📚 核心组件

### [SubAgent 中间件](/middleware/subagent)

**任务委派机制** - Agent 内部的复杂任务拆分

- 通过 `task` 工具调用
- 上下文隔离、并行执行
- Token 优化、短生命周期

### [Room](/multi-agent/room)

**协作空间** - 多 Agent 之间的消息通信

- 成员管理（Join/Leave）
- 消息路由（@mention 支持）
- 广播和点对点通信

### [Pool](/multi-agent/pool)

**生命周期管理器** - 统一管理所有 Agent 实例

- 创建、获取、删除 Agent
- 资源池管理和容量控制
- 监控和统计

### [NomadOS](/multi-agent/nomados)

**统一运行时系统** - 对外提供 Agent 服务

- 自动生成 REST API
- 多接口支持（HTTP/A2A/AGUI）
- 资源注册和管理

### [Scheduler](/multi-agent/scheduler)

**任务调度器** - 智能分配任务给合适的 Agent

- 基于能力的任务路由
- 负载均衡
- 优先级调度

## 🚀 快速开始

```go
import (
    "github.com/nabob-ai/nomad/pkg/core"
)

// 创建 Pool
pool := core.NewPool(&core.PoolOptions{
    Dependencies: deps,
    MaxAgents:    50,
})

// 创建 Agents
pool.Create(ctx, &types.AgentConfig{
    AgentID:    "agent-1",
    TemplateID: "assistant",
})

pool.Create(ctx, &types.AgentConfig{
    AgentID:    "agent-2",
    TemplateID: "assistant",
})

// 创建 Room 协作空间
room := core.NewRoom(pool)
room.Join("alice", "agent-1")
room.Join("bob", "agent-2")

// 发送消息
room.Say(ctx, "alice", "Hello @bob!")
room.Broadcast(ctx, "Meeting starts now")
```

## 🎯 概念对比

| 概念     | 层次     | 生命周期 | 主要职责                |
| -------- | -------- | -------- | ----------------------- |
| SubAgent | 中间件层 | 任务级   | 任务委派、上下文隔离    |
| Room     | 协作层   | 会话级   | 多 Agent 消息通信、路由 |
| Pool     | 管理层   | 应用级   | Agent 生命周期管理      |
| NomadOS  | 运行时层 | 系统级   | 统一运行时、API 网关    |

详细对比请查看 [多Agent概念对比](/multi-agent/comparison)

## 📖 相关文档

- [SubAgent 中间件](/middleware/subagent)
- [Room 协作空间](/multi-agent/room)
- [Pool 生命周期管理](/multi-agent/pool)
- [NomadOS 统一运行时](/multi-agent/nomados)
- [多Agent概念对比](/multi-agent/comparison)
- [多Agent示例](/examples/multi-agent)
