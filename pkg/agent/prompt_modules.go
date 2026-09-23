package agent

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/nabob-ai/nomad/pkg/types"
)

// BasePromptModule 基础 Prompt（来自模板）
type BasePromptModule struct{}

func (m *BasePromptModule) Name() string                      { return "base" }
func (m *BasePromptModule) Priority() int                     { return 0 }
func (m *BasePromptModule) Condition(ctx *PromptContext) bool { return true }
func (m *BasePromptModule) Build(ctx *PromptContext) (string, error) {
	return ctx.Template.SystemPrompt, nil
}

// EnvironmentModule 环境信息模块
type EnvironmentModule struct{}

func (m *EnvironmentModule) Name() string  { return "environment" }
func (m *EnvironmentModule) Priority() int { return 10 }
func (m *EnvironmentModule) Condition(ctx *PromptContext) bool {
	return ctx.Environment != nil
}
func (m *EnvironmentModule) Build(ctx *PromptContext) (string, error) {
	env := ctx.Environment

	var lines []string
	lines = append(lines, "## Environment Information")
	lines = append(lines, "")
	lines = append(lines, "- Working Directory: "+env.WorkingDir)
	lines = append(lines, "- Platform: "+env.Platform)
	lines = append(lines, "- Date: "+env.Date.Format("2006-01-02"))

	// 精简 Git 信息，只保留关键内容以减少 token 消耗
	if env.GitRepo != nil && env.GitRepo.IsRepo {
		lines = append(lines, "- Git Branch: "+env.GitRepo.CurrentBranch)
		// 不再输出 git status 和 recent commits，这些可以通过工具获取
	}

	return strings.Join(lines, "\n"), nil
}

// ToolsManualModule 工具手册模块
type ToolsManualModule struct {
	Config *types.ToolsManualConfig
}

func (m *ToolsManualModule) Name() string  { return "tools_manual" }
func (m *ToolsManualModule) Priority() int { return 20 }
func (m *ToolsManualModule) Condition(ctx *PromptContext) bool {
	if m.Config != nil && m.Config.Mode == "none" {
		return false
	}
	return len(ctx.Tools) > 0
}
func (m *ToolsManualModule) Build(ctx *PromptContext) (string, error) {
	// 根据 Config 决定注入哪些工具
	var toolsToInclude []string

	if m.Config == nil || m.Config.Mode == "" || m.Config.Mode == "all" {
		// 默认：所有工具（除了 Exclude）
		for name := range ctx.Tools {
			if m.Config != nil && contains(m.Config.Exclude, name) {
				continue
			}
			toolsToInclude = append(toolsToInclude, name)
		}
	} else if m.Config.Mode == "listed" {
		// 仅包含 Include 列表中的工具
		if m.Config.Include != nil {
			for _, name := range m.Config.Include {
				if _, exists := ctx.Tools[name]; exists {
					toolsToInclude = append(toolsToInclude, name)
				}
			}
		}
	}

	if len(toolsToInclude) == 0 {
		return "", nil
	}

	sort.Strings(toolsToInclude)

	var lines []string
	lines = append(lines, "## Tools Manual")
	lines = append(lines, "")
	lines = append(lines, "The following tools are available for your use. Use them when appropriate instead of doing everything in natural language.")
	lines = append(lines, "")

	for _, name := range toolsToInclude {
		tool := ctx.Tools[name]
		summary := tool.Description()
		if summary == "" {
			summary = "No detailed manual; infer from tool name and input schema."
		}
		lines = append(lines, fmt.Sprintf("- `%s`: %s", name, summary))
	}

	return strings.Join(lines, "\n"), nil
}

// SandboxModule 沙箱信息模块
type SandboxModule struct{}

func (m *SandboxModule) Name() string  { return "sandbox" }
func (m *SandboxModule) Priority() int { return 15 }
func (m *SandboxModule) Condition(ctx *PromptContext) bool {
	return ctx.Sandbox != nil
}
func (m *SandboxModule) Build(ctx *PromptContext) (string, error) {
	sb := ctx.Sandbox

	var lines []string
	lines = append(lines, "## Sandbox Environment")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("- Type: %s", sb.Kind))
	lines = append(lines, "- Working Directory: "+sb.WorkDir)

	if len(sb.AllowPaths) > 0 {
		lines = append(lines, "- Allowed Paths:")
		for _, path := range sb.AllowPaths {
			lines = append(lines, "  - "+path)
		}
	}

	return strings.Join(lines, "\n"), nil
}

