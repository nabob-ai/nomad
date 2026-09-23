package project

import (
	"context"
)

// Store 项目记忆存储接口
type Store interface {
	// Load 加载项目记忆
	Load(ctx context.Context, projectID string) (*ProjectMemory, error)

	// Save 保存项目记忆
	Save(ctx context.Context, memory *ProjectMemory) error

	// Delete 删除项目记忆
	Delete(ctx context.Context, projectID string) error

	// Exists 检查项目记忆是否存在
	Exists(ctx context.Context, projectID string) (bool, error)

	// AppendEntry 追加条目到指定章节
	AppendEntry(ctx context.Context, projectID string, section MemorySection, entry *Entry) error

	// UpdateWorkflowStep 更新工作流步骤状态
	UpdateWorkflowStep(ctx context.Context, projectID string, step *WorkflowStep) error

	// RecordParams 记录生成参数
	RecordParams(ctx context.Context, projectID string, params *GenerationParams) error

	// GetSummary 获取摘要（用于 AI 上下文注入）
	GetSummary(ctx context.Context, projectID string) (string, error)

	// GetVersionHistory 获取版本历史
	GetVersionHistory(ctx context.Context, projectID string, limit int) ([]*ProjectMemory, error)
}

// StoreConfig 存储配置
type StoreConfig struct {
	// BasePath 基础路径（用于文件存储）
	BasePath string

	// FileName 文件名（默认 AGENTS.md）
	FileName string

	// Template 模板内容
	Template string

	// AutoCreate 是否自动创建
	AutoCreate bool
}

// DefaultStoreConfig 返回默认配置
func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		BasePath:   "workspaces",
		FileName:   "AGENTS.md",
		AutoCreate: true,
	}
}
