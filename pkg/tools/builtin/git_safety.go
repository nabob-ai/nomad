package builtin

import (
	"fmt"
	"regexp"
	"strings"
)

// GitSafetyLevel Git 安全级别
type GitSafetyLevel int

const (
	// GitSafetyLevelStrict 严格模式：所有危险命令都需要用户确认
	GitSafetyLevelStrict GitSafetyLevel = iota
	// GitSafetyLevelNormal 正常模式：只有最危险的命令需要确认
	GitSafetyLevelNormal
	// GitSafetyLevelPermissive 宽松模式：仅阻止极端危险的命令
	GitSafetyLevelPermissive
)

// GitCommandRisk Git 命令风险级别
type GitCommandRisk int

const (
	// GitRiskSafe 安全命令（只读）
	GitRiskSafe GitCommandRisk = iota
	// GitRiskLow 低风险命令（本地修改）
	GitRiskLow
	// GitRiskMedium 中等风险（远程操作）
	GitRiskMedium
	// GitRiskHigh 高风险（可能丢失数据）
	GitRiskHigh
	// GitRiskCritical 极高风险（不可逆操作）
	GitRiskCritical
)

// GitSafetyCheck Git 安全检查结果
type GitSafetyCheck struct {
	IsGitCommand     bool           `json:"is_git_command"`
	Risk             GitCommandRisk `json:"risk"`
	RiskName         string         `json:"risk_name"`
	Command          string         `json:"command"`
	RequiresApproval bool           `json:"requires_approval"`
	Blocked          bool           `json:"blocked"`
	Reason           string         `json:"reason"`
	Warnings         []string       `json:"warnings"`
	Recommendations  []string       `json:"recommendations"`
}

// GitSafetyValidator Git 安全验证器
type GitSafetyValidator struct {
	level    GitSafetyLevel
	patterns map[string]*gitCommandPattern
}

type gitCommandPattern struct {
	pattern *regexp.Regexp
	risk    GitCommandRisk
	reason  string
	blocked bool
}

// NewGitSafetyValidator 创建 Git 安全验证器
func NewGitSafetyValidator(level GitSafetyLevel) *GitSafetyValidator {
	v := &GitSafetyValidator{
		level:    level,
		patterns: make(map[string]*gitCommandPattern),
	}
	v.initPatterns()
	return v
}