// TodoReminderModule Todo 提醒模块
type TodoReminderModule struct {
	Config *types.TodoConfig
}

func (m *TodoReminderModule) Name() string  { return "todo_reminder" }
func (m *TodoReminderModule) Priority() int { return 25 }
func (m *TodoReminderModule) Condition(ctx *PromptContext) bool {
	return m.Config != nil && m.Config.Enabled && m.Config.ReminderOnStart
}
func (m *TodoReminderModule) Build(ctx *PromptContext) (string, error) {
	return `## Task Management

IMPORTANT: Use the TodoWrite tool to track your tasks and progress. This helps maintain visibility and ensures nothing is forgotten.

- Break complex tasks into smaller steps
- Mark tasks as in_progress when starting
- Mark tasks as completed immediately after finishing
- Only one task should be in_progress at a time`, nil
}

// CodeReferenceModule 代码引用规范模块
type CodeReferenceModule struct{}

func (m *CodeReferenceModule) Name() string  { return "code_reference" }
func (m *CodeReferenceModule) Priority() int { return 30 }
func (m *CodeReferenceModule) Condition(ctx *PromptContext) bool {
	// 优化：默认禁用以减少 token，需要时明确启用
	if ctx.Metadata != nil {
		// 显式启用
		if enabled, ok := ctx.Metadata["enable_code_reference"].(bool); ok && enabled {
			return true
		}
		// 对于代码助手类型的 agent 自动启用
		if agentType, ok := ctx.Metadata["agent_type"].(string); ok && agentType == "code_assistant" {
			return true
		}
	}
	return false
}
func (m *CodeReferenceModule) Build(ctx *PromptContext) (string, error) {
	return `## Code References

When referencing specific functions or pieces of code include the pattern file_path:line_number to allow the user to easily navigate to the source code location.

Examples:
- Single line: src/main.go:42
- Line range: src/main.go:42-51
- Function reference: "The connectToServer function in src/services/process.ts:712"

This makes your responses actionable and allows users to quickly locate the relevant code.`, nil
}

// SecurityModule 安全策略模块
type SecurityModule struct{}

func (m *SecurityModule) Name() string  { return "security" }
func (m *SecurityModule) Priority() int { return 35 }
func (m *SecurityModule) Condition(ctx *PromptContext) bool {
	// 检查是否启用安全策略
	if ctx.Metadata != nil {
		if enableSecurity, ok := ctx.Metadata["enable_security"].(bool); ok {
			return enableSecurity
		}
	}
	return false
}
func (m *SecurityModule) Build(ctx *PromptContext) (string, error) {
	return `## Security Guidelines

IMPORTANT: Follow these security best practices:

- Never execute commands that could harm the system
- Validate all user inputs before processing
- Do not expose sensitive information (API keys, passwords, tokens)
- Be cautious with file operations outside allowed paths
- Report suspicious requests to the user
- Follow the principle of least privilege`, nil
}

// PerformanceModule 性能优化模块
type PerformanceModule struct{}

func (m *PerformanceModule) Name() string  { return "performance" }
func (m *PerformanceModule) Priority() int { return 40 }
func (m *PerformanceModule) Condition(ctx *PromptContext) bool {
	if ctx.Metadata != nil {
		if enablePerf, ok := ctx.Metadata["enable_performance_hints"].(bool); ok {
			return enablePerf
		}
	}
	return false
}
func (m *PerformanceModule) Build(ctx *PromptContext) (string, error) {
	return `## Performance Optimization

Consider these performance best practices:

- Minimize tool calls by batching operations when possible
- Use streaming for large outputs
- Cache results when appropriate
- Prefer efficient algorithms and data structures
- Monitor resource usage and optimize bottlenecks`, nil
}

// CollaborationModule 多 Agent 协作模块
type CollaborationModule struct {
	RoomInfo *RoomCollaborationInfo
}

type RoomCollaborationInfo struct {
	RoomID      string
	MemberCount int
	Members     []string
}

