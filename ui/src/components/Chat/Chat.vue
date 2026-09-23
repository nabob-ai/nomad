<template>
  <div class="nomad-chat">
    <!-- Header -->
    <div v-if="config.showHeader" class="chat-header">
      <div class="header-content">
        <div class="header-left">
          <div class="header-avatar">
            <svg class="avatar-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
              />
            </svg>
          </div>
          <div class="header-info">
            <h3 class="header-title">{{ config.title || "AI Agent" }}</h3>
            <p v-if="config.subtitle" class="header-subtitle">{{ config.subtitle }}</p>
          </div>
        </div>

        <div class="header-right">
          <div v-if="isConnected" class="status-indicator connected">
            <span class="status-dot"></span>
            <span class="status-text">在线</span>
          </div>
          <div v-else class="status-indicator disconnected">
            <span class="status-dot"></span>
            <span class="status-text">离线</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Messages -->
    <MessageList
      ref="messageListRef"
      :messages="messages"
      :is-typing="isTyping"
      :config="config"
      class="chat-messages"
    />

    <!-- Quick Replies -->
    <QuickReplies
      v-if="quickReplies.length > 0"
      :replies="quickReplies"
      class="chat-quick-replies"
      @select="handleQuickReply"
    />

    <!-- Composer -->
    <Composer
      v-model="inputText"
      :placeholder="config.placeholder"
      :disabled="!isConnected || isTyping"
      :enable-voice="config.enableVoice"
      :enable-image="config.enableImage"
      class="chat-composer"
      @send="handleSend"
      @voice="handleVoice"
      @image="handleImage"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from "vue";
import MessageList from "./MessageList.vue";
import QuickReplies from "./QuickReplies.vue";
import Composer from "./Composer.vue";
import { useChat } from "@/composables/useChat";
import type { ChatConfig, Message, QuickReply } from "@/types";

const props = withDefaults(defineProps<{
  config?: ChatConfig;
}>(), {
  config: () => ({
    showHeader: true,
    placeholder: "输入消息...",
    enableVoice: false,
    enableImage: false,
  }),
});

const emit = defineEmits<{
  send: [message: Message];
  receive: [message: Message];
  error: [error: Error];
}>();

// Use chat composable
const { messages, isTyping, isConnected, sendMessage, sendImage } = useChat(props.config);

const inputText = ref("");
const messageListRef = ref<{ scrollToBottom: () => void } | null>(null);
const quickReplies = ref<QuickReply[]>(props.config.quickReplies || []);

// Handle send message
const handleSend = async () => {
  if (!inputText.value.trim()) return;

  const text = inputText.value;
  inputText.value = "";

  try {
    await sendMessage(text);
    if (messages.value.length === 0) {
      quickReplies.value = [];
    }
  } catch (error) {
    console.error("Send message error:", error);
    emit("error", error as Error);
  }
};

// Handle quick reply
const handleQuickReply = async (reply: QuickReply) => {
  inputText.value = reply.text;
  await handleSend();
};

// Handle voice input
const handleVoice = (blob: Blob) => {
  console.log("Voice input:", blob);
};

// Handle image upload
const handleImage = async (file: File) => {
  try {
    await sendImage(file);
  } catch (error) {
    console.error("Send image error:", error);
    emit("error", error as Error);
  }
};

// Auto scroll to bottom
watch(
  messages,
  () => {
    messageListRef.value?.scrollToBottom();
  },
  { deep: true },
);

// Initialize
onMounted(() => {
  // Welcome message will be added by useChat
});
</script>

<style scoped>
.nomad-chat {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: white;
  font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

/* Header */
.chat-header {
  flex-shrink: 0;
  padding: 12px 16px;
  border-bottom: 1px solid #f3f4f6;
  background: white;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-icon {
  width: 20px;
  height: 20px;
  color: #3b82f6;
}

.header-info {
  display: flex;
  flex-direction: column;
}

.header-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin: 0;
}

.header-subtitle {
  font-size: 12px;
  color: #6b7280;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.connected .status-dot {
  background: #10b981;
}

.status-indicator.disconnected .status-dot {
  background: #ef4444;
  animation: pulse 1.5s ease-in-out infinite;
}

.status-text {
  color: #6b7280;
}

/* Messages */
.chat-messages {
  flex: 1;
  overflow-y: auto;
}

/* Quick Replies */
.chat-quick-replies {
  flex-shrink: 0;
}

/* Composer */
.chat-composer {
  flex-shrink: 0;
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
