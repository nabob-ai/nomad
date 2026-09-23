package session

import (
	"testing"

	"github.com/nabob-ai/nomad/pkg/types"
)

func TestNewEvent(t *testing.T) {
	invocationID := "test-invocation"
	event := NewEvent(invocationID)

	if event == nil {
		t.Fatal("expected Event, got nil")
	}

	if event.InvocationID != invocationID {
		t.Errorf("expected InvocationID '%s', got '%s'", invocationID, event.InvocationID)
	}

	if event.ID == "" {
		t.Error("expected event ID to be generated")
	}

	if event.Actions.StateDelta == nil {
		t.Error("expected StateDelta to be initialized")
	}

	if event.Actions.ArtifactDelta == nil {
		t.Error("expected ArtifactDelta to be initialized")
	}

	if event.Metadata == nil {
		t.Error("expected Metadata to be initialized")
	}
}

func TestEvent_IsFinalResponse_WithSkipSummarization(t *testing.T) {
	event := &Event{
		Actions: EventActions{
			SkipSummarization: true,
		},
		Content: types.Message{
			Role: types.RoleAssistant,
		},
	}

	if !event.IsFinalResponse() {
		t.Error("expected IsFinalResponse to be true when SkipSummarization is true")
	}
}

func TestEvent_IsFinalResponse_WithLongRunningTools(t *testing.T) {
	event := &Event{
		LongRunningToolIDs: []string{"tool1", "tool2"},
		Content: types.Message{
			Role: types.RoleAssistant,
		},
	}

	if !event.IsFinalResponse() {
		t.Error("expected IsFinalResponse to be true when LongRunningToolIDs is not empty")
	}
}

func TestEvent_IsFinalResponse_WithToolCalls(t *testing.T) {
	event := &Event{
		Content: types.Message{
			Role: types.RoleAssistant,
			ToolCalls: []types.ToolCall{
				{ID: "call1", Name: "tool1"},
			},
		},
	}

	if event.IsFinalResponse() {
		t.Error("expected IsFinalResponse to be false when ToolCalls is present")
	}
}

func TestEvent_IsFinalResponse_NoToolCalls(t *testing.T) {
	event := &Event{
		Content: types.Message{
			Role:    types.RoleAssistant,
			Content: "Final response",
		},
	}

	if !event.IsFinalResponse() {
		t.Error("expected IsFinalResponse to be true when no ToolCalls")
	}
}

func TestIsAppKey(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"app:setting1", true},
		{"app:", false}, // 长度不够
		{"user:setting1", false},
		{"temp:setting1", false},
		{"session:setting1", false},
		{"setting1", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := IsAppKey(tt.key)
			if result != tt.expected {
				t.Errorf("IsAppKey(%q) = %v, expected %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestIsUserKey(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"user:preference1", true},
		{"user:", false}, // 长度不够
		{"app:setting1", false},
		{"temp:setting1", false},
		{"session:setting1", false},
		{"preference1", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := IsUserKey(tt.key)
			if result != tt.expected {
				t.Errorf("IsUserKey(%q) = %v, expected %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestIsTempKey(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"temp:data1", true},
		{"temp:", false}, // 长度不够
		{"app:setting1", false},
		{"user:setting1", false},
		{"session:setting1", false},
		{"data1", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := IsTempKey(tt.key)
			if result != tt.expected {
				t.Errorf("IsTempKey(%q) = %v, expected %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestIsSessionKey(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"session:state1", true},
		{"session:", false}, // 长度不够
		{"app:setting1", false},
		{"user:setting1", false},
		{"temp:setting1", false},
		{"state1", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := IsSessionKey(tt.key)
			if result != tt.expected {
				t.Errorf("IsSessionKey(%q) = %v, expected %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestGenerateEventID(t *testing.T) {
	id1 := generateEventID()
	id2 := generateEventID()

	if id1 == "" {
		t.Error("expected non-empty event ID")
	}

	if id2 == "" {
		t.Error("expected non-empty event ID")
	}

	// IDs 应该不同（虽然在同一秒内可能相同，但这个测试主要验证格式）
	if len(id1) < 4 || id1[:4] != "evt_" {
		t.Errorf("expected event ID to start with 'evt_', got %s", id1)
	}
}

func TestEventActions(t *testing.T) {
	actions := EventActions{
		StateDelta:        map[string]any{"key1": "value1"},
		ArtifactDelta:     map[string]int64{"file1": 1},
		SkipSummarization: true,
		TransferToAgent:   "agent2",
		Escalate:          true,
		CustomActions:     map[string]any{"action1": "data1"},
	}

	if actions.StateDelta["key1"] != "value1" {
		t.Error("expected StateDelta to contain key1")
	}

	if actions.ArtifactDelta["file1"] != 1 {
		t.Error("expected ArtifactDelta to contain file1")
	}

	if !actions.SkipSummarization {
		t.Error("expected SkipSummarization to be true")
	}

	if actions.TransferToAgent != "agent2" {
		t.Error("expected TransferToAgent to be 'agent2'")
	}

	if !actions.Escalate {
		t.Error("expected Escalate to be true")
	}

	if actions.CustomActions["action1"] != "data1" {
		t.Error("expected CustomActions to contain action1")
	}
}

func TestCreateRequest(t *testing.T) {
	req := &CreateRequest{
		AppName:  "test-app",
		UserID:   "user1",
		AgentID:  "agent1",
		Metadata: map[string]any{"key": "value"},
	}

	if req.AppName != "test-app" {
		t.Errorf("expected AppName 'test-app', got '%s'", req.AppName)
	}

	if req.UserID != "user1" {
		t.Errorf("expected UserID 'user1', got '%s'", req.UserID)
	}

	if req.AgentID != "agent1" {
		t.Errorf("expected AgentID 'agent1', got '%s'", req.AgentID)
	}

	if req.Metadata["key"] != "value" {
		t.Error("expected Metadata to contain key")
	}
}

func TestGetRequest(t *testing.T) {
	req := &GetRequest{
		AppName:   "test-app",
		UserID:    "user1",
		SessionID: "session1",
	}

	if req.AppName != "test-app" {
		t.Errorf("expected AppName 'test-app', got '%s'", req.AppName)
	}

	if req.UserID != "user1" {
		t.Errorf("expected UserID 'user1', got '%s'", req.UserID)
	}

	if req.SessionID != "session1" {
		t.Errorf("expected SessionID 'session1', got '%s'", req.SessionID)
	}
}

func TestSessionData(t *testing.T) {
	data := &SessionData{
		ID:       "session1",
		AppName:  "test-app",
		UserID:   "user1",
		AgentID:  "agent1",
		Metadata: map[string]any{"key": "value"},
	}

	if data.ID != "session1" {
		t.Errorf("expected ID 'session1', got '%s'", data.ID)
	}

	if data.AppName != "test-app" {
		t.Errorf("expected AppName 'test-app', got '%s'", data.AppName)
	}

	if data.UserID != "user1" {
		t.Errorf("expected UserID 'user1', got '%s'", data.UserID)
	}

	if data.AgentID != "agent1" {
		t.Errorf("expected AgentID 'agent1', got '%s'", data.AgentID)
	}

	if data.Metadata["key"] != "value" {
		t.Error("expected Metadata to contain key")
	}
}
