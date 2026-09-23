<template>
  <div class="card-message">
    <!-- Card Header -->
    <div v-if="message.content.title || message.content.subtitle" class="card-header">
      <h3 v-if="message.content.title" class="card-title">
        {{ message.content.title }}
      </h3>
      <p v-if="message.content.subtitle" class="card-subtitle">
        {{ message.content.subtitle }}
      </p>
    </div>

    <!-- Card Image -->
    <div v-if="message.content.image" class="card-image">
      <img :src="message.content.image.url" :alt="message.content.image.alt || '卡片图片'" loading="lazy" />
    </div>

    <!-- Card Body -->
    <div v-if="message.content.description" class="card-body">
      <p class="card-description">{{ message.content.description }}</p>
    </div>

    <!-- Card Fields -->
    <div v-if="message.content.fields && message.content.fields.length > 0" class="card-fields">
      <div v-for="(field, index) in message.content.fields" :key="index" class="card-field">
        <span class="field-label">{{ field.label }}</span>
        <span class="field-value">{{ field.value }}</span>
      </div>
    </div>

    <!-- Card Actions -->
    <div v-if="message.content.actions && message.content.actions.length > 0" class="card-actions">
      <button v-for="(action, index) in message.content.actions" :key="index" @click="handleAction(action)" :class="['card-action-btn', action.style === 'primary' ? 'btn-primary' : 'btn-secondary']">
        <svg v-if="action.icon" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getIconPath(action.icon)"></path>
        </svg>
        {{ action.label }}
      </button>
    </div>

    <!-- Card Footer -->
    <div v-if="message.content.footer" class="card-footer">
      <span class="footer-text">{{ message.content.footer }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CardMessage as CardMessageType } from "@/types";

defineProps<{
  message: CardMessageType;
}>();

const emit = defineEmits<{
  action: [action: unknown];
}>();

function handleAction(action: unknown) {
  emit("action", action);
}

function getIconPath(icon: string): string {
  const icons: Record<string, string> = {
    link: "M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14",
    download: "M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4",
    share: "M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z",
  };
  return icons[icon] || "";
}
</script>

<style scoped>
.card-message {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  max-width: 400px;
}

.card-header {
  padding: 16px 16px 8px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 4px;
}

.card-subtitle {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.card-image {
  width: 100%;
}

.card-image img {
  width: 100%;
  height: auto;
  object-fit: cover;
  max-height: 200px;
}

.card-body {
  padding: 12px 16px;
}

.card-description {
  font-size: 14px;
  color: #374151;
  line-height: 1.6;
  margin: 0;
}

.card-fields {
  padding: 8px 16px;
  border-top: 1px solid #f3f4f6;
}

.card-field {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  padding: 4px 0;
}

.field-label {
  color: #6b7280;
  font-weight: 500;
}

.field-value {
  color: #111827;
}

.card-actions {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid #f3f4f6;
}

.card-action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  border: none;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: #f9fafb;
  color: #374151;
  border: 1px solid #e5e7eb;
}

.btn-secondary:hover {
  background: #f3f4f6;
}

.card-footer {
  padding: 8px 16px;
  background: #f9fafb;
  border-top: 1px solid #f3f4f6;
}

.footer-text {
  font-size: 12px;
  color: #6b7280;
}
</style>