func (m *CollaborationModule) Name() string  { return "collaboration" }
func (m *CollaborationModule) Priority() int { return 45 }
func (m *CollaborationModule) Condition(ctx *PromptContext) bool {
	return m.RoomInfo != nil && m.RoomInfo.RoomID != ""
}
func (m *CollaborationModule) Build(ctx *PromptContext) (string, error) {
	var lines []string
	lines = append(lines, "## Multi-Agent Collaboration")
	lines = append(lines, "")
	lines = append(lines, "You are working in a collaborative room: "+m.RoomInfo.RoomID)
	lines = append(lines, fmt.Sprintf("Total members: %d", m.RoomInfo.MemberCount))

	if len(m.RoomInfo.Members) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Room members:")
		for _, member := range m.RoomInfo.Members {
			lines = append(lines, "- "+member)
		}
	}

	lines = append(lines, "")
	lines = append(lines, "Collaboration guidelines:")
	lines = append(lines, "- Use @mention to address specific members")
	lines = append(lines, "- Coordinate tasks to avoid duplication")
	lines = append(lines, "- Share progress and findings with the team")
	lines = append(lines, "- Ask for help when needed")

	return strings.Join(lines, "\n"), nil
}

// WorkflowModule 工作流上下文模块
type WorkflowModule struct {
	WorkflowInfo *WorkflowContextInfo
}

type WorkflowContextInfo struct {
	WorkflowID   string
	CurrentStep  string
	TotalSteps   int
	StepIndex    int
	PreviousStep string
	NextStep     string
}

func (m *WorkflowModule) Name() string  { return "workflow" }
func (m *WorkflowModule) Priority() int { return 50 }
func (m *WorkflowModule) Condition(ctx *PromptContext) bool {
	return m.WorkflowInfo != nil && m.WorkflowInfo.WorkflowID != ""
}
func (m *WorkflowModule) Build(ctx *PromptContext) (string, error) {
	info := m.WorkflowInfo

	var lines []string
	lines = append(lines, "## Workflow Context")
	lines = append(lines, "")
	lines = append(lines, "Workflow ID: "+info.WorkflowID)
	lines = append(lines, fmt.Sprintf("Current Step: %s (Step %d of %d)", info.CurrentStep, info.StepIndex+1, info.TotalSteps))

	if info.PreviousStep != "" {
		lines = append(lines, "Previous Step: "+info.PreviousStep)
	}

	if info.NextStep != "" {
		lines = append(lines, "Next Step: "+info.NextStep)
	}

	lines = append(lines, "")
	lines = append(lines, "Focus on completing the current step efficiently before moving to the next.")

	return strings.Join(lines, "\n"), nil
}

// CustomInstructionsModule 用户自定义指令模块
type CustomInstructionsModule struct {
	Instructions string
}

func (m *CustomInstructionsModule) Name() string  { return "custom_instructions" }
func (m *CustomInstructionsModule) Priority() int { return 55 }
func (m *CustomInstructionsModule) Condition(ctx *PromptContext) bool {
	return m.Instructions != ""
}
func (m *CustomInstructionsModule) Build(ctx *PromptContext) (string, error) {
	return "## Custom Instructions\n\n" + m.Instructions, nil
}

// CapabilitiesModule Agent 能力说明模块
type CapabilitiesModule struct{}

func (m *CapabilitiesModule) Name() string  { return "capabilities" }
func (m *CapabilitiesModule) Priority() int { return 5 }
func (m *CapabilitiesModule) Condition(ctx *PromptContext) bool {
	if ctx.Metadata != nil {
		if showCaps, ok := ctx.Metadata["show_capabilities"].(bool); ok {
			return showCaps
		}
	}
	return false
}
func (m *CapabilitiesModule) Build(ctx *PromptContext) (string, error) {
	var capabilities []string

	// 基于可用工具推断能力
	if ctx.Tools != nil {
		if _, hasRead := ctx.Tools["Read"]; hasRead {
			capabilities = append(capabilities, "Read and analyze files")
		}
		if _, hasWrite := ctx.Tools["Write"]; hasWrite {
			capabilities = append(capabilities, "Create and modify files")
		}
		if _, hasBash := ctx.Tools["Bash"]; hasBash {
			capabilities = append(capabilities, "Execute shell commands")
		}
		if _, hasWebSearch := ctx.Tools["WebSearch"]; hasWebSearch {
			capabilities = append(capabilities, "Search the web for information")
		}
		if _, hasTodo := ctx.Tools["TodoWrite"]; hasTodo {
			capabilities = append(capabilities, "Manage tasks and track progress")
		}
	}

	if len(capabilities) == 0 {
		return "", nil
	}

	var lines []string
	lines = append(lines, "## Your Capabilities")
	lines = append(lines, "")
	lines = append(lines, "You can:")
	for _, cap := range capabilities {
		lines = append(lines, "- "+cap)
	}

	return strings.Join(lines, "\n"), nil
}

