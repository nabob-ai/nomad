<template>
  <div class="tool-call-display">
    <!-- 工具调用说明 -->
    <p v-if="reasoning" class="tool-reasoning">
      {{ reasoning }}
    </p>

    <!-- 工具调用卡片 -->
    <div v-if="isToolCall && toolCall" class="tool-card tool-call-card">
      <div class="tool-header">
        <div class="tool-title">
          <svg class="tool-icon spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
          <span>工具调用</span>
        </div>
        <span class="tool-name">{{ toolCall.name }}</span>
      </div>
      <div class="tool-content">
        <pre class="tool-json">{{ formatJSON(toolCall.input) }}</pre>
      </div>
      <!-- 执行中状态指示器 -->
      <div v-if="toolCall.state === 'executing'" class="tool-executing">
        <span class="executing-dot"></span>
        <span class="executing-text">{{ toolCall.progressMessage || '执行中...' }}</span>
      </div>
      <!-- 进度条 -->
      <div v-if="toolCall.progress && toolCall.progress > 0" class="tool-progress">
        <div class="progress-bar" :style="{ width: `${toolCall.progress * 100}%` }"></div>
        <span v-if="toolCall.progressMessage" class="progress-message">{{ toolCall.progressMessage }}</span>
      </div>
    </div>

    <!-- 工具结果卡片 -->
    <div v-else-if="isToolResult && toolResult" :class="['tool-card', toolResult.error ? 'tool-error-card' : 'tool-result-card']">
      <div class="tool-header" :class="{ error: toolResult.error }">
        <div class="tool-title">
          <svg v-if="toolResult.error" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <span>{{ toolResult.error ? '执行失败' : '执行结果' }}</span>
        </div>
        <span v-if="toolResult.durationMs" class="tool-duration">
          {{ toolResult.durationMs }}ms
        </span>
      </div>
      <div class="tool-content">
        <pre :class="['tool-json', { error: toolResult.error }]">{{ toolResult.error || formatJSON(toolResult.result) }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ThinkAloudEvent } from '@/types/thinking'

const props = defineProps<{
  event: ThinkAloudEvent
}>()

const isToolCall = computed(() => props.event.type === 'tool_call')
const isToolResult = computed(() => props.event.type === 'tool_result')
const toolCall = computed(() => props.event.toolCall)
const toolResult = computed(() => props.event.toolResult)
const reasoning = computed(() => props.event.reasoning)

const formatJSON = (obj: any): string => {
  if (!obj) return ''
  try {
    return JSON.stringify(obj, null, 2)
  } catch {
    return String(obj)
  }
}
</script>

<style scoped>
.tool-call-display {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
  width: 100%;
}

.tool-reasoning {
  font-size: 10px;
  color: #64748b;
  font-style: italic;
  margin: 0;
  padding-bottom: 4px;
  border-bottom: 1px solid #f1f5f9;
}

.tool-card {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  min-width: 0;
  width: 100%;
}

.tool-call-card {
  background: #1e293b;
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.tool-result-card {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
}

.tool-error-card {
  background: #fef2f2;
  border: 1px solid #fecaca;
}

.tool-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  gap: 8px;
  flex-wrap: nowrap;
  min-width: 0;
}

.tool-call-card .tool-header {
  background: #334155;
}

.tool-result-card .tool-header {
  background: #dcfce7;
  border-bottom-color: #bbf7d0;
}

.tool-error-card .tool-header {
  background: #fee2e2;
  border-bottom-color: #fecaca;
}

.tool-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  flex-shrink: 0;
}

.tool-call-card .tool-title {
  color: #a5b4fc;
}

.tool-result-card .tool-title {
  color: #047857;
}

.tool-error-card .tool-title {
  color: #b91c1c;
}

.tool-name {
  font-size: 10px;
  font-family: ui-monospace, monospace;
  padding: 2px 8px;
  border-radius: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}

.tool-call-card .tool-name {
  background: rgba(99, 102, 241, 0.2);
  color: #c7d2fe;
  border: 1px solid rgba(99, 102, 241, 0.3);
}

.tool-duration {
  font-size: 10px;
  color: #10b981;
  font-weight: 400;
}

.tool-content {
  padding: 12px;
  overflow-x: auto;
}

.tool-call-card .tool-content {
  background: #1e293b;
}

.tool-result-card .tool-content,
.tool-error-card .tool-content {
  background: white;
}

.tool-json {
  font-family: ui-monospace, monospace;
  font-size: 10px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  max-height: 128px;
  overflow-y: auto;
}

.tool-call-card .tool-json {
  color: #c7d2fe;
}

.tool-result-card .tool-json {
  color: #065f46;
}

.tool-json.error {
  color: #b91c1c;
}

.tool-progress {
  padding: 0 12px 12px;
  background: #1e293b;
}

.progress-bar {
  height: 3px;
  background: linear-gradient(90deg, #6366f1, #8b5cf6);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.progress-message {
  display: block;
  margin-top: 4px;
  font-size: 10px;
  color: #94a3b8;
}

.tool-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}

.tool-icon.spin {
  animation: spin 2s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.tool-executing {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(99, 102, 241, 0.1);
  border-top: 1px solid rgba(99, 102, 241, 0.2);
}

.executing-dot {
  width: 8px;
  height: 8px;
  background: #6366f1;
  border-radius: 50%;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.8); }
}

.executing-text {
  font-size: 11px;
  color: #a5b4fc;
  font-weight: 500;
}
</style>
