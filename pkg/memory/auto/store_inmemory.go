package auto

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// InMemoryStore 内存存储实现
type InMemoryStore struct {
	mu sync.RWMutex

	// memories 按 ID 存储
	memories map[string]*Memory

	// projectIndex 项目索引
	projectIndex map[string][]string

	// tagIndex 标签索引
	tagIndex map[string][]string

	// config 配置
	config *InMemoryStoreConfig
}

// InMemoryStoreConfig 配置
type InMemoryStoreConfig struct {
	// MaxMemories 最大记忆数
	MaxMemories int
}

// NewInMemoryStore 创建内存存储
func NewInMemoryStore(config *InMemoryStoreConfig) *InMemoryStore {
	if config == nil {
		config = &InMemoryStoreConfig{
			MaxMemories: 10000,
		}
	}
	return &InMemoryStore{
		memories:     make(map[string]*Memory),
		projectIndex: make(map[string][]string),
		tagIndex:     make(map[string][]string),
		config:       config,
	}
}

// Save 保存记忆
func (s *InMemoryStore) Save(ctx context.Context, memory *Memory) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查容量
	if len(s.memories) >= s.config.MaxMemories {
		// 删除最旧的记忆
		s.removeOldest()
	}

	// 保存
	s.memories[memory.ID] = memory

	// 更新项目索引
	if memory.ProjectID != "" {
		s.projectIndex[memory.ProjectID] = append(s.projectIndex[memory.ProjectID], memory.ID)
	}

	// 更新标签索引
	for _, tag := range memory.Tags {
		s.tagIndex[tag] = append(s.tagIndex[tag], memory.ID)
	}

	return nil
}

// Load 加载记忆
func (s *InMemoryStore) Load(ctx context.Context, id string) (*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	memory, ok := s.memories[id]
	if !ok {
		return nil, fmt.Errorf("memory not found: %s", id)
	}

	memory.MarkAccessed()
	return memory, nil
}

// List 列出记忆
func (s *InMemoryStore) List(ctx context.Context, scope MemoryScope, projectID string, limit int) ([]*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Memory

	// 按项目过滤
	if projectID != "" {
		ids := s.projectIndex[projectID]
		for _, id := range ids {
			if memory, ok := s.memories[id]; ok {
				result = append(result, memory)
			}
		}
	} else {
		// 返回所有匹配 scope 的记忆
		for _, memory := range s.memories {
			if memory.Scope == scope || scope == "" {
				result = append(result, memory)
			}
		}
	}

	// 按创建时间排序（最新的在前）
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	// 限制数量
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// Search 搜索记忆
func (s *InMemoryStore) Search(ctx context.Context, query string, tags []string, limit int) ([]*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Memory
	queryLower := strings.ToLower(query)

	// 如果有标签过滤，先获取标签匹配的 ID
	var tagMatchIDs map[string]bool
	if len(tags) > 0 {
		tagMatchIDs = make(map[string]bool)
		for _, tag := range tags {
			if !strings.HasPrefix(tag, "#") {
				tag = "#" + tag
			}
			for _, id := range s.tagIndex[tag] {
				tagMatchIDs[id] = true
			}
		}
	}

	// 遍历所有记忆
	for _, memory := range s.memories {
		// 标签过滤
		if tagMatchIDs != nil && !tagMatchIDs[memory.ID] {
			continue
		}

		// 关键词搜索
		if query != "" {
			titleLower := strings.ToLower(memory.Title)
			contentLower := strings.ToLower(memory.Content)
			if !strings.Contains(titleLower, queryLower) && !strings.Contains(contentLower, queryLower) {
				continue
			}
		}

		result = append(result, memory)
	}

	// 按相关度排序（简单实现：按置信度和访问次数）
	sort.Slice(result, func(i, j int) bool {
		scoreI := result[i].Confidence + float64(result[i].AccessCount)*0.01
		scoreJ := result[j].Confidence + float64(result[j].AccessCount)*0.01
		return scoreI > scoreJ
	})

	// 限制数量
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// Delete 删除记忆
func (s *InMemoryStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	memory, ok := s.memories[id]
	if !ok {
		return nil // 不存在也不报错
	}

	// 从项目索引中删除
	if memory.ProjectID != "" {
		s.removeFromIndex(s.projectIndex, memory.ProjectID, id)
	}

	// 从标签索引中删除
	for _, tag := range memory.Tags {
		s.removeFromIndex(s.tagIndex, tag, id)
	}

	// 删除记忆
	delete(s.memories, id)

	return nil
}

// removeFromIndex 从索引中删除
func (s *InMemoryStore) removeFromIndex(index map[string][]string, key, id string) {
	ids := index[key]
	for i, existingID := range ids {
		if existingID == id {
			index[key] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
}

// removeOldest 删除最旧的记忆
func (s *InMemoryStore) removeOldest() {
	var oldestID string
	var oldestTime = s.memories[oldestID].CreatedAt

	for id, memory := range s.memories {
		if oldestID == "" || memory.CreatedAt.Before(oldestTime) {
			oldestID = id
			oldestTime = memory.CreatedAt
		}
	}

	if oldestID != "" {
		_ = s.Delete(context.Background(), oldestID)
	}
}

// GetStats 获取统计信息
func (s *InMemoryStore) GetStats() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 统计标签分布
	tagCounts := make(map[string]int)
	for tag, ids := range s.tagIndex {
		tagCounts[tag] = len(ids)
	}

	// 统计项目分布
	projectCounts := make(map[string]int)
	for projectID, ids := range s.projectIndex {
		projectCounts[projectID] = len(ids)
	}

	return map[string]any{
		"total_memories": len(s.memories),
		"tag_counts":     tagCounts,
		"project_counts": projectCounts,
	}
}

// 确保实现接口
var _ Store = (*InMemoryStore)(nil)
