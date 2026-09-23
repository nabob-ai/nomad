// MemorySemanticSession 演示如何将 Session 对话内容保存到 SemanticMemory，
// 实现跨会话的长期语义记忆和检索能力。
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nabob-ai/nomad/pkg/memory"
	"github.com/nabob-ai/nomad/pkg/session"
	"github.com/nabob-ai/nomad/pkg/types"
	"github.com/nabob-ai/nomad/pkg/vector"
)

func main() {
	ctx := context.Background()

	// 1. 准备 Session 服务（短期记忆）
	svc := session.NewInMemoryService()
	appName := "ltm_demo"
	userID := "alice"

	teachingSessionID := "teaching-session"

	// 创建教学会话
	teachingSess, err := svc.Create(ctx, &session.CreateRequest{
		AppName:  appName,
		UserID:   userID,
		AgentID:  "teacher-agent",
		Metadata: map[string]any{"kind": "teaching"},
	})
	if err != nil {
		log.Fatalf("create teaching session: %v", err)
	}
	if teachingSess.ID() != teachingSessionID {
		// 使用固定 sessionID 方便演示
		if err := svc.Update(ctx, &session.UpdateRequest{
			SessionID: teachingSess.ID(),
			Metadata:  map[string]any{"_alias": teachingSessionID},
		}); err != nil {
			log.Printf("update teaching session metadata: %v", err)
		}
	}

	// 写入教学事件
	teachEvents := []*session.Event{
		{
			ID:           "evt-teach-1",
			Timestamp:    time.Now(),
			InvocationID: "inv-001",
			AgentID:      "teacher-agent",
			Author:       "user",
			Content: types.Message{
				Role:    types.RoleUser,
				Content: "My secret is 0xabcd",
			},
		},
		{
			ID:           "evt-teach-2",
			Timestamp:    time.Now().Add(1 * time.Second),
			InvocationID: "inv-001",
			AgentID:      "teacher-agent",
			Author:       "assistant",
			Content: types.Message{
				Role:    types.RoleAssistant,
				Content: "Got it, I will remember your secret.",
			},
		},
	}
	for _, ev := range teachEvents {
		if err := svc.AppendEvent(ctx, teachingSess.ID(), ev); err != nil {
			log.Fatalf("append teaching event: %v", err)
		}
	}

	fmt.Printf("✅ Created teaching session: %s\n", teachingSess.ID())

	// 2. 准备语义记忆（长期记忆）
	vecStore := vector.NewMemoryStore()
	embedder := vector.NewMockEmbedder(16)
	semMem := memory.NewSemanticMemory(memory.SemanticMemoryConfig{
		Store:          vecStore,
		Embedder:       embedder,
		NamespaceScope: "user",
		TopK:           3,
	})

	bridge := &memory.LongTermBridge{
		Sessions:       svc,
		SemanticMemory: semMem,
	}

	// 3. 将教学会话保存到长期语义记忆中
	fmt.Println("\n💾 Saving teaching session to long-term semantic memory...")
	if err := bridge.SaveSessionToSemanticMemory(
		ctx,
		appName,
		userID,
		teachingSess.ID(),
		map[string]any{"user_id": userID},
		&memory.LongTermBridgeConfig{MinTokens: 3},
	); err != nil {
		log.Fatalf("save session to semantic memory: %v", err)
	}
	fmt.Println("✅ Saved.")

	// 4. 在新的会话中进行语义查询
	fmt.Println("\n🔍 Querying long-term memory from a new session...")
	studentSess, err := svc.Create(ctx, &session.CreateRequest{
		AppName:  appName,
		UserID:   userID,
		AgentID:  "student-agent",
		Metadata: map[string]any{"kind": "student"},
	})
	if err != nil {
		log.Fatalf("create student session: %v", err)
	}

	question := "What is my secret?"
	fmt.Printf("Question: %s\n", question)

	hits, err := semMem.Search(ctx, question, map[string]any{"user_id": userID}, 3)
	if err != nil {
		log.Fatalf("semantic search failed: %v", err)
	}

	fmt.Println("\nSemantic search hits:")
	for _, h := range hits {
		fmt.Printf("  ID=%s, score=%.4f\n", h.ID, h.Score)
		if txt, ok := h.Metadata["text"].(string); ok {
			fmt.Printf("    text: %s\n", txt)
		}
	}

	fmt.Printf("\n✅ Student session: %s (can now use long-term memory)\n", studentSess.ID())
}
