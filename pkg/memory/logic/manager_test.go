package logic

import (
	"context"
	"testing"
	"time"

	"github.com/nabob-ai/nomad/pkg/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		store := NewInMemoryStore()
		config := &ManagerConfig{
			Store: store,
		}

		manager, err := NewManager(config)
		require.NoError(t, err)
		require.NotNil(t, manager)
		assert.Equal(t, store, manager.store)
	})

	t.Run("nil config", func(t *testing.T) {
		_, err := NewManager(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config is required")
	})

	t.Run("nil store", func(t *testing.T) {
		config := &ManagerConfig{}
		_, err := NewManager(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "store is required")
	})

	t.Run("default values", func(t *testing.T) {
		store := NewInMemoryStore()
		config := &ManagerConfig{
			Store: store,
		}

		manager, err := NewManager(config)
		require.NoError(t, err)
		assert.Equal(t, 0.85, manager.config.ConsolidationThreshold)
		assert.Equal(t, 0.05, manager.config.ConfidenceBoost)
	})
}

func TestRecordMemory(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStore()
	manager, err := NewManager(&ManagerConfig{Store: store})
	require.NoError(t, err)

	t.Run("create new memory", func(t *testing.T) {
		mem := &LogicMemory{
			Namespace:   "user:123",
			Scope:       ScopeUser,
			Type:        "preference",
			Key:         "tone",
			Value:       "casual",
			Description: "User prefers casual tone",
			Provenance: &memory.MemoryProvenance{
				SourceType: memory.SourceUserInput,
				Confidence: 0.7,
			},
		}

		err := manager.RecordMemory(ctx, mem)
		require.NoError(t, err)

		// 验证 ID 已设置
		assert.NotEmpty(t, mem.ID)

		// 验证存储成功
		retrieved, err := store.Get(ctx, "user:123", "tone")
		require.NoError(t, err)
		assert.Equal(t, "casual", retrieved.Value)
		assert.Equal(t, "User prefers casual tone", retrieved.Description)
	})

	t.Run("update existing memory", func(t *testing.T) {
		// 创建初始 Memory
		mem1 := &LogicMemory{
			Namespace:   "user:456",
			Scope:       ScopeUser,
			Type:        "preference",
			Key:         "style",
			Value:       "formal",
			Description: "User prefers formal style",
			Provenance: &memory.MemoryProvenance{
				SourceType: memory.SourceUserInput,
				Confidence: 0.6,
			},
		}
		err := manager.RecordMemory(ctx, mem1)
		require.NoError(t, err)

		// 更新 Memory（相同 namespace + key）
		mem2 := &LogicMemory{
			Namespace:   "user:456",
			Scope:       ScopeUser,
			Type:        "preference",
			Key:         "style",
			Value:       "casual",
			Description: "User changed to casual style",
			Provenance: &memory.MemoryProvenance{
				SourceType: memory.SourceUserInput,
				Confidence: 0.7,
			},
		}
		err = manager.RecordMemory(ctx, mem2)
		require.NoError(t, err)

		// 验证置信度提升
		retrieved, err := store.Get(ctx, "user:456", "style")
		require.NoError(t, err)
		assert.Greater(t, retrieved.Provenance.Confidence, 0.6)
		assert.Equal(t, "casual", retrieved.Value)
	})
}

