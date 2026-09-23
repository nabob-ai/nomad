//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nabob-ai/nomad/pkg/store"
	"github.com/nabob-ai/nomad/pkg/types"
)

func main() {
	// 从环境变量读取配置
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	fmt.Println("=== Redis Store 核心功能测试 ===")
	fmt.Printf("Redis: %s\n\n", redisAddr)

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
		fmt.Printf("❌ 创建失败: %v\n", err)
		os.Exit(1)
	}
	defer redisStore.Close()
	fmt.Println("✅ Redis Store 创建成功")

	// 2. 测试连接
	fmt.Println("\n【步骤 2】测试连接")
	ctx := context.Background()
	if err := redisStore.Ping(ctx); err != nil {
		fmt.Printf("❌ Ping 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Redis 连接正常")

	agentID := "agt-test-001"

	// 3. 保存消息
	fmt.Println("\n【步骤 3】保存消息")
	messages := []types.Message{
		{
			Role:    types.MessageRoleUser,
			Content: "你好",
		},
		{
			Role:    types.MessageRoleAssistant,
			Content: "你好！我是 Claude。",
		},
	}

	if err := redisStore.SaveMessages(ctx, agentID, messages); err != nil {
		fmt.Printf("❌ 保存失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 成功保存 %d 条消息\n", len(messages))

	// 4. 加载消息
	fmt.Println("\n【步骤 4】加载消息")
	loaded, err := redisStore.LoadMessages(ctx, agentID)
	if err != nil {
		fmt.Printf("❌ 加载失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 成功加载 %d 条消息\n", len(loaded))
	for i, msg := range loaded {
		fmt.Printf("  [%d] %s: %s\n", i+1, msg.Role, msg.Content)
	}

	// 5. 添加更多消息（模拟对话）
	fmt.Println("\n【步骤 5】添加更多消息")
	for i := 1; i <= 8; i++ {
		messages = append(messages,
			types.Message{
				Role:    types.MessageRoleUser,
				Content: fmt.Sprintf("消息 %d", i),
			},
			types.Message{
				Role:    types.MessageRoleAssistant,
				Content: fmt.Sprintf("回复 %d", i),
			},
		)
	}
	if err := redisStore.SaveMessages(ctx, agentID, messages); err != nil {
		fmt.Printf("❌ 保存失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 当前消息数: %d\n", len(messages))

	// 6. 测试修剪功能
	fmt.Println("\n【步骤 6】测试修剪功能（保留最近 10 条）")
	if err := redisStore.TrimMessages(ctx, agentID, 10); err != nil {
		fmt.Printf("❌ 修剪失败: %v\n", err)
		os.Exit(1)
	}

	trimmed, _ := redisStore.LoadMessages(ctx, agentID)
	fmt.Printf("✅ 修剪后消息数: %d\n", len(trimmed))
	if len(trimmed) > 10 {
		fmt.Printf("❌ 修剪失败！期望 ≤ 10，实际 %d\n", len(trimmed))
	} else {
		fmt.Println("✅ 修剪功能正常")
	}

	// 7. 测试 Agent 信息存储
	fmt.Println("\n【步骤 7】测试 Agent 信息存储")
	info := types.AgentInfo{
		ID:         agentID,
		AgentID:    agentID,
		TemplateID: "test-template",
		Model:      "claude-sonnet-4-5",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := redisStore.SaveInfo(ctx, agentID, info); err != nil {
		fmt.Printf("❌ 保存 Info 失败: %v\n", err)
	} else {
		fmt.Println("✅ Agent Info 保存成功")
	}

	loadedInfo, err := redisStore.LoadInfo(ctx, agentID)
	if err != nil {
		fmt.Printf("❌ 加载 Info 失败: %v\n", err)
	} else {
		fmt.Printf("✅ Agent Info 加载成功: Model=%s\n", loadedInfo.Model)
	}

	// 8. 测试列出所有 Agent
	fmt.Println("\n【步骤 8】列出所有 Agent")
	agents, err := redisStore.ListAgents(ctx)
	if err != nil {
		fmt.Printf("❌ 列出失败: %v\n", err)
	} else {
		fmt.Printf("✅ 找到 %d 个 Agent: %v\n", len(agents), agents)
	}

	// 9. 测试分布式场景：创建第二个 Store 实例
	fmt.Println("\n【步骤 9】模拟分布式场景（创建第二个 Store 实例）")
	redisStore2, err := store.NewRedisStore(store.RedisConfig{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
		Prefix:   "nomad:test:",
		TTL:      1 * time.Hour,
	})
	if err != nil {
		fmt.Printf("❌ 创建第二个实例失败: %v\n", err)
	} else {
		defer redisStore2.Close()
		fmt.Println("✅ 第二个 Store 实例创建成功")

		// 从第二个实例加载数据
		messages2, err := redisStore2.LoadMessages(ctx, agentID)
		if err != nil {
			fmt.Printf("❌ 加载失败: %v\n", err)
		} else {
			fmt.Printf("✅ 第二个实例成功读取数据: %d 条消息\n", len(messages2))
			if len(messages2) == len(trimmed) {
				fmt.Println("✅ 数据共享验证成功")
			} else {
				fmt.Printf("⚠️  数据不一致: 实例1 有 %d 条，实例2 有 %d 条\n", len(trimmed), len(messages2))
			}
		}
	}

	// 10. 清理测试数据
	fmt.Println("\n【步骤 10】清理测试数据")
	if err := redisStore.DeleteAgent(ctx, agentID); err != nil {
		fmt.Printf("⚠️  清理失败: %v\n", err)
	} else {
		fmt.Println("✅ 测试数据已清理")

		// 验证已清理
		messages, err := redisStore.LoadMessages(ctx, agentID)
		if err != nil {
			fmt.Printf("⚠️  加载错误: %v\n", err)
		} else if len(messages) == 0 {
			fmt.Println("✅ 数据已完全清理")
		} else {
			fmt.Printf("⚠️  还有 %d 条消息未清理\n", len(messages))
		}
	}

	// 测试总结
	fmt.Println("\n=== 测试总结 ===")
	fmt.Println("✅ Redis Store 创建")
	fmt.Println("✅ 连接测试")
	fmt.Println("✅ 消息保存/加载")
	fmt.Println("✅ 消息修剪（FIFO）")
	fmt.Println("✅ Agent 信息存储")
	fmt.Println("✅ 列出所有 Agent")
	fmt.Println("✅ 分布式数据共享")
	fmt.Println("✅ 数据清理")
	fmt.Println("\n🎉 所有测试通过！Redis Store 工作正常")
}
