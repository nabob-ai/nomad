---
title: Human-in-the-Loop (HITL)
description: 实现人工审核和控制敏感 Agent 操作
navigation:
  icon: i-lucide-shield-check
---

# Human-in-the-Loop (HITL) 完整指南

人工在环（Human-in-the-Loop，简称 HITL）是一种关键的 AI 安全机制，允许人类在 Agent 执行敏感操作前进行审核、批准或修改。

## 📖 概述

### 什么是 HITL？

HITL 是一种设计模式，在自动化流程中引入人工决策点，确保：

- **安全性**：防止 Agent 执行危险操作
- **合规性**：满足审计和合规要求
- **可控性**：保持人类对关键决策的控制权
- **透明性**：让用户了解并掌控 Agent 行为

### 适用场景

| 场景     | 示例                   | 风险等级 |
| -------- | ---------------------- | -------- |
| 文件操作 | 删除文件、修改配置     | 🔴 高    |
| 系统命令 | 执行 Shell 命令        | 🔴 高    |
| 外部 API | 发送邮件、调用付费 API | 🟡 中    |
| 数据修改 | 更新数据库、修改记录   | 🟡 中    |
| 资源消耗 | 大规模计算、批量处理   | 🟡 中    |

## 🚀 快速开始

### 基本用法

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/nabob-ai/nomad/pkg/agent"
    "github.com/nabob-ai/nomad/pkg/middleware"
)

