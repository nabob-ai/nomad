package main

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/nabob-ai/nomad/pkg/agent"
	"github.com/nabob-ai/nomad/pkg/types"
)

// =============================================================================
// AgentIntegrationSuite - OpenRouter Agent 集成测试套件
// =============================================================================

type AgentIntegrationSuite struct {
	suite.Suite

	// 每个测试独立的 Agent（避免上下文污染）
	ag        *agent.Agent
	ctx       context.Context
	cancel    context.CancelFunc
	workspace string
	apiKey    string
	eventCh   <-chan types.AgentEventEnvelope
}

// SetupSuite 在所有测试开始前执行一次
func (s *AgentIntegrationSuite) SetupSuite() {
	s.T().Log("🚀 初始化 Agent 集成测试套件...")

	// 检查环境变量
	s.apiKey = os.Getenv("OPENROUTER_API_KEY")
	if s.apiKey == "" {
		s.T().Skip("跳过测试：需要设置 OPENROUTER_API_KEY 环境变量")
	}

	s.workspace = "./workspace"

	// 确保工作目录存在
	err := os.MkdirAll(s.workspace, 0755)
	s.Require().NoError(err, "创建工作目录失败")
}

// TearDownSuite 在所有测试结束后执行一次
func (s *AgentIntegrationSuite) TearDownSuite() {
	s.T().Log("🧹 清理测试套件...")

	// 清理测试文件
	_ = os.RemoveAll(s.workspace)
}

// SetupTest 在每个测试方法前执行 - 创建新的 Agent
func (s *AgentIntegrationSuite) SetupTest() {
	// 清理可能存在的测试文件
	_ = os.Remove(s.workspace + "/test.txt")

	// 为每个测试创建独立的 Agent（避免对话上下文污染）
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 2*time.Minute)

	var err error
	s.ag, err = createTestAgent(s.apiKey)
	s.Require().NoError(err, "创建 Agent 失败")

	// 订阅事件（用于调试）
	s.eventCh = s.ag.Subscribe(
		[]types.AgentChannel{types.ChannelProgress, types.ChannelMonitor},
		nil,
	)

	// 启动事件监听
	go s.handleEvents()
}

// TearDownTest 在每个测试方法后执行 - 关闭 Agent
func (s *AgentIntegrationSuite) TearDownTest() {
	if s.ag != nil {
		_ = s.ag.Close()
		s.ag = nil
	}
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

// handleEvents 处理 Agent 事件（调试用）
func (s *AgentIntegrationSuite) handleEvents() {
	for envelope := range s.eventCh {
		switch e := envelope.Event.(type) {
		case *types.ProgressToolStartEvent:
			s.T().Logf("  [工具开始] %s", e.Call.Name)
		case *types.ProgressToolEndEvent:
			s.T().Logf("  [工具完成] %s", e.Call.Name)
		case *types.ProgressToolErrorEvent:
			s.T().Logf("  [工具错误] %s: %s", e.Call.Name, e.Error)
		}
	}
}

// =============================================================================
// 测试用例
// =============================================================================

func (s *AgentIntegrationSuite) TestCreateFile() {
	result, err := s.ag.Chat(s.ctx, "使用 Write 工具在当前目录创建文件 test.txt，文件内容为: Hello World")

	s.Require().NoError(err, "Chat 调用失败")
	s.Require().NotNil(result, "结果不应为空")
	s.Equal("ok", result.Status, "状态应为 ok")

	// 等待文件操作完成
	time.Sleep(300 * time.Millisecond)

	// 验证文件创建
	data, err := os.ReadFile(s.workspace + "/test.txt")
	s.Require().NoError(err, "文件应该已创建")
	s.Equal("Hello World", strings.TrimSpace(string(data)), "文件内容不匹配")
}

func (s *AgentIntegrationSuite) TestReadFile() {
	// 先创建测试文件
	testContent := "这是测试内容"
	err := os.WriteFile(s.workspace+"/test.txt", []byte(testContent), 0644)
	s.Require().NoError(err, "创建测试文件失败")

	result, err := s.ag.Chat(s.ctx, "使用 Read 工具读取 test.txt 文件的内容")

	s.Require().NoError(err, "Chat 调用失败")
	s.Require().NotNil(result, "结果不应为空")
	s.Equal("ok", result.Status, "状态应为 ok")
}

func (s *AgentIntegrationSuite) TestBashCommand() {
	result, err := s.ag.Chat(s.ctx, "使用 Bash 工具执行命令: ls -la")

	s.Require().NoError(err, "Chat 调用失败")
	s.Require().NotNil(result, "结果不应为空")
	s.Equal("ok", result.Status, "状态应为 ok")
}

func (s *AgentIntegrationSuite) TestAgentStatus() {
	// 先执行一个简单操作确保有步骤记录
	_, err := s.ag.Chat(s.ctx, "你好")
	s.Require().NoError(err)

	status := s.ag.Status()

	s.Equal(types.AgentStateReady, status.State, "Agent 状态应为 Ready")
	assert.Positive(s.T(), status.StepCount, "步骤计数应大于 0")
	s.NotEmpty(status.AgentID, "Agent ID 不应为空")

	s.T().Logf("Agent 状态: ID=%s, State=%s, Steps=%d",
		status.AgentID, status.State, status.StepCount)
}

// TestMultipleBashCommands 使用 Table-Driven 测试多个命令
func (s *AgentIntegrationSuite) TestMultipleBashCommands() {
	tests := []struct {
		name   string
		prompt string
	}{
		{"列出当前目录", "使用 Bash 工具执行命令: pwd"},
		{"显示日期", "使用 Bash 工具执行命令: date"},
		{"回显文本", "使用 Bash 工具执行命令: echo hello world"},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			result, err := s.ag.Chat(s.ctx, tc.prompt)

			s.Require().NoError(err, "Chat 调用失败")
			s.Require().NotNil(result, "结果不应为空")
			s.Equal("ok", result.Status, "状态应为 ok")
		})
	}
}

// =============================================================================
// 测试入口
// =============================================================================

func TestAgentIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（-short 模式）")
	}

	suite.Run(t, new(AgentIntegrationSuite))
}
