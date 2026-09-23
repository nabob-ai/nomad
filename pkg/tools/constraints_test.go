package tools

import (
	"context"
	"testing"
)

func TestNoConstraints(t *testing.T) {
	constraints := &NoConstraints{}
	ctx := context.Background()

	if !constraints.IsAllowed(ctx, "AnyTool") {
		t.Error("NoConstraints should allow all tools")
	}

	if constraints.GetConstraintType() != ConstraintTypeNone {
		t.Errorf("Expected type %s, got %s", ConstraintTypeNone, constraints.GetConstraintType())
	}

	allowed := constraints.GetAllowedTools(ctx)
	if allowed != nil {
		t.Error("NoConstraints should return nil for allowed tools")
	}
}

func TestWhitelistConstraints(t *testing.T) {
	constraints := NewWhitelistConstraints([]string{"Read", "Write"})
	ctx := context.Background()

	if !constraints.IsAllowed(ctx, "Read") {
		t.Error("Read should be allowed")
	}

	if !constraints.IsAllowed(ctx, "Write") {
		t.Error("Write should be allowed")
	}

	if constraints.IsAllowed(ctx, "Bash") {
		t.Error("Bash should not be allowed")
	}

	allowed := constraints.GetAllowedTools(ctx)
	if len(allowed) != 2 {
		t.Errorf("Expected 2 allowed tools, got %d", len(allowed))
	}
}

func TestBlacklistConstraints(t *testing.T) {
	constraints := NewBlacklistConstraints([]string{"Bash", "Write"})
	ctx := context.Background()

	if constraints.IsAllowed(ctx, "Bash") {
		t.Error("Bash should not be allowed")
	}

	if constraints.IsAllowed(ctx, "Write") {
		t.Error("Write should not be allowed")
	}

	if !constraints.IsAllowed(ctx, "Read") {
		t.Error("Read should be allowed")
	}

	allowed := constraints.GetAllowedTools(ctx)
	if allowed != nil {
		t.Error("Blacklist should return nil for allowed tools")
	}
}

func TestRequiredToolConstraints(t *testing.T) {
	constraints := NewRequiredToolConstraints("Read")
	ctx := context.Background()

	if !constraints.IsAllowed(ctx, "Read") {
		t.Error("Read should be allowed")
	}

	if constraints.IsAllowed(ctx, "Write") {
		t.Error("Write should not be allowed")
	}

	allowed := constraints.GetAllowedTools(ctx)
	if len(allowed) != 1 || allowed[0] != "Read" {
		t.Errorf("Expected only Read to be allowed, got %v", allowed)
	}

	if constraints.GetRequiredTool() != "Read" {
		t.Errorf("Expected required tool 'Read', got '%s'", constraints.GetRequiredTool())
	}
}

func TestToolChoice_ToConstraints(t *testing.T) {
	tests := []struct {
		name         string
		choice       *ToolChoice
		expectedType ConstraintType
		expectedTool string
	}{
		{
			name:         "auto",
			choice:       ToolChoiceAuto,
			expectedType: ConstraintTypeNone,
		},
		{
			name:         "any",
			choice:       ToolChoiceAny,
			expectedType: ConstraintTypeNone,
		},
		{
			name:         "required tool",
			choice:       ToolChoiceRequired("Read"),
			expectedType: ConstraintTypeRequired,
			expectedTool: "Read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraints := tt.choice.ToConstraints()
			if constraints.GetConstraintType() != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, constraints.GetConstraintType())
			}

			if tt.expectedTool != "" {
				if req, ok := constraints.(*RequiredToolConstraints); ok {
					if req.GetRequiredTool() != tt.expectedTool {
						t.Errorf("Expected tool %s, got %s", tt.expectedTool, req.GetRequiredTool())
					}
				} else {
					t.Error("Expected RequiredToolConstraints")
				}
			}
		})
	}
}