func TestProcessEvent(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStore()

	t.Run("with valid matcher", func(t *testing.T) {
		// 创建测试 Matcher
		matcher := &testMatcher{
			supportedTypes: []string{"user_action"},
			memories: []*LogicMemory{
				{
					Namespace:   "user:789",
					Scope:       ScopeUser,
					Type:        "behavior",
					Key:         "click_pattern",
					Value:       "fast_clicker",
					Description: "User clicks quickly",
					Provenance: &memory.MemoryProvenance{
						SourceType: memory.SourceUserInput,
						Confidence: 0.75,
					},
				},
			},
		}

		manager, err := NewManager(&ManagerConfig{
			Store:    store,
			Matchers: []PatternMatcher{matcher},
		})
		require.NoError(t, err)

		// 处理事件
		event := Event{
			Type:      "user_action",
			Source:    "user:789",
			Data:      map[string]any{"action": "click"},
			Timestamp: time.Now(),
		}

		err = manager.ProcessEvent(ctx, event)
		require.NoError(t, err)

		// 验证 Memory 已创建
		retrieved, err := store.Get(ctx, "user:789", "click_pattern")
		require.NoError(t, err)
		assert.Equal(t, "fast_clicker", retrieved.Value)
		assert.Equal(t, "User clicks quickly", retrieved.Description)
	})

	t.Run("no matchers", func(t *testing.T) {
		manager, err := NewManager(&ManagerConfig{
			Store:    store,
			Matchers: nil, // 没有 Matcher
		})
		require.NoError(t, err)

		event := Event{
			Type:      "unknown",
			Source:    "user:999",
			Data:      map[string]any{},
			Timestamp: time.Now(),
		}

		// 不应返回错误
		err = manager.ProcessEvent(ctx, event)
		assert.NoError(t, err)
	})

	t.Run("unsupported event type", func(t *testing.T) {
		matcher := &testMatcher{
			supportedTypes: []string{"user_action"},
			memories:       []*LogicMemory{},
		}

		manager, err := NewManager(&ManagerConfig{
			Store:    store,
			Matchers: []PatternMatcher{matcher},
		})
		require.NoError(t, err)

		event := Event{
			Type:      "unsupported_type",
			Source:    "user:999",
			Data:      map[string]any{},
			Timestamp: time.Now(),
		}

		// 不应调用 Matcher，不返回错误
		err = manager.ProcessEvent(ctx, event)
		assert.NoError(t, err)
	})
}

func TestRetrieveMemories(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStore()
	manager, err := NewManager(&ManagerConfig{Store: store})
	require.NoError(t, err)

	// 创建测试数据
	memories := []*LogicMemory{
		{
			ID:          "1",
			Namespace:   "user:123",
			Scope:       ScopeUser,
			Type:        "preference",
			Key:         "tone",
			Value:       "casual",
			Description: "Casual tone",
			Provenance:  &memory.MemoryProvenance{Confidence: 0.9},
		},
		{
			ID:          "2",
			Namespace:   "user:123",
			Scope:       ScopeUser,
			Type:        "preference",
			Key:         "style",
			Value:       "concise",
			Description: "Concise style",
			Provenance:  &memory.MemoryProvenance{Confidence: 0.7},
		},
		{
			ID:          "3",
			Namespace:   "user:456",
			Scope:       ScopeUser,
			Type:        "preference",
			Key:         "language",
			Value:       "en",
			Description: "English",
			Provenance:  &memory.MemoryProvenance{Confidence: 0.8},
		},
	}

	for _, mem := range memories {
		err := store.Save(ctx, mem)
		require.NoError(t, err)
	}

	t.Run("retrieve all for namespace", func(t *testing.T) {
		retrieved, err := manager.RetrieveMemories(ctx, "user:123")
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
	})

	t.Run("retrieve with TopK", func(t *testing.T) {
		retrieved, err := manager.RetrieveMemories(ctx, "user:123", WithTopK(1))
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		// 应该返回置信度最高的
		assert.Equal(t, "tone", retrieved[0].Key)
	})

	t.Run("retrieve with MinConfidence", func(t *testing.T) {
		retrieved, err := manager.RetrieveMemories(ctx, "user:123", WithMinConfidence(0.8))
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, "tone", retrieved[0].Key)
	})

	t.Run("retrieve increments access count", func(t *testing.T) {
		// 先获取当前访问次数
		mem, err := store.Get(ctx, "user:123", "tone")
		require.NoError(t, err)
		oldCount := mem.AccessCount

		// 检索 Memory
		_, err = manager.RetrieveMemories(ctx, "user:123")
		require.NoError(t, err)

		// 等待异步更新
		time.Sleep(100 * time.Millisecond)

		// 验证访问次数增加
		mem, err = store.Get(ctx, "user:123", "tone")
		require.NoError(t, err)
		assert.Greater(t, mem.AccessCount, oldCount)
	})
}