func (v *GitSafetyValidator) initPatterns() {
	// 极高风险：强制推送到主分支（支持 --force, -f, --force-with-lease）
	v.addPattern("force_push_main", `git\s+push\s+.*(-f|--force)\s+.*\s*(main|mnomad)`, GitRiskCritical,
		"Force push to main/mnomad branch can cause irreversible data loss", true)
	v.addPattern("force_push_main_alt", `git\s+push\s+.*\s*(main|mnomad).*(-f|--force)`, GitRiskCritical,
		"Force push to main/mnomad branch can cause irreversible data loss", true)
	v.addPattern("force_push_main_short", `git\s+push\s+-f\s+\S+\s+(main|mnomad)`, GitRiskCritical,
		"Force push to main/mnomad branch can cause irreversible data loss", true)

	// 高风险：强制推送（支持简写 -f）
	v.addPattern("force_push", `git\s+push\s+.*--force`, GitRiskHigh,
		"Force push can overwrite remote history and cause data loss", false)
	v.addPattern("force_push_short", `git\s+push\s+.*-f`, GitRiskHigh,
		"Force push can overwrite remote history and cause data loss", false)
	v.addPattern("force_push_lease", `git\s+push\s+.*--force-with-lease`, GitRiskHigh,
		"Force push (even with lease) can overwrite remote history", false)

	// 高风险：硬重置
	v.addPattern("hard_reset", `git\s+reset\s+--hard`, GitRiskHigh,
		"Hard reset discards all uncommitted changes permanently", false)
	v.addPattern("hard_reset_remote", `git\s+reset\s+--hard\s+origin`, GitRiskHigh,
		"Hard reset to remote can discard local commits", false)

	// 高风险：清理未跟踪文件
	v.addPattern("clean_force", `git\s+clean\s+-[dDfFxX]*f`, GitRiskHigh,
		"Git clean -f permanently removes untracked files", false)
	v.addPattern("clean_all", `git\s+clean\s+-[dDfFxX]*d`, GitRiskHigh,
		"Git clean -d removes untracked directories", false)

	// 中等风险：修改历史
	v.addPattern("rebase_interactive", `git\s+rebase\s+-i`, GitRiskMedium,
		"Interactive rebase modifies commit history (not supported in non-interactive mode)", false)
	v.addPattern("commit_amend", `git\s+commit\s+.*--amend`, GitRiskMedium,
		"Amending commits modifies history, verify authorship first", false)
	v.addPattern("cherry_pick", `git\s+cherry-pick`, GitRiskMedium,
		"Cherry-pick can cause conflicts and duplicate commits", false)

	// 中等风险：跳过验证
	v.addPattern("no_verify", `git\s+(commit|push)\s+.*--no-verify`, GitRiskMedium,
		"Skipping hooks bypasses quality checks", false)
	v.addPattern("no_gpg_sign", `git\s+commit\s+.*--no-gpg-sign`, GitRiskLow,
		"Skipping GPG signing", false)

	// 中等风险：远程操作
	v.addPattern("remote_delete", `git\s+push\s+.*--delete`, GitRiskMedium,
		"Deleting remote branches/tags is irreversible without backup", false)
	v.addPattern("remote_prune", `git\s+remote\s+prune`, GitRiskMedium,
		"Pruning removes remote-tracking references", false)

	// 中等风险：分支删除
	v.addPattern("branch_delete_force", `git\s+branch\s+-[dD]\s+`, GitRiskMedium,
		"Deleting branches can lose unmerged work", false)

	// 中等风险：配置修改
	v.addPattern("config_global", `git\s+config\s+--global`, GitRiskMedium,
		"Modifying global git config affects all repositories", false)
	v.addPattern("config_system", `git\s+config\s+--system`, GitRiskHigh,
		"Modifying system git config requires admin and affects all users", false)

	// 中等风险：子模块操作
	v.addPattern("submodule_deinit", `git\s+submodule\s+deinit`, GitRiskMedium,
		"Deinitializing submodules removes their contents", false)

	// 低风险：一般推送
	v.addPattern("push", `git\s+push`, GitRiskLow,
		"Pushing changes to remote repository", false)

	// 低风险：拉取和合并
	v.addPattern("pull_rebase", `git\s+pull\s+.*--rebase`, GitRiskLow,
		"Pull with rebase modifies local history", false)
	v.addPattern("merge", `git\s+merge`, GitRiskLow,
		"Merging branches", false)

	// 安全：只读操作
	v.addPattern("status", `git\s+status`, GitRiskSafe, "", false)
	v.addPattern("log", `git\s+log`, GitRiskSafe, "", false)
	v.addPattern("diff", `git\s+diff`, GitRiskSafe, "", false)
	v.addPattern("show", `git\s+show`, GitRiskSafe, "", false)
	v.addPattern("branch_list", `git\s+branch\s*$`, GitRiskSafe, "", false)
	v.addPattern("remote_list", `git\s+remote\s*$`, GitRiskSafe, "", false)
	v.addPattern("fetch", `git\s+fetch`, GitRiskSafe, "", false)
	v.addPattern("stash_list", `git\s+stash\s+list`, GitRiskSafe, "", false)
}

func (v *GitSafetyValidator) addPattern(name, pattern string, risk GitCommandRisk, reason string, blocked bool) {
	regex, err := regexp.Compile(`(?i)` + pattern)
	if err != nil {
		return
	}
	v.patterns[name] = &gitCommandPattern{
		pattern: regex,
		risk:    risk,
		reason:  reason,
		blocked: blocked,
	}
}

// Check 检查 Git 命令的安全性
func (v *GitSafetyValidator) Check(command string) *GitSafetyCheck {
	result := &GitSafetyCheck{
		Command:         command,
		Risk:            GitRiskSafe,
		RiskName:        "safe",
		Warnings:        []string{},
		Recommendations: []string{},
	}

	// 检查是否是 git 命令
	// 必须是 "git " 开头或包含 "git " 子字符串（排除 github, gitk 等）
	trimmedCmd := strings.TrimSpace(command)
	lowerCmd := strings.ToLower(trimmedCmd)
	isGit := strings.HasPrefix(lowerCmd, "git ") ||
		lowerCmd == "git" ||
		strings.Contains(lowerCmd, " git ") ||
		strings.HasSuffix(lowerCmd, " git")

	if !isGit {
		result.IsGitCommand = false
		return result
	}
	result.IsGitCommand = true

	// 检查所有模式
	for _, p := range v.patterns {
		if p.pattern.MatchString(command) {
			// 更新为更高的风险级别
			if p.risk > result.Risk {
				result.Risk = p.risk
				result.RiskName = riskToName(p.risk)
			}

			// 添加警告
			if p.reason != "" {
				result.Warnings = append(result.Warnings, p.reason)
			}

			// 检查是否被阻止
			if p.blocked {
				result.Blocked = true
				result.Reason = p.reason
			}
		}
	}

	// 根据安全级别和风险确定是否需要批准
	result.RequiresApproval = v.requiresApproval(result.Risk)

	// 添加建议
	result.Recommendations = v.getRecommendations(command, result.Risk)

	return result
}

