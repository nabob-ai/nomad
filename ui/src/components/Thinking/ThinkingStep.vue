<template>
  <div class="thinking-step">
    <!-- 步骤指示器 -->
    <div class="step-indicator">
      <div :class="['step-dot', stepDotClass]">
        <!-- 推理图标 -->
        <svg v-if="step.type === 'reasoning'" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
          />
        </svg>

        <!-- 工具调用图标 -->
        <svg v-else-if="step.type === 'tool_call'" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
        </svg>

        <!-- 工具结果图标 -->
        <svg v-else-if="step.type === 'tool_result'" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>

        <!-- 决策图标 -->
        <svg v-else-if="step.type === 'decision'" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>

        <!-- 审批图标 -->
        <svg v-else-if="step.type === 'approval'" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>

        <!-- 会话摘要图标 -->
        <svg v-else-if="step.type === 'session_summarized'" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>

        <!-- 默认图标 -->
        <svg v-else class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>

      <!-- 连接线 -->
      <div v-if="!isLast" class="step-line"></div>
    </div>

    <!-- 步骤内容 -->
    <div class="step-content">
      <div class="step-header">
        <span :class="['step-type', stepTypeClass]">
          {{ stepLabel }}
        </span>
        <span class="step-time">{{ formatTime(step.timestamp) }}</span>
      </div>

      <!-- 推理内容 -->
      <div v-if="step.type === 'reasoning' && step.content" class="step-body">
        <p class="text-sm text-slate-700 dark:text-slate-300 leading-relaxed">
          {{ step.content }}
        </p>
      </div>

      <!-- 决策内容 -->
      <div v-if="step.type === 'decision' && step.content" class="step-body">
        <div class="decision-content">
          <svg class="w-4 h-4 text-purple-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          <p class="text-sm text-slate-800 dark:text-slate-200 font-medium">
            {{ step.content }}
          </p>
        </div>
      </div>

      <!-- 工具调用 -->
      <div v-if="step.type === 'tool_call' && step.tool" class="step-body">
        <div class="tool-call">
          <div class="tool-header">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
            </svg>
            <span class="font-mono font-semibold">{{ step.tool.name }}</span>
          </div>
          <pre class="tool-args">{{ formatToolArgs(step.tool.args) }}</pre>
        </div>
      </div>

      <!-- 工具结果 -->
      <div v-if="step.type === 'tool_result' && step.result" class="step-body">
        <div class="tool-result">
          <div class="result-header">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span class="font-semibold">执行结果</span>
          </div>
          <pre class="result-content">{{ formatResult(step.result) }}</pre>
        </div>
      </div>

      <!-- 审批步骤 -->
      <div v-if="step.type === 'approval' && step.tool" class="step-body">
        <div class="approval-step">
          <div class="approval-step-header">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <span class="font-semibold">等待审批</span>
          </div>
          <div class="approval-step-content">
            <span class="font-mono text-sm">{{ step.tool.name }}</span>
          </div>
        </div>
      </div>

      <!-- 会话摘要步骤 -->
      <div v-if="step.type === 'session_summarized'" class="step-body">
        <div class="session-summarized">
          <div class="session-summarized-header">
            <svg class="w-4 h-4 text-cyan-600 dark:text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            <span class="font-semibold text-cyan-700 dark:text-cyan-300">已汇总会话历史记录</span>
          </div>
          <div v-if="step.sessionSummarized" class="session-summarized-stats">
            <div class="stat-item">
              <span class="stat-label">消息</span>
              <span class="stat-value">{{ step.sessionSummarized.messagesBefore }} → {{ step.sessionSummarized.messagesAfter }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">节省</span>
              <span class="stat-value text-emerald-600 dark:text-emerald-400">{{ formatTokens(step.sessionSummarized.tokensSaved) }} tokens</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">压缩率</span>
              <span class="stat-value">{{ formatPercent(step.sessionSummarized.compressionRatio) }}</span>
            </div>
          </div>
          <div v-if="step.sessionSummarized?.summaryPreview" class="session-summarized-preview">
            <p class="text-xs text-slate-500 dark:text-slate-400 line-clamp-2">
              {{ step.sessionSummarized.summaryPreview }}...
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from "vue";
import type { ThinkingStep } from "@/types/thinking";
import type { PropType } from "vue";

export default defineComponent({
  name: "ThinkingStep",
  props: {
    step: {
      type: Object as PropType<ThinkingStep>,
      required: true,
    },
    isLast: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const stepDotClass = computed(() => {
      const classes: Record<string, string> = {
        reasoning: "step-dot-reasoning",
        tool_call: "step-dot-tool-call",
        tool_result: "step-dot-tool-result",
        decision: "step-dot-decision",
        approval: "step-dot-approval",
        session_summarized: "step-dot-session-summarized",
      };
      return classes[props.step.type] || "";
    });

    const stepTypeClass = computed(() => {
      const classes: Record<string, string> = {
        reasoning: "step-type-reasoning",
        tool_call: "step-type-tool-call",
        tool_result: "step-type-tool-result",
        decision: "step-type-decision",
        approval: "step-type-approval",
        session_summarized: "step-type-session-summarized",
      };
      return classes[props.step.type] || "";
    });

    const stepLabel = computed(() => {
      const labels: Record<string, string> = {
        reasoning: "推理",
        tool_call: "工具调用",
        tool_result: "执行结果",
        decision: "决策",
        approval: "审批请求",
        session_summarized: "历史汇总",
      };
      return labels[props.step.type] || props.step.type;
    });

    const formatTime = (timestamp: number): string => {
      const date = new Date(timestamp);
      return date.toLocaleTimeString("zh-CN", {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      });
    };

    const formatToolArgs = (args: any): string => {
      try {
        return JSON.stringify(args, null, 2);
      } catch (e) {
        return String(args);
      }
    };

    const formatResult = (result: any): string => {
      try {
        if (typeof result === "string") return result;
        return JSON.stringify(result, null, 2);
      } catch (e) {
        return String(result);
      }
    };

    const formatTokens = (tokens: number): string => {
      return tokens.toLocaleString();
    };

    const formatPercent = (ratio: number): string => {
      return `${Math.round(ratio * 100)}%`;
    };

    return {
      stepDotClass,
      stepTypeClass,
      stepLabel,
      formatTime,
      formatToolArgs,
      formatResult,
      formatTokens,
      formatPercent,
    };
  },
});
</script>

<style scoped>
.thinking-step {
  @apply flex gap-3;
}

.step-indicator {
  @apply flex flex-col items-center;
}

.step-dot {
  @apply w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0;
}

.step-dot-reasoning {
  @apply bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400;
}

.step-dot-tool-call {
  @apply bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400;
}

.step-dot-tool-result {
  @apply bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400;
}

.step-dot-decision {
  @apply bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400;
}

.step-dot-approval {
  @apply bg-amber-100 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400;
}

.step-dot-session-summarized {
  @apply bg-cyan-100 dark:bg-cyan-900/30 text-cyan-600 dark:text-cyan-400;
}

.step-line {
  @apply w-0.5 flex-1 bg-slate-200 dark:bg-slate-700 mt-1;
}

.step-content {
  @apply flex-1 pb-4;
}

.step-header {
  @apply flex items-center justify-between mb-2;
}

.step-type {
  @apply text-xs font-semibold px-2 py-1 rounded;
}

.step-type-reasoning {
  @apply bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300;
}

.step-type-tool-call {
  @apply bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300;
}

.step-type-tool-result {
  @apply bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300;
}

.step-type-decision {
  @apply bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300;
}

.step-type-approval {
  @apply bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300;
}

.step-type-session-summarized {
  @apply bg-cyan-100 dark:bg-cyan-900/30 text-cyan-700 dark:text-cyan-300;
}

.step-time {
  @apply text-xs text-slate-500 dark:text-slate-400 tabular-nums;
}

.step-body {
  @apply text-sm;
}

.decision-content {
  @apply flex gap-2 items-start p-2 bg-purple-50 dark:bg-purple-900/10 rounded border border-purple-200 dark:border-purple-800;
}

.tool-call,
.tool-result {
  @apply rounded-lg overflow-hidden border;
}

.tool-call {
  @apply border-indigo-200 dark:border-indigo-800;
}

.tool-result {
  @apply border-emerald-200 dark:border-emerald-800;
}

.tool-header,
.result-header {
  @apply flex items-center gap-2 px-3 py-2 text-sm font-medium;
}

.tool-header {
  @apply bg-indigo-50 dark:bg-indigo-900/20 text-indigo-700 dark:text-indigo-300;
}

.result-header {
  @apply bg-emerald-50 dark:bg-emerald-900/20 text-emerald-700 dark:text-emerald-300;
}

.tool-args,
.result-content {
  @apply p-3 text-xs font-mono bg-slate-50 dark:bg-slate-900 text-slate-700 dark:text-slate-300 overflow-x-auto;
}

.approval-step {
  @apply rounded-lg overflow-hidden border border-amber-200 dark:border-amber-800;
}

.approval-step-header {
  @apply flex items-center gap-2 px-3 py-2 text-sm font-medium bg-amber-50 dark:bg-amber-900/20 text-amber-700 dark:text-amber-300;
}

.approval-step-content {
  @apply p-3 bg-white dark:bg-slate-800;
}

/* 会话摘要样式 */
.session-summarized {
  @apply rounded-lg overflow-hidden border border-cyan-200 dark:border-cyan-800 bg-gradient-to-r from-cyan-50 to-blue-50 dark:from-cyan-900/20 dark:to-blue-900/20;
}

.session-summarized-header {
  @apply flex items-center gap-2 px-3 py-2 text-sm font-medium;
}

.session-summarized-stats {
  @apply flex flex-wrap gap-4 px-3 py-2 border-t border-cyan-100 dark:border-cyan-800/50;
}

.stat-item {
  @apply flex flex-col;
}

.stat-label {
  @apply text-xs text-slate-500 dark:text-slate-400;
}

.stat-value {
  @apply text-sm font-medium text-slate-700 dark:text-slate-200;
}

.session-summarized-preview {
  @apply px-3 py-2 border-t border-cyan-100 dark:border-cyan-800/50 bg-white/50 dark:bg-slate-900/50;
}
</style>
