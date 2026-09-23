# Providers 总览

nomad 支持 10+ 个主流 AI Provider，覆盖国际和中国市场。

## 支持的 Providers

### 🌍 国际主流

| Provider       | 特点                    | 配置名称           | 文档                        |
| -------------- | ----------------------- | ------------------ | --------------------------- |
| **OpenAI**     | 最流行，GPT-4/o1/o3     | `openai`           | [详细文档](./openai.md)     |
| **Anthropic**  | Claude 系列             | `anthropic`        | [详细文档](./anthropic.md)  |
| **Gemini**     | 超长上下文 1M，视频理解 | `gemini`, `google` | [详细文档](./gemini.md)     |
| **Groq**       | 超快推理速度            | `groq`             | [详细文档](./groq.md)       |
| **OpenRouter** | 聚合平台，数百模型      | `openrouter`       | [详细文档](./openrouter.md) |
| **Mistral**    | 欧洲主流，开源友好      | `mistral`          | [详细文档](./mistral.md)    |
| **Ollama**     | 本地部署首选            | `ollama`           | [详细文档](./ollama.md)     |

### 🇨🇳 中国市场

| Provider          | 特点           | 配置名称              | 文档                      |
| ----------------- | -------------- | --------------------- | ------------------------- |
| **DeepSeek**      | R1 推理模型    | `deepseek`            | [详细文档](./deepseek.md) |
| **智谱 GLM**      | ChatGLM 系列   | `glm`, `zhipu`        | [详细文档](./glm.md)      |
| **豆包 Doubao**   | 字节跳动企业级 | `doubao`, `bytedance` | [详细文档](./doubao.md)   |
| **月之暗面 Kimi** | 长上下文 200K  | `moonshot`, `kimi`    | [详细文档](./moonshot.md) |

## 快速开始

### 1. 基础使用

```go
import (
	"github.com/nabob-ai/nomad/pkg/provider"
	"github.com/nabob-ai/nomad/pkg/types"
)

// 创建配置
config := &types.ModelConfig{
	Provider: "openai",    // 选择 Provider
	Model:    "gpt-4o",    // 选择模型
	APIKey:   "your-key",  // API Key
}

// 使用工厂创建 Provider
factory := provider.NewMultiProviderFactory()
p, err := factory.Create(config)

// 发送消息
messages := []types.Message{
	{Role: types.RoleUser, Content: "Hello!"},
}

response, err := p.Complete(ctx, messages, nil)
fmt.Println(response.Message.Content)
```

### 2. 切换 Provider

只需修改 `Provider` 字段，代码无需改动：

```go
// 从 OpenAI 切换到 Groq
config.Provider = "groq"
config.Model = "llama-3.3-70b-versatile"

// 切换到本地 Ollama
config.Provider = "ollama"
config.Model = "llama3.2"
config.BaseURL = "http://localhost:11434/v1"
config.APIKey = "" // Ollama 不需要 API Key
```

## 功能对比

| 功能         | OpenAI   | Anthropic | Gemini      | Groq | OpenRouter | Mistral | Ollama | DeepSeek | GLM  | Doubao | Moonshot |
| ------------ | -------- | --------- | ----------- | ---- | ---------- | ------- | ------ | -------- | ---- | ------ | -------- |
| 流式输出     | ✅       | ✅        | ✅          | ✅   | ✅         | ✅      | ✅     | ✅       | ✅   | ✅     | ✅       |
| 工具调用     | ✅       | ✅        | ✅          | ✅   | ✅         | ✅      | ✅     | ✅       | ✅   | ✅     | ✅       |
| 视觉输入     | ✅       | ✅        | ✅          | ❌   | ✅         | ✅      | ⚠️     | ❌       | ⚠️   | ⚠️     | ❌       |
| 音频输入     | ✅       | ❌        | ✅          | ❌   | ✅         | ❌      | ❌     | ❌       | ❌   | ❌     | ❌       |
| 视频输入     | ❌       | ❌        | ✅          | ❌   | ❌         | ❌      | ❌     | ❌       | ❌   | ❌     | ❌       |
| 推理模型     | ✅ o1/o3 | ❌        | ✅ Thinking | ❌   | ✅         | ✅      | ❌     | ✅ R1    | ❌   | ❌     | ❌       |
| Prompt Cache | ✅       | ✅        | ✅          | ❌   | ✅         | ❌      | ❌     | ❌       | ❌   | ❌     | ❌       |
| 超长上下文   | 128K     | 200K      | 1M-2M       | 128K | -          | 128K    | 128K   | 64K      | 128K | 128K   | 200K     |
| 本地部署     | ❌       | ❌        | ❌          | ❌   | ❌         | ❌      | ✅     | ❌       | ❌   | ❌     | ❌       |
| 无需 API Key | ❌       | ❌        | ❌          | ❌   | ❌         | ❌      | ✅     | ❌       | ❌   | ❌     | ❌       |

