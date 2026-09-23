// Package project 提供项目级别的外部记忆系统
// 用于存储项目上下文、用户偏好、工作流状态等信息
// 典型实现：AGENTS.md 文件模式
package project

import (
	"time"
)

// MemorySection 记忆章节类型
type MemorySection string

const (
	// SectionPreferences 用户偏好（从对话中提取）
	SectionPreferences MemorySection = "preferences"
	// SectionChoices 关键选择（用户在工作流中的选择）
	SectionChoices MemorySection = "choices"
	// SectionWorkflow 工作流状态
	SectionWorkflow MemorySection = "workflow"
	// SectionParams 生成参数
	SectionParams MemorySection = "params"
	// SectionHistory 项目历史
	SectionHistory MemorySection = "history"
	// SectionCustom 自定义章节
	SectionCustom MemorySection = "custom"
)

// ProjectMemory 项目记忆
type ProjectMemory struct {
	// ProjectID 项目标识
	ProjectID string `json:"project_id"`

	// Title 项目标题
	Title string `json:"title"`

	// Description 项目描述
	Description string `json:"description"`

	// Sections 章节内容
	Sections map[MemorySection]*Section `json:"sections"`

	// Metadata 元数据
	Metadata map[string]any `json:"metadata,omitempty"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`

	// Version 版本号
	Version int `json:"version"`
}

// Section 章节
type Section struct {
	// Name 章节名称
	Name string `json:"name"`

	// Protected 是否受保护（不可自动删除）
	Protected bool `json:"protected"`

	// Entries 条目列表
	Entries []*Entry `json:"entries"`
}

// Entry 条目
type Entry struct {
	// ID 条目 ID
	ID string `json:"id"`

	// Category 分类（如 preference, constraint, style）
	Category string `json:"category"`

	// Content 内容
	Content string `json:"content"`

	// Value 结构化值（可选）
	Value any `json:"value,omitempty"`

	// Source 来源（dialog, choice, system）
	Source string `json:"source"`

	// Confidence 置信度
	Confidence float64 `json:"confidence"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// Metadata 元数据
	Metadata map[string]any `json:"metadata,omitempty"`
}

// WorkflowStep 工作流步骤
type WorkflowStep struct {
	// ID 步骤 ID
	ID string `json:"id"`

	// Name 步骤名称
	Name string `json:"name"`

	// Status 状态: pending, in_progress, completed, skipped
	Status string `json:"status"`

	// StartedAt 开始时间
	StartedAt *time.Time `json:"started_at,omitempty"`

	// CompletedAt 完成时间
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Note 备注
	Note string `json:"note,omitempty"`
}

// GenerationParams 生成参数
type GenerationParams struct {
	// Model 模型名称
	Model string `json:"model"`

	// Temperature 温度
	Temperature float64 `json:"temperature"`

	// MaxTokens 最大 tokens
	MaxTokens int `json:"max_tokens"`

	// Extra 额外参数
	Extra map[string]any `json:"extra,omitempty"`

	// Timestamp 记录时间
	Timestamp time.Time `json:"timestamp"`
}

// NewProjectMemory 创建新的项目记忆
func NewProjectMemory(projectID, title, description string) *ProjectMemory {
	now := time.Now()
	return &ProjectMemory{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Sections: map[MemorySection]*Section{
			SectionPreferences: {Name: "用户偏好", Protected: true, Entries: []*Entry{}},
			SectionChoices:     {Name: "关键选择", Protected: true, Entries: []*Entry{}},
			SectionWorkflow:    {Name: "工作流状态", Protected: false, Entries: []*Entry{}},
			SectionParams:      {Name: "生成参数", Protected: false, Entries: []*Entry{}},
			SectionHistory:     {Name: "项目历史", Protected: false, Entries: []*Entry{}},
		},
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}
}

// AddEntry 添加条目到指定章节
func (pm *ProjectMemory) AddEntry(section MemorySection, entry *Entry) {
	if pm.Sections[section] == nil {
		pm.Sections[section] = &Section{
			Name:    string(section),
			Entries: []*Entry{},
		}
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	pm.Sections[section].Entries = append(pm.Sections[section].Entries, entry)
	pm.UpdatedAt = time.Now()
	pm.Version++
}

// GetEntries 获取指定章节的所有条目
func (pm *ProjectMemory) GetEntries(section MemorySection) []*Entry {
	if pm.Sections[section] == nil {
		return []*Entry{}
	}
	return pm.Sections[section].Entries
}

// GetSummary 获取摘要（用于注入 AI 上下文）
func (pm *ProjectMemory) GetSummary() string {
	// 由 Store 实现具体的格式化逻辑
	return ""
}
