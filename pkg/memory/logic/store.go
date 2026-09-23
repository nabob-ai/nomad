package logic

import (
	"context"
)

// LogicMemoryStore 存储接口（类似现有 BackendProtocol）
// 应用层可以提供不同的实现（PostgreSQL, Redis, InMemory 等）
type LogicMemoryStore interface {
	// ===== 基础 CRUD =====

	// Save 保存或更新 Memory
	// 如果 Namespace + Key 已存在，则更新；否则创建新记录
	Save(ctx context.Context, memory *LogicMemory) error

	// Get 获取单个 Memory
	// 返回 ErrMemoryNotFound 如果不存在
	Get(ctx context.Context, namespace, key string) (*LogicMemory, error)

	// Delete 删除 Memory
	// 如果不存在不返回错误
	Delete(ctx context.Context, namespace, key string) error

	// List 列出符合条件的 Memory
	// filters 用于过滤和排序
	List(ctx context.Context, namespace string, filters ...Filter) ([]*LogicMemory, error)

	// ===== 高级查询 =====

	// SearchByType 按类型搜索
	SearchByType(ctx context.Context, namespace, memoryType string) ([]*LogicMemory, error)

	// SearchByScope 按作用域搜索
	SearchByScope(ctx context.Context, namespace string, scope MemoryScope) ([]*LogicMemory, error)

	// GetTopK 获取 TopK Memory
	// 按 orderBy 排序后返回前 k 个
	GetTopK(ctx context.Context, namespace string, k int, orderBy OrderBy) ([]*LogicMemory, error)

	// ===== 统计 =====

	// IncrementAccessCount 增加访问计数
	// 同时更新 LastAccessed 时间
	IncrementAccessCount(ctx context.Context, namespace, key string) error

	// GetStats 获取统计信息
	GetStats(ctx context.Context, namespace string) (*MemoryStats, error)

	// ===== 生命周期 =====

	// Prune 清理低价值 Memory
	// 返回清理的数量
	Prune(ctx context.Context, criteria PruneCriteria) (int, error)

	// ===== 连接管理 =====

	// Close 关闭连接
	Close() error
}

// StoreConfig 存储配置（通用）
type StoreConfig struct {
	// Type 存储类型（"postgres", "redis", "inmemory"）
	Type string

	// ConnectionString 连接字符串（PostgreSQL/Redis）
	ConnectionString string

	// TableName 表名（PostgreSQL，可选，默认 "logic_memories"）
	TableName string

	// KeyPrefix 键前缀（Redis，可选）
	KeyPrefix string

	// MaxConnections 最大连接数（可选）
	MaxConnections int
}

// 错误定义
var (
	// ErrMemoryNotFound Memory 不存在
	ErrMemoryNotFound = &StoreError{Code: "MEMORY_NOT_FOUND", Message: "logic memory not found"}

	// ErrDuplicateKey 键冲突
	ErrDuplicateKey = &StoreError{Code: "DUPLICATE_KEY", Message: "logic memory with this key already exists"}

	// ErrInvalidNamespace 无效的 Namespace
	ErrInvalidNamespace = &StoreError{Code: "INVALID_NAMESPACE", Message: "invalid namespace"}

	// ErrStoreClosed 存储已关闭
	ErrStoreClosed = &StoreError{Code: "STORE_CLOSED", Message: "logic memory store is closed"}
)

// StoreError 存储错误
type StoreError struct {
	Code    string
	Message string
	Err     error
}

func (e *StoreError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *StoreError) Unwrap() error {
	return e.Err
}

// NewStoreError 创建新的存储错误
func NewStoreError(code, message string, err error) *StoreError {
	return &StoreError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
