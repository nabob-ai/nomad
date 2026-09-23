<template>
  <div :class="['message-bubble', `message-${message.role}`, message.status === 'error' && 'message-error']">
    <!-- Avatar -->
    <div v-if="showAvatar && message.role === 'assistant'" class="message-avatar">
      <div class="avatar-icon">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
          ></path>
        </svg>
      </div>
    </div>

    <!-- Content -->
    <div class="message-content-wrapper">
      <!-- Bubble -->
      <div :class="['message-bubble-content', message.role === 'user' ? 'user-bubble' : 'assistant-bubble']">
        <!-- Text Content -->
        <div v-if="message.type === 'text'" class="message-text" v-html="renderedContent"></div>

        <!-- Image Content -->
        <div v-else-if="message.type === 'image'" class="message-image">
          <img :src="message.content.url" :alt="message.content.alt" class="max-w-full rounded" loading="lazy" />
        </div>

        <!-- System Message -->
        <div v-else-if="message.type === 'system'" class="message-system">
          {{ message.content.text }}
        </div>

        <!-- Status Indicator -->
        <div v-if="message.status" class="message-status">
          <svg v-if="message.status === 'pending'" class="w-3 h-3 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
          <svg v-else-if="message.status === 'sent'" class="w-3 h-3 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
          <svg v-else-if="message.status === 'error'" class="w-3 h-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </div>
      </div>

      <!-- ThinkingBlock (仅显示助手消息) -->
      <ThinkingBlock
        v-if="message.role === 'assistant' && thinkingSteps.length > 0"
        :thoughts="convertedThoughts"
        :is-finished="!isThinking"
        @approve="handleApprove"
        @reject="handleReject"
      />

      <!-- WorkflowProgressView (仅当有激活工作流时显示) -->
      <WorkflowProgressView
        v-if="showWorkflow"
        :steps="workflowSteps"
        title="工作流进度"
        :show-progress="true"
        :show-steps="true"
        :max-visible-steps="3"
      />

      <!-- Timestamp -->
      <div v-if="showTimestamp" class="message-timestamp">
        {{ formatTime(message.createdAt) }}
      </div>

      <!-- Actions (on hover) -->
      <div v-if="showActions" class="message-actions">
        <button @click="$emit('copy', message)" class="action-button" title="复制">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
          </svg>
        </button>
        <button v-if="message.status === 'error'" @click="$emit('retry', message)" class="action-button" title="重试">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
        </button>
        <button @click="$emit('delete', message)" class="action-button" title="删除">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { renderMarkdown } from "@/utils/markdown";
import { formatTime } from "@/utils/format";
import { useThinkingStore } from "@/stores/thinking";
import { useWorkflowStore } from "@/stores/workflow";
import { ThinkingBlock } from "@/components/Thinking";
import { WorkflowProgressView } from "@/components/Workflow";
import type { Message, TextMessage } from "@/types";
import type { ThinkAloudEvent } from "@/types/thinking";

const props = defineProps<{
  message: Message;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  showActions?: boolean;
}>();

const emit = defineEmits<{
  copy: [message: Message];
  retry: [message: Message];
  delete: [message: Message];
  approve: [request: any];
  reject: [request: any];
}>();

const thinkingStore = useThinkingStore();
const workflowStore = useWorkflowStore();

const renderedContent = computed(() => {
  if (props.message.type === "text") {
    const textMessage = props.message as TextMessage;
    return renderMarkdown(textMessage.content.text);
  }
  return "";
});

// 获取该消息的思维步骤
const thinkingSteps = computed(() => {
  return thinkingStore.getSteps(props.message.id);
});

// 转换为新的 ThinkAloudEvent 格式
const convertedThoughts = computed((): ThinkAloudEvent[] => {
  return thinkingSteps.value.map((step, index) => ({
    id: step.id || `step-${index}`,
    type: step.type === 'decision' ? 'reasoning' : step.type === 'approval' ? 'approval_required' : step.type,
    timestamp: new Date(step.timestamp).toISOString(),
    reasoning: step.content,
    toolCall: step.tool ? {
      id: step.id || `tool-${index}`,
      name: step.tool.name,
      input: step.tool.args,
      state: 'completed' as const,
      progress: 1
    } : undefined,
    toolResult: step.result ? {
      id: step.id || `result-${index}`,
      name: step.tool?.name || 'unknown',
      result: step.result
    } : undefined
  }));
});

// 判断是否正在思考
const isThinking = computed(() => {
  return thinkingStore.isThinking && thinkingStore.currentMessageId === props.message.id;
});

// 是否显示工作流进度
const showWorkflow = computed(() => {
  return props.message.role === "assistant" && workflowStore.hasActiveWorkflow;
});

// 工作流步骤
const workflowSteps = computed(() => {
  return workflowStore.steps;
});

const handleApprove = (request: any) => {
  emit('approve', request);
};

const handleReject = (request: any) => {
  emit('reject', request);
};
</script>

<style scoped>
.message-bubble {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-user {
  flex-direction: row-reverse;
}

.message-assistant {
  flex-direction: row;
}

.message-avatar {
  flex-shrink: 0;
}

.avatar-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e0f2fe 0%, #bae6fd 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0284c7;
}

.message-content-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-width: 80%;
}

.message-user .message-content-wrapper {
  align-items: flex-end;
}

.message-assistant .message-content-wrapper {
  align-items: flex-start;
}

.message-bubble-content {
  border-radius: 12px;
  padding: 12px 16px;
  position: relative;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.user-bubble {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border-bottom-right-radius: 4px;
}

.assistant-bubble {
  background: #f8fafc;
  color: #1e293b;
  border: 1px solid #e2e8f0;
  border-bottom-left-radius: 4px;
}

.message-text {
  font-size: 14px;
  line-height: 1.6;
}

.message-text :deep(p) {
  margin-bottom: 8px;
}

.message-text :deep(p:last-child) {
  margin-bottom: 0;
}

.message-text :deep(code) {
  background: rgba(0, 0, 0, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
  font-family: ui-monospace, monospace;
}

.user-bubble .message-text :deep(code) {
  background: rgba(255, 255, 255, 0.2);
}

.message-text :deep(pre) {
  background: #1e293b;
  color: #f1f5f9;
  padding: 12px;
  border-radius: 8px;
  margin: 8px 0;
  overflow-x: auto;
  font-size: 13px;
}

.message-text :deep(pre code) {
  background: transparent;
  padding: 0;
}

.message-text :deep(a) {
  color: inherit;
  text-decoration: underline;
}

.message-text :deep(a:hover) {
  text-decoration: none;
}

.message-image img {
  max-height: 256px;
  object-fit: contain;
  border-radius: 8px;
}

.message-system {
  font-size: 12px;
  text-align: center;
  color: #64748b;
  font-style: italic;
}

.message-status {
  position: absolute;
  bottom: 4px;
  right: 4px;
  opacity: 0.7;
}

.message-timestamp {
  font-size: 11px;
  color: #94a3b8;
  padding: 0 8px;
}

.message-actions {
  display: none;
  gap: 4px;
  position: absolute;
  top: -32px;
  right: 0;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  padding: 4px;
}

.message-bubble:hover .message-actions {
  display: flex;
}

.action-button {
  padding: 6px;
  color: #64748b;
  border-radius: 6px;
  transition: all 0.15s;
  background: transparent;
  border: none;
  cursor: pointer;
}

.action-button:hover {
  background: #f1f5f9;
  color: #3b82f6;
}

.message-error .message-bubble-content {
  border-color: #fecaca;
  background: #fef2f2;
}
</style>