func TestMergeMemory(t *testing.T) {
	store := NewInMemoryStore()
	manager, err := NewManager(&ManagerConfig{
		Store:           store,
		ConfidenceBoost: 0.1,
	})
	require.NoError(t, err)

	t.Run("boost confidence", func(t *testing.T) {
		existing := &LogicMemory{
			AccessCount: 5,
			Provenance:  &memory.MemoryProvenance{Confidence: 0.7},
		}

		new := &LogicMemory{
			Provenance: &memory.MemoryProvenance{Confidence: 0.8},
		}

		manager.mergeMemory(existing, new)

		// 访问次数应该增加
		assert.Equal(t, 6, existing.AccessCount)

		// 置信度应该提升（使用 InDelta 处理浮点数精度）
		assert.InDelta(t, 0.8, existing.Provenance.Confidence, 0.0001)
	})

	t.Run("merge sources", func(t *testing.T) {
		existing := &LogicMemory{
			Provenance: &memory.MemoryProvenance{
				Confidence: 0.7,
				Sources:    []string{"source1"},
			},
		}

		new := &LogicMemory{
			Provenance: &memory.MemoryProvenance{
				Confidence: 0.8,
				Sources:    []string{"source2", "source3"},
			},
		}

		manager.mergeMemory(existing, new)

		// Sources 应该合并
		assert.Len(t, existing.Provenance.Sources, 3)
		assert.Contains(t, existing.Provenance.Sources, "source1")
		assert.Contains(t, existing.Provenance.Sources, "source2")
	})

	t.Run("update description if longer", func(t *testing.T) {
		existing := &LogicMemory{
			Description: "Short",
			Provenance:  &memory.MemoryProvenance{Confidence: 0.7},
		}

		new := &LogicMemory{
			Description: "Much longer description",
			Provenance:  &memory.MemoryProvenance{Confidence: 0.8},
		}

		manager.mergeMemory(existing, new)

		assert.Equal(t, "Much longer description", existing.Description)
	})

	t.Run("merge metadata", func(t *testing.T) {
		existing := &LogicMemory{
			Metadata:   map[string]any{"key1": "value1"},
			Provenance: &memory.MemoryProvenance{Confidence: 0.7},
		}

		new := &LogicMemory{
			Metadata:   map[string]any{"key2": "value2", "key3": "value3"},
			Provenance: &memory.MemoryProvenance{Confidence: 0.8},
		}

		manager.mergeMemory(existing, new)

		assert.Len(t, existing.Metadata, 3)
		assert.Equal(t, "value1", existing.Metadata["key1"])
		assert.Equal(t, "value2", existing.Metadata["key2"])
	})
}

func TestPruneMemories(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStore()
	manager, err := NewManager(&ManagerConfig{Store: store})
	require.NoError(t, err)

	// 创建测试数据
	now := time.Now()
	memories := []*LogicMemory{
		{
			ID:           "1",
			Namespace:    "user:123",
			Key:          "high_confidence",
			Provenance:   &memory.MemoryProvenance{Confidence: 0.9},
			AccessCount:  10,
			LastAccessed: now,
			CreatedAt:    now.Add(-24 * time.Hour),
		},
		{
			ID:           "2",
			Namespace:    "user:123",
			Key:          "low_confidence",
			Provenance:   &memory.MemoryProvenance{Confidence: 0.3},
			AccessCount:  1,
			LastAccessed: now.Add(-48 * time.Hour),
			CreatedAt:    now.Add(-72 * time.Hour),
		},
		{
			ID:           "3",
			Namespace:    "user:123",
			Key:          "old_unused",
			Provenance:   &memory.MemoryProvenance{Confidence: 0.6},
			AccessCount:  0,
			LastAccessed: now.Add(-96 * time.Hour),
			CreatedAt:    now.Add(-120 * time.Hour),
		},
	}

	for _, mem := range memories {
		err := store.Save(ctx, mem)
		require.NoError(t, err)
	}

	t.Run("prune by low confidence", func(t *testing.T) {
		criteria := PruneCriteria{
			MinConfidence: 0.5,
		}

		count, err := manager.PruneMemories(ctx, criteria)
		require.NoError(t, err)
		assert.Equal(t, 1, count) // 应该清理 low_confidence

		// 验证
		_, err = store.Get(ctx, "user:123", "low_confidence")
		assert.Error(t, err)
	})

	t.Run("prune by last access", func(t *testing.T) {
		// 重新创建数据
		for _, mem := range memories {
			err := store.Save(ctx, mem)
			require.NoError(t, err)
		}

		criteria := PruneCriteria{
			SinceLastAccess: 72 * time.Hour,
		}

		count, err := manager.PruneMemories(ctx, criteria)
		require.NoError(t, err)
		assert.Equal(t, 1, count) // 应该清理 old_unused
	})
}

// testMatcher 测试用的 PatternMatcher
type testMatcher struct {
	supportedTypes []string
	memories       []*LogicMemory
	err            error
}

func (m *testMatcher) MatchEvent(ctx context.Context, event Event) ([]*LogicMemory, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.memories, nil
}

func (m *testMatcher) SupportedEventTypes() []string {
	return m.supportedTypes
}
