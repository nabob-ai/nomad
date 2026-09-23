// Package examples 包含 Nomad 框架的示例代码和集成演示。
package examples

import (
	"context"
	"fmt"

	"github.com/nabob-ai/nomad/pkg/backends"
	"github.com/nabob-ai/nomad/pkg/middleware"
	"github.com/nabob-ai/nomad/pkg/types"
)

// Phase4IntegrationExample 展示如何集成 SummarizationMiddleware 和 AgentMemoryMiddleware
func Phase4IntegrationExample() error {
	_ = context.Background()

	// 1. 创建后端存储
	// 使用 CompositeBackend 支持多路径路由
	filesystemBackend := createFilesystemBackend()
	composite := backends.NewCompositeBackend(filesystemBackend, []backends.RouteConfig{
		{Prefix: "/memories/", Backend: createStoreBackend()},
	})

	// 2. 创建自定义 Summarizer (使用真实 LLM)
	customSummarizer := func(ctx context.Context, messages []types.Message) (string, error) {
		// 这里应该调用真实的 LLM API 生成总结
		// 示例使用 Provider 接口
		return generateSummaryWithLLM(ctx, messages)
	}

	// 3. 创建 SummarizationMiddleware
	summarizationMW, err := middleware.NewSummarizationMiddleware(&middleware.SummarizationMiddlewareConfig{
		Summarizer:             customSummarizer,
		MaxTokensBeforeSummary: 170000, // Claude 3.5 Sonnet 的 200k 上下文的 85%
		MessagesToKeep:         6,      // 保留最近 3 轮对话
		SummaryPrefix:          "## Previous conversation summary:",
	})
	if err != nil {
		return fmt.Errorf("failed to create summarization middleware: %w", err)
	}

	// 4. 创建 AgentMemoryMiddleware
	memoryMW, err := middleware.NewAgentMemoryMiddleware(&middleware.AgentMemoryMiddlewareConfig{
		Backend:    composite,
		MemoryPath: "/memories/",
		// 使用默认模板: <agent_memory>...</agent_memory>
	})
	if err != nil {
		return fmt.Errorf("failed to create agent memory middleware: %w", err)
	}

	// 5. 构建中间件栈
	middlewares := []middleware.Middleware{
		memoryMW,        // 优先级 5: 最早执行,注入记忆
		summarizationMW, // 优先级 40: 在调用模型前检查并总结
		// ... 其他中间件
	}

	// 6. 创建 Agent 并使用中间件
	fmt.Printf("Created agent with %d middlewares\n", len(middlewares))
	fmt.Printf("- SummarizationMiddleware: max_tokens=%d, keep_messages=%d\n",
		170000, 6)
	fmt.Printf("- AgentMemoryMiddleware: memory_path=%s\n", "/memories/")

	return nil
}

// generateSummaryWithLLM 使用 LLM 生成对话总结
func generateSummaryWithLLM(_ context.Context, messages []types.Message) (string, error) {
	// 这是一个示例实现
	// 实际应用中应该:
	// 1. 使用 Provider 接口调用 LLM
	// 2. 使用专门的总结提示词
	// 3. 设置较低的 temperature (如 0.3)

	// 构建总结请求
	_ = `Please provide a concise summary of the following conversation, capturing:
1. Main topics discussed
2. Important decisions or conclusions
3. Action items or next steps
4. Relevant technical details or constraints

Keep the summary focused and informative, around 200-300 words.`

	// TODO: 实际应该组装消息并调用LLM
	// Example:
	// summaryMessages := []types.Message{
	//     {Role: types.MessageRoleSystem, Content: summaryPrompt},
	// }
	// summaryMessages = append(summaryMessages, messages...)
	// resp, err := provider.Stream(ctx, summaryMessages, &provider.StreamOptions{
	//     Temperature: 0.3,
	//     MaxTokens:   500,
	// })
	_ = messages

	// 示例返回
	return "Summary of conversation: [Topics discussed, decisions made, next steps...]", nil
}

// createStoreBackend 创建 Store Backend (持久化存储)
func createStoreBackend() backends.BackendProtocol {
	// 实际应用中应该使用真实的 Store Backend
	// 这里返回 nil 仅用于示例
	return nil
}

// createFilesystemBackend 创建 Filesystem Backend
func createFilesystemBackend() backends.BackendProtocol {
	// 实际应用中应该使用真实的 Filesystem Backend
	// 这里返回 nil 仅用于示例
	return nil
}

