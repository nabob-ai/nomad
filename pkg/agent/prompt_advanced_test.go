package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/nabob-ai/nomad/pkg/types"
)

func TestPromptBuilder_AllModules(t *testing.T) {
	deps := setupPromptTestDeps(t)

	// 注册完整功能的模板
	deps.TemplateRegistry.Register(&types.AgentTemplateDefinition{
		ID:           "full-featured",
		SystemPrompt: "You are a full-featured assistant.",
		Tools:        []any{"Read", "Write", "Bash"},
		Runtime: &types.AgentTemplateRuntime{
			Todo: &types.TodoConfig{
				Enabled:         true,
				ReminderOnStart: true,
			},
		},
	})

	config := &types.AgentConfig{
		TemplateID: "full-featured",
		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5",
			APIKey:   "test-key",
		},
		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindMock,
			WorkDir: "/tmp/test",
		},
		Context: &types.ContextManagerOptions{
			MaxTokens: 100000,
		},
		Metadata: map[string]any{
			"agent_type":               "code_assistant",
			"show_capabilities":        true,
			"show_limitations":         true,
			"enable_security":          true,
			"enable_performance_hints": true,
			"custom_instructions":      "Always write tests for your code.",
		},
	}

	ag, err := Create(context.Background(), config, deps)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	defer func() { _ = ag.Close() }()

	systemPrompt := ag.GetSystemPrompt()

	// 验证所有模块都被注入
	expectedSections := []string{
		"You are a full-featured assistant",
		"## Your Capabilities",
		"## Environment Information",
		"## Sandbox Environment",
		"## Tools Manual",
		"## Task Management",
		"## Code References",
		"## Security Guidelines",
		"## Performance Optimization",
		"## Custom Instructions",
		"## Important Limitations",
		"## Context Window Management",
	}

	for _, section := range expectedSections {
		if !strings.Contains(systemPrompt, section) {
			t.Errorf("System prompt should contain: %s", section)
		}
	}

	t.Logf("Full-featured System Prompt length: %d", len(systemPrompt))
}

func TestPromptBuilder_RoomCollaboration(t *testing.T) {
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
		Metadata: map[string]any{
			"room_id":           "room-123",
			"room_member_count": 3,
			"room_members":      []string{"alice", "bob", "charlie"},
		},
	}

	ag, err := Create(context.Background(), config, deps)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	defer func() { _ = ag.Close() }()

	systemPrompt := ag.GetSystemPrompt()

	// 验证协作模块
	if !strings.Contains(systemPrompt, "## Multi-Agent Collaboration") {
		t.Error("System prompt should contain collaboration section")
	}

	if !strings.Contains(systemPrompt, "room-123") {
		t.Error("System prompt should contain room ID")
	}

	if !strings.Contains(systemPrompt, "alice") {
		t.Error("System prompt should contain member names")
	}

	t.Logf("Room Collaboration System Prompt:\n%s", systemPrompt)
}

func TestPromptBuilder_WorkflowContext(t *testing.T) {
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
		Metadata: map[string]any{
			"workflow_id":            "wf-456",
			"workflow_current_step":  "data_processing",
			"workflow_total_steps":   5,
			"workflow_step_index":    2,
			"workflow_previous_step": "data_collection",
			"workflow_next_step":     "data_analysis",
		},
	}

	ag, err := Create(context.Background(), config, deps)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	defer func() { _ = ag.Close() }()

	systemPrompt := ag.GetSystemPrompt()

	// 验证工作流模块
	if !strings.Contains(systemPrompt, "## Workflow Context") {
		t.Error("System prompt should contain workflow section")
	}

	if !strings.Contains(systemPrompt, "wf-456") {
		t.Error("System prompt should contain workflow ID")
	}

	if !strings.Contains(systemPrompt, "data_processing") {
		t.Error("System prompt should contain current step")
	}

	if !strings.Contains(systemPrompt, "Step 3 of 5") {
		t.Error("System prompt should contain step progress")
	}

	t.Logf("Workflow Context System Prompt:\n%s", systemPrompt)
}

