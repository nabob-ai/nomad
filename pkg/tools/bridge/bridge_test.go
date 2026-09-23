package bridge

import (
	"context"
	"testing"

	"github.com/nabob-ai/nomad/pkg/tools"
)

// mockTool 用于测试的模拟工具
type mockTool struct {
	name        string
	description string
	executeFunc func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error)
}

func (m *mockTool) Name() string        { return m.name }
func (m *mockTool) Description() string { return m.description }
func (m *mockTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}
func (m *mockTool) Prompt() string { return "" }
func (m *mockTool) Execute(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, input, tc)
	}
	return map[string]any{"result": "ok"}, nil
}

func TestToolBridge_CallTool(t *testing.T) {
	registry := tools.NewRegistry()

	// 注册模拟工具
	registry.Register("TestTool", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{
			name:        "TestTool",
			description: "A test tool",
			executeFunc: func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
				return map[string]any{
					"input_received": input,
				}, nil
			},
		}, nil
	})

	bridge := NewToolBridge(registry)
	ctx := context.Background()

	// 测试调用工具
	result, err := bridge.CallTool(ctx, "TestTool", map[string]any{"key": "value"}, nil)
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success, got error: %s", result.Error)
	}

	if result.Name != "TestTool" {
		t.Errorf("expected name 'TestTool', got %s", result.Name)
	}
}

func TestToolBridge_CallToolNotFound(t *testing.T) {
	registry := tools.NewRegistry()
	bridge := NewToolBridge(registry)
	ctx := context.Background()

	result, err := bridge.CallTool(ctx, "NonExistent", nil, nil)
	if err != nil {
		t.Fatalf("CallTool returned error: %v", err)
	}

	if result.Success {
		t.Error("expected failure for non-existent tool")
	}
}

func TestToolBridge_CallToolJSON(t *testing.T) {
	registry := tools.NewRegistry()

	registry.Register("EchoTool", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{
			name: "EchoTool",
			executeFunc: func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
				return input, nil
			},
		}, nil
	})

	bridge := NewToolBridge(registry)
	ctx := context.Background()

	result, err := bridge.CallToolJSON(ctx, "EchoTool", `{"message": "hello"}`, nil)
	if err != nil {
		t.Fatalf("CallToolJSON failed: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success, got error: %s", result.Error)
	}
}

func TestToolBridge_CallToolJSON_InvalidJSON(t *testing.T) {
	registry := tools.NewRegistry()
	bridge := NewToolBridge(registry)
	ctx := context.Background()

	result, err := bridge.CallToolJSON(ctx, "AnyTool", "not valid json", nil)
	if err != nil {
		t.Fatalf("CallToolJSON returned error: %v", err)
	}

	if result.Success {
		t.Error("expected failure for invalid JSON")
	}
}

func TestToolBridge_BatchCall(t *testing.T) {
	registry := tools.NewRegistry()

	callCount := 0
	registry.Register("CountTool", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{
			name: "CountTool",
			executeFunc: func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
				callCount++
				return map[string]any{"count": callCount}, nil
			},
		}, nil
	})

	bridge := NewToolBridge(registry)
	ctx := context.Background()

	calls := []CallToolInput{
		{Name: "CountTool", Input: nil},
		{Name: "CountTool", Input: nil},
		{Name: "CountTool", Input: nil},
	}

	result := bridge.CallToolsBatch(ctx, calls, nil)

	if result.Succeeded != 3 {
		t.Errorf("expected 3 succeeded, got %d", result.Succeeded)
	}

	if result.Failed != 0 {
		t.Errorf("expected 0 failed, got %d", result.Failed)
	}

	if len(result.Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(result.Results))
	}
}

func TestToolBridge_ParallelCall(t *testing.T) {
	registry := tools.NewRegistry()

	registry.Register("SlowTool", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{
			name: "SlowTool",
			executeFunc: func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
				return map[string]any{"done": true}, nil
			},
		}, nil
	})

	bridge := NewToolBridge(registry)
	ctx := context.Background()

	calls := []CallToolInput{
		{Name: "SlowTool", Input: nil},
		{Name: "SlowTool", Input: nil},
		{Name: "SlowTool", Input: nil},
	}

	result := bridge.CallToolsParallel(ctx, calls, nil)

	if result.Succeeded != 3 {
		t.Errorf("expected 3 succeeded, got %d", result.Succeeded)
	}
}

func TestToolBridge_ListAvailableTools(t *testing.T) {
	registry := tools.NewRegistry()

	registry.Register("Tool1", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{name: "Tool1"}, nil
	})
	registry.Register("Tool2", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{name: "Tool2"}, nil
	})

	bridge := NewToolBridge(registry)
	toolNames := bridge.ListAvailableTools()

	if len(toolNames) != 2 {
		t.Errorf("expected 2 tools, got %d", len(toolNames))
	}
}

func TestToolChain_Execute(t *testing.T) {
	registry := tools.NewRegistry()

	registry.Register("AddTool", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{
			name: "AddTool",
			executeFunc: func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
				a := input["a"].(float64)
				b := input["b"].(float64)
				return map[string]any{"result": a + b}, nil
			},
		}, nil
	})

	registry.Register("MultiplyTool", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{
			name: "MultiplyTool",
			executeFunc: func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
				x := input["x"].(float64)
				y := input["y"].(float64)
				return map[string]any{"result": x * y}, nil
			},
		}, nil
	})

	bridge := NewToolBridge(registry)
	chain := NewToolChain(bridge)

	chain.AddStep(ChainStep{
		Name:  "AddTool",
		Input: map[string]any{"a": float64(5), "b": float64(3)},
	})

	chain.AddStep(ChainStep{
		Name: "MultiplyTool",
		InputMapper: func(prevResult any) map[string]any {
			prev := prevResult.(map[string]any)
			return map[string]any{
				"x": prev["result"],
				"y": float64(2),
			}
		},
	})

	ctx := context.Background()
	result := chain.Execute(ctx, nil)

	if !result.Success {
		t.Fatalf("chain execution failed: %s", result.Error)
	}

	if len(result.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(result.Steps))
	}

	// 最终结果应该是 (5+3)*2 = 16
	finalResult := result.FinalResult.(map[string]any)
	if finalResult["result"] != float64(16) {
		t.Errorf("expected final result 16, got %v", finalResult["result"])
	}
}

func TestToolChain_FailureStopsChain(t *testing.T) {
	registry := tools.NewRegistry()

	registry.Register("FailTool", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{
			name: "FailTool",
			executeFunc: func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
				return nil, context.DeadlineExceeded
			},
		}, nil
	})

	registry.Register("NeverCalled", func(config map[string]any) (tools.Tool, error) {
		return &mockTool{
			name: "NeverCalled",
			executeFunc: func(ctx context.Context, input map[string]any, tc *tools.ToolContext) (any, error) {
				t.Error("this tool should not be called")
				return nil, nil
			},
		}, nil
	})

	bridge := NewToolBridge(registry)
	chain := NewToolChain(bridge)

	chain.AddStep(ChainStep{Name: "FailTool"})
	chain.AddStep(ChainStep{Name: "NeverCalled"})

	ctx := context.Background()
	result := chain.Execute(ctx, nil)

	if result.Success {
		t.Error("expected chain to fail")
	}

	// 只应该有 1 个步骤结果
	if len(result.Steps) != 1 {
		t.Errorf("expected 1 step result, got %d", len(result.Steps))
	}
}
