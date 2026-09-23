package logic

import (
	"time"

	"github.com/nabob-ai/nomad/pkg/memory"
)

// MemoryScope 记忆作用域
type MemoryScope string

const (
	// ScopeSession 单次会话级别（短期）
	ScopeSession MemoryScope = "session"
	// ScopeUser 用户级别（中期）
	ScopeUser MemoryScope = "user"
	// ScopeGlobal 全局级别（长期）
	ScopeGlobal MemoryScope = "global"
)

// LogicMemory 逻辑记忆（通用结构）
// Logic Memory 用于存储用户偏好、行为模式等可学习的记忆
type LogicMemory struct {
	// ===== 基础字段 =====

	// ID 唯一标识符
	ID string `json:"id"`

	// Namespace 租户隔离（如 user:123, team:456, global）
	Namespace string `json:"namespace"`

	// Scope 作用域（Session/User/Global）
	Scope MemoryScope `json:"scope"`

	// ===== Memory 类型（应用层定义）=====

	// Type Memory 类型（如 "user_preference", "behavior_pattern" 等）
	Type string `json:"type"`

	// Category 分类（可选，用于进一步分组）
	Category string `json:"category,omitempty"`

	// ===== Memory 内容 =====

	// Key 唯一标识（如 "writing_tone_preference"）
	// Namespace + Key 构成全局唯一键
	Key string `json:"key"`

	// Value 值（结构化数据）
	// 应用层可以存储任意 JSON 可序列化的数据
	Value any `json:"value"`

	// Description 人类可读描述（用于 Prompt 注入）
	// 例如："用户偏好口语化表达，避免使用书面语"
	Description string `json:"description"`

	// ===== 溯源（复用现有 Provenance）=====

	// Provenance 记忆溯源信息
	Provenance *memory.MemoryProvenance `json:"provenance"`

	// ===== 统计信息 =====

	// AccessCount 访问次数（用于 LRU 淘汰）
	AccessCount int `json:"access_count"`

	// LastAccessed 最后访问时间
	LastAccessed time.Time `json:"last_accessed"`

	// ===== 元信息 =====

	// Metadata 扩展字段（应用层自定义）
	Metadata map[string]any `json:"metadata,omitempty"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// Event 通用事件结构
// 用于 PatternMatcher 从事件中识别 Memory
type Event struct {
	// Type 事件类型（如 "user_message", "tool_result", "user_feedback" 等）
	Type string `json:"type"`

	// Source 来源（agent_id, user_id 等）
	Source string `json:"source"`

	// Data 事件数据（结构化）
	Data map[string]any `json:"data"`

	// Timestamp 时间戳
	Timestamp time.Time `json:"timestamp"`
}

// Filter Logic Memory 查询过滤器
type Filter func(*FilterOptions)

// FilterOptions 过滤选项
type FilterOptions struct {
	// Type 过滤类型
	Type string

	// Scope 过滤作用域
	Scope MemoryScope

	// MinConfidence 最低置信度
	MinConfidence float64

	// MaxResults TopK 限制
	MaxResults int

	// OrderBy 排序字段
	OrderBy OrderBy

	// SinceLastAccess 最后访问时间过滤
	SinceLastAccess time.Duration
}

// WithType 按类型过滤
func WithType(memoryType string) Filter {
	return func(opts *FilterOptions) {
		opts.Type = memoryType
	}
}

// WithScope 按作用域过滤
func WithScope(scope MemoryScope) Filter {
	return func(opts *FilterOptions) {
		opts.Scope = scope
	}
}

// WithMinConfidence 按最低置信度过滤
func WithMinConfidence(confidence float64) Filter {
	return func(opts *FilterOptions) {
		opts.MinConfidence = confidence
	}
}

// WithTopK 限制返回数量
func WithTopK(k int) Filter {
	return func(opts *FilterOptions) {
		opts.MaxResults = k
	}
}

// WithOrderBy 指定排序方式
func WithOrderBy(orderBy OrderBy) Filter {
	return func(opts *FilterOptions) {
		opts.OrderBy = orderBy
	}
}

// OrderBy 排序方式
type OrderBy string

const (
	// OrderByConfidence 按置信度降序
	OrderByConfidence OrderBy = "confidence DESC"
	// OrderByLastAccessed 按最后访问时间降序
	OrderByLastAccessed OrderBy = "last_accessed DESC"
	// OrderByCreatedAt 按创建时间降序
	OrderByCreatedAt OrderBy = "created_at DESC"
	// OrderByAccessCount 按访问次数降序
	OrderByAccessCount OrderBy = "access_count DESC"
)

// PruneCriteria 清理条件
type PruneCriteria struct {
	// MinConfidence 最低置信度（低于此值将被清理）
	MinConfidence float64

	// MaxAge 最大年龄（超过此时长将被清理）
	MaxAge time.Duration

	// MinAccessCount 最少访问次数（低于此值将被清理）
	MinAccessCount int

	// SinceLastAccess 最后访问时间（超过此时长未访问将被清理）
	SinceLastAccess time.Duration
}

// MemoryStats Logic Memory 统计信息
type MemoryStats struct {
	// TotalCount 总记忆数
	TotalCount int

	// CountByType 按类型统计
	CountByType map[string]int

	// CountByScope 按作用域统计
	CountByScope map[MemoryScope]int

	// AverageConfidence 平均置信度
	AverageConfidence float64

	// LastUpdated 最后更新时间
	LastUpdated time.Time
}

// ApplyFilters 应用过滤器到 FilterOptions
func ApplyFilters(filters ...Filter) *FilterOptions {
	opts := &FilterOptions{
		MinConfidence: 0.0,
		MaxResults:    0, // 0 表示不限制
		OrderBy:       OrderByConfidence,
	}

	for _, filter := range filters {
		filter(opts)
	}

	return opts
}
