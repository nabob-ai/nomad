<script setup lang="ts">
/**
 * AgentLoopDemo - 演示完整 Agent Loop + HITL 集成
 *
 * 功能:
 * - 重试逻辑 (通过后端 ModelFallbackManager)
 * - Human-in-the-Loop 审批流程
 * - 真实工具执行
 * - 流式响应
 */

import { ref } from "vue";
import { useAgentLoop } from "@/composables/useAgentLoop";
import type { ThinkAloudEvent } from "@/composables/useAgentLoop";
import { ThinkingBlock } from "@/components/Thinking";
import { StreamingText } from "@/components/Common";

// Props
const props = defineProps<{
  modelConfig?: {
    provider?: string;
    model?: string;
  };
}>();

// 思考事件列表（转换为新格式）
import type { ThinkAloudEvent as NewThinkAloudEvent } from "@/types/thinking";
const thinkEvents = ref<NewThinkAloudEvent[]>([]);

// Agent Loop
const {
  isRunning,
  isPaused,
  currentOutput,
  pendingApproval,
  isConnected,
  execute,
  approveAndResume,
  rejectTool,
  cancel
} = useAgentLoop({
  modelConfig: props.modelConfig,
  sensitiveTools: ["Edit", "Write", "bash", "fs_write"],
  maxRetries: 3,
  maxLoops: 10,
  onThink: (event) => {
    // 转换为新的 ThinkAloudEvent 格式
    const newEvent: NewThinkAloudEvent = {
      id: event.id,
      timestamp: event.timestamp,
      stage: event.stage as any,
      reasoning: event.reasoning,
      decision: event.decision,
      type: 'reasoning' as const,
    };

    if (event.toolCall) {
      newEvent.type = 'tool_call';
      newEvent.toolCall = {
        id: event.id,
        name: event.toolCall.toolName,
        input: event.toolCall.args,
        state: 'executing',
        progress: 0
      };
    } else if (event.toolResult) {
      newEvent.type = 'tool_result';
      newEvent.toolResult = {
        id: event.id,
        name: event.toolResult.toolName,
        result: event.toolResult.result
      };
    } else if (event.approvalRequest) {
      newEvent.type = 'approval_required';
      newEvent.approvalRequest = {
        id: event.approvalRequest.id,
        toolName: event.approvalRequest.toolName,
        input: event.approvalRequest.args,
        reason: '敏感操作需要人工审批'
      };
    }

    thinkEvents.value = [...thinkEvents.value, newEvent];
  },
  onApprovalRequired: (request) => {
    console.log("Approval required:", request);
  },
  onToolStart: (toolName, args) => {
    console.log("Tool started:", toolName, args);
  },
  onToolEnd: (toolName, result) => {
    console.log("Tool ended:", toolName, result);
  },
  onTextDelta: () => {
    // 已通过 currentOutput 响应式更新
  },
  onComplete: (result) => {
    console.log("Execution complete:", result.status);
  },
  onError: (error) => {
    console.error("Execution error:", error);
  },
});

// 用户输入
const userInput = ref("");

// 发送消息
const sendMessage = async () => {
  if (!userInput.value.trim() || isRunning.value) return;

  thinkEvents.value = [];
  const input = userInput.value;
  userInput.value = "";

  await execute(input);
};

// 批准工具
const handleApprove = async (request: any) => {
  if (!request) return;
  await approveAndResume(request.id);
};

// 拒绝工具
const handleReject = (request: any) => {
  if (!request) return;
  rejectTool(request.id, "用户拒绝");
};

// 取消执行
const handleCancel = () => {
  cancel();
};

// 计算是否完成
const isFinished = ref(false);

// 监听运行状态
import { watch } from 'vue';
watch(isRunning, (running) => {
  if (!running && !isPaused.value) {
    isFinished.value = true;
  } else {
    isFinished.value = false;
  }
});
</script>

<template>
  <div class="agent-loop-demo">
    <!-- 连接状态 -->
    <div class="connection-status" :class="{ connected: isConnected }">
      <span class="status-dot"></span>
      {{ isConnected ? "已连接" : "未连接" }}
    </div>

    <!-- 思考过程 - 使用新的 ThinkingBlock 组件 -->
    <ThinkingBlock
      v-if="thinkEvents.length > 0"
      :thoughts="thinkEvents"
      :is-finished="isFinished"
      :pending-approval="pendingApproval as any"
      :is-approving="isRunning && isPaused"
      @approve="handleApprove"
      @reject="handleReject"
    />

    <!-- 输出面板 -->
    <div class="output-panel">
      <h3>📝 输出</h3>
      <div class="output-content">
        <StreamingText
          v-if="currentOutput"
          :content="currentOutput"
          :is-streaming="isRunning && !isPaused"
        />
        <em v-else class="placeholder">等待输出...</em>
      </div>
    </div>

    <!-- 输入面板 -->
    <div class="input-panel">
      <textarea
        v-model="userInput"
        placeholder="输入你的请求..."
        :disabled="isRunning"
        @keydown.enter.ctrl="sendMessage"
        rows="3"
      ></textarea>
      <div class="input-actions">
        <button
          class="btn btn-primary"
          @click="sendMessage"
          :disabled="!userInput.trim() || isRunning"
        >
          {{ isRunning ? (isPaused ? "等待审批..." : "执行中...") : "发送" }}
        </button>
        <button
          class="btn btn-secondary"
          @click="handleCancel"
          :disabled="!isRunning"
        >
          取消
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.agent-loop-demo {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
  max-width: 800px;
  margin: 0 auto;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}

.connection-status.connected {
  color: #22c55e;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ef4444;
}

.connection-status.connected .status-dot {
  background: #22c55e;
}

.output-panel {
  background: #f8fafc;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid #e2e8f0;
}

.output-panel h3 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #475569;
}

.output-content {
  min-height: 100px;
  line-height: 1.6;
  color: #1e293b;
}

.output-content .placeholder {
  color: #94a3b8;
}

.input-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.input-panel textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.input-panel textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.input-panel textarea:disabled {
  background: #f8fafc;
  color: #94a3b8;
}

.input-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.btn {
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-primary:active:not(:disabled) {
  transform: scale(0.98);
}

.btn-secondary {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.btn-secondary:hover:not(:disabled) {
  background: #e2e8f0;
}
</style>