func TestDefaultToolSelector_SelectTools(t *testing.T) {
	selector := &DefaultToolSelector{}
	ctx := context.Background()

	// 创建模拟工具
	allTools := []Tool{
		&mockTool{name: "Read"},
		&mockTool{name: "Write"},
		&mockTool{name: "Bash"},
	}

	tests := []struct {
		name          string
		constraints   ToolConstraints
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "no constraints",
			constraints:   &NoConstraints{},
			expectedCount: 3,
		},
		{
			name:          "whitelist",
			constraints:   NewWhitelistConstraints([]string{"Read", "Write"}),
			expectedCount: 2,
			expectedNames: []string{"Read", "Write"},
		},
		{
			name:          "blacklist",
			constraints:   NewBlacklistConstraints([]string{"Bash"}),
			expectedCount: 2,
			expectedNames: []string{"Read", "Write"},
		},
		{
			name:          "required",
			constraints:   NewRequiredToolConstraints("Read"),
			expectedCount: 1,
			expectedNames: []string{"Read"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selected, err := selector.SelectTools(ctx, allTools, tt.constraints)
			if err != nil {
				t.Fatalf("SelectTools failed: %v", err)
			}

			if len(selected) != tt.expectedCount {
				t.Errorf("Expected %d tools, got %d", tt.expectedCount, len(selected))
			}

			if tt.expectedNames != nil {
				names := make(map[string]bool)
				for _, tool := range selected {
					names[tool.Name()] = true
				}
				for _, expected := range tt.expectedNames {
					if !names[expected] {
						t.Errorf("Expected tool %s not found", expected)
					}
				}
			}
		})
	}
}

func TestDefaultToolSelector_ShouldUseToolChoice(t *testing.T) {
	selector := &DefaultToolSelector{}
	ctx := context.Background()

	tests := []struct {
		name        string
		constraints ToolConstraints
		shouldUse   bool
		choiceType  string
	}{
		{
			name:        "no constraints",
			constraints: &NoConstraints{},
			shouldUse:   false,
		},
		{
			name:        "required tool",
			constraints: NewRequiredToolConstraints("Read"),
			shouldUse:   true,
			choiceType:  "tool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			choice, shouldUse := selector.ShouldUseToolChoice(ctx, tt.constraints)
			if shouldUse != tt.shouldUse {
				t.Errorf("Expected shouldUse=%v, got %v", tt.shouldUse, shouldUse)
			}

			if tt.shouldUse && choice != nil {
				if choice.Type != tt.choiceType {
					t.Errorf("Expected choice type %s, got %s", tt.choiceType, choice.Type)
				}
			}
		})
	}
}

func TestConstraintsBuilder(t *testing.T) {
	tests := []struct {
		name         string
		buildFunc    func(*ConstraintsBuilder) *ConstraintsBuilder
		expectError  bool
		expectedType ConstraintType
	}{
		{
			name: "whitelist",
			buildFunc: func(b *ConstraintsBuilder) *ConstraintsBuilder {
				return b.WithWhitelist("Read", "Write")
			},
			expectError:  false,
			expectedType: ConstraintTypeWhitelist,
		},
		{
			name: "blacklist",
			buildFunc: func(b *ConstraintsBuilder) *ConstraintsBuilder {
				return b.WithBlacklist("Bash")
			},
			expectError:  false,
			expectedType: ConstraintTypeBlacklist,
		},
		{
			name: "required",
			buildFunc: func(b *ConstraintsBuilder) *ConstraintsBuilder {
				return b.WithRequired("Read")
			},
			expectError:  false,
			expectedType: ConstraintTypeRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewConstraintsBuilder()
			builder = tt.buildFunc(builder)
			constraints, err := builder.Build()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if constraints.GetConstraintType() != tt.expectedType {
					t.Errorf("Expected type %s, got %s", tt.expectedType, constraints.GetConstraintType())
				}
			}
		})
	}
}

// mockTool 用于测试的模拟工具
type mockTool struct {
	name string
}

func (m *mockTool) Name() string {
	return m.name
}

func (m *mockTool) Description() string {
	return "Mock tool for testing"
}

func (m *mockTool) InputSchema() map[string]any {
	return map[string]any{}
}

func (m *mockTool) Execute(ctx context.Context, input map[string]any, toolCtx *ToolContext) (any, error) {
	return "mock result", nil
}

func (m *mockTool) Prompt() string {
	return ""
}
