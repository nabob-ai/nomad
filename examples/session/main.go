// Session 演示会话管理系统，包括分层状态管理、事件溯源、事件过滤和会话列表。
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nabob-ai/nomad/pkg/session"
	"github.com/nabob-ai/nomad/pkg/types"
)

func main() {
	// 1. 创建 Session 服务
	sessionService := session.NewInMemoryService()
	ctx := context.Background()

	// 2. 创建新会话
	sess, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: "my-app",
		UserID:  "user-123",
		AgentID: "agent-001",
		Metadata: map[string]any{
			"source": "web",
			"region": "us-west",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	sessionID := sess.ID()
	fmt.Printf("✅ Created session: %s\n\n", sessionID)

	// 3. 演示分层状态管理
	demonstrateStateManagement(ctx, sessionService, sessionID)

	// 4. 演示事件管理
	demonstrateEventManagement(ctx, sessionService, sessionID)

	// 5. 演示事件过滤
	demonstrateEventFiltering(ctx, sessionService, sessionID)

	// 6. 演示会话列表
	demonstrateSessionListing(ctx, sessionService)
}

// 演示分层状态管理
func demonstrateStateManagement(ctx context.Context, service *session.InMemoryService, sessionID string) {
	fmt.Println("📊 State Management Demo")
	fmt.Println("========================")

	// 获取会话
	sess, _ := service.Get(ctx, &session.GetRequest{
		AppName:   "my-app",
		UserID:    "user-123",
		SessionID: sessionID,
	})

	state := sess.State()

	// 设置不同作用域的状态
	states := map[string]any{
		"app:version":           "1.0.0",                            // 应用级
		"user:preferences":      map[string]string{"theme": "dark"}, // 用户级
		"temp:current_task":     "processing",                       // 临时
		"session:message_count": 0,                                  // 会话级
	}

	for key, value := range states {
		if err := state.Set(key, value); err != nil {
			log.Printf("Error setting %s: %v", key, err)
		}
	}

	// 读取状态
	fmt.Println("\n📖 Reading states:")
	for key := range states {
		val, err := state.Get(key)
		if err != nil {
			log.Printf("Error getting %s: %v", key, err)
			continue
		}
		fmt.Printf("  %s = %v\n", key, val)
	}

	// 使用迭代器遍历所有状态
	fmt.Println("\n🔄 Iterating all states:")
	for key, value := range state.All() {
		scope := getScope(key)
		fmt.Printf("  [%s] %s = %v\n", scope, key, value)
	}

	fmt.Println()
}

// 演示事件管理
func demonstrateEventManagement(ctx context.Context, service *session.InMemoryService, sessionID string) {
	fmt.Println("📝 Event Management Demo")
	fmt.Println("========================")

	// 创建多个事件
	events := []*session.Event{
		{
			ID:           "evt-001",
			Timestamp:    time.Now(),
			InvocationID: "inv-001",
			AgentID:      "agent-001",
			Branch:       "main",
			Author:       "user",
			Content: types.Message{
				Role:    types.RoleUser,
				Content: "Hello, can you help me?",
			},
		},
		{
			ID:           "evt-002",
			Timestamp:    time.Now().Add(1 * time.Second),
			InvocationID: "inv-001",
			AgentID:      "agent-001",
			Branch:       "main",
			Author:       "agent-001",
			Content: types.Message{
				Role:    types.RoleAssistant,
				Content: "Of course! What do you need help with?",
			},
		},
		{
			ID:           "evt-003",
			Timestamp:    time.Now().Add(2 * time.Second),
			InvocationID: "inv-001",
			AgentID:      "agent-001",
			Branch:       "main",
			Author:       "agent-001",
			Content: types.Message{
				Role: types.RoleAssistant,
				ToolCalls: []types.ToolCall{
					{
						ID:   "call-001",
						Name: "search",
						Arguments: map[string]any{
							"query": "golang best practices",
						},
					},
				},
			},
			LongRunningToolIDs: []string{"call-001"},
		},
	}

	// 添加事件
	for _, evt := range events {
		if err := service.AppendEvent(ctx, sessionID, evt); err != nil {
			log.Printf("Error appending event: %v", err)
		}
	}

	// 获取会话并查看事件
	sess, _ := service.Get(ctx, &session.GetRequest{
		AppName:   "my-app",
		UserID:    "user-123",
		SessionID: sessionID,
	})

	eventList := sess.Events()
	fmt.Printf("\n📊 Total events: %d\n\n", eventList.Len())

	// 遍历事件
	fmt.Println("📜 Event timeline:")
	for evt := range eventList.All() {
		fmt.Printf("  [%s] %s (%s): ", evt.Timestamp.Format("15:04:05"), evt.Author, evt.ID)
		if evt.Content.Content != "" {
			fmt.Printf("%s\n", evt.Content.Content)
		} else if len(evt.Content.ToolCalls) > 0 {
			fmt.Printf("Tool call: %s\n", evt.Content.ToolCalls[0].Name)
		}

		if evt.IsFinalResponse() {
			fmt.Println("    ✓ Final response")
		}
	}

	fmt.Println()
}

// 演示事件过滤
func demonstrateEventFiltering(ctx context.Context, service *session.InMemoryService, sessionID string) {
	fmt.Println("🔍 Event Filtering Demo")
	fmt.Println("=======================")

	// 按作者过滤
	filter := &session.EventFilter{
		Author: "agent-001",
		Limit:  10,
	}

	events, err := service.GetEvents(ctx, sessionID, filter)
	if err != nil {
		log.Printf("Error filtering events: %v", err)
		return
	}

	fmt.Printf("\n📋 Agent events (filtered): %d\n", len(events))
	for _, evt := range events {
		fmt.Printf("  - %s: %s\n", evt.ID, evt.Content.Content)
	}

	// 按时间范围过滤
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()
	timeFilter := &session.EventFilter{
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	events, _ = service.GetEvents(ctx, sessionID, timeFilter)
	fmt.Printf("\n📅 Events in last hour: %d\n\n", len(events))
}

// 演示会话列表
func demonstrateSessionListing(ctx context.Context, service *session.InMemoryService) {
	fmt.Println("📚 Session Listing Demo")
	fmt.Println("=======================")

	// 创建多个会话
	for i := 1; i <= 3; i++ {
		_, _ = service.Create(ctx, &session.CreateRequest{
			AppName: "my-app",
			UserID:  "user-123",
			AgentID: fmt.Sprintf("agent-%03d", i),
		})
	}

	// 列出会话
	sessions, err := service.List(ctx, &session.ListRequest{
		AppName: "my-app",
		UserID:  "user-123",
		Limit:   10,
	})
	if err != nil {
		log.Printf("Error listing sessions: %v", err)
		return
	}

	fmt.Printf("\n📊 Total sessions for user-123: %d\n", len(sessions))
	for i, sess := range sessions {
		s := *sess
		fmt.Printf("  %d. %s (Agent: %s, Updated: %s)\n",
			i+1,
			s.ID(),
			s.AgentID(),
			s.LastUpdateTime().Format("15:04:05"),
		)
	}

	fmt.Println()
}

// 辅助函数：获取状态作用域
func getScope(key string) string {
	if session.IsAppKey(key) {
		return "APP"
	}
	if session.IsUserKey(key) {
		return "USER"
	}
	if session.IsTempKey(key) {
		return "TEMP"
	}
	if session.IsSessionKey(key) {
		return "SESSION"
	}
	return "UNKNOWN"
}