// AgentMemoryExample 展示如何管理 Agent 记忆
func AgentMemoryExample() {
	// agent.md 文件示例内容:
	agentMemoryExample := `# Agent Personality

You are Claude, a helpful AI coding assistant created by Anthropic.

## Core Principles

1. **Code Quality First**: Always write clean, well-documented code
2. **Test-Driven**: Write tests before implementing features
3. **Security Conscious**: Check for common vulnerabilities (SQL injection, XSS, etc.)

## User Preferences

- **Language**: Prefers Go over Python for backend services
- **Testing**: Uses table-driven tests in Go
- **Documentation**: Likes detailed inline comments
- **Code Style**: Follows official Go style guide

## Project Context

- Working on nomad, a Go-based agent framework
- Recently implemented Phase 4 features (Summarization + Memory)
- Focus on production-ready, well-tested code

## Learning from Feedback

- User prefers progressive enhancement over big-bang changes
- User values incremental commits with clear messages
- User wants all tests to pass before committing`

	fmt.Println("=== Agent Memory Example ===")
	fmt.Println("\nCreate /agent.md with the following content:")
	fmt.Println(agentMemoryExample)
	fmt.Println("\nThe AgentMemoryMiddleware will:")
	fmt.Println("1. Load this file on agent start")
	fmt.Println("2. Inject it into every system prompt")
	fmt.Println("3. Add usage guidelines for long-term memory")
}

// AdvancedSummarizationExample 展示高级总结策略
func AdvancedSummarizationExample() {
	fmt.Println("=== Advanced Summarization Strategies ===")

	fmt.Println("\n1. Dynamic Threshold Based on Task Complexity:")
	fmt.Println("   - Simple tasks: 100k tokens")
	fmt.Println("   - Complex tasks: 150k tokens")
	fmt.Println("   - Code review: 180k tokens (need more context)")

	fmt.Println("\n2. Smart Message Retention:")
	fmt.Println("   - Keep recent 6 messages by default")
	fmt.Println("   - Keep more if containing important context:")
	fmt.Println("     * Error messages and stack traces")
	fmt.Println("     * User requirements or specifications")
	fmt.Println("     * Critical decisions")

	fmt.Println("\n3. Multi-Stage Summarization:")
	fmt.Println("   - Stage 1: Summarize every 50k tokens")
	fmt.Println("   - Stage 2: Re-summarize summaries at 150k tokens")
	fmt.Println("   - Keeps context compressed while preserving key info")

	fmt.Println("\n4. Custom Token Counting:")
	fmt.Println("   - Use model's official tokenizer for accuracy")
	fmt.Println("   - Account for tool definitions in token count")
	fmt.Println("   - Include system prompt in calculation")
}

// MonitoringExample 展示如何监控中间件性能
func MonitoringExample(summarizationMW *middleware.SummarizationMiddleware) {
	fmt.Println("=== Monitoring Middleware Performance ===")

	// 获取配置
	config := summarizationMW.GetConfig()
	fmt.Printf("\nConfiguration:\n")
	fmt.Printf("  Max Tokens: %v\n", config["max_tokens_before_summary"])
	fmt.Printf("  Messages to Keep: %v\n", config["messages_to_keep"])
	fmt.Printf("  Summary Prefix: %v\n", config["summary_prefix"])

	// 获取统计信息
	count := summarizationMW.GetSummarizationCount()
	fmt.Printf("\nStatistics:\n")
	fmt.Printf("  Total Summarizations: %d\n", count)

	// 如果总结次数过多,可能需要调整阈值
	if count > 10 {
		fmt.Println("\n⚠️  Warning: High summarization count detected")
		fmt.Println("   Consider increasing max_tokens_before_summary")
		fmt.Println("   or reducing messages_to_keep")
	}
}

// BestPractices 最佳实践指南
func BestPractices() {
	fmt.Println("=== Phase 4 Middleware Best Practices ===")

	fmt.Println("\n📋 SummarizationMiddleware:")
	fmt.Println("  1. Use real LLM for summarization in production")
	fmt.Println("  2. Set max_tokens to 85% of model's context window")
	fmt.Println("  3. Monitor summarization count to tune parameters")
	fmt.Println("  4. Consider using cheaper models for summarization")
	fmt.Println("  5. Keep summaries concise (200-300 words)")

	fmt.Println("\n🧠 AgentMemoryMiddleware:")
	fmt.Println("  1. Keep agent.md focused and well-structured")
	fmt.Println("  2. Update regularly based on user feedback")
	fmt.Println("  3. Version control agent.md with git")
	fmt.Println("  4. Use CompositeBackend for flexible storage")
	fmt.Println("  5. Test agent behavior after updating memory")

	fmt.Println("\n🔧 Integration:")
	fmt.Println("  1. Place AgentMemoryMiddleware early (low priority number)")
	fmt.Println("  2. Place SummarizationMiddleware before model call")
	fmt.Println("  3. Test with long conversations to verify behavior")
	fmt.Println("  4. Monitor token usage and adjust thresholds")
	fmt.Println("  5. Implement graceful degradation on errors")

	fmt.Println("\n🎯 Production Checklist:")
	fmt.Println("  [ ] Custom summarizer implemented with real LLM")
	fmt.Println("  [ ] agent.md created and tested")
	fmt.Println("  [ ] Token thresholds tuned for your model")
	fmt.Println("  [ ] Monitoring and alerting set up")
	fmt.Println("  [ ] Error handling tested (LLM failure, file not found)")
	fmt.Println("  [ ] Performance benchmarked under load")
}
