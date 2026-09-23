<template>
  <div v-if="replies.length > 0" class="quick-replies">
    <div class="quick-replies-inner">
      <div class="quick-replies-scroll">
        <button v-for="(reply, index) in replies" :key="index" @click="handleSelect(reply)" :class="['quick-reply-button', reply.isHighlight && 'highlight']">
          <svg v-if="reply.icon" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="reply.icon"></path>
          </svg>
          <span>{{ reply.text }}</span>
          <span v-if="reply.isNew" class="new-badge">NEW</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { QuickReply } from "@/types";

defineProps<{
  replies: QuickReply[];
}>();

const emit = defineEmits<{
  select: [reply: QuickReply];
}>();

function handleSelect(reply: QuickReply) {
  emit("select", reply);
}
</script>

<style scoped>
.quick-replies {
  border-top: 1px solid #f3f4f6;
  background: white;
}

.quick-replies-inner {
  max-width: 800px;
  margin: 0 auto;
  padding: 12px 16px;
}

.quick-replies-scroll {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
  scrollbar-width: thin;
}

.quick-replies-scroll::-webkit-scrollbar {
  height: 4px;
}

.quick-replies-scroll::-webkit-scrollbar-track {
  background: transparent;
}

.quick-replies-scroll::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 9999px;
}

.quick-reply-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #f9fafb;
  border: 1px solid #f3f4f6;
  border-radius: 9999px;
  font-size: 13px;
  font-weight: 500;
  color: #4b5563;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.15s;
  flex-shrink: 0;
}

.quick-reply-button:hover {
  background: #f3f4f6;
  border-color: #e5e7eb;
  color: #111827;
}

.quick-reply-button:active {
  transform: scale(0.95);
}

.quick-reply-button.highlight {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.quick-reply-button.highlight:hover {
  background: #2563eb;
  border-color: #2563eb;
}

.new-badge {
  margin-left: 4px;
  padding: 2px 6px;
  background: #ef4444;
  color: white;
  font-size: 10px;
  font-weight: 700;
  border-radius: 4px;
}
</style>
