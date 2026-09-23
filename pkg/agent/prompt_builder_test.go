package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/nabob-ai/nomad/pkg/provider"
	"github.com/nabob-ai/nomad/pkg/sandbox"
	"github.com/nabob-ai/nomad/pkg/store"
	"github.com/nabob-ai/nomad/pkg/tools"
	"github.com/nabob-ai/nomad/pkg/tools/builtin"
	"github.com/nabob-ai/nomad/pkg/types"
)

func TestPromptBuilder_Basic(t *testing.T) {
	deps := setupPromptTestDeps(t)

	config := &types.AgentConfig{
		TemplateID: "test-template",
		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5",
			APIKey:   "test-key",
		},
		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindMock,
			WorkDir: "/tmp/test",
		},
	}

	ag, err := Create(context.Background(), config, deps)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	defer func() { _ = ag.Close() }()

	// 获取 System Prompt
	systemPrompt := ag.GetSystemPrompt()

	// 验证基础 Prompt 存在
	if !strings.Contains(systemPrompt, "You are a test assistant") {
		t.Error("System prompt should contain base prompt")
	}

	// 验证环境信息存在
	if !strings.Contains(systemPrompt, "## Environment Information") {
		t.Error("System prompt should contain environment information")
	}

	// 验证工具手册存在
	if !strings.Contains(systemPrompt, "## Tools Manual") {
		t.Error("System prompt should contain tools manual")
	}

	t.Logf("System Prompt length: %d", len(systemPrompt))
}

func TestPromptBuilder_CodeAssistant(t *testing.T) {
	deps := setupPromptTestDeps(t)

	// 注册代码助手模板
	deps.TemplateRegistry.Register(&types.AgentTemplateDefinition{
		ID:           "code-assistant",
		SystemPrompt: "You are a professional code assistant.",
		Tools:        []any{"Read", "Write"},
	})

	config := &types.AgentConfig{
		TemplateID: "code-assistant",
		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5",
			APIKey:   "test-key",
		},
		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindMock,
			WorkDir: "/tmp/test",
		},
		Metadata: map[string]any{
			"agent_type": "code_assistant",
		},
	}

	ag, err := Create(context.Background(), config, deps)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	defer func() { _ = ag.Close() }()

	// 获取 System Prompt
	systemPrompt := ag.GetSystemPrompt()

	// 验证代码引用规范存在
	if !strings.Contains(systemPrompt, "## Code References") {
		t.Error("Code assistant should have code reference guidelines")
	}

	if !strings.Contains(systemPrompt, "file_path:line_number") {
		t.Error("Code assistant should mention file_path:line_number format")
	}

	t.Logf("Code Assistant System Prompt:\n%s", systemPrompt)
}

func TestPromptBuilder_TodoReminder(t *testing.T) {
	deps := setupPromptTestDeps(t)

	// 注册带 Todo 提醒的模板
	deps.TemplateRegistry.Register(&types.AgentTemplateDefinition{
		ID:           "todo-assistant",
		SystemPrompt: "You are a task management assistant.",
		Tools:        []any{"Read", "Write"},
		Runtime: &types.AgentTemplateRuntime{
			Todo: &types.TodoConfig{
				Enabled:         true,
				ReminderOnStart: true,
			},
		},
	})

	config := &types.AgentConfig{
		TemplateID: "todo-assistant",
		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5",
			APIKey:   "test-key",
		},
		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindMock,
			WorkDir: "/tmp/test",
		},
	}

	ag, err := Create(context.Background(), config, deps)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	defer func() { _ = ag.Close() }()

	// 获取 System Prompt
	systemPrompt := ag.GetSystemPrompt()

	// 验证 Todo 提醒存在
	if !strings.Contains(systemPrompt, "## Task Management") {
		t.Error("Todo assistant should have task management section")
	}

	if !strings.Contains(systemPrompt, "TodoWrite") {
		t.Error("Todo assistant should mention TodoWrite tool")
	}

	t.Logf("Todo Assistant System Prompt:\n%s", systemPrompt)
}

func TestPromptBuilder_ToolsManualConfig(t *testing.T) {
	deps := setupPromptTestDeps(t)

	// 注册模板，只包含部分工具
	deps.TemplateRegistry.Register(&types.AgentTemplateDefinition{
		ID:           "selective-tools",
		SystemPrompt: "You are a selective assistant.",
		Tools:        []any{"Read", "Write"},
		Runtime: &types.AgentTemplateRuntime{
			ToolsManual: &types.ToolsManualConfig{
				Mode:    "listed",
				Include: []string{"Read"}, // 只包含 Read 工具
			},
		},
	})

	config := &types.AgentConfig{
		TemplateID: "selective-tools",
		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5",
			APIKey:   "test-key",
		},
		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindMock,
			WorkDir: "/tmp/test",
		},
	}

	ag, err := Create(context.Background(), config, deps)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	defer func() { _ = ag.Close() }()

	// 获取 System Prompt
	systemPrompt := ag.GetSystemPrompt()

	// 验证只包含 Read 工具
	if !strings.Contains(systemPrompt, "`Read`") {
		t.Error("System prompt should contain Read tool")
	}

	// 验证不包含 Write 工具
	if strings.Contains(systemPrompt, "`Write`") {
		t.Error("System prompt should not contain Write tool (excluded by config)")
	}

	t.Logf("Selective Tools System Prompt:\n%s", systemPrompt)
}

// setupPromptTestDeps 创建测试依赖
func setupPromptTestDeps(t *testing.T) *Dependencies {
	// 创建工具注册表
	toolRegistry := tools.NewRegistry()
	builtin.RegisterAll(toolRegistry)

	// 创建Sandbox工厂
	sandboxFactory := sandbox.NewFactory()

	// 创建Provider工厂
	providerFactory := &provider.AnthropicFactory{}

	// 创建Store (使用临时目录)
	jsonStore, err := store.NewJSONStore(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 创建模板注册表
	templateRegistry := NewTemplateRegistry()
	templateRegistry.Register(&types.AgentTemplateDefinition{
		ID:           "test-template",
		SystemPrompt: "You are a test assistant.",
		Model:        "claude-sonnet-4-5",
		Tools:        []any{"Read", "Write"},
	})

	return &Dependencies{
		Store:            jsonStore,
		SandboxFactory:   sandboxFactory,
		ToolRegistry:     toolRegistry,
		ProviderFactory:  providerFactory,
		TemplateRegistry: templateRegistry,
	}
}
