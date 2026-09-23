package search

import (
	"testing"

	"github.com/nabob-ai/nomad/pkg/tools"
)

func TestToolIndex_BasicOperations(t *testing.T) {
	index := NewToolIndex()

	// 测试添加工具
	entry := ToolIndexEntry{
		Name:        "Read",
		Description: "Read file from filesystem",
		Category:    "filesystem",
		Keywords:    []string{"file", "read", "content"},
		Deferred:    false,
		Source:      "builtin",
	}

	if err := index.IndexToolEntry(entry); err != nil {
		t.Fatalf("IndexToolEntry failed: %v", err)
	}

	// 验证工具被添加
	if index.Count() != 1 {
		t.Errorf("expected 1 tool, got %d", index.Count())
	}

	// 测试获取工具
	got := index.GetTool("Read")
	if got == nil {
		t.Fatal("expected tool to exist")
	}
	if got.Name != "Read" {
		t.Errorf("expected name 'Read', got %s", got.Name)
	}
}

func TestToolIndex_Search(t *testing.T) {
	index := NewToolIndex()

	// 添加多个工具
	entries := []ToolIndexEntry{
		{
			Name:        "Read",
			Description: "Read file from filesystem",
			Category:    "filesystem",
			Keywords:    []string{"file", "read"},
		},
		{
			Name:        "Write",
			Description: "Write content to file",
			Category:    "filesystem",
			Keywords:    []string{"file", "write"},
		},
		{
			Name:        "Bash",
			Description: "Execute shell commands",
			Category:    "execution",
			Keywords:    []string{"shell", "command"},
		},
	}

	for _, e := range entries {
		if err := index.IndexToolEntry(e); err != nil {
			t.Fatalf("IndexToolEntry failed: %v", err)
		}
	}

	// 搜索文件相关工具
	results := index.Search("file read", 5)

	if len(results) == 0 {
		t.Fatal("expected search results")
	}

	// Read 应该排在最前面
	if results[0].Entry.Name != "Read" {
		t.Errorf("expected 'Read' as top result, got %s", results[0].Entry.Name)
	}
}

func TestToolIndex_CategoryFilter(t *testing.T) {
	index := NewToolIndex()

	entries := []ToolIndexEntry{
		{Name: "Read", Category: "filesystem"},
		{Name: "Write", Category: "filesystem"},
		{Name: "Bash", Category: "execution"},
	}

	for _, e := range entries {
		_ = index.IndexToolEntry(e)
	}

	// 获取特定分类的工具
	fsTools := index.SearchByCategory("filesystem")
	if len(fsTools) != 2 {
		t.Errorf("expected 2 filesystem tools, got %d", len(fsTools))
	}

	execTools := index.SearchByCategory("execution")
	if len(execTools) != 1 {
		t.Errorf("expected 1 execution tool, got %d", len(execTools))
	}
}

func TestToolIndex_DeferredTools(t *testing.T) {
	index := NewToolIndex()

	entries := []ToolIndexEntry{
		{Name: "Read", Deferred: false},
		{Name: "MCPTool1", Deferred: true},
		{Name: "MCPTool2", Deferred: true},
	}

	for _, e := range entries {
		_ = index.IndexToolEntry(e)
	}

	// 获取延迟加载的工具
	deferred := index.GetDeferredTools()
	if len(deferred) != 2 {
		t.Errorf("expected 2 deferred tools, got %d", len(deferred))
	}

	// 获取活跃的工具
	active := index.GetActiveTools()
	if len(active) != 1 {
		t.Errorf("expected 1 active tool, got %d", len(active))
	}
}

func TestToolIndex_WithExamples(t *testing.T) {
	index := NewToolIndex()

	entry := ToolIndexEntry{
		Name:        "HttpRequest",
		Description: "Make HTTP requests",
		Examples: []tools.ToolExample{
			{
				Description: "GET request",
				Input:       map[string]any{"url": "https://example.com"},
			},
		},
	}

	if err := index.IndexToolEntry(entry); err != nil {
		t.Fatalf("IndexToolEntry failed: %v", err)
	}

	got := index.GetTool("HttpRequest")
	if got == nil {
		t.Fatal("expected tool to exist")
	}
	if len(got.Examples) != 1 {
		t.Errorf("expected 1 example, got %d", len(got.Examples))
	}
}

func TestToolIndex_RemoveTool(t *testing.T) {
	index := NewToolIndex()

	_ = index.IndexToolEntry(ToolIndexEntry{Name: "Test"})

	if index.Count() != 1 {
		t.Fatal("expected 1 tool after add")
	}

	removed := index.RemoveTool("Test")
	if !removed {
		t.Error("expected removal to succeed")
	}

	if index.Count() != 0 {
		t.Error("expected 0 tools after remove")
	}

	got := index.GetTool("Test")
	if got != nil {
		t.Error("tool should not exist after removal")
	}
}

func TestToolIndex_GetAllTools(t *testing.T) {
	index := NewToolIndex()

	entries := []ToolIndexEntry{
		{Name: "Tool1"},
		{Name: "Tool2"},
		{Name: "Tool3"},
	}

	for _, e := range entries {
		_ = index.IndexToolEntry(e)
	}

	all := index.GetAllTools()
	if len(all) != 3 {
		t.Errorf("expected 3 tools, got %d", len(all))
	}
}

func TestToolIndex_UpdateExisting(t *testing.T) {
	index := NewToolIndex()

	// 添加初始条目
	_ = index.IndexToolEntry(ToolIndexEntry{
		Name:        "Test",
		Description: "Original description",
	})

	// 更新条目
	_ = index.IndexToolEntry(ToolIndexEntry{
		Name:        "Test",
		Description: "Updated description",
	})

	// 验证只有一个条目
	if index.Count() != 1 {
		t.Errorf("expected 1 tool after update, got %d", index.Count())
	}

	// 验证描述被更新
	got := index.GetTool("Test")
	if got.Description != "Updated description" {
		t.Errorf("expected updated description, got %s", got.Description)
	}
}
