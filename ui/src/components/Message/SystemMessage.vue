<template>
  <div :class="['system-message', `type-${message.content.type || 'info'}`]">
    <div class="system-icon">
      <svg v-if="message.content.type === 'success'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <svg v-else-if="message.content.type === 'error'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <svg v-else-if="message.content.type === 'warning'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
      </svg>
      <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
    </div>

    <div class="system-content">
      <p class="system-text">{{ message.content.text }}</p>
      <span v-if="message.createdAt" class="system-time">
        {{ formatTime(message.createdAt) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SystemMessage as SystemMessageType } from "@/types";

defineProps<{
  message: SystemMessageType;
}>();

function formatTime(timestamp: number): string {
  const date = new Date(timestamp);
  return date.toLocaleTimeString("zh-CN", {
    hour: "2-digit",
    minute: "2-digit",
  });
}
</script>

<style scoped>
.system-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 14px;
  margin: 16px auto;
  max-width: 400px;
}

.type-info {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
}

.type-success {
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
}

.type-warning {
  background: #fffbeb;
  border: 1px solid #fde68a;
}

.type-error {
  background: #fef2f2;
  border: 1px solid #fecaca;
}

.system-icon {
  flex-shrink: 0;
}

.system-icon svg {
  width: 16px;
  height: 16px;
}

.type-info .system-icon {
  color: #3b82f6;
}

.type-success .system-icon {
  color: #10b981;
}

.type-warning .system-icon {
  color: #f59e0b;
}

.type-error .system-icon {
  color: #ef4444;
}

.system-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.system-text {
  margin: 0;
}

.type-info .system-text {
  color: #1e40af;
}

.type-success .system-text {
  color: #065f46;
}

.type-warning .system-text {
  color: #92400e;
}

.type-error .system-text {
  color: #991b1b;
}

.system-time {
  font-size: 12px;
  color: #6b7280;
}
</style>