func main() {
    ctx := context.Background()

    // 1. 创建 HITL 中间件
    hitlMW, err := middleware.NewHumanInTheLoopMiddleware(&middleware.HumanInTheLoopMiddlewareConfig{
        // 配置需要审核的工具
        InterruptOn: map[string]interface{}{
            "Bash":     true,  // 启用默认审核
            "fs_delete":    true,  // 文件删除需要审核
            "HttpRequest": true,  // HTTP 请求需要审核
        },

        // 审核处理器
        ApprovalHandler: func(ctx context.Context, req *middleware.ReviewRequest) ([]middleware.Decision, error) {
            // 显示待审核的操作
            for _, action := range req.ActionRequests {
                fmt.Printf("\n🚨 需要审核的操作:\n")
                fmt.Printf("  工具: %s\n", action.ToolName)
                fmt.Printf("  参数: %+v\n", action.Input)
                fmt.Printf("  说明: %s\n\n", action.Message)

                // 获取人工决策
                fmt.Print("请选择操作 (approve/reject): ")
                var choice string
                fmt.Scanln(&choice)

                if choice == "approve" || choice == "y" {
                    return []middleware.Decision{{
                        Type:   middleware.DecisionApprove,
                        Reason: "用户批准",
                    }}, nil
                }

                return []middleware.Decision{{
                    Type:   middleware.DecisionReject,
                    Reason: "用户拒绝",
                }}, nil
            }

            return nil, fmt.Errorf("no decision")
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // 2. 注册中间件
    stack := middleware.NewStack()
    stack.Use(hitlMW)

    // 3. 创建 Agent
    ag, err := agent.Create(ctx, config, &agent.Dependencies{
        MiddlewareStack: stack,
        // ... 其他依赖
    })
    if err != nil {
        log.Fatal(err)
    }
    defer ag.Close()

    // 4. 使用 Agent（敏感操作会触发审核）
    result, _ := ag.Chat(ctx, "请删除 /tmp/test.txt 文件")
    fmt.Println(result.Text)
}
```

## ⚙️ 配置详解

### InterruptOn 配置

`InterruptOn` 字段定义哪些工具需要审核，支持三种配置方式：

#### 1. 布尔值配置（简单模式）

```go
InterruptOn: map[string]interface{}{
    "Bash":     true,   // 启用默认审核
    "HttpRequest": false,  // 不需要审核
}
```

#### 2. 详细配置对象

```go
InterruptOn: map[string]interface{}{
    "Write": map[string]interface{}{
        "message": "文件写入操作需要审核，请确认参数正确",
        "allowed_decisions": []string{"approve", "reject", "edit"},
    },
}
```

#### 3. 使用 InterruptConfig 结构体

```go
InterruptOn: map[string]interface{}{
    "database_update": &middleware.InterruptConfig{
        Enabled:          true,
        Message:          "数据库更新操作需要审核",
        AllowedDecisions: []middleware.DecisionType{
            middleware.DecisionApprove,
            middleware.DecisionReject,
            middleware.DecisionEdit,
        },
    },
}
```

### 决策类型

HITL 支持三种决策类型：

| 决策类型          | 说明                 | 使用场景           |
| ----------------- | -------------------- | ------------------ |
| `DecisionApprove` | 批准：按原参数执行   | 操作合理，可以执行 |
| `DecisionReject`  | 拒绝：取消执行       | 操作不安全或不合理 |
| `DecisionEdit`    | 编辑：修改参数后执行 | 参数需要调整       |

```go
const (
    DecisionApprove DecisionType = "approve"  // 批准执行
    DecisionReject  DecisionType = "reject"   // 拒绝执行
    DecisionEdit    DecisionType = "edit"     // 编辑参数后执行
)
```

### 审核处理器实现

#### 命令行交互

```go
ApprovalHandler: func(ctx context.Context, req *middleware.ReviewRequest) ([]middleware.Decision, error) {
    for _, action := range req.ActionRequests {
        fmt.Printf("工具: %s\n参数: %+v\n", action.ToolName, action.Input)
        fmt.Print("批准? (y/n): ")

        var answer string
        fmt.Scanln(&answer)

        if answer == "y" {
            return []middleware.Decision{{Type: middleware.DecisionApprove}}, nil
        }
        return []middleware.Decision{{Type: middleware.DecisionReject}}, nil
    }
    return nil, fmt.Errorf("no decision")
}
```

#### 参数编辑示例

```go
ApprovalHandler: func(ctx context.Context, req *middleware.ReviewRequest) ([]middleware.Decision, error) {
    action := req.ActionRequests[0]

    fmt.Printf("工具: %s\n", action.ToolName)
    fmt.Printf("当前参数: %+v\n", action.Input)
    fmt.Print("选择操作 (approve/reject/edit): ")

    var choice string
    fmt.Scanln(&choice)

    switch choice {
    case "approve":
        return []middleware.Decision{{Type: middleware.DecisionApprove}}, nil

    case "reject":
        return []middleware.Decision{{Type: middleware.DecisionReject}}, nil

    case "edit":
        editedInput := make(map[string]interface{})
        for key, value := range action.Input {
            fmt.Printf("编辑 %s (当前: %v, 按回车保持不变): ", key, value)
            var newValue string
            fmt.Scanln(&newValue)
            if newValue != "" {
                editedInput[key] = newValue
            } else {
                editedInput[key] = value
            }
        }

        return []middleware.Decision{{
            Type:        middleware.DecisionEdit,
            EditedInput: editedInput,
        }}, nil
    }

    return nil, fmt.Errorf("invalid choice")
}
```

## 📝 实战示例

### 示例 1: 文件操作审核

创建一个针对文件操作的 HITL 中间件，自动检测危险路径：

```go
func NewFileOperationHITL() (*middleware.HumanInTheLoopMiddleware, error) {
    return middleware.NewHumanInTheLoopMiddleware(&middleware.HumanInTheLoopMiddlewareConfig{
        InterruptOn: map[string]interface{}{
            "fs_delete": true,
            "Write":  true,
        },
        ApprovalHandler: func(ctx context.Context, req *middleware.ReviewRequest) ([]middleware.Decision, error) {
            action := req.ActionRequests[0]

            // 检查危险路径
            if path, ok := action.Input["path"].(string); ok {
                dangerousPaths := []string{"/", "/etc", "/usr", "/bin", "/home"}
                for _, dp := range dangerousPaths {
                    if strings.HasPrefix(path, dp) {
                        fmt.Printf("⛔ 危险路径: %s\n", path)
                        fmt.Println("自动拒绝删除系统重要目录")
                        return []middleware.Decision{{
                            Type:   middleware.DecisionReject,
                            Reason: fmt.Sprintf("拒绝操作系统重要路径: %s", path),
                        }}, nil
                    }
                }
            }

            // 正常审核流程
            fmt.Printf("操作: %s\n参数: %+v\n", action.ToolName, action.Input)
            fmt.Print("批准? (y/n): ")

            var answer string
            fmt.Scanln(&answer)

            if answer == "y" {
                return []middleware.Decision{{Type: middleware.DecisionApprove}}, nil
            }
            return []middleware.Decision{{Type: middleware.DecisionReject}}, nil
        },
    })
}
```

### 示例 2: 基于风险级别的智能审核

```go
type RiskLevel int

const (
    RiskLow    RiskLevel = 1  // 低风险，自动批准
    RiskMedium RiskLevel = 2  // 中风险，需要审核
    RiskHigh   RiskLevel = 3  // 高风险，严格审核
)

func assessRisk(action middleware.ActionRequest) RiskLevel {
    if action.ToolName == "Bash" {
        if cmd, ok := action.Input["command"].(string); ok {
            if strings.Contains(cmd, "rm -rf") || strings.Contains(cmd, "mkfs") {
                return RiskHigh
            }
            if strings.Contains(cmd, "rm ") || strings.Contains(cmd, "chmod") {
                return RiskMedium
            }
        }
    }
    return RiskLow
}

func smartApprovalHandler(ctx context.Context, req *middleware.ReviewRequest) ([]middleware.Decision, error) {
    action := req.ActionRequests[0]
    risk := assessRisk(action)

    switch risk {
    case RiskLow:
        fmt.Printf("✅ 自动批准低风险操作: %s\n", action.ToolName)
        return []middleware.Decision{{
            Type:   middleware.DecisionApprove,
            Reason: "低风险操作自动批准",
        }}, nil

    case RiskMedium:
        fmt.Printf("⚠️  中风险操作: %s\n", action.ToolName)
        fmt.Print("是否批准? (y/n): ")
        var answer string
        fmt.Scanln(&answer)
        if answer == "y" {
            return []middleware.Decision{{Type: middleware.DecisionApprove}}, nil
        }
        return []middleware.Decision{{Type: middleware.DecisionReject}}, nil

    case RiskHigh:
        fmt.Printf("🚨 高风险操作: %s\n", action.ToolName)
        fmt.Printf("参数: %+v\n", action.Input)
        fmt.Print("输入 'CONFIRM' 确认执行: ")
        var confirm string
        fmt.Scanln(&confirm)
        if confirm == "CONFIRM" {
            return []middleware.Decision{{Type: middleware.DecisionApprove}}, nil
        }
        return []middleware.Decision{{Type: middleware.DecisionReject}}, nil
    }

    return nil, fmt.Errorf("unknown risk level")
}
```

## 🎯 最佳实践

### 1. 只审核敏感操作

```go
// ✅ 推荐：只审核敏感操作
InterruptOn: map[string]interface{}{
    "Bash":     true,  // Shell 命令
    "fs_delete":    true,  // 文件删除
    "api_payment":  true,  // 付费 API
}

// ❌ 不推荐：审核所有操作
InterruptOn: map[string]interface{}{
    "Read":   true,  // 读取通常安全
    "calculate": true,  // 计算无副作用
}
```

### 2. 提供清晰的审核信息

```go
InterruptOn: map[string]interface{}{
    "email_send": map[string]interface{}{
        "message": "📧 邮件发送审核\n请确认收件人和内容是否正确",
    },
}
```

### 3. 实现超时机制

```go
ApprovalHandler: func(ctx context.Context, req *middleware.ReviewRequest) ([]middleware.Decision, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()

    decisionCh := make(chan middleware.Decision, 1)
    go func() {
        decisionCh <- getUserDecision(req)
    }()

    select {
    case decision := <-decisionCh:
        return []middleware.Decision{decision}, nil
    case <-ctx.Done():
        return nil, fmt.Errorf("审核超时")
    }
}
```

### 4. 记录审核日志

```go
type AuditLog struct {
    Timestamp time.Time
    Tool      string
    Input     map[string]interface{}
    Decision  string
    Reason    string
}

func logApproval(action middleware.ActionRequest, decision middleware.Decision) {
    log := AuditLog{
        Timestamp: time.Now(),
        Tool:      action.ToolName,
        Input:     action.Input,
        Decision:  string(decision.Type),
        Reason:    decision.Reason,
    }
    // 保存到文件或数据库
    saveToAuditLog(log)
}
```

### 5. 支持批量审核

对于多个操作，可以一次性审核：

```go
ApprovalHandler: func(ctx context.Context, req *middleware.ReviewRequest) ([]middleware.Decision, error) {
    fmt.Printf("有 %d 个操作需要审核:\n", len(req.ActionRequests))

    for i, action := range req.ActionRequests {
        fmt.Printf("%d. %s: %+v\n", i+1, action.ToolName, action.Input)
    }

    fmt.Print("批准所有? (y/n/详细审核): ")
    var choice string
    fmt.Scanln(&choice)

    if choice == "y" {
        // 批准所有
        decisions := make([]middleware.Decision, len(req.ActionRequests))
        for i := range decisions {
            decisions[i] = middleware.Decision{Type: middleware.DecisionApprove}
        }
        return decisions, nil
    }

    // 逐个审核...
    return nil, nil
}
```

## 🔒 安全建议

### 1. 默认拒绝策略

当无法获取决策时，默认拒绝而不是批准：

```go
ApprovalHandler: func(ctx context.Context, req *middleware.ReviewRequest) ([]middleware.Decision, error) {
    decision, err := tryGetDecision(ctx, req)
    if err != nil {
        // 默认拒绝
        return []middleware.Decision{{
            Type:   middleware.DecisionReject,
            Reason: "无法获取审核决策，默认拒绝",
        }}, nil
    }
    return []middleware.Decision{decision}, nil
}
```

### 2. 审核权限控制

实现基于角色的审核权限：

```go
type ApprovalPolicy struct {
    AllowedRoles map[string][]string  // tool -> roles
}

func (p *ApprovalPolicy) CheckPermission(user string, toolName string) bool {
    roles := getUserRoles(user)
    allowedRoles := p.AllowedRoles[toolName]

    for _, role := range roles {
        for _, allowed := range allowedRoles {
            if role == allowed {
                return true
            }
        }
    }
    return false
}
```

### 3. 审核记录不可篡改

使用只追加的日志系统记录所有审核决策：

```go
func appendOnlyAuditLog(entry AuditLog) {
    f, _ := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()

    data, _ := json.Marshal(entry)
    f.WriteString(string(data) + "\n")
}
```

## 📚 进阶主题

### 与 System Prompt 配合

可以在 System Prompt 中告知 Agent 哪些操作需要审核：

```go
const HITL_SYSTEM_PROMPT = `## Human-in-the-Loop (HITL)

某些敏感操作需要人工批准才能执行。当你调用这些工具时:

1. **系统会暂停执行**，等待人工审核
2. **人工审核员可以**：
   - 批准(approve): 执行操作
   - 拒绝(reject): 取消操作
   - 编辑(edit): 修改参数后执行

3. **如果操作被拒绝**：
   - 你会收到拒绝原因
   - 可以调整策略或尝试其他方法
   - 不要重复尝试被拒绝的操作

4. **最佳实践**：
   - 清楚解释为什么需要执行该操作
   - 提供足够的上下文帮助审核
   - 尊重人工审核决策`

// 将这段提示添加到 Agent 配置中
config.SystemPrompt = basePrompt + "\n\n" + middleware.HITL_SYSTEM_PROMPT
```

### 动态调整审核策略

根据运行时状态动态调整哪些操作需要审核：

```go
type DynamicHITL struct {
    *middleware.HumanInTheLoopMiddleware
    policy *ApprovalPolicy
}

func (d *DynamicHITL) UpdatePolicy(toolName string, requiresApproval bool) {
    d.policy.Update(toolName, requiresApproval)
}

// 使用场景：在测试环境禁用审核，生产环境启用
if isProduction {
    hitl.UpdatePolicy("Bash", true)
} else {
    hitl.UpdatePolicy("Bash", false)
}
```

## 🔗 相关资源

- [中间件系统概览](/middleware) - 了解中间件架构
- [内置中间件文档](/middleware/builtin) - 所有内置中间件
- [安全最佳实践](/best-practices/security) - Agent 安全指南
- [HITL 实现代码](https://github.com/nabob-ai/nomad/blob/main/pkg/middleware/hitl.go)

## ❓ 常见问题

### Q: HITL 会影响性能吗？

A: 是的，HITL 会引入延迟，因为需要等待人工决策。建议：

- 只对真正敏感的操作启用审核
- 实现超时机制避免无限等待
- 对低风险操作使用自动批准

### Q: 如何在 Agent 被拒绝后恢复？

A: Agent 会收到拒绝信息，可以：

- 尝试其他方法完成任务
- 向用户解释原因并请求指导
- 调整参数后重试

### Q: 支持异步审核吗？

A: 支持。可以通过消息队列或数据库实现异步审核，让 Agent 等待决策。

### Q: 如何处理批量操作审核？

A: 可以在 `ApprovalHandler` 中一次性处理所有 `ActionRequests`，返回对应的决策列表。
