package provider

import (
	"testing"

	"github.com/nabob-ai/nomad/pkg/types"
)

func TestNewMultiProviderFactory(t *testing.T) {
	factory := NewMultiProviderFactory()
	if factory == nil {
		t.Fatal("expected MultiProviderFactory, got nil")
	}
}

func TestMultiProviderFactory_Create_Anthropic(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "anthropic",
		APIKey:   "test-key",
		Model:    "claude-3-opus",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected provider, got nil")
	}
}

func TestMultiProviderFactory_Create_OpenAI(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "openai",
		APIKey:   "test-key",
		Model:    "gpt-4",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected provider, got nil")
	}
}

func TestMultiProviderFactory_Create_Deepseek(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "deepseek",
		APIKey:   "test-key",
		Model:    "deepseek-chat",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected provider, got nil")
	}
}

func TestMultiProviderFactory_Create_GLM(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "glm",
		APIKey:   "test-key",
		Model:    "glm-4",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected provider, got nil")
	}
}

func TestMultiProviderFactory_Create_Groq(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "groq",
		APIKey:   "test-key",
		Model:    "mixtral-8x7b",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected provider, got nil")
	}
}

func TestMultiProviderFactory_Create_Ollama(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "ollama",
		Model:    "llama2",
		BaseURL:  "http://localhost:11434",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected provider, got nil")
	}
}

func TestMultiProviderFactory_Create_DefaultProvider(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		APIKey: "test-key",
		Model:  "claude-3-opus",
		// Provider 为空，应该默认为 anthropic
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected default provider (anthropic), got nil")
	}
}

func TestMultiProviderFactory_Create_CustomProvider(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "openai_compatible",
		APIKey:   "test-key",
		Model:    "custom-model",
		BaseURL:  "https://api.custom.com/v1",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected custom provider, got nil")
	}
}

func TestMultiProviderFactory_Create_UnsupportedProvider(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "unsupported_provider",
		APIKey:   "test-key",
		Model:    "some-model",
	}

	_, err := factory.Create(config)
	if err == nil {
		t.Error("expected error for unsupported provider, got nil")
	}
}

func TestMultiProviderFactory_Create_WithBaseURL(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "unknown",
		APIKey:   "test-key",
		Model:    "some-model",
		BaseURL:  "https://api.unknown.com/v1", // 提供了 BaseURL，应该作为自定义 provider
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected custom provider for unknown provider with BaseURL, got nil")
	}
}

func TestMultiProviderFactory_Create_Gemini(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "gemini",
		APIKey:   "test-key",
		Model:    "gemini-pro",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected gemini provider, got nil")
	}
}

func TestMultiProviderFactory_Create_Doubao(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "doubao",
		APIKey:   "test-key",
		Model:    "doubao-pro",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected doubao provider, got nil")
	}
}

func TestMultiProviderFactory_Create_Moonshot(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "moonshot",
		APIKey:   "test-key",
		Model:    "moonshot-v1",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected moonshot provider, got nil")
	}
}

func TestMultiProviderFactory_Create_Mistral(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "mistral",
		APIKey:   "test-key",
		Model:    "mistral-medium",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected mistral provider, got nil")
	}
}

func TestMultiProviderFactory_Create_OpenRouter(t *testing.T) {
	factory := NewMultiProviderFactory()

	config := &types.ModelConfig{
		Provider: "openrouter",
		APIKey:   "test-key",
		Model:    "openai/gpt-4",
	}

	provider, err := factory.Create(config)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if provider == nil {
		t.Error("expected openrouter provider, got nil")
	}
}
