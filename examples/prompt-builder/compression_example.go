// Package main 演示如何使用 Prompt 压缩功能
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nabob-ai/nomad/pkg/agent"
	"github.com/nabob-ai/nomad/pkg/provider"
	"github.com/nabob-ai/nomad/pkg/types"
)

// RunCompressionExamples 运行压缩示例
func RunCompressionExamples() {
	// 示例 1: 直接使用压缩器
	directCompressionExample()

	// 示例 2: 通过模板配置使用压缩
	templateCompressionExample()
}

// directCompressionExample 直接使用 EnhancedPromptCompressor
func directCompressionExample() {
	fmt.Println("=== 直接压缩示例 ===")

	// 创建 Provider (这里用 nil 因为我们只测试简单压缩)
	var prov provider.Provider // 实际使用时需要创建真实的 Provider

	// 创建压缩器
	compressor := agent.NewEnhancedPromptCompressor(prov, "zh")

	// 生成一个长 prompt
	longPrompt := generateLongPrompt()
	fmt.Printf("原始长度: %d 字符\n", len(longPrompt))

	// 测试不同压缩级别
	levels := []struct {
		name  string
		level agent.CompressionLevel
	}{
		{"轻度压缩", agent.CompressionLevelLight},
		{"中度压缩", agent.CompressionLevelModerate},
		{"激进压缩", agent.CompressionLevelAggressive},
	}

	for _, l := range levels {
		result, err := compressor.Compress(context.Background(), longPrompt, &agent.CompressOptions{
			Mode:             agent.CompressionModeSimple, // 使用简单模式（不需要 LLM）
			Level:            l.level,
			PreserveSections: []string{"Tools Manual", "Security"},
		})
		if err != nil {
			log.Printf("压缩失败: %v", err)
			continue
		}

		fmt.Printf("\n%s:\n", l.name)
		fmt.Printf("  压缩后长度: %d 字符\n", result.CompressedLength)
		fmt.Printf("  压缩率: %.2f%%\n", result.CompressionRatio*100)
		fmt.Printf("  节省 Token: %d\n", result.TokensSaved)
	}
}

// templateCompressionExample 通过模板配置使用压缩
func templateCompressionExample() {
	fmt.Println("\n=== 模板配置压缩示例 ===")

	// 创建带压缩配置的模板
	compressionConfig := &types.PromptCompressionConfig{
		Enabled:          true,
		MaxLength:        5000, // 超过 5000 字符触发压缩
		TargetLength:     3000, // 压缩到 3000 字符
		Mode:             "hybrid",
		Level:            2, // 中度压缩
		PreserveSections: []string{"Tools Manual", "Security Guidelines"},
		CacheEnabled:     true,
		Language:         "zh",
	}

	templateID := "code-assistant-compressed"

	fmt.Printf("模板 ID: %s\n", templateID)
	fmt.Printf("压缩配置:\n")
	fmt.Printf("  启用: %v\n", compressionConfig.Enabled)
	fmt.Printf("  ���发阈值: %d 字符\n", compressionConfig.MaxLength)
	fmt.Printf("  目标长度: %d 字符\n", compressionConfig.TargetLength)
	fmt.Printf("  压缩模式: %s\n", compressionConfig.Mode)
	fmt.Printf("  压缩级别: %d\n", compressionConfig.Level)
	fmt.Printf("  保留段落: %v\n", compressionConfig.PreserveSections)
	fmt.Printf("  缓存启用: %v\n", compressionConfig.CacheEnabled)
	fmt.Printf("  语言: %s\n", compressionConfig.Language)
}

// generateLongPrompt 生成测试用的长 prompt
func generateLongPrompt() string {
	return `# System Prompt

You are a helpful AI assistant.

## Tools Manual

This section describes the available tools:

### ReadFile
Read contents of a file from the filesystem.
Parameters:
- path: The file path to read

### WriteFile
Write contents to a file.
Parameters:
- path: The file path to write
- content: The content to write

### Bash
Execute a shell command.
Parameters:
- command: The command to execute

## Security Guidelines

IMPORTANT: Follow these security rules at all times:

1. Never expose API keys or credentials
2. Always validate user input before processing
3. Use secure connections for network operations
4. Do not execute arbitrary code without user confirmation
5. Protect sensitive data in transit and at rest

## General Instructions

When responding to user queries:
1. Be helpful and accurate
2. Provide clear explanations
3. Use examples when helpful
4. Ask clarifying questions when needed

## Code Style Guidelines

Follow these coding conventions:
- Use meaningful variable names
- Add comments for complex logic
- Follow the project's existing patterns
- Write tests for new functionality

## Error Handling

When errors occur:
1. Log the error with context
2. Provide user-friendly error messages
3. Suggest possible solutions
4. Never expose internal details

## Performance Tips

For optimal performance:
- Cache frequently used data
- Minimize network requests
- Use efficient algorithms
- Profile before optimizing

## Additional Guidelines

Section A: More detailed instructions about specific tasks.
This section contains supplementary information that may be useful
in certain scenarios but is not critical for basic operation.

Section B: Extended documentation about advanced features.
This content provides additional context but can be compressed
if token limits are a concern.

Section C: Historical notes and version information.
These details are less critical for day-to-day operations
and can be summarized or removed during compression.

## Closing Notes

Remember to always be helpful and follow the guidelines above.
If in doubt, ask the user for clarification.`
}
