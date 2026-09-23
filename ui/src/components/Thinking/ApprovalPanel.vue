<template>
  <div class="approval-panel">
    <div class="approval-header">
      <div class="approval-icon">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
      </div>
      <div class="approval-info">
        <h4 class="approval-title">需要人工审批</h4>
        <p class="approval-description">
          Agent 请求执行敏感操作
          <span class="tool-name-badge">{{ pendingApproval.name }}</span>
          。请确认是否允许。
        </p>
      </div>
    </div>

    <!-- 参数详情 -->
    <div class="params-section">
      <div class="params-header">
        <span>参数详情</span>
        <span class="params-tool-name">{{ pendingApproval.name }}</span>
      </div>
      <pre class="params-content">{{ formatJSON(pendingApproval.input) }}</pre>
    </div>

    <!-- 操作按钮 -->
    <div class="action-buttons">
      <button
        class="btn btn-approve"
        :disabled="isLoading"
        @click="handleApprove"
      >
        <svg v-if="!isLoading" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
        <span v-else class="loading-spinner"></span>
        {{ isLoading ? '处理中...' : '批准执行' }}
      </button>
      <button
        class="btn btn-reject"
        :disabled="isLoading"
        @click="handleReject"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
        拒绝
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ToolCallRecord } from '@/types/thinking'

const props = defineProps<{
  pendingApproval: ToolCallRecord
  isLoading?: boolean
}>()

const emit = defineEmits<{
  approve: [request: ToolCallRecord]
  reject: [request: ToolCallRecord]
}>()

// 格式化 JSON
const formatJSON = (obj: any): string => {
  if (!obj) return ''
  try {
    return JSON.stringify(obj, null, 2)
  } catch {
    return String(obj)
  }
}

// 处理批准
const handleApprove = () => {
  emit('approve', props.pendingApproval)
}

// 处理拒绝
const handleReject = () => {
  emit('reject', props.pendingApproval)
}
</script>

<style scoped>
.approval-panel {
  margin-top: 16px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.approval-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.approval-icon {
  padding: 8px;
  background: #fef3c7;
  border-radius: 50%;
  color: #d97706;
  flex-shrink: 0;
}

.approval-info {
  flex: 1;
}

.approval-title {
  font-size: 14px;
  font-weight: 700;
  color: #92400e;
  margin: 0 0 4px 0;
}

.approval-description {
  font-size: 12px;
  color: #b45309;
  margin: 0;
  line-height: 1.5;
}

.tool-name-badge {
  font-family: ui-monospace, monospace;
  font-weight: 700;
  background: #fef3c7;
  padding: 1px 6px;
  border-radius: 4px;
}

.params-section {
  background: white;
  border: 1px solid #fde68a;
  border-radius: 8px;
  margin-bottom: 16px;
  overflow: hidden;
}

.params-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #fefce8;
  border-bottom: 1px solid #fef3c7;
  font-size: 12px;
  font-weight: 600;
  color: #92400e;
}

.params-tool-name {
  font-family: ui-monospace, monospace;
  opacity: 0.7;
}

.params-content {
  padding: 12px;
  font-family: ui-monospace, monospace;
  font-size: 12px;
  color: #475569;
  margin: 0;
  overflow-x: auto;
  background: #fafaf9;
  white-space: pre-wrap;
  word-break: break-word;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-approve {
  background: #059669;
  color: white;
  box-shadow: 0 2px 4px rgba(5, 150, 105, 0.3);
}

.btn-approve:hover:not(:disabled) {
  background: #047857;
}

.btn-approve:active:not(:disabled) {
  transform: scale(0.98);
}

.btn-reject {
  background: white;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.btn-reject:hover:not(:disabled) {
  background: #fef2f2;
  border-color: #fca5a5;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