// LimitationsModule 限制说明模块
type LimitationsModule struct{}

func (m *LimitationsModule) Name() string  { return "limitations" }
func (m *LimitationsModule) Priority() int { return 60 }
func (m *LimitationsModule) Condition(ctx *PromptContext) bool {
	if ctx.Metadata != nil {
		if showLimits, ok := ctx.Metadata["show_limitations"].(bool); ok {
			return showLimits
		}
	}
	return false
}
func (m *LimitationsModule) Build(ctx *PromptContext) (string, error) {
	var lines []string
	lines = append(lines, "## Important Limitations")
	lines = append(lines, "")
	lines = append(lines, "Be aware of these limitations:")
	lines = append(lines, "- You cannot access the internet directly (unless WebSearch tool is available)")
	lines = append(lines, "- You cannot execute code outside the sandbox environment")
	lines = append(lines, "- You have limited context window - be concise")
	lines = append(lines, "- You cannot remember information across different sessions")

	// 基于沙箱类型添加特定限制
	if ctx.Sandbox != nil {
		if ctx.Sandbox.Kind == types.SandboxKindMock {
			lines = append(lines, "- Running in mock sandbox - file operations are simulated")
		}
		if len(ctx.Sandbox.AllowPaths) > 0 {
			lines = append(lines, "- File access is restricted to allowed paths only")
		}
	}

	return strings.Join(lines, "\n"), nil
}

// ContextWindowModule 上下文窗口管理模块
type ContextWindowModule struct {
	MaxTokens int
	Strategy  string
}

func (m *ContextWindowModule) Name() string  { return "context_window" }
func (m *ContextWindowModule) Priority() int { return 65 }
func (m *ContextWindowModule) Condition(ctx *PromptContext) bool {
	return m.MaxTokens > 0
}
func (m *ContextWindowModule) Build(ctx *PromptContext) (string, error) {
	var lines []string
	lines = append(lines, "## Context Window Management")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Maximum context tokens: %d", m.MaxTokens))

	if m.Strategy != "" {
		lines = append(lines, "Compression strategy: "+m.Strategy)
	}

	lines = append(lines, "")
	lines = append(lines, "To manage context efficiently:")
	lines = append(lines, "- Summarize long outputs")
	lines = append(lines, "- Reference files by path instead of including full content")
	lines = append(lines, "- Use tools to retrieve information on-demand")
	lines = append(lines, "- Focus on relevant information only")

	return strings.Join(lines, "\n"), nil
}

// ProfessionalObjectivityModule 专业客观性模块
type ProfessionalObjectivityModule struct{}

func (m *ProfessionalObjectivityModule) Name() string  { return "professional_objectivity" }
func (m *ProfessionalObjectivityModule) Priority() int { return 8 }
func (m *ProfessionalObjectivityModule) Condition(ctx *PromptContext) bool {
	// 优化：默认禁用以减少 token，这些原则可以内嵌到 base prompt
	if ctx.Metadata != nil {
		if enabled, ok := ctx.Metadata["enable_objectivity"].(bool); ok && enabled {
			return true
		}
	}
	return false
}
func (m *ProfessionalObjectivityModule) Build(ctx *PromptContext) (string, error) {
	return `## Professional Objectivity

Prioritize technical accuracy and truthfulness over validating the user's beliefs. Focus on facts and problem-solving, providing direct, objective technical info without any unnecessary superlatives, praise, or emotional validation.

Guidelines:
- Apply the same rigorous standards to all ideas
- Disagree when necessary, even if it's not what the user wants to hear
- Provide respectful correction over false agreement
- Investigate to find the truth rather than confirming assumptions
- Avoid over-the-top validation or excessive praise like "You're absolutely right"
- Be direct and honest about limitations and tradeoffs`, nil
}

// ConcisenessModule 简洁性模块
// 强调简洁、高效的沟通风格
type ConcisenessModule struct{}