**图例**：

- ✅ 完全支持
- ⚠️ 部分模型支持
- ❌ 不支持

## 使用场景推荐

### 🚀 追求速度

**Groq** - 业界最快的推理速度

```go
config := &types.ModelConfig{
	Provider: "groq",
	Model:    "llama-3.3-70b-versatile",
	APIKey:   "gsk-xxx",
}
```

**适合**：实时对话、客服系统、快速原型

### 💰 成本优化

**Ollama** - 本地部署，零 API 成本

```go
config := &types.ModelConfig{
	Provider: "ollama",
	Model:    "llama3.2",
	BaseURL:  "http://localhost:11434/v1",
}
```

**适合**：开发测试、离线应用、隐私敏感

### 🎯 通用场景

**OpenAI GPT-4o** - 性价比最高

```go
config := &types.ModelConfig{
	Provider: "openai",
	Model:    "gpt-4o",
	APIKey:   "sk-xxx",
}
```

**适合**：大多数应用场景

### 🧠 复杂推理

**OpenAI o1-preview** 或 **DeepSeek R1**

```go
// OpenAI o1
config := &types.ModelConfig{
	Provider: "openai",
	Model:    "o1-preview",
	APIKey:   "sk-xxx",
}

// DeepSeek R1（更便宜）
config := &types.ModelConfig{
	Provider: "deepseek",
	Model:    "deepseek-reasoner",
	APIKey:   "sk-xxx",
}
```

**适合**：数学问题、代码调试、复杂规划

### 📚 长上下文

**Gemini** - 1M-2M 超长上下文（最长）

```go
config := &types.ModelConfig{
	Provider: "gemini",
	Model:    "gemini-1.5-pro",    // 2M tokens
	// Model:    "gemini-2.0-flash-exp", // 1M tokens
	APIKey:   "your-key",
}
```

**适合**：完整代码仓库分析、长篇文档处理、视频理解

**Moonshot Kimi** - 200K 超长上下文

```go
config := &types.ModelConfig{
	Provider: "moonshot",
	Model:    "moonshot-v1-128k", // 或 moonshot-v1-32k
	APIKey:   "sk-xxx",
}
```

**适合**：文档分析、长文本理解

### 🌐 模型聚合

**OpenRouter** - 一次接入，数百模型

```go
config := &types.ModelConfig{
	Provider: "openrouter",
	Model:    "openai/gpt-4o",      // 或任何支持的模型
	APIKey:   "sk-or-xxx",
}

// 随时切换模型，无需更改代码
config.Model = "anthropic/claude-3-opus"
config.Model = "google/gemini-pro"
config.Model = "meta-llama/llama-3-70b"
```

**适合**：需要多模型支持、模型 A/B 测试

## 环境变量配置

可以通过环境变量配置 API Keys：

```bash
# 国际 Providers
export OPENAI_API_KEY="sk-xxx"
export ANTHROPIC_API_KEY="sk-ant-xxx"
export GEMINI_API_KEY="your-key"
export GROQ_API_KEY="gsk-xxx"
export OPENROUTER_API_KEY="sk-or-xxx"
export MISTRAL_API_KEY="xxx"

# 中国 Providers
export DEEPSEEK_API_KEY="sk-xxx"
export ZHIPU_API_KEY="xxx"
export DOUBAO_API_KEY="xxx"
export MOONSHOT_API_KEY="sk-xxx"

# Ollama（本地部署，可选）
export OLLAMA_BASE_URL="http://localhost:11434/v1"
```

