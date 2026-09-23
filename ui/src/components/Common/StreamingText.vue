<template>
  <div class="streaming-text">
    <div class="markdown-content" v-html="renderedContent"></div>
    <span v-if="isStreaming && displayedContent.length < content.length" class="cursor">▊</span>
  </div>
</template>

<script setup lang="ts">
/**
 * StreamingText - 流式文字显示组件
 * 实现打字机效果，每次显示几个字符，营造快速流畅的感觉
 * 支持 Markdown 渲染
 */
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
// @ts-ignore - markdown-it 可能没有类型声明
import MarkdownIt from 'markdown-it'

// 初始化 markdown-it
const md = new (MarkdownIt as any)({
  html: true,
  breaks: true,
  linkify: true,
  typographer: true
})

interface Props {
  content: string
  isStreaming?: boolean
  chunkSize?: number
  interval?: number
}

const props = withDefaults(defineProps<Props>(), {
  isStreaming: false,
  chunkSize: 4,
  interval: 30
})

const displayedContent = ref('')

const renderedContent = computed(() => {
  if (!displayedContent.value) return ''
  return md.render(displayedContent.value)
})

let animationFrameId: number | null = null
let lastUpdateTime = 0

const streamText = (targetContent: string) => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }

  const animate = (currentTime: number) => {
    if (!lastUpdateTime) {
      lastUpdateTime = currentTime
    }

    const elapsed = currentTime - lastUpdateTime

    if (elapsed >= props.interval) {
      if (displayedContent.value.length < targetContent.length) {
        const nextLength = Math.min(
          displayedContent.value.length + props.chunkSize,
          targetContent.length
        )
        displayedContent.value = targetContent.substring(0, nextLength)
        lastUpdateTime = currentTime
      } else {
        displayedContent.value = targetContent
        if (animationFrameId) {
          cancelAnimationFrame(animationFrameId)
          animationFrameId = null
        }
        return
      }
    }

    animationFrameId = requestAnimationFrame(animate)
  }

  animationFrameId = requestAnimationFrame(animate)
}

watch(() => props.content, (newContent) => {
  if (props.isStreaming) {
    if (displayedContent.value.length < newContent.length) {
      streamText(newContent)
    } else {
      displayedContent.value = newContent
    }
  } else {
    displayedContent.value = newContent
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId)
      animationFrameId = null
    }
  }
}, { immediate: true })

watch(() => props.isStreaming, (isStreaming) => {
  if (isStreaming && displayedContent.value.length < props.content.length) {
    streamText(props.content)
  } else if (!isStreaming) {
    displayedContent.value = props.content
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId)
      animationFrameId = null
    }
  }
})

onMounted(() => {
  if (props.isStreaming && props.content) {
    streamText(props.content)
  } else {
    displayedContent.value = props.content
  }
})

onUnmounted(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
})
</script>

<style scoped>
.streaming-text {
  word-break: break-word;
}

.markdown-content {
  line-height: 1.7;
}

.markdown-content :deep(h1) {
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
  margin: 1rem 0 0.5rem;
}

.markdown-content :deep(h2) {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0.75rem 0 0.5rem;
}

.markdown-content :deep(h3) {
  font-size: 1.125rem;
  font-weight: 600;
  color: #374151;
  margin: 0.5rem 0 0.25rem;
}

.markdown-content :deep(p) {
  margin: 0.5rem 0;
  color: #374151;
}

.markdown-content :deep(strong) {
  font-weight: 600;
  color: #111827;
}

.markdown-content :deep(em) {
  font-style: italic;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin: 0.5rem 0;
  padding-left: 1.25rem;
}

.markdown-content :deep(li) {
  margin: 0.25rem 0;
  color: #374151;
}

.markdown-content :deep(ul > li) {
  list-style-type: disc;
}

.markdown-content :deep(ol > li) {
  list-style-type: decimal;
}

.markdown-content :deep(blockquote) {
  border-left: 3px solid #e5e7eb;
  padding-left: 0.75rem;
  margin: 0.5rem 0;
  color: #6b7280;
  font-style: italic;
}

.markdown-content :deep(code) {
  background: #f3f4f6;
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
  font-size: 0.875em;
  font-family: ui-monospace, SFMono-Regular, monospace;
  color: #1f2937;
}

.markdown-content :deep(pre) {
  background: #1f2937;
  color: #f9fafb;
  padding: 0.75rem;
  border-radius: 0.375rem;
  overflow-x: auto;
  margin: 0.5rem 0;
  font-size: 0.875rem;
}

.markdown-content :deep(pre code) {
  background: transparent;
  padding: 0;
  color: inherit;
}

.markdown-content :deep(a) {
  color: #2563eb;
  text-decoration: none;
}

.markdown-content :deep(a:hover) {
  text-decoration: underline;
}

.markdown-content :deep(hr) {
  border: none;
  border-top: 1px solid #e5e7eb;
  margin: 1rem 0;
}

.cursor {
  display: inline-block;
  width: 0.6em;
  height: 1em;
  background: #3b82f6;
  margin-left: 2px;
  animation: blink 1s ease-in-out infinite;
  vertical-align: baseline;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}
</style>
