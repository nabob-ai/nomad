package tools

import (
	"context"
	"testing"
)

func TestToolAnnotations_Defaults(t *testing.T) {
	a := &ToolAnnotations{}

	if a.ReadOnly {
		t.Error("ReadOnly should be false by default")
	}
	if a.Destructive {
		t.Error("Destructive should be false by default")
	}
	if a.Idempotent {
		t.Error("Idempotent should be false by default")
	}
	if a.OpenWorld {
		t.Error("OpenWorld should be false by default")
	}
	if a.RiskLevel != 0 {
		t.Errorf("RiskLevel should be 0 by default, got %d", a.RiskLevel)
	}
}

func TestToolAnnotations_Clone(t *testing.T) {
	original := &ToolAnnotations{
		ReadOnly:    true,
		Destructive: false,
		RiskLevel:   2,
		Category:    "test",
	}

	cloned := original.Clone()

	if cloned.ReadOnly != original.ReadOnly {
		t.Error("ReadOnly not cloned correctly")
	}
	if cloned.Destructive != original.Destructive {
		t.Error("Destructive not cloned correctly")
	}
	if cloned.RiskLevel != original.RiskLevel {
		t.Error("RiskLevel not cloned correctly")
	}
	if cloned.Category != original.Category {
		t.Error("Category not cloned correctly")
	}

	// 修改克隆不应影响原始
	cloned.RiskLevel = 5
	if original.RiskLevel == 5 {
		t.Error("Modifying clone should not affect original")
	}
}

func TestToolAnnotations_IsSafeForAutoApproval(t *testing.T) {
	tests := []struct {
		name     string
		ann      *ToolAnnotations
		expected bool
	}{
		{
			name:     "safe read-only",
			ann:      AnnotationsSafeReadOnly,
			expected: true,
		},
		{
			name:     "safe write",
			ann:      AnnotationsSafeWrite,
			expected: false,
		},
		{
			name:     "destructive write",
			ann:      AnnotationsDestructiveWrite,
			expected: false,
		},
		{
			name:     "execution",
			ann:      AnnotationsExecution,
			expected: false,
		},
		{
			name:     "network read",
			ann:      AnnotationsNetworkRead,
			expected: false, // OpenWorld = true
		},
		{
			name:     "read-only but open world",
			ann:      &ToolAnnotations{ReadOnly: true, OpenWorld: true},
			expected: false,
		},
		{
			name:     "read-only internal",
			ann:      &ToolAnnotations{ReadOnly: true, OpenWorld: false},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ann.IsSafeForAutoApproval(); got != tt.expected {
				t.Errorf("IsSafeForAutoApproval() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestToolAnnotations_RiskLevelName(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{RiskLevelSafe, "safe"},
		{RiskLevelLow, "low"},
		{RiskLevelMedium, "medium"},
		{RiskLevelHigh, "high"},
		{RiskLevelCritical, "critical"},
		{99, "unknown"},
		{-1, "unknown"},
	}

	for _, tt := range tests {
		ann := &ToolAnnotations{RiskLevel: tt.level}
		if got := ann.RiskLevelName(); got != tt.expected {
			t.Errorf("RiskLevelName() for level %d = %v, want %v", tt.level, got, tt.expected)
		}
	}
}

func TestPredefinedAnnotations(t *testing.T) {
	// 测试预定义注解的正确性
	tests := []struct {
		name        string
		ann         *ToolAnnotations
		readOnly    bool
		destructive bool
		openWorld   bool
		minRisk     int
	}{
		{"SafeReadOnly", AnnotationsSafeReadOnly, true, false, false, RiskLevelSafe},
		{"SafeWrite", AnnotationsSafeWrite, false, false, false, RiskLevelLow},
		{"DestructiveWrite", AnnotationsDestructiveWrite, false, true, false, RiskLevelHigh},
		{"Execution", AnnotationsExecution, false, true, true, RiskLevelHigh}, // OpenWorld = true (可能访问网络)
		{"NetworkRead", AnnotationsNetworkRead, true, false, true, RiskLevelLow},
		{"NetworkWrite", AnnotationsNetworkWrite, false, false, true, RiskLevelMedium},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ann.ReadOnly != tt.readOnly {
				t.Errorf("ReadOnly = %v, want %v", tt.ann.ReadOnly, tt.readOnly)
			}
			if tt.ann.Destructive != tt.destructive {
				t.Errorf("Destructive = %v, want %v", tt.ann.Destructive, tt.destructive)
			}
			if tt.ann.OpenWorld != tt.openWorld {
				t.Errorf("OpenWorld = %v, want %v", tt.ann.OpenWorld, tt.openWorld)
			}
			if tt.ann.RiskLevel < tt.minRisk {
				t.Errorf("RiskLevel = %d, want >= %d", tt.ann.RiskLevel, tt.minRisk)
			}
		})
	}
}

// mockToolWithAnnotations 带注解的模拟工具
type mockToolWithAnnotations struct {
	name        string
	annotations *ToolAnnotations
}

func (m *mockToolWithAnnotations) Name() string        { return m.name }
func (m *mockToolWithAnnotations) Description() string { return "mock tool with annotations" }
func (m *mockToolWithAnnotations) InputSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}
func (m *mockToolWithAnnotations) Execute(ctx context.Context, input map[string]any, tc *ToolContext) (any, error) {
	return nil, nil
}
func (m *mockToolWithAnnotations) Prompt() string                { return "mock prompt" }
func (m *mockToolWithAnnotations) Annotations() *ToolAnnotations { return m.annotations }