func (v *GitSafetyValidator) requiresApproval(risk GitCommandRisk) bool {
	switch v.level {
	case GitSafetyLevelStrict:
		// 严格模式：中等风险及以上都需要批准
		return risk >= GitRiskMedium
	case GitSafetyLevelNormal:
		// 正常模式：高风险及以上需要批准
		return risk >= GitRiskHigh
	case GitSafetyLevelPermissive:
		// 宽松模式：只有极高风险需要批准
		return risk >= GitRiskCritical
	default:
		return risk >= GitRiskHigh
	}
}

func (v *GitSafetyValidator) getRecommendations(command string, risk GitCommandRisk) []string {
	var recs []string

	lowerCmd := strings.ToLower(command)

	// 针对特定命令的建议
	if strings.Contains(lowerCmd, "push") && strings.Contains(lowerCmd, "--force") {
		recs = append(recs, "Consider using --force-with-lease instead of --force for safer force pushes")
		recs = append(recs, "Ensure no one else has pushed to this branch before force pushing")
	}

	if strings.Contains(lowerCmd, "reset") && strings.Contains(lowerCmd, "--hard") {
		recs = append(recs, "Consider using git stash to save uncommitted changes first")
		recs = append(recs, "Use git reflog to recover lost commits if needed")
	}

	if strings.Contains(lowerCmd, "commit") && strings.Contains(lowerCmd, "--amend") {
		recs = append(recs, "Check authorship before amending: git log -1 --format='%an %ae'")
		recs = append(recs, "Only amend commits that haven't been pushed to remote")
		recs = append(recs, "Verify branch is ahead of remote: git status")
	}

	if strings.Contains(lowerCmd, "clean") {
		recs = append(recs, "Run git clean -n first to preview what will be deleted")
	}

	if strings.Contains(lowerCmd, "rebase") && strings.Contains(lowerCmd, "-i") {
		recs = append(recs, "Interactive rebase is not supported in non-interactive environments")
		recs = append(recs, "Consider using non-interactive rebase commands instead")
	}

	// 通用建议
	if risk >= GitRiskMedium {
		recs = append(recs, "Consider creating a backup branch before this operation")
	}

	return recs
}

func riskToName(risk GitCommandRisk) string {
	switch risk {
	case GitRiskSafe:
		return "safe"
	case GitRiskLow:
		return "low"
	case GitRiskMedium:
		return "medium"
	case GitRiskHigh:
		return "high"
	case GitRiskCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// FormatCheckResult 格式化检查结果为用户可读的消息
func (c *GitSafetyCheck) FormatCheckResult() string {
	if !c.IsGitCommand {
		return ""
	}

	if c.Blocked {
		return fmt.Sprintf("🚫 BLOCKED: This git command is not allowed.\nReason: %s\nCommand: %s", c.Reason, c.Command)
	}

	if c.RequiresApproval {
		msg := fmt.Sprintf("⚠️ GIT SAFETY WARNING [%s risk]\n", strings.ToUpper(c.RiskName))
		msg += fmt.Sprintf("Command: %s\n", c.Command)

		if len(c.Warnings) > 0 {
			msg += "\nWarnings:\n"
			var msgSb314 strings.Builder
			for _, w := range c.Warnings {
				msgSb314.WriteString(fmt.Sprintf("  • %s\n", w))
			}
			msg += msgSb314.String()
		}

		if len(c.Recommendations) > 0 {
			msg += "\nRecommendations:\n"
			var msgSb321 strings.Builder
			for _, r := range c.Recommendations {
				msgSb321.WriteString(fmt.Sprintf("  • %s\n", r))
			}
			msg += msgSb321.String()
		}

		msg += "\nThis command requires user approval before execution."
		return msg
	}

	return ""
}

// 全局 Git 安全验证器
var globalGitSafetyValidator *GitSafetyValidator

// GetGlobalGitSafetyValidator 获取全局 Git 安全验证器
func GetGlobalGitSafetyValidator() *GitSafetyValidator {
	if globalGitSafetyValidator == nil {
		// 默认使用严格模式
		globalGitSafetyValidator = NewGitSafetyValidator(GitSafetyLevelStrict)
	}
	return globalGitSafetyValidator
}

// SetGlobalGitSafetyLevel 设置全局 Git 安全级别
func SetGlobalGitSafetyLevel(level GitSafetyLevel) {
	globalGitSafetyValidator = NewGitSafetyValidator(level)
}
