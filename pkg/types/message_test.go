package types

import (
	"encoding/json"
	"testing"
)

func TestMessageMetadata_Defaults(t *testing.T) {
	m := NewMessageMetadata()
	if !m.UserVisible {
		t.Error("UserVisible should be true by default")
	}
	if !m.AgentVisible {
		t.Error("AgentVisible should be true by default")
	}
}

func TestMessageMetadata_AgentOnly(t *testing.T) {
	m := NewMessageMetadata().AgentOnly()
	if m.UserVisible {
		t.Error("UserVisible should be false for AgentOnly")
	}
	if !m.AgentVisible {
		t.Error("AgentVisible should be true for AgentOnly")
	}
}

func TestMessageMetadata_UserOnly(t *testing.T) {
	m := NewMessageMetadata().UserOnly()
	if !m.UserVisible {
		t.Error("UserVisible should be true for UserOnly")
	}
	if m.AgentVisible {
		t.Error("AgentVisible should be false for UserOnly")
	}
}

func TestMessageMetadata_Invisible(t *testing.T) {
	m := NewMessageMetadata().Invisible()
	if m.UserVisible {
		t.Error("UserVisible should be false for Invisible")
	}
	if m.AgentVisible {
		t.Error("AgentVisible should be false for Invisible")
	}
}

func TestMessageMetadata_WithSource(t *testing.T) {
	m := NewMessageMetadata().WithSource("summary")
	if m.Source != "summary" {
		t.Errorf("Source should be 'summary', got '%s'", m.Source)
	}
}

func TestMessageMetadata_WithTags(t *testing.T) {
	m := NewMessageMetadata().WithTags("important", "review")
	if len(m.Tags) != 2 {
		t.Errorf("Tags length should be 2, got %d", len(m.Tags))
	}
	if m.Tags[0] != "important" || m.Tags[1] != "review" {
		t.Error("Tags content mismatch")
	}
}

func TestMessageMetadata_IsVisible(t *testing.T) {
	tests := []struct {
		name     string
		meta     *MessageMetadata
		forAgent bool
		expected bool
	}{
		{"default for agent", NewMessageMetadata(), true, true},
		{"default for user", NewMessageMetadata(), false, true},
		{"agent only for agent", NewMessageMetadata().AgentOnly(), true, true},
		{"agent only for user", NewMessageMetadata().AgentOnly(), false, false},
		{"user only for agent", NewMessageMetadata().UserOnly(), true, false},
		{"user only for user", NewMessageMetadata().UserOnly(), false, true},
		{"invisible for agent", NewMessageMetadata().Invisible(), true, false},
		{"invisible for user", NewMessageMetadata().Invisible(), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.meta.IsVisible(tt.forAgent); got != tt.expected {
				t.Errorf("IsVisible(%v) = %v, want %v", tt.forAgent, got, tt.expected)
			}
		})
	}
}

func TestFilterMessagesForAgent(t *testing.T) {
	messages := []Message{
		{Role: RoleUser, Content: "Hello"},
		{Role: RoleAssistant, Content: "Hi", Metadata: NewMessageMetadata().AgentOnly()},
		{Role: RoleSystem, Content: "Summary", Metadata: NewMessageMetadata().UserOnly()},
		{Role: RoleUser, Content: "World", Metadata: NewMessageMetadata().Invisible()},
	}

	filtered := FilterMessagesForAgent(messages)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 messages for agent, got %d", len(filtered))
	}
	if filtered[0].Content != "Hello" {
		t.Error("First message should be 'Hello'")
	}
	if filtered[1].Content != "Hi" {
		t.Error("Second message should be 'Hi'")
	}
}

func TestFilterMessagesForUser(t *testing.T) {
	messages := []Message{
		{Role: RoleUser, Content: "Hello"},
		{Role: RoleAssistant, Content: "Hi", Metadata: NewMessageMetadata().AgentOnly()},
		{Role: RoleSystem, Content: "Summary", Metadata: NewMessageMetadata().UserOnly()},
		{Role: RoleUser, Content: "World", Metadata: NewMessageMetadata().Invisible()},
	}

	filtered := FilterMessagesForUser(messages)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 messages for user, got %d", len(filtered))
	}
	if filtered[0].Content != "Hello" {
		t.Error("First message should be 'Hello'")
	}
	if filtered[1].Content != "Summary" {
		t.Error("Second message should be 'Summary'")
	}
}