func TestPromptOptimizer_RemoveDuplicates(t *testing.T) {
	optimizer := &PromptOptimizer{
		RemoveDuplicates: true,
	}

	input := `Line 1


Line 2


Line 3`

	expected := `Line 1

Line 2

Line 3`

	result := optimizer.Optimize(input)

	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestPromptOptimizer_Truncate(t *testing.T) {
	optimizer := &PromptOptimizer{
		MaxLength: 100,
	}

	input := strings.Repeat("a", 200)

	result := optimizer.Optimize(input)

	if len(result) > 150 { // 允许一些截断提示的额外字符
		t.Errorf("Expected length <= 150, got %d", len(result))
	}
}

func TestAnalyzePrompt(t *testing.T) {
	prompt := `## Section 1

Content 1

## Section 2

Content 2`

	stats := AnalyzePrompt(prompt)

	if stats.SectionCount != 2 {
		t.Errorf("Expected 2 sections, got %d", stats.SectionCount)
	}

	if stats.TotalLength != len(prompt) {
		t.Errorf("Expected length %d, got %d", len(prompt), stats.TotalLength)
	}

	t.Logf("Stats: %s", FormatStats(stats))
}

func TestPromptTemplatePresets(t *testing.T) {
	presets := []struct {
		name   string
		preset *PromptTemplatePreset
	}{
		{"CodeAssistant", CodeAssistantPreset},
		{"ResearchAssistant", ResearchAssistantPreset},
		{"DataAnalyst", DataAnalystPreset},
		{"DevOpsEngineer", DevOpsEngineerPreset},
		{"TechnicalWriter", TechnicalWriterPreset},
		{"ProjectManager", ProjectManagerPreset},
		{"SecurityAuditor", SecurityAuditorPreset},
		{"GeneralAssistant", GeneralAssistantPreset},
	}

	for _, tc := range presets {
		t.Run(tc.name, func(t *testing.T) {
			if tc.preset == nil {
				t.Fatal("Preset is nil")
			}

			if tc.preset.ID == "" {
				t.Error("Preset ID should not be empty")
			}

			if tc.preset.Template == nil {
				t.Error("Preset template should not be nil")
			}

			if tc.preset.Template.SystemPrompt == "" {
				t.Error("Preset system prompt should not be empty")
			}

			t.Logf("Preset %s: %s", tc.preset.ID, tc.preset.Description)
		})
	}
}

func TestGetPreset(t *testing.T) {
	preset := GetPreset("code-assistant")
	if preset == nil {
		t.Fatal("Should find code-assistant preset")
	}

	if preset.ID != "code-assistant" {
		t.Errorf("Expected ID 'code-assistant', got '%s'", preset.ID)
	}

	notFound := GetPreset("non-existent")
	if notFound != nil {
		t.Error("Should return nil for non-existent preset")
	}
}

func TestRegisterAllPresets(t *testing.T) {
	registry := NewTemplateRegistry()
	RegisterAllPresets(registry)

	// 验证所有预设都已注册
	for _, preset := range AllPresets {
		template, err := registry.Get(preset.ID)
		if err != nil {
			t.Errorf("Failed to get preset %s: %v", preset.ID, err)
		}

		if template.ID != preset.ID {
			t.Errorf("Expected ID %s, got %s", preset.ID, template.ID)
		}
	}
}

func TestPromptBuilder_WithPreset(t *testing.T) {
	deps := setupPromptTestDeps(t)

	// 注册所有预设
	RegisterAllPresets(deps.TemplateRegistry)

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

	systemPrompt := ag.GetSystemPrompt()

	// 验证代码助手预设的特征
	if !strings.Contains(systemPrompt, "professional code assistant") {
		t.Error("Should contain code assistant description")
	}

	if !strings.Contains(systemPrompt, "## Code References") {
		t.Error("Should contain code reference section")
	}

	t.Logf("Code Assistant Preset System Prompt length: %d", len(systemPrompt))
}
