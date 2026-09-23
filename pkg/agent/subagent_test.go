package agent

import (
	"testing"
)

func TestSubAgentManager_RegisterBuiltinSpecs(t *testing.T) {
	// 创建一个简单的 Dependencies（不需要完整初始化）
	deps := &Dependencies{
		TemplateRegistry: NewTemplateRegistry(),
	}

	mgr := NewSubAgentManager(deps)

	// 验证内置规格已注册
	specs := mgr.ListSpecs()
	if len(specs) == 0 {
		t.Fatal("expected builtin specs to be registered")
	}

	// 验证 Explore 规格
	exploreSpec, err := mgr.GetSpec("Explore")
	if err != nil {
		t.Fatalf("failed to get Explore spec: %v", err)
	}
	if exploreSpec.Name != "Explore" {
		t.Errorf("expected name 'Explore', got '%s'", exploreSpec.Name)
	}
	if len(exploreSpec.Tools) == 0 {
		t.Error("expected Explore to have tools")
	}

	// 验证 Plan 规格
	planSpec, err := mgr.GetSpec("Plan")
	if err != nil {
		t.Fatalf("failed to get Plan spec: %v", err)
	}
	if planSpec.Name != "Plan" {
		t.Errorf("expected name 'Plan', got '%s'", planSpec.Name)
	}

	// 验证 general-purpose 规格
	gpSpec, err := mgr.GetSpec("general-purpose")
	if err != nil {
		t.Fatalf("failed to get general-purpose spec: %v", err)
	}
	if gpSpec.Name != "general-purpose" {
		t.Errorf("expected name 'general-purpose', got '%s'", gpSpec.Name)
	}
}

func TestSubAgentManager_GetSpec_NotFound(t *testing.T) {
	deps := &Dependencies{
		TemplateRegistry: NewTemplateRegistry(),
	}

	mgr := NewSubAgentManager(deps)

	_, err := mgr.GetSpec("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent spec")
	}
}

func TestSubAgentManagerFactory(t *testing.T) {
	deps := &Dependencies{
		TemplateRegistry: NewTemplateRegistry(),
	}

	mgr := NewSubAgentManager(deps)
	factory := NewSubAgentManagerFactory(mgr)

	// 测试 ListTypes
	types := factory.ListTypes()
	if len(types) == 0 {
		t.Error("expected factory to list types")
	}

	// 验证包含内置类型
	hasExplore := false
	hasPlan := false
	for _, typ := range types {
		if typ == "Explore" {
			hasExplore = true
		}
		if typ == "Plan" {
			hasPlan = true
		}
	}
	if !hasExplore {
		t.Error("expected 'Explore' in types")
	}
	if !hasPlan {
		t.Error("expected 'Plan' in types")
	}

	// 测试 Create
	executor, err := factory.Create("Explore")
	if err != nil {
		t.Fatalf("failed to create executor: %v", err)
	}

	spec := executor.GetSpec()
	if spec.Name != "Explore" {
		t.Errorf("expected spec name 'Explore', got '%s'", spec.Name)
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "he..."},
		{"", 5, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "..."},
	}

	for _, tt := range tests {
		result := truncateString(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}
