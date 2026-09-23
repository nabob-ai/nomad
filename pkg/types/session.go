package types

// SessionState 会话状态信息
type SessionState struct {
	// HasHistory 是否有历史消息
	HasHistory bool `json:"has_history"`
	// MessageCount 历史消息数量
	MessageCount int `json:"message_count"`
	// IsResumed 是否是恢复的会话（Agent 实例是新创建的，但有历史数据）
	IsResumed bool `json:"is_resumed"`
	// LastMessageTime 最后一条消息的时间
	LastMessageTime *string `json:"last_message_time,omitempty"`
}

// RecoveryConfig 会话恢复配置
type RecoveryConfig struct {
	// Enabled 是否启用恢复机制
	Enabled bool `json:"enabled"`
	// TriggerPatterns 触发恢复的消息模式（正则表达式）
	// 当用户消息匹配这些模式时，认为用户想继续之前的工作
	TriggerPatterns []string `json:"trigger_patterns,omitempty"`
	// MaxTriggerLength 触发恢复的最大消息长度
	// 超过此长度的消息不被视为恢复触发（因为可能是新的详细指令）
	MaxTriggerLength int `json:"max_trigger_length,omitempty"`
}

// RecoveryContext 恢复上下文（传递给 RecoveryHook）
type RecoveryContext struct {
	// AgentID Agent 标识
	AgentID string `json:"agent_id"`
	// SessionState 会话状态
	SessionState *SessionState `json:"session_state"`
	// OriginalMessage 用户原始消息
	OriginalMessage string `json:"original_message"`
	// WorkDir 工作目录（如果有）
	WorkDir string `json:"work_dir,omitempty"`
	// Metadata 自定义元数据（应用层可以传入额外信息）
	Metadata map[string]any `json:"metadata,omitempty"`
}

// RecoveryResult 恢复结果
type RecoveryResult struct {
	// ShouldRecover 是否应该执行恢复
	ShouldRecover bool `json:"should_recover"`
	// EnhancedMessage 增强后的消息（如果需要恢复）
	EnhancedMessage string `json:"enhanced_message,omitempty"`
	// Instructions 恢复指令（可选，会被添加到消息前）
	Instructions string `json:"instructions,omitempty"`
}

// RecoveryHook 恢复钩子函数类型
// 应用层实现这个函数来定义自己的恢复逻辑
type RecoveryHook func(ctx *RecoveryContext) *RecoveryResult
