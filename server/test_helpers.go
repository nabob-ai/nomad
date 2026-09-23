package server

import (
	"context"

	"github.com/nabob-ai/nomad/pkg/provider"
	"github.com/nabob-ai/nomad/pkg/types"
)

// MockProvider 模拟 Provider 用于测试
type MockProvider struct {
	systemPrompt string
}

func (m *MockProvider) Complete(ctx context.Context, messages []types.Message, opts *provider.StreamOptions) (*provider.CompleteResponse, error) {
	return &provider.CompleteResponse{
		Message: types.Message{
			Role:    types.RoleAssistant,
			Content: "Mock response",
		},
		Usage: &provider.TokenUsage{
			InputTokens:  10,
			OutputTokens: 20,
			TotalTokens:  30,
		},
	}, nil
}

func (m *MockProvider) Stream(ctx context.Context, messages []types.Message, opts *provider.StreamOptions) (<-chan provider.StreamChunk, error) {
	ch := make(chan provider.StreamChunk, 1)
	ch <- provider.StreamChunk{
		Type:         "content_block_delta",
		TextDelta:    "Mock response",
		FinishReason: "end_turn",
	}
	close(ch)
	return ch, nil
}

func (m *MockProvider) Config() *types.ModelConfig {
	return &types.ModelConfig{
		Provider: "mock",
		Model:    "test-model",
	}
}

func (m *MockProvider) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{
		SupportStreaming:    true,
		SupportToolCalling:  true,
		SupportSystemPrompt: true,
		SupportVision:       false,
	}
}

func (m *MockProvider) SetSystemPrompt(prompt string) error {
	m.systemPrompt = prompt
	return nil
}

func (m *MockProvider) GetSystemPrompt() string {
	return m.systemPrompt
}

func (m *MockProvider) Close() error {
	return nil
}

// MockProviderFactory 模拟 ProviderFactory
type MockProviderFactory struct {
	providers map[string]*MockProvider
}

func NewMockProviderFactory() *MockProviderFactory {
	return &MockProviderFactory{
		providers: make(map[string]*MockProvider),
	}
}

func (f *MockProviderFactory) Create(config *types.ModelConfig) (provider.Provider, error) {
	key := config.Provider + "/" + config.Model

	if p, exists := f.providers[key]; exists {
		return p, nil
	}

	// 创建新的 mock provider
	p := &MockProvider{}
	f.providers[key] = p
	return p, nil
}

func (f *MockProviderFactory) SetProvider(key string, p *MockProvider) {
	f.providers[key] = p
}
