package logic

import (
	"context"
)

// PatternMatcher 模式识别接口（应用层实现）
// 用于从事件中识别和提取 Logic Memory
type PatternMatcher interface {
	// MatchEvent 从事件中识别 Memory
	// 返回识别出的 Memory 列表（可能为空）
	MatchEvent(ctx context.Context, event Event) ([]*LogicMemory, error)

	// SupportedEventTypes 返回支持的事件类型列表
	// 用于 Manager 筛选事件，避免不必要的调用
	SupportedEventTypes() []string
}

// NoopMatcher 空实现（用于测试或默认场景）
type NoopMatcher struct{}

// MatchEvent 实现 PatternMatcher 接口
func (m *NoopMatcher) MatchEvent(ctx context.Context, event Event) ([]*LogicMemory, error) {
	return nil, nil
}

// SupportedEventTypes 实现 PatternMatcher 接口
func (m *NoopMatcher) SupportedEventTypes() []string {
	return []string{}
}

// 确保 NoopMatcher 实现 PatternMatcher 接口
var _ PatternMatcher = (*NoopMatcher)(nil)