然后在代码中使用：

```go
import "os"

config := &types.ModelConfig{
	Provider: "openai",
	Model:    "gpt-4o",
	APIKey:   os.Getenv("OPENAI_API_KEY"),
}
```

## 统一接口

所有 Provider 都实现相同的接口，确保代码兼容性：

```go
type Provider interface {
	// 流式对话
	Stream(ctx context.Context, messages []Message, opts *StreamOptions) (<-chan StreamChunk, error)

	// 非流式对话
	Complete(ctx context.Context, messages []Message, opts *StreamOptions) (*CompleteResponse, error)

	// 获取配置
	Config() *ModelConfig

	// 获取能力
	Capabilities() ProviderCapabilities

	// 设置/获取系统提示词
	SetSystemPrompt(prompt string) error
	GetSystemPrompt() string

	// 关闭连接
	Close() error
}
```

## 架构优势

### OpenAI 兼容层

大多数 Provider 基于 OpenAI 兼容层实现，具有以下优势：

1. **快速开发**：新增 Provider 只需 50-100 行代码
2. **统一体验**：API 调用方式完全一致
3. **自动重试**：内置 429/5xx 错误重试机制
4. **流式优化**：高效的 SSE 流式解析
5. **多模态支持**：统一的图片/音频处理

### 可扩展性

```go
// 轻松添加新 Provider
type NewProvider struct {
	*OpenAICompatibleProvider
}

func NewNewProvider(config *types.ModelConfig) (Provider, error) {
	options := &OpenAICompatibleOptions{
		RequireAPIKey: true,
		DefaultModel:  "model-name",
	}

	return NewOpenAICompatibleProvider(config, "https://api.example.com/v1", "NewProvider", options)
}
```

## 性能优化

### 连接复用

Provider 实例会复用 HTTP 连接：

```go
// ✅ 推荐：复用 Provider 实例
provider, _ := factory.Create(config)
for i := 0; i < 100; i++ {
	response, _ := provider.Complete(ctx, messages, nil)
}

// ❌ 避免：每次创建新实例
for i := 0; i < 100; i++ {
	provider, _ := factory.Create(config)
	response, _ := provider.Complete(ctx, messages, nil)
}
```

### 并发请求

```go
// 并发发送多个请求
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
	wg.Add(1)
	go func(index int) {
		defer wg.Done()
		response, _ := provider.Complete(ctx, messages[index], nil)
		// 处理响应
	}(i)
}
wg.Wait()
```

## 故障转移

```go
// 主 Provider + 备用 Provider
primaryConfig := &types.ModelConfig{
	Provider: "openai",
	Model:    "gpt-4o",
	APIKey:   os.Getenv("OPENAI_API_KEY"),
}

fallbackConfig := &types.ModelConfig{
	Provider: "groq",
	Model:    "llama-3.3-70b-versatile",
	APIKey:   os.Getenv("GROQ_API_KEY"),
}

// 尝试主 Provider
primaryProvider, _ := factory.Create(primaryConfig)
response, err := primaryProvider.Complete(ctx, messages, nil)

if err != nil {
	// 失败时使用备用 Provider
	fallbackProvider, _ := factory.Create(fallbackConfig)
	response, err = fallbackProvider.Complete(ctx, messages, nil)
}
```

## 下一步

- 查看具体 Provider 的详细文档
- 了解 [Agent 集成](../core-concepts/agent.md)
- 学习 [工具调用](../examples/tools/custom.md)
- 探索 [中间件](../core-concepts/middleware.md)

## 相关资源

- [Provider 接口定义](https://pkg.go.dev/github.com/nabob-ai/nomad/pkg/provider)
- [类型定义](https://pkg.go.dev/github.com/nabob-ai/nomad/pkg/types)
- [GitHub 示例](https://github.com/nabob-ai/nomad/tree/main/examples)
