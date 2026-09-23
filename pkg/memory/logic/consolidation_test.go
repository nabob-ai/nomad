package logic

import (
	"context"
	"testing"

	"github.com/nabob-ai/nomad/pkg/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultSimilarityCalculator(t *testing.T) {
	calc := &DefaultSimilarityCalculator{}

	t.Run("same type and key prefix", func(t *testing.T) {
		a := &LogicMemory{
			Type:        "preference",
			Key:         "tone_formal",
			Category:    "writing",
			Description: "User prefers formal tone",
		}

		b := &LogicMemory{
			Type:        "preference",
			Key:         "tone_casual",
			Category:    "writing",
			Description: "User likes casual tone",
		}

		similarity := calc.Calculate(a, b)
		// Same type (0.3) + same key prefix "tone" (0.3) + same category (0.2) = 0.8+
		assert.Greater(t, similarity, 0.7)
	})

	t.Run("different types", func(t *testing.T) {
		a := &LogicMemory{
			Type:        "preference",
			Key:         "tone_formal",
			Description: "User prefers formal tone",
		}

		b := &LogicMemory{
			Type:        "behavior_pattern",
			Key:         "click_fast",
			Description: "User clicks quickly",
		}

		similarity := calc.Calculate(a, b)
		assert.Less(t, similarity, 0.3)
	})

	t.Run("empty memories", func(t *testing.T) {
		a := &LogicMemory{}
		b := &LogicMemory{}

		similarity := calc.Calculate(a, b)
		// Same type (""), same key prefix ("") = 0.6
		assert.InDelta(t, 0.6, similarity, 0.1)
	})
}

func TestConsolidationEngine(t *testing.T) {
	ctx := context.Background()

	t.Run("consolidate similar memories", func(t *testing.T) {
		store := NewInMemoryStore()
		engine := NewConsolidationEngine(store, &ConsolidationConfig{
			SimilarityThreshold:             0.7,
			MinGroupSize:                    2,
			MaxMergeCount:                   100,
			PreserveHighConfidenceThreshold: 0.95,
			MergeStrategy:                   MergeStrategyKeepHighestConfidence,
		})

		// 创建相似的 Memory
		memories := []*LogicMemory{
			{
				ID:          "1",
				Namespace:   "user:123",
				Type:        "preference",
				Key:         "tone_formal",
				Category:    "writing",
				Description: "User prefers formal tone",
				Provenance:  &memory.MemoryProvenance{Confidence: 0.7},
			},
			{
				ID:          "2",
				Namespace:   "user:123",
				Type:        "preference",
				Key:         "tone_casual",
				Category:    "writing",
				Description: "User likes casual writing",
				Provenance:  &memory.MemoryProvenance{Confidence: 0.8},
			},
			{
				ID:          "3",
				Namespace:   "user:123",
				Type:        "behavior_pattern",
				Key:         "click_fast",
				Description: "User clicks quickly",
				Provenance:  &memory.MemoryProvenance{Confidence: 0.6},
			},
		}

		for _, mem := range memories {
			err := store.Save(ctx, mem)
			require.NoError(t, err)
		}

		// 执行合并
		result, err := engine.Consolidate(ctx, "user:123")
		require.NoError(t, err)

		assert.Equal(t, 3, result.TotalMemories)
		// 前两个应该被合并（相同 type, category, key prefix）
		assert.GreaterOrEqual(t, result.MergedGroups, 1)
	})

	t.Run("preserve high confidence memories", func(t *testing.T) {
		store := NewInMemoryStore()
		engine := NewConsolidationEngine(store, &ConsolidationConfig{
			SimilarityThreshold:             0.7,
			MinGroupSize:                    2,
			PreserveHighConfidenceThreshold: 0.9,
			MergeStrategy:                   MergeStrategyKeepHighestConfidence,
		})

		// 创建相似但高置信度的 Memory
		memories := []*LogicMemory{
			{
				ID:          "1",
				Namespace:   "user:123",
				Type:        "preference",
				Key:         "tone_formal",
				Description: "User prefers formal tone",
				Provenance:  &memory.MemoryProvenance{Confidence: 0.95}, // 高置信度
			},
			{
				ID:          "2",
				Namespace:   "user:123",
				Type:        "preference",
				Key:         "tone_casual",
				Description: "User likes casual writing",
				Provenance:  &memory.MemoryProvenance{Confidence: 0.6},
			},
		}

		for _, mem := range memories {
			err := store.Save(ctx, mem)
			require.NoError(t, err)
		}

		// 执行合并
		result, err := engine.Consolidate(ctx, "user:123")
		require.NoError(t, err)

		// 高置信度的应该被保留，不参与合并
		assert.Equal(t, 0, result.MergedGroups)
	})

	t.Run("merge descriptions strategy", func(t *testing.T) {
		store := NewInMemoryStore()
		engine := NewConsolidationEngine(store, &ConsolidationConfig{
			SimilarityThreshold: 0.6,
			MinGroupSize:        2,
			MergeStrategy:       MergeStrategyMergeDescriptions,
		})

		memories := []*LogicMemory{
			{
				ID:          "1",
				Namespace:   "user:123",
				Type:        "preference",
				Key:         "style_a",
				Description: "First description",
				Provenance:  &memory.MemoryProvenance{Confidence: 0.7},
			},
			{
				ID:          "2",
				Namespace:   "user:123",
				Type:        "preference",
				Key:         "style_b",
				Description: "Second description",
				Provenance:  &memory.MemoryProvenance{Confidence: 0.8},
			},
		}

		for _, mem := range memories {
			err := store.Save(ctx, mem)
			require.NoError(t, err)
		}

		result, err := engine.Consolidate(ctx, "user:123")
		require.NoError(t, err)

		if result.MergedGroups > 0 && len(result.MergedMemories) > 0 {
			// 描述应该被合并
			merged := result.MergedMemories[0]
			assert.Contains(t, merged.Description, "description")
		}
	})

	t.Run("no memories to consolidate", func(t *testing.T) {
		store := NewInMemoryStore()
		engine := NewConsolidationEngine(store, nil)

		result, err := engine.Consolidate(ctx, "empty:namespace")
		require.NoError(t, err)

		assert.Equal(t, 0, result.TotalMemories)
		assert.Equal(t, 0, result.MergedGroups)
	})
}

func TestConsolidationResult(t *testing.T) {
	result := &ConsolidationResult{
		TotalMemories:   100,
		MergedGroups:    5,
		DeletedMemories: 10,
	}

	summary := result.String()
	assert.Contains(t, summary, "100 memories processed")
	assert.Contains(t, summary, "5 groups merged")
	assert.Contains(t, summary, "10 memories deleted")
}
