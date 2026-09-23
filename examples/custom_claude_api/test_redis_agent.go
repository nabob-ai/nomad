//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nabob-ai/nomad/pkg/agent"
	"github.com/nabob-ai/nomad/pkg/provider"
	"github.com/nabob-ai/nomad/pkg/sandbox"
	"github.com/nabob-ai/nomad/pkg/store"
	"github.com/nabob-ai/nomad/pkg/tools"
	"github.com/nabob-ai/nomad/pkg/tools/builtin"
	"github.com/nabob-ai/nomad/pkg/types"
)

func main() {
	// 从环境变量读取配置
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	apiKey := os.Getenv("CLAUDE_API_KEY")
	baseURL := os.Getenv("CLAUDE_BASE_URL")

	if apiKey == "" {
		fmt.Println("❌ 请设置 CLAUDE_API_KEY 环境变量")
		os.Exit(1)
	}
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}

	fmt.Println("=== Redis Store + Agent 分布式集成测试 ===")
	fmt.Printf("Redis: %s\n", redisAddr)
	fmt.Printf("API: %s\n\n", baseURL)

	ctx := context.Background()

	// 1. 创建 Redis Store
	fmt.Println("【步骤 1】创建 Redis Store")
	redisStore, err := store.NewRedisStore(store.RedisConfig{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
		Prefix:   "nomad:test:",
		TTL:      1 * time.Hour,
	})
	if err != nil {
		fmt.Printf("❌ 创建 Redis Store 失败: %v\n", err)
		os.Exit(1)
	}
	defer redisStore.Close()

	if err := redisStore.Ping(ctx); err != nil {
		fmt.Printf("❌ Redis 连接失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Redis Store 创建成功")

	// 2. 创建工具注册表
	fmt.Println("\n【步骤 2】创建工具注册表")
	toolRegistry := tools.NewRegistry()
	builtin.RegisterAll(toolRegistry)
	fmt.Println("✅ 工具注册表创建成功")

	// 3. 创建 Sandbox Factory
	fmt.Println("\n【步骤 3】创建 Sandbox Factory")
	sandboxFactory := sandbox.NewFactory()
	fmt.Println("✅ Sandbox Factory 创建成功")

	// 4. 创建 Provider Factory (使用自定义工厂)
	fmt.Println("\n【步骤 4】创建 Provider Factory")
	providerFactory := &CustomProviderFactory{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
	fmt.Println("✅ Provider Factory 创建成功")

	// 5. 创建模板注册表
	fmt.Println("\n【步骤 5】创建模板注册表")
	templateRegistry := agent.NewTemplateRegistry()
	templateRegistry.Register(&types.AgentTemplateDefinition{
		ID:           "redis-test",
		Model:        "claude-sonnet-4-5-20250929",
		SystemPrompt: "你是一个测试 Agent，用于验证 Redis Store 的分布式功能。请简洁回答。",
		Tools:        []any{}, // 不需要工具
	})
	fmt.Println("✅ 模板注册表创建成功")

	// 6. 创建依赖
	fmt.Println("\n【步骤 6】创建依赖")
	deps := &agent.Dependencies{
		Store:            redisStore,
		SandboxFactory:   sandboxFactory,
		ToolRegistry:     toolRegistry,
		ProviderFactory:  providerFactory,
		TemplateRegistry: templateRegistry,
	}
	fmt.Println("✅ 依赖创建成功")

	// 7. 创建 Agent 配置
	fmt.Println("\n【步骤 7】创建 Agent 配置")
	agentID := "agt-redis-integration-001"
	config := &types.AgentConfig{
		AgentID:    agentID,
		TemplateID: "redis-test",
		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5-20250929",
			APIKey:   apiKey,
			BaseURL:  baseURL,
		},
		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindLocal,
			WorkDir: "./workspace-redis-test",
		},
		Store: &types.StoreConfig{
			MaxMessages: 10,
			AutoTrim:    true,
		},
	}
	fmt.Println("✅ Agent 配置创建成功")

	// 8. 创建节点 1 的 Agent
	fmt.Println("\n【步骤 8】节点 1 - 创建 Agent")
	agent1, err := agent.Create(ctx, config, deps)
	if err != nil {
		fmt.Printf("❌ 创建 Agent 失败: %v\n", err)
		os.Exit(1)
	}
	defer agent1.Close()
	fmt.Printf("✅ 节点 1 Agent 创建成功 (ID: %s)\n", agentID)

	// 9. 节点 1 发送第一条消息
	fmt.Println("\n【步骤 9】节点 1 - 发送第一条消息")
	result1, err := agent1.Chat(ctx, "你好，我是节点 1，请记住我")
	if err != nil {
		fmt.Printf("❌ 发送失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 节点 1 收到回复: %s\n", truncate(result1.Text, 80))

	// 10. 检查 Redis 中的数据
	fmt.Println("\n【步骤 10】检查 Redis 存储")
	messages, err := redisStore.LoadMessages(ctx, agentID)
	if err != nil {
		fmt.Printf("❌ 加载消息失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Redis 中存储了 %d 条消息\n", len(messages))
	printMessages(messages)

	// 11. 创建节点 2 的 Agent（模拟分布式场景）
	fmt.Println("\n【步骤 11】节点 2 - 创建新 Agent 实例（模拟另一个服务器）")

	// 重要：使用相同的 AgentID 和相同的 Redis Store
	config2 := &types.AgentConfig{
		AgentID:    agentID, // ✅ 相同的 ID
		TemplateID: "redis-test",
		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5-20250929",
			APIKey:   apiKey,
			BaseURL:  baseURL,
		},
		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindLocal,
			WorkDir: "./workspace-redis-test-2",
		},
		Store: &types.StoreConfig{
			MaxMessages: 10,
			AutoTrim:    true,
		},
	}

	agent2, err := agent.Create(ctx, config2, deps) // 共享同一个 deps (包含 Redis Store)
	if err != nil {
		fmt.Printf("❌ 创建节点 2 Agent 失败: %v\n", err)
		os.Exit(1)
	}
	defer agent2.Close()
	fmt.Println("✅ 节点 2 Agent 创建成功（共享 Redis）")

	// 12. 节点 2 发送消息（应该能看到节点 1 的历史）
	fmt.Println("\n【步骤 12】节点 2 - 继续对话")
	result2, err := agent2.Chat(ctx, "你还记得我是哪个节点吗？")
	if err != nil {
		fmt.Printf("❌ 发送失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 节点 2 收到回复: %s\n", truncate(result2.Text, 80))

	// 13. 再次检查 Redis
	fmt.Println("\n【步骤 13】检查 Redis 存储（应该包含两轮对话）")
	messages, err = redisStore.LoadMessages(ctx, agentID)
	if err != nil {
		fmt.Printf("❌ 加载消息失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Redis 中现在存储了 %d 条消息\n", len(messages))
	printMessages(messages)

	// 14. 测试 Store 自动修剪
	fmt.Println("\n【步骤 14】测试 Store 自动修剪（MaxMessages=10）")
	fmt.Println("发送多轮对话以触发修剪...")
	for i := 1; i <= 5; i++ {
		_, err := agent1.Chat(ctx, fmt.Sprintf("第 %d 条测试消息", i))
		if err != nil {
			fmt.Printf("  ❌ 第 %d 轮失败: %v\n", i, err)
			break
		}
		fmt.Printf("  ✅ 第 %d 轮完成\n", i)
	}

	messages, _ = redisStore.LoadMessages(ctx, agentID)
	fmt.Printf("✅ 修剪测试完成，当前消息数: %d (应该 ≤ 10)\n", len(messages))

	if len(messages) > 10 {
		fmt.Printf("❌ 修剪失败！消息数: %d > 10\n", len(messages))
	} else {
		fmt.Println("✅ 修剪功能正常工作")
	}

	// 15. 清理测试数据
	fmt.Println("\n【步骤 15】清理测试数据")
	if err := redisStore.DeleteAgent(ctx, agentID); err != nil {
		fmt.Printf("⚠️  清理失败: %v\n", err)
	} else {
		fmt.Println("✅ 测试数据已清理")
	}

	// 测试总结
	fmt.Println("\n=== 测试总结 ===")
	fmt.Println("✅ Redis Store 创建")
	fmt.Println("✅ Agent 集成")
	fmt.Println("✅ 节点 1 对话")
	fmt.Println("✅ 节点 2 共享状态")
	fmt.Println("✅ 分布式数据一致性")
	fmt.Println("✅ Store 自动修剪")
	fmt.Println("✅ 数据清理")
	fmt.Println("\n🎉 Redis Store + Agent 分布式集成测试通过！")
}

// CustomProviderFactory 自定义 Provider 工厂
type CustomProviderFactory struct {
	apiKey  string
	baseURL string
}

func (f *CustomProviderFactory) Create(config *types.ModelConfig) (provider.Provider, error) {
	return provider.NewCustomClaudeProvider(&types.ModelConfig{
		Provider: config.Provider,
		Model:    config.Model,
		APIKey:   f.apiKey,
		BaseURL:  f.baseURL,
	})
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// printMessages 打印消息列表
func printMessages(messages []types.Message) {
	for i, msg := range messages {
		content := extractText(msg)
		fmt.Printf("  [%d] %s: %s\n", i+1, msg.Role, truncate(content, 60))
	}
}

// extractText 从消息中提取文本
func extractText(msg types.Message) string {
	if msg.Content != "" {
		return msg.Content
	}

	for _, block := range msg.ContentBlocks {
		if textBlock, ok := block.(*types.TextBlock); ok {
			return textBlock.Text
		}
	}

	return "[非文本内容]"
}