func (m *ConcisenessModule) Name() string  { return "conciseness" }
func (m *ConcisenessModule) Priority() int { return 9 }
func (m *ConcisenessModule) Condition(ctx *PromptContext) bool {
	// 优化：默认禁用以减少 token，这些原则可以内嵌到 base prompt
	if ctx.Metadata != nil {
		if enabled, ok := ctx.Metadata["enable_conciseness"].(bool); ok && enabled {
			return true
		}
	}
	return false
}
func (m *ConcisenessModule) Build(ctx *PromptContext) (string, error) {
	return `## Tone and Style

- Your responses should be short and concise
- You can use markdown for formatting
- Output text to communicate with the user; all text outside of tool use is displayed to the user
- Only use tools to complete tasks, never as a means to communicate
- NEVER create files unless absolutely necessary. Prefer editing existing files
- Only use emojis if the user explicitly requests them`, nil
}

// AvoidOverEngineeringModule 避免过度工程化模块
type AvoidOverEngineeringModule struct{}

func (m *AvoidOverEngineeringModule) Name() string  { return "avoid_over_engineering" }
func (m *AvoidOverEngineeringModule) Priority() int { return 12 }
func (m *AvoidOverEngineeringModule) Condition(ctx *PromptContext) bool {
	// 优化：默认禁用以减少 token，这些原则可以内嵌到 base prompt
	if ctx.Metadata != nil {
		if enabled, ok := ctx.Metadata["enable_avoid_over_engineering"].(bool); ok && enabled {
			return true
		}
	}
	return false
}
func (m *AvoidOverEngineeringModule) Build(ctx *PromptContext) (string, error) {
	return `## Avoid Over-Engineering

Only make changes that are directly requested or clearly necessary. Keep solutions simple and focused.

Principles:
- Don't add features, refactor code, or make "improvements" beyond what was asked
- A bug fix doesn't need surrounding code cleaned up
- A simple feature doesn't need extra configurability
- Don't add docstrings, comments, or type annotations to code you didn't change
- Don't add error handling for scenarios that can't happen
- Don't create helpers or abstractions for one-time operations
- Don't design for hypothetical future requirements
- Three similar lines of code is better than a premature abstraction`, nil
}

// PlanningWithoutTimelinesModule 无时间线规划模块
type PlanningWithoutTimelinesModule struct{}

func (m *PlanningWithoutTimelinesModule) Name() string  { return "planning_no_timelines" }
func (m *PlanningWithoutTimelinesModule) Priority() int { return 13 }
func (m *PlanningWithoutTimelinesModule) Condition(ctx *PromptContext) bool {
	if ctx.Metadata != nil {
		if enabled, ok := ctx.Metadata["enable_planning_guidelines"].(bool); ok {
			return enabled
		}
	}
	return false
}
func (m *PlanningWithoutTimelinesModule) Build(ctx *PromptContext) (string, error) {
	return `## Planning Guidelines

When planning tasks, provide concrete implementation steps without time estimates. Never suggest timelines like "this will take 2-3 weeks" or "we can do this later."

Focus on:
- What needs to be done
- Break work into actionable steps
- Let users decide scheduling
- Avoid committing to deadlines or durations`, nil
}

// GitSafetyModule Git 安全协议模块
type GitSafetyModule struct{}

func (m *GitSafetyModule) Name() string  { return "git_safety" }
func (m *GitSafetyModule) Priority() int { return 32 }
func (m *GitSafetyModule) Condition(ctx *PromptContext) bool {
	// 检查是否有 Git 环境
	if ctx.Environment != nil && ctx.Environment.GitRepo != nil && ctx.Environment.GitRepo.IsRepo {
		return true
	}
	return false
}
func (m *GitSafetyModule) Build(ctx *PromptContext) (string, error) {
	return `## Git Safety Protocol

CRITICAL: Follow these git safety rules:

### NEVER do:
- Update git config without explicit user request
- Run destructive/irreversible commands (push --force, hard reset) without explicit request
- Skip hooks (--no-verify, --no-gpg-sign) unless explicitly requested
- Force push to main/mnomad - warn user if they request it
- Use git commit --amend on commits you didn't create
- Commit changes unless the user explicitly asks

### ALWAYS do:
- Before amending: check authorship with git log -1 --format='%an %ae'
- Before amending: verify not pushed with git status
- Create descriptive commit messages explaining "why" not just "what"
- Run git status and git diff before committing
- Never use interactive flags (-i) as they require interactive input

### Commit Message Format:
- Summarize the change (new feature, bug fix, refactor, etc.)
- Focus on "why" rather than "what"
- End with: 🤖 Generated with AI assistance`, nil
}

// 辅助函数
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
