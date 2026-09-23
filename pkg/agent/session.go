package agent

import (
	"context"
	"time"

	"github.com/nabob-ai/nomad/pkg/logging"
	"github.com/nabob-ai/nomad/pkg/types"
)

var sessionLog = logging.ForComponent("AgentSession")

// GetSessionState 获取会话状态信息
// 返回当前会话的状态，用于判断是否需要恢复
func (a *Agent) GetSessionState(ctx context.Context) *types.SessionState {
	a.mu.RLock()
	defer a.mu.RUnlock()

	state := &types.SessionState{
		HasHistory:   len(a.messages) > 0,
		MessageCount: len(a.messages),
		IsResumed:    false, // 默认不是恢复的会话
	}

	return state
}

// ProcessRecovery 处理会话恢复
// 如果注册了 RecoveryHook，则调用它来决定是否需要恢复以及如何恢复
func (a *Agent) ProcessRecovery(ctx context.Context, message string, workDir string, metadata map[string]any) (string, bool) {
	if a.deps.RecoveryHook == nil {
		return message, false
	}

	// 获取会话状态
	sessionState := a.GetSessionState(ctx)

	// 检查是否是恢复的会话（有历史消息且 Agent 是刚创建的）
	// 判断标准：Agent 创建时间在最近 5 秒内，但有历史消息
	isResumed := sessionState.HasHistory && time.Since(a.createdAt) < 5*time.Second
	sessionState.IsResumed = isResumed

	if !isResumed {
		// 不是恢复的会话，不需要处理
		return message, false
	}

	// 构建恢复上下文
	recoveryCtx := &types.RecoveryContext{
		AgentID:         a.id,
		SessionState:    sessionState,
		OriginalMessage: message,
		WorkDir:         workDir,
		Metadata:        metadata,
	}

	// 调用恢复钩子
	result := a.deps.RecoveryHook(recoveryCtx)
	if result == nil || !result.ShouldRecover {
		return message, false
	}

	sessionLog.Info(ctx, "session recovery triggered", map[string]any{"agent_id": a.id})

	// 返回增强后的消息
	if result.EnhancedMessage != "" {
		return result.EnhancedMessage, true
	}

	// 如果只提供了 Instructions，则拼接到原始消息前
	if result.Instructions != "" {
		return result.Instructions + "\n\n---\n\n" + message, true
	}

	return message, false
}

// IsResumedSession 检查当前是否是恢复的会话
func (a *Agent) IsResumedSession() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// 判断标准：Agent 创建时间在最近 5 秒内，但有历史消息
	return len(a.messages) > 0 && time.Since(a.createdAt) < 5*time.Second
}

// GetWorkDir 获取工作目录
func (a *Agent) GetWorkDir() string {
	if a.sandbox != nil {
		return a.sandbox.WorkDir()
	}
	return ""
}