func TestFilterMessages(t *testing.T) {
	messages := []Message{
		{Role: RoleUser, Content: "A"},
		{Role: RoleUser, Content: "B", Metadata: &MessageMetadata{UserVisible: false, AgentVisible: false}},
		{Role: RoleUser, Content: "C"},
	}

	filtered := FilterMessages(messages, func(m *Message) bool {
		return m.Metadata == nil || m.Metadata.UserVisible
	})

	if len(filtered) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(filtered))
	}
}

func TestFilterMessagesBySource(t *testing.T) {
	messages := []Message{
		{Role: RoleUser, Content: "A"},
		{Role: RoleSystem, Content: "B", Metadata: NewMessageMetadata().WithSource("summary")},
		{Role: RoleSystem, Content: "C", Metadata: NewMessageMetadata().WithSource("user")},
	}

	filtered := FilterMessagesBySource(messages, "summary")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 message with source 'summary', got %d", len(filtered))
	}
	if filtered[0].Content != "B" {
		t.Error("Filtered message content mismatch")
	}
}

func TestFilterMessagesByTag(t *testing.T) {
	messages := []Message{
		{Role: RoleUser, Content: "A"},
		{Role: RoleUser, Content: "B", Metadata: NewMessageMetadata().WithTags("important")},
		{Role: RoleUser, Content: "C", Metadata: NewMessageMetadata().WithTags("review")},
	}

	filtered := FilterMessagesByTag(messages, "important")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 message with tag 'important', got %d", len(filtered))
	}
	if filtered[0].Content != "B" {
		t.Error("Filtered message content mismatch")
	}
}

func TestMessage_JSONSerialization(t *testing.T) {
	original := Message{
		Role:    RoleAssistant,
		Content: "Hello",
		Metadata: &MessageMetadata{
			UserVisible:  true,
			AgentVisible: false,
			Source:       "test",
			Tags:         []string{"tag1", "tag2"},
		},
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var restored Message
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify
	if restored.Role != original.Role {
		t.Errorf("Role mismatch: %v != %v", restored.Role, original.Role)
	}
	if restored.Content != original.Content {
		t.Errorf("Content mismatch: %v != %v", restored.Content, original.Content)
	}
	if restored.Metadata == nil {
		t.Fatal("Metadata should not be nil")
	}
	if restored.Metadata.UserVisible != original.Metadata.UserVisible {
		t.Errorf("UserVisible mismatch")
	}
	if restored.Metadata.AgentVisible != original.Metadata.AgentVisible {
		t.Errorf("AgentVisible mismatch")
	}
	if restored.Metadata.Source != original.Metadata.Source {
		t.Errorf("Source mismatch")
	}
	if len(restored.Metadata.Tags) != len(original.Metadata.Tags) {
		t.Errorf("Tags length mismatch")
	}
}

func TestMessage_JSONSerialization_NoMetadata(t *testing.T) {
	// 测试没有 Metadata 的消息序列化（向后兼容）
	original := Message{
		Role:    RoleUser,
		Content: "Test",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var restored Message
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if restored.Role != original.Role {
		t.Errorf("Role mismatch")
	}
	if restored.Content != original.Content {
		t.Errorf("Content mismatch")
	}
	// Metadata 可以为 nil，向后兼容
	if restored.Metadata != nil {
		t.Error("Metadata should be nil for backward compatibility")
	}
}

func TestMessage_NoMetadata_DefaultVisible(t *testing.T) {
	// 没有 Metadata 的消息应该默认对双方可见
	msg := Message{Role: RoleUser, Content: "Test"}

	// 测试过滤器行为：无 Metadata 的消息应该通过所有过滤器
	messages := []Message{msg}

	agentFiltered := FilterMessagesForAgent(messages)
	if len(agentFiltered) != 1 {
		t.Error("Message without metadata should be visible to agent")
	}

	userFiltered := FilterMessagesForUser(messages)
	if len(userFiltered) != 1 {
		t.Error("Message without metadata should be visible to user")
	}
}
