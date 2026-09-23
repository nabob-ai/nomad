package logic

import (
	"context"
	"sort"
	"sync"
	"time"
)

// InMemoryStore 内存存储实现（用于测试和简单场景）
type InMemoryStore struct {
	mu       sync.RWMutex
	memories map[string]*LogicMemory // key: namespace:key
	closed   bool
}

// NewInMemoryStore 创建内存存储
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		memories: make(map[string]*LogicMemory),
	}
}

// makeKey 生成存储键
func makeKey(namespace, key string) string {
	return namespace + ":" + key
}

// Save 保存或更新 Memory
func (s *InMemoryStore) Save(ctx context.Context, memory *LogicMemory) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStoreClosed
	}

	if memory.Namespace == "" {
		return ErrInvalidNamespace
	}

	key := makeKey(memory.Namespace, memory.Key)

	// 如果是新记录，设置创建时间
	if _, exists := s.memories[key]; !exists {
		if memory.CreatedAt.IsZero() {
			memory.CreatedAt = time.Now()
		}
	}

	// 更新时间
	memory.UpdatedAt = time.Now()
	if memory.LastAccessed.IsZero() {
		memory.LastAccessed = memory.UpdatedAt
	}

	// 深拷贝存储
	stored := *memory
	s.memories[key] = &stored

	return nil
}

// Get 获取单个 Memory
func (s *InMemoryStore) Get(ctx context.Context, namespace, key string) (*LogicMemory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrStoreClosed
	}

	storeKey := makeKey(namespace, key)
	memory, exists := s.memories[storeKey]
	if !exists {
		return nil, ErrMemoryNotFound
	}

	// 返回拷贝
	result := *memory
	return &result, nil
}

// Delete 删除 Memory
func (s *InMemoryStore) Delete(ctx context.Context, namespace, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStoreClosed
	}

	storeKey := makeKey(namespace, key)
	delete(s.memories, storeKey)
	return nil
}

// List 列出符合条件的 Memory
func (s *InMemoryStore) List(ctx context.Context, namespace string, filters ...Filter) ([]*LogicMemory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrStoreClosed
	}

	opts := ApplyFilters(filters...)

	var result []*LogicMemory
	for _, memory := range s.memories {
		// 过滤 namespace
		if namespace != "" && memory.Namespace != namespace {
			continue
		}

		// 过滤类型
		if opts.Type != "" && memory.Type != opts.Type {
			continue
		}

		// 过滤作用域
		if opts.Scope != "" && memory.Scope != opts.Scope {
			continue
		}

		// 过滤置信度
		if memory.Provenance != nil && memory.Provenance.Confidence < opts.MinConfidence {
			continue
		}

		// 返回拷贝
		copied := *memory
		result = append(result, &copied)
	}

	// 排序
	sortMemories(result, opts.OrderBy)

	// 限制数量
	if opts.MaxResults > 0 && len(result) > opts.MaxResults {
		result = result[:opts.MaxResults]
	}

	return result, nil
}

// SearchByType 按类型搜索
func (s *InMemoryStore) SearchByType(ctx context.Context, namespace, memoryType string) ([]*LogicMemory, error) {
	return s.List(ctx, namespace, WithType(memoryType))
}

// SearchByScope 按作用域搜索
func (s *InMemoryStore) SearchByScope(ctx context.Context, namespace string, scope MemoryScope) ([]*LogicMemory, error) {
	return s.List(ctx, namespace, WithScope(scope))
}

// GetTopK 获取 TopK Memory
func (s *InMemoryStore) GetTopK(ctx context.Context, namespace string, k int, orderBy OrderBy) ([]*LogicMemory, error) {
	return s.List(ctx, namespace, WithTopK(k), WithOrderBy(orderBy))
}

// IncrementAccessCount 增加访问计数
func (s *InMemoryStore) IncrementAccessCount(ctx context.Context, namespace, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStoreClosed
	}

	storeKey := makeKey(namespace, key)
	memory, exists := s.memories[storeKey]
	if !exists {
		return ErrMemoryNotFound
	}

	memory.AccessCount++
	memory.LastAccessed = time.Now()
	return nil
}

// GetStats 获取统计信息
func (s *InMemoryStore) GetStats(ctx context.Context, namespace string) (*MemoryStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrStoreClosed
	}

	stats := &MemoryStats{
		CountByType:  make(map[string]int),
		CountByScope: make(map[MemoryScope]int),
	}

	var totalConfidence float64
	var lastUpdated time.Time

	for _, memory := range s.memories {
		if namespace != "" && memory.Namespace != namespace {
			continue
		}

		stats.TotalCount++
		stats.CountByType[memory.Type]++
		stats.CountByScope[memory.Scope]++

		if memory.Provenance != nil {
			totalConfidence += memory.Provenance.Confidence
		}

		if memory.UpdatedAt.After(lastUpdated) {
			lastUpdated = memory.UpdatedAt
		}
	}

	if stats.TotalCount > 0 {
		stats.AverageConfidence = totalConfidence / float64(stats.TotalCount)
	}
	stats.LastUpdated = lastUpdated

	return stats, nil
}

// Prune 清理低价值 Memory
func (s *InMemoryStore) Prune(ctx context.Context, criteria PruneCriteria) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return 0, ErrStoreClosed
	}

	var toDelete []string
	now := time.Now()

	for storeKey, memory := range s.memories {
		shouldPrune := false

		// 置信度过低
		if memory.Provenance != nil && memory.Provenance.Confidence < criteria.MinConfidence {
			shouldPrune = true
		}

		// 太久未访问
		if criteria.SinceLastAccess > 0 && now.Sub(memory.LastAccessed) > criteria.SinceLastAccess {
			shouldPrune = true
		}

		// 访问次数过少且年龄过大
		if criteria.MinAccessCount > 0 && criteria.MaxAge > 0 {
			if memory.AccessCount < criteria.MinAccessCount && now.Sub(memory.CreatedAt) > criteria.MaxAge {
				shouldPrune = true
			}
		}

		if shouldPrune {
			toDelete = append(toDelete, storeKey)
		}
	}

	for _, key := range toDelete {
		delete(s.memories, key)
	}

	return len(toDelete), nil
}

// Close 关闭存储
func (s *InMemoryStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.closed = true
	s.memories = nil
	return nil
}

// sortMemories 排序 Memory 列表
func sortMemories(memories []*LogicMemory, orderBy OrderBy) {
	switch orderBy {
	case OrderByConfidence:
		sort.Slice(memories, func(i, j int) bool {
			ci, cj := 0.0, 0.0
			if memories[i].Provenance != nil {
				ci = memories[i].Provenance.Confidence
			}
			if memories[j].Provenance != nil {
				cj = memories[j].Provenance.Confidence
			}
			return ci > cj
		})
	case OrderByLastAccessed:
		sort.Slice(memories, func(i, j int) bool {
			return memories[i].LastAccessed.After(memories[j].LastAccessed)
		})
	case OrderByCreatedAt:
		sort.Slice(memories, func(i, j int) bool {
			return memories[i].CreatedAt.After(memories[j].CreatedAt)
		})
	case OrderByAccessCount:
		sort.Slice(memories, func(i, j int) bool {
			return memories[i].AccessCount > memories[j].AccessCount
		})
	}
}

// 确保 InMemoryStore 实现 LogicMemoryStore 接口
var _ LogicMemoryStore = (*InMemoryStore)(nil)