// mockToolWithoutAnnotations 不带注解的模拟工具
type mockToolWithoutAnnotations struct {
	name string
}

func (m *mockToolWithoutAnnotations) Name() string        { return m.name }
func (m *mockToolWithoutAnnotations) Description() string { return "mock tool without annotations" }
func (m *mockToolWithoutAnnotations) InputSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}
func (m *mockToolWithoutAnnotations) Execute(ctx context.Context, input map[string]any, tc *ToolContext) (any, error) {
	return nil, nil
}
func (m *mockToolWithoutAnnotations) Prompt() string { return "mock prompt" }

func TestGetAnnotations(t *testing.T) {
	tests := []struct {
		name         string
		tool         Tool
		expectedSafe bool
		expectedRisk int
	}{
		{
			name:         "tool with annotations",
			tool:         &mockToolWithAnnotations{name: "test", annotations: AnnotationsSafeReadOnly},
			expectedSafe: true,
			expectedRisk: RiskLevelSafe,
		},
		{
			name:         "tool without annotations",
			tool:         &mockToolWithoutAnnotations{name: "test"},
			expectedSafe: false,           // 默认不安全
			expectedRisk: RiskLevelMedium, // 未知工具保守处理
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ann := GetAnnotations(tt.tool)
			if ann == nil {
				t.Fatal("GetAnnotations returned nil")
			}
			if ann.IsSafeForAutoApproval() != tt.expectedSafe {
				t.Errorf("IsSafeForAutoApproval() = %v, want %v", ann.IsSafeForAutoApproval(), tt.expectedSafe)
			}
			if ann.RiskLevel != tt.expectedRisk {
				t.Errorf("RiskLevel = %d, want %d", ann.RiskLevel, tt.expectedRisk)
			}
		})
	}
}

func TestIsToolSafeForAutoApproval(t *testing.T) {
	safeTool := &mockToolWithAnnotations{name: "safe", annotations: AnnotationsSafeReadOnly}
	unsafeTool := &mockToolWithAnnotations{name: "unsafe", annotations: AnnotationsDestructiveWrite}
	noAnnotationTool := &mockToolWithoutAnnotations{name: "none"}

	if !IsToolSafeForAutoApproval(safeTool) {
		t.Error("Safe read-only tool should be safe for auto approval")
	}
	if IsToolSafeForAutoApproval(unsafeTool) {
		t.Error("Destructive tool should not be safe for auto approval")
	}
	if IsToolSafeForAutoApproval(noAnnotationTool) {
		t.Error("Tool without annotations should not be safe for auto approval by default")
	}
}

func TestGetToolRiskLevel(t *testing.T) {
	highRiskTool := &mockToolWithAnnotations{name: "high", annotations: AnnotationsDestructiveWrite}
	lowRiskTool := &mockToolWithAnnotations{name: "low", annotations: AnnotationsSafeReadOnly}
	noAnnotationTool := &mockToolWithoutAnnotations{name: "none"}

	if GetToolRiskLevel(highRiskTool) != RiskLevelHigh {
		t.Errorf("High risk tool level = %d, want %d", GetToolRiskLevel(highRiskTool), RiskLevelHigh)
	}
	if GetToolRiskLevel(lowRiskTool) != RiskLevelSafe {
		t.Errorf("Low risk tool level = %d, want %d", GetToolRiskLevel(lowRiskTool), RiskLevelSafe)
	}
	if GetToolRiskLevel(noAnnotationTool) != RiskLevelMedium {
		t.Errorf("No annotation tool level = %d, want %d", GetToolRiskLevel(noAnnotationTool), RiskLevelMedium)
	}
}

func TestAnnotationsNilSafety(t *testing.T) {
	// 测试 nil annotations 的安全性
	var ann *ToolAnnotations = nil

	// GetAnnotations 应该对 nil 工具返回默认注解
	defaultAnn := GetAnnotations(&mockToolWithoutAnnotations{name: "test"})
	if defaultAnn == nil {
		t.Error("GetAnnotations should never return nil")
	}

	// 测试 nil annotations 的方法（应该 panic 或返回默认值）
	// 这里我们验证非 nil 注解的行为
	ann = &ToolAnnotations{}
	if ann.IsSafeForAutoApproval() {
		// 默认空注解不应该被认为是安全的（ReadOnly = false）
		t.Error("Empty annotations should not be safe for auto approval")
	}
}
