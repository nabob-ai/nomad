package tools

// ToolAnnotations 工具安全注解
// 用于描述工具的安全特性，帮助权限系统做出智能决策
type ToolAnnotations struct {
	// ReadOnly 工具是否只读（不修改任何状态）
	// 只读工具在 SmartApprove 模式下可以自动批准
	ReadOnly bool `json:"read_only"`

	// Destructive 工具是否具有破坏性（可能导致数据丢失）
	// 破坏性工具在大多数模式下都需要用户确认
	Destructive bool `json:"destructive"`

	// Idempotent 工具是否幂等（多次执行结果相同）
	// 幂等工具更安全，可以重试
	Idempotent bool `json:"idempotent"`

	// OpenWorld 工具是否涉及外部系统（网络、第三方 API）
	// 涉及外部系统的工具即使只读也可能有安全风险
	OpenWorld bool `json:"open_world"`

	// RiskLevel 风险级别 (0-4)
	// 0: safe - 完全安全的只读操作
	// 1: low - 低风险，可逆操作
	// 2: medium - 中等风险，需要注意
	// 3: high - 高风险，可能导致数据丢失
	// 4: critical - 极高风险，可能导致不可逆损失
	RiskLevel int `json:"risk_level"`

	// Category 工具分类
	// 例如: "filesystem", "execution", "network", "database", "system"
	Category string `json:"category,omitempty"`

	// RequiresConfirmation 是否需要用户确认
	// 设为 true 时，无论权限模式如何都需要确认
	RequiresConfirmation bool `json:"requires_confirmation,omitempty"`
}

// RiskLevel 风险级别常量
const (
	RiskLevelSafe     = 0 // 完全安全
	RiskLevelLow      = 1 // 低风险
	RiskLevelMedium   = 2 // 中等风险
	RiskLevelHigh     = 3 // 高风险
	RiskLevelCritical = 4 // 极高风险
)

// Category 工具分类常量
const (
	CategoryFilesystem = "filesystem" // 文件系统操作
	CategoryExecution  = "execution"  // 命令执行
	CategoryNetwork    = "network"    // 网络请求
	CategoryDatabase   = "database"   // 数据库操作
	CategorySystem     = "system"     // 系统操作
	CategoryMCP        = "mcp"        // MCP 工具
	CategoryCustom     = "custom"     // 自定义工具
)

// 预定义注解模板

// AnnotationsSafeReadOnly 安全只读操作
// 适用于: Read, Glob, Grep 等只读工具
var AnnotationsSafeReadOnly = &ToolAnnotations{
	ReadOnly:    true,
	Destructive: false,
	Idempotent:  true,
	OpenWorld:   false,
	RiskLevel:   RiskLevelSafe,
	Category:    CategoryFilesystem,
}

// AnnotationsSafeWrite 安全写操作
// 适用于: Write, Edit 等文件写入工具
var AnnotationsSafeWrite = &ToolAnnotations{
	ReadOnly:    false,
	Destructive: false,
	Idempotent:  true, // 覆盖写是幂等的
	OpenWorld:   false,
	RiskLevel:   RiskLevelLow,
	Category:    CategoryFilesystem,
}

// AnnotationsDestructiveWrite 破坏性写操作
// 适用于: 删除文件、清空目录等操作
var AnnotationsDestructiveWrite = &ToolAnnotations{
	ReadOnly:             false,
	Destructive:          true,
	Idempotent:           false,
	OpenWorld:            false,
	RiskLevel:            RiskLevelHigh,
	Category:             CategoryFilesystem,
	RequiresConfirmation: true,
}

// AnnotationsExecution 命令执行
// 适用于: Bash 等命令执行工具
var AnnotationsExecution = &ToolAnnotations{
	ReadOnly:    false,
	Destructive: true, // 可能执行任何命令
	Idempotent:  false,
	OpenWorld:   true, // 可能访问网络
	RiskLevel:   RiskLevelHigh,
	Category:    CategoryExecution,
}

// AnnotationsNetworkRead 网络只读操作
// 适用于: WebFetch, WebSearch 等网络读取工具
var AnnotationsNetworkRead = &ToolAnnotations{
	ReadOnly:    true,
	Destructive: false,
	Idempotent:  true,
	OpenWorld:   true,
	RiskLevel:   RiskLevelLow,
	Category:    CategoryNetwork,
}

