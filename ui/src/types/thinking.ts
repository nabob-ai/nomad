/**
 * 思维过程相关类型定义
 * 参考 yunjin 的 types/writing.ts
 */

// 工具调用状态
export type ToolCallState =
  | 'pending'     // 待执行
  | 'queued'      // 已排队
  | 'executing'   // 执行中
  | 'completed'   // 已完成
  | 'failed'      // 失败
  | 'cancelling'  // 取消中
  | 'cancelled'   // 已取消

// Think Aloud 事件类型
export type ThinkAloudEventType =
  | 'reasoning'           // 推理思考
  | 'tool_call'          // 工具调用开始
  | 'tool_result'        // 工具调用结果
  | 'approval_required'  // 需要审批
  | 'session_summarized' // 会话历史已汇总

// 思考阶段类型
export type ThinkingStage =
  | '任务规划'      // 收到请求后分析和规划
  | '推理分析'      // 模型推理过程
  | '工具规划'      // 决定调用工具
  | '工具执行'      // 执行工具中
  | '结果总结'      // 汇总结果
  | '人工介入'      // HITL 审批

// 工具调用数据
export interface ToolCallData {
  id: string
  name: string
  input: Record<string, any>
  state: ToolCallState
  progress?: number
  progressMessage?: string
}

// 工具结果数据
export interface ToolResultData {
  id: string
  name: string
  result: any
  error?: string
  durationMs?: number
}

// 审批请求数据
export interface ApprovalRequestData {
  id: string
  toolName: string
  input: Record<string, any>
  reason: string
}

// Think Aloud 事件
export interface ThinkAloudEvent {
  id: string
  type: ThinkAloudEventType
  timestamp: string
  stage?: ThinkingStage | string
  reasoning?: string
  decision?: string
  toolCall?: ToolCallData
  toolResult?: ToolResultData
  approvalRequest?: ApprovalRequestData
}

// 工具调用审批信息
export interface ToolCallApproval {
  callId: string
  required: boolean
  approved: boolean
  reason?: string
  timestamp: string
  approvedBy?: string
}

// 工具调用审计条目
export interface ToolCallAuditEntry {
  state: ToolCallState
  timestamp: string
  note: string
}

// 工具调用记录（完整版）
export interface ToolCallRecord {
  id: string
  name: string
  input: Record<string, any>
  result?: any
  error?: string
  isError: boolean
  progress: number
  progressMessage?: string
  intermediate?: Record<string, any>
  startedAt?: string
  completedAt?: string
  durationMs?: number
  state: ToolCallState
  approval: ToolCallApproval
  createdAt: string
  updatedAt: string
  auditTrail: ToolCallAuditEntry[]
}

// 旧版类型（保持向后兼容）
export type ThinkingStepType = 'reasoning' | 'tool_call' | 'tool_result' | 'decision' | 'approval' | 'session_summarized'

// 会话摘要数据
export interface SessionSummarizedData {
  messagesBefore: number
  messagesAfter: number
  tokensBefore: number
  tokensAfter: number
  tokensSaved: number
  compressionRatio: number
  summaryPreview: string
}

export interface ThinkingStep {
  id?: string
  type: ThinkingStepType
  content?: string
  tool?: {
    name: string
    args: any
  }
  result?: any
  timestamp: number
  messageId?: string
  // 会话摘要相关数据
  sessionSummarized?: SessionSummarizedData
}

export interface ThinkingState {
  stepsByMessage: Map<string, ThinkingStep[]>
  currentThought: string
  currentMessageId: string | null
  isThinking: boolean
}

// 类型守卫
export function isValidThinkAloudEvent(event: any): event is ThinkAloudEvent {
  return (
    event &&
    typeof event === 'object' &&
    typeof event.type === 'string' &&
    ['reasoning', 'tool_call', 'tool_result', 'approval_required'].includes(event.type)
  )
}
