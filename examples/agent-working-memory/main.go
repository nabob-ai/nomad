// AgentWorkingMemory 演示带 Working Memory 的 Agent，支持自动加载到
// system prompt、LLM 主动更新和跨会话持久化。
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nabob-ai/nomad/pkg/agent"
	"github.com/nabob-ai/nomad/pkg/provider"
	"github.com/nabob-ai/nomad/pkg/sandbox"
	"github.com/nabob-ai/nomad/pkg/store"
	"github.com/nabob-ai/nomad/pkg/tools"
	"github.com/nabob-ai/nomad/pkg/tools/builtin"
	"github.com/nabob-ai/nomad/pkg/types"
)

//
// 使用场景:
// - 记住用户偏好和设置
// - 跟踪任务进度
// - 维护会话状态

func main() {
	// 检查 API Key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is required")
	}

	ctx := context.Background()

	// 1. 工具注册表
	toolRegistry := tools.NewRegistry()
	builtin.RegisterAll(toolRegistry)

	// 2. Sandbox 工厂
	sandboxFactory := sandbox.NewFactory()

	// 3. Provider 工厂
	providerFactory := &provider.AnthropicFactory{}

	// 4. Store
	jsonStore, err := store.NewJSONStore(".nomad-working-memory")
	if err != nil {
		log.Fatalf("create store failed: %v", err)
	}

	// 5. 模板注册表
	templateRegistry := agent.NewTemplateRegistry()
	templateRegistry.Register(&types.AgentTemplateDefinition{
		ID: "task-assistant",
		SystemPrompt: `You are a helpful task assistant with Working Memory.

Your capabilities:
- Remember user information and preferences
- Track task progress across multiple conversations
- Maintain context between sessions

When you learn something important about the user or their tasks, use the update_working_memory tool to store it.

Remember: Working Memory is automatically loaded at the start of each conversation, so you always have access to it.`,
		Model: "claude-sonnet-4-5",
		// 基础工具，working memory 工具由中间件自动注入
		Tools: []any{"Read", "Write", "Bash"},
	})

	deps := &agent.Dependencies{
		Store:            jsonStore,
		SandboxFactory:   sandboxFactory,
		ToolRegistry:     toolRegistry,
		ProviderFactory:  providerFactory,
		TemplateRegistry: templateRegistry,
	}

	// 6. 工作目录
	workDir := "./workspace"
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		log.Fatalf("create workspace dir failed: %v", err)
	}

	// 7. 创建 Agent 配置，启用 working_memory 中间件
	config := &types.AgentConfig{
		TemplateID: "task-assistant",
		ModelConfig: &types.ModelConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-5",
			APIKey:   apiKey,
		},
		Sandbox: &types.SandboxConfig{
			Kind:    types.SandboxKindLocal,
			WorkDir: workDir,
		},
		// 启用中间件: filesystem + agent_memory + working_memory
		Middlewares: []string{"filesystem", "agent_memory", "working_memory"},
		// 设置 threadID 用于 Working Memory 隔离
		Metadata: map[string]any{
			"thread_id":   "demo-thread-001",
			"resource_id": "demo-task",
		},
	}

	// 8. 创建 Agent
	ag, err := agent.Create(ctx, config, deps)
	if err != nil {
		log.Fatalf("create agent failed: %v", err)
	}
	defer func() { _ = ag.Close() }()

	fmt.Printf("✅ Working Memory Agent created: %s\n", ag.ID())
	fmt.Printf("📝 Thread ID: %s\n", config.Metadata["thread_id"])
	fmt.Println("\n=== Working Memory Agent Demo ===")
	fmt.Println("演示场景：任务助手记住用户信息并跟踪任务进度")
	fmt.Println()

	// 9. 订阅事件（简化展示）
	eventCh := ag.Subscribe([]types.AgentChannel{types.ChannelProgress}, nil)
	go func() {
		for envelope := range eventCh {
			if evt, ok := envelope.Event.(types.EventType); ok {
				switch e := evt.(type) {
				case *types.ProgressTextChunkStartEvent:
					fmt.Print("\n[Assistant] ")
				case *types.ProgressTextChunkEvent:
					fmt.Print(e.Delta)
				case *types.ProgressDoneEvent:
					if e.Reason == "tool_use" {
						fmt.Printf("\n💡 使用了工具\n")
					}
				}
			}
		}
	}()

	// 10. 运行交互式演示
	runDemo(ctx, ag)
}

func runDemo(ctx context.Context, ag *agent.Agent) {
	scenarios := []struct {
		title  string
		prompt string
		wait   time.Duration
	}{
		{
			title: "会话 1: 初次见面，Agent 记住用户信息",
			prompt: `你好！我是 Alice，一名软件工程师。

我的偏好：
- 编程语言：TypeScript 和 Go
- 代码风格：简洁、函数式
- 回答风格：直接、技术性强

请记住这些信息。`,
			wait: 3 * time.Second,
		},
		{
			title: "会话 2: 开始一个任务",
			prompt: `我需要创建一个 REST API 项目。请帮我规划一下步骤。

要求：
- 使用我偏好的语言
- 包含基础的 CRUD 操作
- 有完整的错误处理`,
			wait: 3 * time.Second,
		},
		{
			title:  "会话 3: 查询任务进度",
			prompt: `我的 REST API 项目进展如何了？还记得我的要求吗？`,
			wait:   3 * time.Second,
		},
	}

	// 自动演示模式
	if len(os.Args) > 1 && os.Args[1] == "--auto" {
		fmt.Println("🤖 自动演示模式")
		fmt.Println("=" + strings.Repeat("=", 60))
		for i, scenario := range scenarios {
			fmt.Printf("\n\n%s\n", scenario.title)
			fmt.Println(strings.Repeat("-", 60))
			fmt.Printf("[User] %s\n", scenario.prompt)

			_, err := ag.Chat(ctx, scenario.prompt)
			if err != nil {
				log.Printf("Error in scenario %d: %v", i+1, err)
			}

			time.Sleep(scenario.wait)
		}

		fmt.Println("\n\n" + strings.Repeat("=", 60))
		fmt.Println("✅ 演示完成！")
		printAgentStatus(ag)
		return
	}

	// 交互模式
	fmt.Println("🎮 交互模式")
	fmt.Println("输入 'auto' 运行自动演示")
	fmt.Println("输入 'status' 查看 Agent 状态")
	fmt.Println("输入 'quit' 退出")
	fmt.Println(strings.Repeat("=", 60))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n[You] ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch input {
		case "quit", "exit", "q":
			fmt.Println("👋 再见！")
			return

		case "status":
			printAgentStatus(ag)
			continue

		case "auto":
			runDemo(ctx, ag)
			return

		case "help":
			printHelp()
			continue
		}

		// 发送消息
		_, err := ag.Chat(ctx, input)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func printAgentStatus(ag *agent.Agent) {
	status := ag.Status()
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 Agent Status")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("State:      %s\n", status.State)
	fmt.Printf("Steps:      %d\n", status.StepCount)
	fmt.Printf("Cursor:     %d\n", status.Cursor)
	fmt.Println(strings.Repeat("=", 60))
}

func printHelp() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📖 可用命令")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("auto   - 运行自动演示")
	fmt.Println("status - 查看 Agent 状态")
	fmt.Println("help   - 显示此帮助")
	fmt.Println("quit   - 退出程序")
	fmt.Println(strings.Repeat("=", 60))
}
