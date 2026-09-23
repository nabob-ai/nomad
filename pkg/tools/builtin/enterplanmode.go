package builtin

import (
	"context"
	"fmt"
	"time"

	"github.com/nabob-ai/nomad/pkg/tools"
)

// PlanModeManagerInterface Plan 模式管理器接口
// 用于工具与 Agent 的 PlanModeManager 交互
type PlanModeManagerInterface interface {
	EnterPlanMode(planID, planFilePath, reason string)
	ExitPlanMode()
	IsActive() bool
}

// EnterPlanModeTool 进入规划模式工具
// 用于复杂任务的规划阶段，在此模式下只允许只读操作和计划文件写入
type EnterPlanModeTool struct {
	planFileManager *PlanFileManager
}

// EnterPlanModeResult 进入规划模式的结果
type EnterPlanModeResult struct {
	OK           bool     `json:"ok"`
	PlanFilePath string   `json:"plan_file_path"`
	PlanID       string   `json:"plan_id"`
	AllowedTools []string `json:"allowed_tools"`
	Workflow     string   `json:"workflow"`
	Message      string   `json:"message"`
}

// NewEnterPlanModeTool 创建EnterPlanMode工具
func NewEnterPlanModeTool(config map[string]any) (tools.Tool, error) {
	basePath := ".plans" // 默认使用相对路径，会在工作目录下创建
	if bp, ok := config["base_path"].(string); ok && bp != "" {
		basePath = bp
	}

	return &EnterPlanModeTool{
		planFileManager: NewPlanFileManager(basePath),
	}, nil
}

func (t *EnterPlanModeTool) Name() string {
	return "EnterPlanMode"
}

func (t *EnterPlanModeTool) Description() string {
	return "进入规划模式，用于复杂任务的规划阶段。在此模式下只允许只读操作和计划文件写入。"
}

func (t *EnterPlanModeTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"reason": map[string]any{
				"type":        "string",
				"description": "进入规划模式的原因说明",
			},
		},
		"required": []string{},
	}
}

func (t *EnterPlanModeTool) Execute(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
	reason := GetStringParam(input, "reason", "")

	// 获取工作目录，动态设置计划文件存储路径
	// 这样计划文件会存储在 {workDir}/.plans/ 下，支持按项目隔离
	planManager := t.planFileManager
	if tc != nil && tc.Sandbox != nil {
		workDir := tc.Sandbox.WorkDir()
		if workDir != "" {
			// 创建新的 PlanFileManager，使用工作目录下的 .plans 子目录
			planManager = NewPlanFileManagerWithProject(workDir+"/.plans", "")
		}
	}

	// 确保计划目录存在
	if err := planManager.EnsureDir(); err != nil {
		return NewClaudeErrorResponse(fmt.Errorf("failed to create plans directory: %w", err)), nil
	}

	// 生成计划文件路径
	planPath := planManager.GeneratePath()
	planID := planManager.GenerateID()

	// 创建初始计划文件
	initialContent := fmt.Sprintf(`# 实施计划

> 计划 ID: %s
> 创建时间: %s
> 状态: 规划中

## 概述

[在此描述任务目标和实施方案]

## 执行步骤

1. [步骤 1]
2. [步骤 2]
3. [步骤 3]

## 关键文件

| 文件 | 用途 |
|------|------|
| | |

## 风险与应对

-

## 成功标准

-

---
*此计划文件将随着规划进展持续更新。*
`, planID, time.Now().Format(time.RFC3339))

	if err := planManager.Save(planPath, initialContent); err != nil {
		return NewClaudeErrorResponse(fmt.Errorf("failed to create plan file: %w", err)), nil
	}

	// 激活 Agent 级别的 Plan 模式约束
	if tc != nil && tc.Services != nil {
		if pmm, ok := tc.Services["plan_mode_manager"].(PlanModeManagerInterface); ok {
			pmm.EnterPlanMode(planID, planPath, reason)
		}
	}

	// 定义规划模式允许的工具
	allowedTools := []string{
		"Read",
		"Glob",
		"Grep",
		"WebFetch",
		"WebSearch",
		"AskUserQuestion",
		"Write", // 仅限计划文件
		"ExitPlanMode",
		"Task", // 仅限 Explore 子代理
	}

	workflow := `## 规划模式工作流程

### 阶段 1：初步了解
- 启动 Explore 代理了解代码库
- 向用户提问以澄清需求

### 阶段 2：方案设计
- 设计实施方案
- 识别关键文件和依赖关系

### 阶段 3：综合分析
- 收集发现并创建综合计划
- 向用户询问权衡和偏好

### 阶段 4：最终计划
- 更新计划文件，包含最终建议
- 包括实施步骤和成功标准

### 阶段 5：退出
- 规划完成后调用 ExitPlanMode
- 等待用户批准后再开始实施

**重要提示**：在规划模式下，你只能：
- 读取文件（Read, Glob, Grep）
- 搜索网页（WebFetch, WebSearch）
- 向用户提问（AskUserQuestion）
- 写入计划文件（Write - 仅限：` + planPath + `）
- 启动 Explore 子代理（Task - 仅限 Explore 类型）

在计划获得批准之前，你不能编辑代码、运行命令或进行任何修改。`

	// 构建响应
	result := map[string]any{
		"ok":                   true,
		"plan_file_path":       planPath,
		"plan_id":              planID,
		"allowed_tools":        allowedTools,
		"workflow":             workflow,
		"message":              "已进入规划模式。现在可以探索代码库并创建计划。",
		"created_at":           time.Now().Unix(),
		"plan_mode_enforced":   true,
		"agent_level_enforced": true,
	}

	if reason != "" {
		result["reason"] = reason
	}

	// 添加提示信息
	result["next_steps"] = []string{
		"1. 阅读代码库，理解现有模式",
		"2. 向用户提问，澄清需求",
		"3. 更新计划文件，记录发现",
		"4. 调用 ExitPlanMode 请求用户审批",
	}

	result["constraints"] = map[string]any{
		"read_only":             true,
		"plan_file_writable":    planPath,
		"code_editing_disabled": true,
		"bash_disabled":         true,
		"task_explore_only":     true,
	}

	return result, nil
}

func (t *EnterPlanModeTool) Prompt() string {
	return `进入规划模式，用于复杂任务的规划阶段。

## 何时使用此工具

当遇到以下情况时应使用此工具：
1. **多种有效方案** - 任务可以用多种方式解决，需要权衡
2. **重大架构决策** - 需要在架构模式之间做选择
3. **大规模变更** - 任务涉及多个文件或系统
4. **需求不明确** - 需要先探索才能理解完整范围
5. **需要用户确认** - 需要向用户提问以确认方向

## 规划模式的限制

在规划模式下，你只能：
- 读取文件（Read, Glob, Grep）
- 搜索网页（WebFetch, WebSearch）
- 向用户提问（AskUserQuestion）
- 写入计划文件（Write - 仅限指定的计划文件）

你不能：
- 编辑代码文件
- 运行 Bash 命令
- 创建或删除文件（除计划文件外）

## 工作流程

1. 调用此工具进入规划模式
2. 探索代码库，理解现有模式
3. 向用户提问以澄清需求
4. 更新计划文件记录发现和方案
5. 调用 ExitPlanMode 请求用户批准
6. 用户批准后开始实施

## 示例

{
  "reason": "需要设计用户认证系统，有多种方案可选"
}

## 注意事项

- 规划模式需要用户批准才能退出
- 计划文件存储在 .nomad/plans/ 目录
- 完成规划后必须调用 ExitPlanMode`
}
