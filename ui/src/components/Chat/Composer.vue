<template>
  <div class="composer">
    <form class="composer-form" @submit.prevent="handleSend">
      <!-- Textarea -->
      <textarea
        ref="textareaRef"
        v-model="localValue"
        :placeholder="placeholder"
        :disabled="disabled"
        class="composer-textarea"
        rows="1"
        @keydown="handleKeyDown"
        @input="handleInput"
      />

      <!-- Bottom Toolbar -->
      <div class="composer-toolbar">
        <!-- Left Actions -->
        <div class="toolbar-left">
          <button
            v-if="enableImage"
            type="button"
            class="toolbar-btn"
            title="上传图片"
            @click="handleImageClick"
          >
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
            </svg>
          </button>

          <button
            v-if="enableVoice"
            type="button"
            :class="['toolbar-btn', isRecording && 'recording']"
            title="语音输入"
            @click="toggleVoice"
          >
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"/>
            </svg>
          </button>

          <!-- Character Count -->
          <span v-if="showCharCount && maxLength" class="char-count">
            {{ localValue.length }} / {{ maxLength }}
          </span>
        </div>

        <!-- Right Actions -->
        <div class="toolbar-right">
          <button
            type="submit"
            :disabled="!canSend"
            class="send-btn"
            title="发送 (Enter)"
          >
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
            </svg>
          </button>
        </div>
      </div>
    </form>

    <!-- Hidden File Input -->
    <input
      ref="fileInputRef"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleFileChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from "vue";

const props = withDefaults(defineProps<{
  modelValue: string;
  placeholder?: string;
  disabled?: boolean;
  enableVoice?: boolean;
  enableImage?: boolean;
  maxLength?: number;
  showCharCount?: boolean;
}>(), {
  placeholder: "输入消息...",
  disabled: false,
  enableVoice: false,
  enableImage: false,
  showCharCount: false,
});

const emit = defineEmits<{
  "update:modelValue": [value: string];
  send: [];
  voice: [blob: Blob];
  image: [file: File];
}>();

const textareaRef = ref<HTMLTextAreaElement>();
const fileInputRef = ref<HTMLInputElement>();
const localValue = ref(props.modelValue);
const isRecording = ref(false);

// Can send if has content and not disabled
const canSend = computed(() => {
  return localValue.value.trim().length > 0 && !props.disabled;
});

// Handle input
function handleInput() {
  emit("update:modelValue", localValue.value);
  adjustTextareaHeight();
}

// Handle key down
function handleKeyDown(e: KeyboardEvent) {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault();
    if (canSend.value) {
      handleSend();
    }
  }
}

// Handle send
function handleSend() {
  if (!canSend.value) return;
  emit("send");
}

// Adjust textarea height
function adjustTextareaHeight() {
  if (!textareaRef.value) return;
  textareaRef.value.style.height = "auto";
  const newHeight = Math.min(textareaRef.value.scrollHeight, 200);
  textareaRef.value.style.height = `${newHeight}px`;
}

// Handle image upload
function handleImageClick() {
  fileInputRef.value?.click();
}

function handleFileChange(e: Event) {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];
  if (file) {
    emit("image", file);
    target.value = "";
  }
}

// Handle voice input
function toggleVoice() {
  if (isRecording.value) {
    stopRecording();
  } else {
    startRecording();
  }
}

function startRecording() {
  isRecording.value = true;
  console.log("Start recording...");
}

function stopRecording() {
  isRecording.value = false;
  console.log("Stop recording...");
}

// Watch modelValue changes from parent
watch(
  () => props.modelValue,
  (newValue) => {
    localValue.value = newValue;
    nextTick(adjustTextareaHeight);
  },
);

// Focus input
function focus() {
  textareaRef.value?.focus();
}

defineExpose({ focus });
</script>

<style scoped>
.composer {
  padding: 16px;
  background: white;
  padding-bottom: 24px;
}

.composer-form {
  position: relative;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  background: white;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.2s;
}

.composer-form:focus-within {
  border-color: #d1d5db;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.12);
}

.composer-textarea {
  width: 100%;
  min-height: 80px;
  max-height: 200px;
  padding: 12px 16px;
  border: none;
  background: transparent;
  font-size: 15px;
  line-height: 1.6;
  color: #111827;
  resize: none;
  outline: none;
}

.composer-textarea::placeholder {
  color: #9ca3af;
}

.composer-textarea:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.composer-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-top: 1px solid #f3f4f6;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-btn {
  padding: 8px;
  border: none;
  background: transparent;
  color: #9ca3af;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.15s;
}

.toolbar-btn:hover {
  background: #f3f4f6;
  color: #4b5563;
}

.toolbar-btn.recording {
  color: #ef4444;
  background: #fef2f2;
  animation: pulse 1.5s ease-in-out infinite;
}

.icon {
  width: 20px;
  height: 20px;
}

.char-count {
  font-size: 12px;
  color: #9ca3af;
}

.send-btn {
  padding: 8px;
  border: none;
  background: #111827;
  color: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
}

.send-btn:hover:not(:disabled) {
  background: #000;
}

.send-btn:disabled {
  background: #e5e7eb;
  color: #9ca3af;
  cursor: not-allowed;
}

.send-btn .icon {
  width: 16px;
  height: 16px;
}

.hidden {
  display: none;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
