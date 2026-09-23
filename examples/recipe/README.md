# Recipe System 示例

本示例演示 Nomad 的 Recipe 系统，一种声明式的 Agent 配置方式。

## 功能特点

- 📖 YAML 格式的声明式配置
- 🔧 Builder 模式创建 Recipe
- 🔌 MCP 扩展集成
- 📝 参数化模板
- 🔐 权限模式配置
- 💾 Recipe 导入导出

## 运行示例

```bash
go run ./examples/recipe/
```

## Recipe 结构

```yaml
version: "1.0"
title: My Assistant
description: 助手描述

# 系统指令
instructions: |
  你是一个专业的助手...

# 初始提示
prompt: 请帮我完成任务

# 启用的工具
tools:
  - Read
  - Write
  - Bash

# MCP 扩展
extensions:
  - type: stdio
    name: github
    cmd: npx
    args: ["-y", "@anthropics/mcp-github"]

# 可配置参数
parameters:
  - key: project_name
    input_type: string
    requirement: required
    description: 项目名称

# 权限模式
permission_mode: smart_approve

# 建议活动
activities:
  - 审查代码
  - 写测试

# 作者信息
author:
  name: Your Name
  url: https://github.com/you
```

## 使用方式

### 1. 从 YAML 文件加载

```go
import "github.com/nabob-ai/nomad/pkg/recipe"

// 加载 Recipe
r, err := recipe.Load("./recipes/code-review.yaml")
if err != nil {
    log.Fatal(err)
}

// 验证 Recipe
if err := r.Validate(); err != nil {
    log.Printf("验证警告: %v", err)
}

// 使用 Recipe 创建 Agent 配置
agentConfig := r.ToAgentConfig()
```

### 2. 使用 Builder 创建

```go
r := recipe.NewBuilder().
    WithTitle("Code Review Assistant").
    WithDescription("代码审查助手").
    WithVersion("1.0").
    WithInstructions("你是专业的代码审查专家...").
    WithTools("Read", "List", "Search").
    WithPermissionMode(recipe.PermissionModeSmartApprove).
    Build()

// 保存为 YAML
r.Save("./my-recipe.yaml")
```

### 3. MCP 扩展配置

```yaml
extensions:
  # stdio 类型 - 启动外部进程
  - type: stdio
    name: github
    cmd: npx
    args: ["-y", "@anthropics/mcp-github"]
    env:
      GITHUB_TOKEN: "${GITHUB_TOKEN}"
    timeout: 30

  # sse 类型 - 连接远程服务
  - type: sse
    name: search
    url: http://localhost:3000/mcp
    timeout: 10

  # builtin 类型 - 内置扩展
  - type: builtin
    name: filesystem
```

### 4. 参数化模板

```yaml
parameters:
  - key: language
    input_type: select
    requirement: required
    description: 编程语言
    default: go
    options:
      - go
      - python
      - typescript

prompt: |
  请为我创建一个 {{language}} 项目。
```

支持的参数类型：

| 类型 | 说明 | 示例 |
|------|------|------|
| `string` | 文本输入 | 项目名称 |
| `number` | 数字输入 | 并发数 |
| `boolean` | 布尔开关 | 是否启用测试 |
| `select` | 下拉选择 | 编程语言 |
| `file` | 文件选择 | 配置文件路径 |

### 5. 权限模式

```yaml
# 自动批准所有操作
permission_mode: auto_approve

# 智能审批（推荐）
permission_mode: smart_approve

# 所有操作都需确认
permission_mode: always_ask
```

## 与 Agent 集成

```go
import (
    "github.com/nabob-ai/nomad/pkg/agent"
    "github.com/nabob-ai/nomad/pkg/recipe"
)

// 加载 Recipe
r, _ := recipe.Load("./my-recipe.yaml")

// 转换为 Agent 配置
agentConfig := r.ToAgentConfig()
agentConfig.ModelConfig = &types.ModelConfig{
    Provider: "anthropic",
    Model:    "claude-sonnet-4-5",
    APIKey:   os.Getenv("ANTHROPIC_API_KEY"),
}

// 创建 Agent
ag, err := agent.Create(ctx, agentConfig, deps)
```

## 内置 Recipe 示例

查看 `examples/recipes/` 目录获取更多示例：

- `code-review.yaml` - 代码审查助手
- `writing-assistant.yaml` - 写作助手

## Recipe 分享

Recipe 可以通过以下方式分享：

1. **文件分享** - 直接分享 YAML 文件
2. **Git 仓库** - 创建 Recipe 集合仓库
3. **URL 加载** - 从远程 URL 加载 Recipe

```go
// 从 URL 加载
r, err := recipe.LoadFromURL("https://example.com/recipes/my-recipe.yaml")
```

## 相关示例

- [permission](../permission/) - 权限系统详解
- [mcp](../mcp/) - MCP 扩展详解
- [agent](../agent/) - Agent 基础示例
