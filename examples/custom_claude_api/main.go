package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nabob-ai/nomad/pkg/agent"
	"github.com/nabob-ai/nomad/pkg/provider"
	"github.com/nabob-ai/nomad/pkg/sandbox"
	"github.com/nabob-ai/nomad/pkg/store"
	"github.com/nabob-ai/nomad/pkg/tools"
	"github.com/nabob-ai/nomad/pkg/tools/builtin"
	"github.com/nabob-ai/nomad/pkg/types"
)

// Config 应用配置
type Config struct {
	APIKey         string
	BaseURL        string
	Model          string
	StoreDir       string
	StoreMaxMsg    int
	StoreAutoTrim  bool
	SandboxWorkDir string
}

func main() {
	// 加载配置
	config := loadConfig()
	validateConfig(config)

	fmt.Println("===Nabob AI Agent ===")
	fmt.Printf("Model: %s\n", config.Model)
	fmt.Printf("Base URL: %s\n", config.BaseURL)
	fmt.Printf("Store: %s (max: %d messages)\n\n", config.StoreDir, config.StoreMaxMsg)

	// 初始化组件
	ctx := context.Background()
	deps, err := initializeDependencies(config)
	if err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 Agent
	ag, err := createAgent(ctx, config, deps)
	if err != nil {
		fmt.Printf("❌ 创建 Agent 失败: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := ag.Close(); err != nil {
			fmt.Printf("⚠️  关闭 Agent 失败: %v\n", err)
		}
	}()

	// 设置事件监听
	setupEventHandlers(ag)

	// 运行示例
	runExamples(ag)

	fmt.Println("\n✅ 程序执行完成")
}

// loadConfig 从环境变量加载配置
func loadConfig() *Config {
	return &Config{
		APIKey:  getEnv("CLAUDE_API_KEY", ""),
		BaseURL: getEnv("CLAUDE_BASE_URL", "https://api.anthropic.com"),
		Model:   getEnv("CLAUDE_MODEL", "claude-sonnet-4-5-20250929"),

		StoreDir:      getEnv("STORE_DIR", ".nomad"),
		StoreMaxMsg:   getEnvInt("STORE_MAX_MESSAGES", 20),
		StoreAutoTrim: getEnvBool("STORE_AUTO_TRIM", true),

		SandboxWorkDir: getEnv("SANDBOX_WORK_DIR", "./workspace"),
	}
}

// validateConfig 验证配置
func validateConfig(config *Config) {
	if config.APIKey == "" {
		fmt.Println("❌ 错误: 未配置 CLAUDE_API_KEY")
		fmt.Println("请设置环境变量或创建 .env 文件")
		fmt.Println("示例: export CLAUDE_API_KEY=your-key-here")
		os.Exit(1)
	}

	if config.Model == "" {
		fmt.Println("❌ 错误: 未配置 CLAUDE_MODEL")
		os.Exit(1)
	}
}

// initializeDependencies 初始化依赖组件
func initializeDependencies(config *Config) (*agent.Dependencies, error) {
	// 工具注册表
	toolRegistry := tools.NewRegistry()
	builtin.RegisterAll(toolRegistry)

	// 沙箱工厂
	sandboxFactory := sandbox.NewFactory()

	// Provider 工厂
	providerFactory := &provider.CustomClaudeFactory{}

	// Store
	jsonStore, err := store.NewJSONStore(config.StoreDir)
	if err != nil {
		return nil, fmt.Errorf("create store: %w", err)
	}

	// 模板注册表
	templateRegistry := agent.NewTemplateRegistry()
	templateRegistry.Register(&types.AgentTemplateDefinition{
		ID:           "assistant",
		SystemPrompt: "You are Claude, an AI assistant created by Anthropic. You are helpful, harmless, and honest.",
		Tools:        "*",
	})

	return &agent.Dependencies{
		Store:            jsonStore,
		SandboxFactory:   sandboxFactory,
		ToolRegistry:     toolRegistry,
		ProviderFactory:  providerFactory,
		TemplateRegistry: templateRegistry,
	}, nil
}

// createAgent 创建 Agent
func createAgent(ctx context.Context, config *Config, deps *agent.Dependencies) (*agent.Agent, error) {
	agentConfig := &types.AgentConfig{
		TemplateID: "assistant",

		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    config.Model,
			APIKey:   config.APIKey,
			BaseURL:  config.BaseURL,
		},

		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindLocal,
			WorkDir: config.SandboxWorkDir,
		},

		Store: &types.StoreConfig{
			MaxMessages: config.StoreMaxMsg,
			AutoTrim:    config.StoreAutoTrim,
		},

		Context: &types.ContextManagerOptions{
			MaxTokens: 200000,
		},
	}

	return agent.Create(ctx, agentConfig, deps)
}

// setupEventHandlers 设置事件处理器
func setupEventHandlers(ag *agent.Agent) {
	// EventBus 是私有字段，Agent API 没有提供公开的事件订阅方法
	// 事件处理功能已移除，因为无法访问私有的 eventBus 字段
	// 如果需要事件处理，应该在 Agent 包中提供公开的事件订阅 API
	_ = ag
}

// runExamples 运行示例
func runExamples(ag *agent.Agent) {
	fmt.Println("【示例 1】基础对话")
	fmt.Println("----------------------------------------")
	chat(ag, "你好，请简单介绍一下你自己")

	fmt.Println("\n【示例 2】多轮对话")
	fmt.Println("----------------------------------------")
	chat(ag, "我最喜欢的颜色是蓝色")
	chat(ag, "我刚才说什么了？")
}

// chat 辅助函数
func chat(ag *agent.Agent, message string) {
	fmt.Printf("\n💬 用户: %s\n", message)
	fmt.Print("🤖 助手: ")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := ag.Chat(ctx, message)
	if err != nil {
		fmt.Printf("\n❌ 错误: %v\n", err)
		return
	}

	fmt.Printf("%s\n", result.Text)
}

// 辅助函数

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
