package agent

import (
	"testing"
)

func TestPlanModeManager_EnterAndExit(t *testing.T) {
	pmm := NewPlanModeManager()

	// 初始状态应该是非活跃
	if pmm.IsActive() {
		t.Error("expected plan mode to be inactive initially")
	}

	// 进入 Plan 模式
	pmm.EnterPlanMode("plan-123", ".nomad/plans/test.md", "testing")

	if !pmm.IsActive() {
		t.Error("expected plan mode to be active after EnterPlanMode")
	}

	state := pmm.GetState()
	if state.PlanID != "plan-123" {
		t.Errorf("expected plan ID 'plan-123', got '%s'", state.PlanID)
	}
	if state.PlanFilePath != ".nomad/plans/test.md" {
		t.Errorf("expected plan file path '.nomad/plans/test.md', got '%s'", state.PlanFilePath)
	}
	if state.Reason != "testing" {
		t.Errorf("expected reason 'testing', got '%s'", state.Reason)
	}

	// 退出 Plan 模式
	pmm.ExitPlanMode()

	if pmm.IsActive() {
		t.Error("expected plan mode to be inactive after ExitPlanMode")
	}
}

func TestPlanModeManager_ValidateToolCall_NotInPlanMode(t *testing.T) {
	pmm := NewPlanModeManager()

	// 不在 Plan 模式时，所有工具都应该被允许
	allowed, reason := pmm.ValidateToolCall("Bash", nil)
	if !allowed {
		t.Errorf("expected Bash to be allowed when not in plan mode, reason: %s", reason)
	}

	allowed, reason = pmm.ValidateToolCall("Edit", nil)
	if !allowed {
		t.Errorf("expected Edit to be allowed when not in plan mode, reason: %s", reason)
	}
}

func TestPlanModeManager_ValidateToolCall_InPlanMode(t *testing.T) {
	pmm := NewPlanModeManager()
	pmm.EnterPlanMode("plan-123", ".nomad/plans/test.md", "testing")

	// 允许的工具
	allowedTools := []string{"Read", "Glob", "Grep", "WebFetch", "WebSearch", "AskUserQuestion", "ExitPlanMode"}
	for _, tool := range allowedTools {
		allowed, reason := pmm.ValidateToolCall(tool, nil)
		if !allowed {
			t.Errorf("expected %s to be allowed in plan mode, reason: %s", tool, reason)
		}
	}

	// 禁止的工具
	disallowedTools := []string{"Bash", "Edit", "Delete", "Move", "Copy"}
	for _, tool := range disallowedTools {
		allowed, _ := pmm.ValidateToolCall(tool, nil)
		if allowed {
			t.Errorf("expected %s to be disallowed in plan mode", tool)
		}
	}
}

func TestPlanModeManager_ValidateWriteCall(t *testing.T) {
	pmm := NewPlanModeManager()
	pmm.EnterPlanMode("plan-123", ".nomad/plans/test.md", "testing")

	// 允许写入计划文件
	allowed, reason := pmm.ValidateToolCall("Write", map[string]any{
		"path": ".nomad/plans/test.md",
	})
	if !allowed {
		t.Errorf("expected Write to plan file to be allowed, reason: %s", reason)
	}

	// 允许写入 .nomad/plans/ 目录下的任何文件
	allowed, reason = pmm.ValidateToolCall("Write", map[string]any{
		"path": ".nomad/plans/another.md",
	})
	if !allowed {
		t.Errorf("expected Write to .nomad/plans/ to be allowed, reason: %s", reason)
	}

	// 禁止写入其他文件
	allowed, _ = pmm.ValidateToolCall("Write", map[string]any{
		"path": "src/main.go",
	})
	if allowed {
		t.Error("expected Write to src/main.go to be disallowed in plan mode")
	}
}

func TestPlanModeManager_ValidateTaskCall(t *testing.T) {
	pmm := NewPlanModeManager()
	pmm.EnterPlanMode("plan-123", ".nomad/plans/test.md", "testing")

	// 允许 Explore 子代理
	allowed, reason := pmm.ValidateToolCall("Task", map[string]any{
		"subagent_type": "Explore",
	})
	if !allowed {
		t.Errorf("expected Task with Explore to be allowed, reason: %s", reason)
	}

	// 禁止其他类型的子代理
	allowed, _ = pmm.ValidateToolCall("Task", map[string]any{
		"subagent_type": "Plan",
	})
	if allowed {
		t.Error("expected Task with Plan to be disallowed in plan mode")
	}

	allowed, _ = pmm.ValidateToolCall("Task", map[string]any{
		"subagent_type": "general-purpose",
	})
	if allowed {
		t.Error("expected Task with general-purpose to be disallowed in plan mode")
	}
}

func TestPlanModeManager_GetState_ReturnsCopy(t *testing.T) {
	pmm := NewPlanModeManager()
	pmm.EnterPlanMode("plan-123", ".nomad/plans/test.md", "testing")

	state1 := pmm.GetState()
	state2 := pmm.GetState()

	// 修改 state1 不应该影响 state2
	state1.PlanID = "modified"

	if state2.PlanID == "modified" {
		t.Error("GetState should return a copy, not a reference")
	}
}