// AnnotationsNetworkWrite 网络写操作
// 适用于: HTTP POST/PUT/DELETE 等修改操作
var AnnotationsNetworkWrite = &ToolAnnotations{
	ReadOnly:    false,
	Destructive: false,
	Idempotent:  false,
	OpenWorld:   true,
	RiskLevel:   RiskLevelMedium,
	Category:    CategoryNetwork,
}

// AnnotationsDatabaseRead 数据库只读操作
// 适用于: SELECT 查询等
var AnnotationsDatabaseRead = &ToolAnnotations{
	ReadOnly:    true,
	Destructive: false,
	Idempotent:  true,
	OpenWorld:   false,
	RiskLevel:   RiskLevelSafe,
	Category:    CategoryDatabase,
}

// AnnotationsDatabaseWrite 数据库写操作
// 适用于: INSERT, UPDATE 等
var AnnotationsDatabaseWrite = &ToolAnnotations{
	ReadOnly:    false,
	Destructive: false,
	Idempotent:  false,
	OpenWorld:   false,
	RiskLevel:   RiskLevelMedium,
	Category:    CategoryDatabase,
}

// AnnotationsDatabaseDestructive 数据库破坏性操作
// 适用于: DELETE, DROP 等
var AnnotationsDatabaseDestructive = &ToolAnnotations{
	ReadOnly:             false,
	Destructive:          true,
	Idempotent:           false,
	OpenWorld:            false,
	RiskLevel:            RiskLevelCritical,
	Category:             CategoryDatabase,
	RequiresConfirmation: true,
}

// AnnotationsMCPTool MCP 工具默认注解
var AnnotationsMCPTool = &ToolAnnotations{
	ReadOnly:    false,
	Destructive: false,
	Idempotent:  false,
	OpenWorld:   true, // MCP 工具通常涉及外部系统
	RiskLevel:   RiskLevelMedium,
	Category:    CategoryMCP,
}

// AnnotationsUserInteraction 用户交互工具
// 适用于: AskUserQuestion 等
var AnnotationsUserInteraction = &ToolAnnotations{
	ReadOnly:    true,
	Destructive: false,
	Idempotent:  true,
	OpenWorld:   false,
	RiskLevel:   RiskLevelSafe,
	Category:    CategoryCustom,
}

// Clone 克隆注解（用于修改）
func (a *ToolAnnotations) Clone() *ToolAnnotations {
	if a == nil {
		return nil
	}
	return &ToolAnnotations{
		ReadOnly:             a.ReadOnly,
		Destructive:          a.Destructive,
		Idempotent:           a.Idempotent,
		OpenWorld:            a.OpenWorld,
		RiskLevel:            a.RiskLevel,
		Category:             a.Category,
		RequiresConfirmation: a.RequiresConfirmation,
	}
}

// WithCategory 设置分类（链式调用）
func (a *ToolAnnotations) WithCategory(category string) *ToolAnnotations {
	a.Category = category
	return a
}

// WithRiskLevel 设置风险级别（链式调用）
func (a *ToolAnnotations) WithRiskLevel(level int) *ToolAnnotations {
	a.RiskLevel = level
	return a
}

// WithRequiresConfirmation 设置是否需要确认（链式调用）
func (a *ToolAnnotations) WithRequiresConfirmation(requires bool) *ToolAnnotations {
	a.RequiresConfirmation = requires
	return a
}

// IsSafeForAutoApproval 判断是否可以自动批准
// 在 SmartApprove 模式下，只读且不涉及外部系统的工具可以自动批准
func (a *ToolAnnotations) IsSafeForAutoApproval() bool {
	if a == nil {
		return false
	}
	return a.ReadOnly && !a.OpenWorld && !a.RequiresConfirmation
}

// RiskLevelName 获取风险级别名称
func (a *ToolAnnotations) RiskLevelName() string {
	if a == nil {
		return "unknown"
	}
	switch a.RiskLevel {
	case RiskLevelSafe:
		return "safe"
	case RiskLevelLow:
		return "low"
	case RiskLevelMedium:
		return "medium"
	case RiskLevelHigh:
		return "high"
	case RiskLevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}
